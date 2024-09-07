package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const chunkSize = 1 * 1024 * 1024 // 1MB chunks // hardcoded purposefully

func (h *Handler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Generate a unique file ID
	fileID := h.generateFileID(header.Filename)

	// Create a temporary directory to store chunks
	tempDir, err := os.MkdirTemp("", "file-chunks-")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create temporary directory"})
		return
	}
	defer os.RemoveAll(tempDir)

	// Split file into chunks
	chunks, err := h.splitFile(file, tempDir)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to split file"})
		return
	}

	// Upload chunks to database in parallel
	var wg sync.WaitGroup
	for i, chunk := range chunks {
		wg.Add(1)
		go func(i int, chunkPath string) {
			defer wg.Done()
			if err := h.uploadChunkToDB(fileID, i, chunkPath); err != nil {
				// done purposefully, should be handled better
				log.Fatalf("Failed to upload chunk %d: %v", i, err)
			}
		}(i, chunk)
	}
	wg.Wait()

	// Store file metadata
	err = h.db.Exec("INSERT INTO file_metadata (file_id, filename, mime_type) VALUES ($1, $2, $3)",
		fileID, header.Filename, header.Header.Get("Content-Type")).Error
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to store file metadata"})
		return
	}

	c.JSON(200, gin.H{"file_id": fileID})
}

func (h *Handler) generateFileID(filename string) string {
	hash := sha256.Sum256([]byte(filename + fmt.Sprint(time.Now().UnixNano())))
	return hex.EncodeToString(hash[:])
}

func (h *Handler) splitFile(file io.Reader, tempDir string) ([]string, error) {
	var chunks []string
	buffer := make([]byte, chunkSize)
	chunkIndex := 0

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		chunkFilename := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", chunkIndex))
		chunkFile, err := os.Create(chunkFilename)
		if err != nil {
			return nil, err
		}

		_, err = chunkFile.Write(buffer[:n])
		chunkFile.Close()
		if err != nil {
			return nil, err
		}

		chunks = append(chunks, chunkFilename)
		chunkIndex++
	}

	return chunks, nil
}

func (h *Handler) uploadChunkToDB(fileID string, chunkIndex int, chunkPath string) error {
	chunkData, err := os.ReadFile(chunkPath)
	if err != nil {
		return err
	}

	tx := h.db.Exec("INSERT INTO file_chunks (file_id, chunk_index, chunk_data) VALUES ($1, $2, $3)",
		fileID, chunkIndex, chunkData)
	return tx.Error
}
