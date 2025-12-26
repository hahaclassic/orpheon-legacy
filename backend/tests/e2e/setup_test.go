package e2e

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
)

const (
	migrationsDir = "../../db/migrations"
)

var (
	postgresConnStr string
	redisURL        string
	minIOURL        string
	appURL          string
)

var (
	pgxPool     *pgxpool.Pool
	redisClient *redis.Client
	minioClient *minio.Client

	minioAudioBucketName string = "audio"
)

func TestMain(m *testing.M) {
	var (
		err  error
		code int

		dockerPool     *dockertest.Pool
		dockerPostgres *dockertest.Resource
		dockerMinIO    *dockertest.Resource
		dockerRedis    *dockertest.Resource
		dockerApp      *dockertest.Resource
	)

	defer func() {
		isErr := false
		if r := recover(); r != nil {
			isErr = true
			fmt.Println("[ERR]", r)
		}

		teardown(dockerPool, []*dockertest.Resource{dockerPostgres, dockerMinIO, dockerRedis})

		if isErr {
			os.Exit(1)
		}
		os.Exit(code)
	}()

	dockerPool, err = dockertest.NewPool("")
	if err != nil {
		panic(fmt.Sprintf("failed to start docker: %v", err))
	}
	dockerPool.MaxWait = 2 * time.Minute

	dockerPostgres, err = setupPostgres(dockerPool)
	if err != nil {
		panic(fmt.Sprintf("failed to start postgres: %v", err))
	}

	dockerMinIO, err = setupMinIO(dockerPool)
	if err != nil {
		panic(fmt.Sprintf("failed to start minio: %v", err))
	}

	dockerRedis, err = setupRedis(dockerPool)
	if err != nil {
		panic(fmt.Sprintf("failed to start redis: %v", err))
	}

	dockerApp, err = setupApp(dockerPool)
	if err != nil {
		panic(fmt.Sprintf("failed to start app: %v", err))
	}

	code = m.Run()
}

func teardown(pool *dockertest.Pool, resources []*dockertest.Resource) {
	if pgxPool != nil {
		pgxPool.Close()
	}

	if redisClient != nil {
		err := redisClient.Close()
		if err != nil {
			slog.Error("err", "redis close error", err)
		}
	}

	for i := range resources {
		if resources[i] == nil {
			continue
		}
		if err := pool.Purge(resources[i]); err != nil {
			slog.Error("failed to purge docker resource", "err", err)
		}
	}
}

func setupPostgres(dockerPool *dockertest.Pool) (*dockertest.Resource, error) {
	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17",
		Env: []string{
			"POSTGRES_USER=testuser",
			"POSTGRES_PASSWORD=password",
			"POSTGRES_DB=testdb",
		},
	})
	if err != nil {
		return nil, err
	}

	port := resource.GetPort("5432/tcp")
	connString := fmt.Sprintf("postgres://testuser:password@localhost:%s/testdb?sslmode=disable", port)
	postgresConnStr = connString

	if err := dockerPool.Retry(func() error {
		var err error
		pgxPool, err = pgxpool.New(context.Background(), connString)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return pgxPool.Ping(ctx)
	}); err != nil {
		return nil, err
	}

	return resource, nil
}

func runMigrationsUp(dbpool *pgxpool.Pool) error {
	db, err := sql.Open("pgx", dbpool.Config().ConnString())
	if err != nil {
		return err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			slog.Error("err", "pgx close error", err)
		}
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, migrationsDir)
}

func runMigrationsDown(dbpool *pgxpool.Pool) error {
	db, err := sql.Open("pgx", dbpool.Config().ConnString())
	if err != nil {
		return err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			slog.Error("err", "pgx close error", err)
		}
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Reset(db, migrationsDir)
}

func setupRedis(dockerPool *dockertest.Pool) (*dockertest.Resource, error) {
	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7",
	})
	if err != nil {
		return nil, err
	}

	port := resource.GetPort("6379/tcp")
	addr := fmt.Sprintf("localhost:%s", port)
	redisURL = addr

	if err := dockerPool.Retry(func() error {
		redisClient = redis.NewClient(&redis.Options{
			Addr: addr,
		})
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return redisClient.Ping(ctx).Err()
	}); err != nil {
		return nil, err
	}

	return resource, nil
}

func clearRedis(rdb *redis.Client) error {
	return rdb.FlushAll(context.Background()).Err()
}

func setupMinIO(dockerPool *dockertest.Pool) (*dockertest.Resource, error) {
	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "minio/minio",
		Tag:        "latest",
		Env: []string{
			"MINIO_ROOT_USER=minioadmin",
			"MINIO_ROOT_PASSWORD=minioadmin",
		},
		Cmd: []string{"server", "/data"},
	})
	if err != nil {
		return nil, err
	}

	port := resource.GetPort("9000/tcp")
	endpoint := fmt.Sprintf("localhost:%s", port)
	minIOURL = endpoint

	if err := dockerPool.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		minioClient, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
			Secure: false,
		})
		if err != nil {
			return err
		}

		err = minioClient.MakeBucket(ctx, minioAudioBucketName, minio.MakeBucketOptions{})
		if err != nil {
			exists, errBucketExists := minioClient.BucketExists(ctx, minioAudioBucketName)
			if errBucketExists != nil || !exists {
				return fmt.Errorf("could not create bucket: %v", err)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return resource, nil
}

func clearMinioBucket(ctx context.Context, client *minio.Client, bucketName string) error {
	opts := minio.ListObjectsOptions{
		Recursive: true,
	}
	for obj := range client.ListObjects(ctx, bucketName, opts) {
		if obj.Err != nil {
			return obj.Err
		}
		err := client.RemoveObject(ctx, bucketName, obj.Key, minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func setupApp(dockerPool *dockertest.Pool) (*dockertest.Resource, error) {
	appRes, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Name:       "orpheon-app",
		Repository: "orpheon-e2e",
		Tag:        "latest",
		Env: []string{
			"DB_URL=" + postgresConnStr,
			"REDIS_URL=" + redisURL,
			"MINIO_URL=" + minIOURL,
			"APP_PORT=8080",
		},
		ExposedPorts: []string{"8080/tcp"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"8080/tcp": {{HostIP: "0.0.0.0", HostPort: "8080"}},
		},
	}, func(cfg *docker.HostConfig) {
		cfg.AutoRemove = true
		cfg.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		panic(fmt.Sprintf("failed to start app container: %v", err))
	}

	time.Sleep(5 * time.Second)

	appPort := appRes.GetPort("8080/tcp")
	appURL = fmt.Sprintf("http://localhost:%s/api/v1", appPort)

	err = dockerPool.Retry(func() error {
		resp, err := http.Get(appURL + "/albums")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("app not ready yet, status: %d", resp.StatusCode)
		}
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("app failed to become ready: %v", err))
	}

	fmt.Printf("Application is ready at %s\n", appURL)
	return appRes, nil
}
