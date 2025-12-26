package main

// Запуск тестов для исследования производительности FS и MinIO

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

const (
	fileSize       = 10 // bytes
	minFiles       = 1
	maxFiles       = 3001
	step           = 250
	testIterations = 4
)

var (
	numGoroutines = runtime.NumCPU()
)

var configPath = ".env"

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

	// Create results directory if it doesn't exist
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		log.Fatalln("Failed to create results directory:", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	resultsFile, err := os.Create(fmt.Sprintf("%s/result_%s.txt", resultsDir, timestamp))
	if err != nil {
		log.Fatalln("Failed to create results file:", err)
	}
	defer resultsFile.Close()

	fmt.Fprintf(resultsFile, "=== Performance Test Results (%s) ===\n\n", timestamp)

	fmt.Fprintln(resultsFile, "=== Configuration Settings ===")
	fmt.Fprintf(resultsFile, "Min Files: %d\n", minFiles)
	fmt.Fprintf(resultsFile, "Max Files: %d\n", maxFiles)
	fmt.Fprintf(resultsFile, "Step: %d\n", step)
	fmt.Fprintf(resultsFile, "Test Iterations: %d\n", testIterations)
	fmt.Fprintf(resultsFile, "Number of Goroutines: %d\n", numGoroutines)

	srcData := make([]byte, fileSize)
	_, err = rand.Read(srcData)
	if err != nil {
		log.Fatalln("Failed to read source file:", err)
	}
	fileSize := len(srcData)
	var sizeStr string
	if fileSize < 1024 {
		sizeStr = fmt.Sprintf("%d bytes", fileSize)
	} else if fileSize < 1024*1024 {
		sizeStr = fmt.Sprintf("%.2f kb", float64(fileSize)/1024)
	} else {
		sizeStr = fmt.Sprintf("%.2f mb", float64(fileSize)/(1024*1024))
	}
	fmt.Fprintf(resultsFile, "File size = %s\n\n\n", sizeStr)

	// Write table header
	header := "N | FS UPLOAD | MinIO UPLOAD | upload ratio | FS READ | MinIO READ | read ratio |"
	separator := "---+-----------+--------------+--------------+---------+------------+------------+"
	fmt.Fprintln(resultsFile, header)
	fmt.Fprintln(resultsFile, separator)

	fileCounts := []int{}
	for i := minFiles; i <= maxFiles; i += step {
		fileCounts = append(fileCounts, i)
	}

	for _, n := range fileCounts {
		var totalFSUpload, totalMinIOUpload, totalFSRead, totalMinIORead time.Duration
		var trackIDs1, trackIDs3 []uuid.UUID

		for i := range testIterations {
			runtime.GC()
			if err := cleanupFS(); err != nil {
				log.Fatalf("Failed to cleanup FS: %v\n", err)
			}
			if err := cleanupMinIO(ctx, minioClient, conf.MinIO.BucketAudio); err != nil {
				log.Fatalf("Failed to cleanup MinIO: %v\n", err)
			}

			fmt.Printf("\nTesting with N=%d (iteration %d/%d):\n", n, i+1, testIterations)

			fmt.Println("Testing FS upload...")
			t1, c1 := testWriteParallel(ctx, fsRepo, n, srcData)
			totalFSUpload += t1
			trackIDs1 = c1
			fmt.Printf("FS upload time: %v\n", t1)

			fmt.Println("Testing FS read...")
			t2 := testReadParallel(ctx, fsRepo, trackIDs1)
			totalFSRead += t2
			fmt.Printf("FS read time: %v\n", t2)

			fmt.Println("Testing MinIO upload...")
			t3, c3 := testWriteParallel(ctx, minioRepo, n, srcData)
			totalMinIOUpload += t3
			trackIDs3 = c3
			fmt.Printf("MinIO upload time: %v\n", t3)

			fmt.Println("Testing MinIO read...")
			t4 := testReadParallel(ctx, minioRepo, trackIDs3)
			totalMinIORead += t4
			fmt.Printf("MinIO read time: %v\n", t4)
		}

		avgFSUpload := totalFSUpload / testIterations
		avgMinIOUpload := totalMinIOUpload / testIterations
		avgFSRead := totalFSRead / testIterations
		avgMinIORead := totalMinIORead / testIterations

		uploadRatio := float64(avgMinIOUpload) / float64(avgFSUpload)
		readRatio := float64(avgMinIORead) / float64(avgFSRead)

		fmt.Fprintf(resultsFile, "%d | %9v | %12v | %11.2fx | %7v | %10v | %10.2fx |\n",
			n, avgFSUpload, avgMinIOUpload, uploadRatio, avgFSRead, avgMinIORead, readRatio)

		fmt.Printf("\nAverage Ratios - Upload: %.2fx, Read: %.2fx\n", uploadRatio, readRatio)
	}

	if err := cleanupFS(); err != nil {
		log.Fatalf("Failed to cleanup FS: %v\n", err)
	}
	if err := cleanupMinIO(ctx, minioClient, conf.MinIO.BucketAudio); err != nil {
		log.Fatalf("Failed to cleanup MinIO: %v\n", err)
	}

	fmt.Printf("\nResults have been saved to results/result_%s.txt\n", timestamp)
}

func testReadParallel(ctx context.Context, repo audio.AudioFileRepository, trackIDs []uuid.UUID) time.Duration {
	start := time.Now()

	numWorkers := numGoroutines
	trackIDChan := make(chan uuid.UUID, len(trackIDs))
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	for range numWorkers {
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

	for range numWorkers {
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
		err := client.RemoveObject(ctx, bucket, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			return fmt.Errorf("failed to remove object %s: %w", object.Key, err)
		}
	}
	return nil
}
