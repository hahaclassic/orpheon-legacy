package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/config"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/track/audio"
	inframinio "github.com/hahaclassic/orpheon/backend/internal/infrastructure/minio"
	audio_fs "github.com/hahaclassic/orpheon/backend/internal/repository/content/track/audio/fs"
	audio_minio "github.com/hahaclassic/orpheon/backend/internal/repository/content/track/audio/minio"
	minio "github.com/minio/minio-go/v7"
)

const (
	fsDir      = "../audio/"
	resultsDir = "../results"
)

var (
	numGoroutines = runtime.NumCPU()
	configPath    = ".env"
)

func init() {
	flag.StringVar(&configPath, "config", ".env", "path to config file")
	flag.Parse()
}

func main() {
	conf := config.MustLoad(configPath)
	ctx := context.Background()

	minioClient, err := inframinio.NewMinioClient(conf.MinIO)
	if err != nil {
		log.Fatalln("Failed to create MinIO client:", err)
	}

	minioRepo, err := audio_minio.NewAudioFileRepository(ctx, minioClient, conf.MinIO.BucketAudio)
	if err != nil {
		log.Fatalln("Failed to create MinIO repository:", err)
	}

	fsRepo, err := audio_fs.NewAudioFileRepository(conf.AudioStorage.BasePath)
	if err != nil {
		log.Fatalln("Failed to create FS repository:", err)
	}

	// Настройки цикла по размеру файла (байты)
	const (
		minFileSize    = 1 * 1024 * 1024
		maxFileSize    = 100 * 1024 * 1024 // 100 MB
		stepFileSize   = 10 * 1024 * 1024
		numFiles       = 100 // количество файлов в каждом тесте фиксированное
		testIterations = 4
	)

	// Создаем директорию для результатов
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		log.Fatalln("Failed to create results directory:", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	resultsFile, err := os.Create(fmt.Sprintf("%s/result_size_%s.txt", resultsDir, timestamp))
	if err != nil {
		log.Fatalln("Failed to create results file:", err)
	}
	defer resultsFile.Close()

	fmt.Fprintf(resultsFile, "=== Performance Test Results by File Size (%s) ===\n\n", timestamp)
	fmt.Fprintln(resultsFile, "=== Configuration Settings ===")
	fmt.Fprintf(resultsFile, "Min File Size: %d bytes\n", minFileSize)
	fmt.Fprintf(resultsFile, "Max File Size: %d bytes\n", maxFileSize)
	fmt.Fprintf(resultsFile, "Step File Size: %d bytes\n", stepFileSize)
	fmt.Fprintf(resultsFile, "Number of Files: %d\n", numFiles)
	fmt.Fprintf(resultsFile, "Test Iterations: %d\n", testIterations)
	fmt.Fprintf(resultsFile, "Number of Goroutines: %d\n\n", numGoroutines)

	fmt.Fprintln(resultsFile, "FileSize | FS UPLOAD | MinIO UPLOAD | upload ratio | FS READ | MinIO READ | read ratio |")
	fmt.Fprintln(resultsFile, "---------+-----------+--------------+--------------+---------+------------+------------+")

	for fileSize := minFileSize; fileSize <= maxFileSize; fileSize += stepFileSize {
		srcData := make([]byte, fileSize)
		_, err := rand.Read(srcData)
		if err != nil {
			log.Fatalf("Failed to generate random data: %v\n", err)
		}

		var totalFSUpload, totalMinIOUpload, totalFSRead, totalMinIORead time.Duration
		var trackIDsFS, trackIDsMinIO []uuid.UUID

		for iter := 0; iter < testIterations; iter++ {
			runtime.GC()
			if err := cleanupFS(); err != nil {
				log.Fatalf("Failed to cleanup FS: %v\n", err)
			}
			if err := cleanupMinIO(ctx, minioClient, conf.MinIO.BucketAudio); err != nil {
				log.Fatalf("Failed to cleanup MinIO: %v\n", err)
			}

			fmt.Printf("Testing file size=%d bytes, iteration %d/%d\n", fileSize, iter+1, testIterations)

			tFSUp, idsFS := testWriteParallel(ctx, fsRepo, numFiles, srcData)
			totalFSUpload += tFSUp
			trackIDsFS = idsFS
			fmt.Printf("FS upload: %v\n", tFSUp)

			tFSRead := testReadParallel(ctx, fsRepo, trackIDsFS)
			totalFSRead += tFSRead
			fmt.Printf("FS read: %v\n", tFSRead)

			tMinIOUp, idsMinIO := testWriteParallel(ctx, minioRepo, numFiles, srcData)
			totalMinIOUpload += tMinIOUp
			trackIDsMinIO = idsMinIO
			fmt.Printf("MinIO upload: %v\n", tMinIOUp)

			tMinIORead := testReadParallel(ctx, minioRepo, trackIDsMinIO)
			totalMinIORead += tMinIORead
			fmt.Printf("MinIO read: %v\n", tMinIORead)
		}

		avgFSUpload := totalFSUpload / testIterations
		avgMinIOUpload := totalMinIOUpload / testIterations
		avgFSRead := totalFSRead / testIterations
		avgMinIORead := totalMinIORead / testIterations

		uploadRatio := float64(avgMinIOUpload) / float64(avgFSUpload)
		readRatio := float64(avgMinIORead) / float64(avgFSRead)

		fmt.Fprintf(resultsFile, "%8d | %9v | %12v | %11.2fx | %7v | %10v | %10.2fx |\n",
			fileSize, avgFSUpload, avgMinIOUpload, uploadRatio, avgFSRead, avgMinIORead, readRatio)

		fmt.Printf("Average ratios for file size %d bytes - Upload: %.2fx, Read: %.2fx\n\n", fileSize, uploadRatio, readRatio)
	}

	if err := cleanupFS(); err != nil {
		log.Fatalf("Failed to cleanup FS: %v\n", err)
	}
	if err := cleanupMinIO(ctx, minioClient, conf.MinIO.BucketAudio); err != nil {
		log.Fatalf("Failed to cleanup MinIO: %v\n", err)
	}

	fmt.Printf("Results saved to %s/result_size_%s.txt\n", resultsDir, timestamp)
}

func testReadParallel(ctx context.Context, repo audio.AudioFileRepository, trackIDs []uuid.UUID) time.Duration {
	start := time.Now()
	numWorkers := numGoroutines
	trackIDChan := make(chan uuid.UUID, len(trackIDs))
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for trackID := range trackIDChan {
				_, err := repo.GetAudioChunk(ctx, &entity.AudioChunk{
					TrackID: trackID,
					Start:   0,
					End:     math.MaxInt64,
				})
				if err != nil {
					log.Fatalf("Failed to get audio chunk: %v\n", err)
				}
			}
		}()
	}

	for _, trackID := range trackIDs {
		trackIDChan <- trackID
	}
	close(trackIDChan)

	wg.Wait()
	return time.Since(start)
}

func testWriteParallel(ctx context.Context, repo audio.AudioFileRepository, n int, fileData []byte) (time.Duration, []uuid.UUID) {
	trackIDs := make([]uuid.UUID, n)
	for i := range trackIDs {
		trackIDs[i] = uuid.New()
	}

	start := time.Now()
	numWorkers := numGoroutines
	trackIDChan := make(chan uuid.UUID, len(trackIDs))
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for trackID := range trackIDChan {
				err := repo.UploadAudioFile(ctx, &entity.AudioChunk{
					TrackID: trackID,
					Data:    fileData,
					Start:   0,
					End:     math.MaxInt64,
				})
				if err != nil {
					log.Fatalf("Failed to upload audio file: %v\n", err)
				}
			}
		}()
	}

	for _, trackID := range trackIDs {
		trackIDChan <- trackID
	}
	close(trackIDChan)

	wg.Wait()
	return time.Since(start), trackIDs
}

func cleanupFS() error {
	if err := os.RemoveAll(fsDir); err != nil {
		return fmt.Errorf("failed to remove FS directory: %w", err)
	}
	if err := os.MkdirAll(fsDir, 0755); err != nil {
		return fmt.Errorf("failed to recreate FS directory: %w", err)
	}
	return nil
}

func cleanupMinIO(ctx context.Context, client *minio.Client, bucket string) error {
	objectsCh := client.ListObjects(ctx, bucket, minio.ListObjectsOptions{})
	for object := range objectsCh {
		if object.Err != nil {
			return fmt.Errorf("failed to list objects: %w", object.Err)
		}
		if err := client.RemoveObject(ctx, bucket, object.Key, minio.RemoveObjectOptions{}); err != nil {
			return fmt.Errorf("failed to remove object %s: %w", object.Key, err)
		}
	}
	return nil
}
