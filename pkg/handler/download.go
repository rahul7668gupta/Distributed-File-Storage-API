package handler

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rahul7668gupta/dfsa/pkg/model"
	"gorm.io/gorm"
)

func (h *Handler) DownloadFile(c *gin.Context) {
	fileID := c.Param("id")

	// Get file metadata
	var metadata model.FileMetadata
	if err := h.db.Where("file_id = ?", fileID).First(&metadata).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "File not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to retrieve file metadata"})
		}
		return
	}

	// Get the total number of chunks
	var totalChunks int64
	if err := h.db.Model(&model.FileChunk{}).Where("file_id = ?", fileID).Count(&totalChunks).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to count file chunks"})
		return
	}

	// Prepare a buffer to store the merged file
	var buffer bytes.Buffer

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(int(totalChunks))

	// Create a channel to receive chunks in order
	chunks := make(chan []byte, int(totalChunks))

	// Fetch chunks in parallel
	for i := 0; i < int(totalChunks); i++ {
		go func(chunkIndex int) {
			defer wg.Done()
			var chunk model.FileChunk
			if err := h.db.Where("file_id = ? AND chunk_index = ?", fileID, chunkIndex).First(&chunk).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					fmt.Printf("Failed to retrieve chunk %d: %v\n", chunkIndex, err)
				}
				return
			}
			chunks <- chunk.ChunkData
		}(i)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(chunks)
	}()

	// Merge chunks in order
	for chunk := range chunks {
		if _, err := buffer.Write(chunk); err != nil {
			c.JSON(500, gin.H{"error": "Failed to merge file chunks"})
			return
		}
	}

	// Set the appropriate headers for file download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", metadata.Filename))
	c.Header("Content-Type", metadata.MimeType)

	// Stream the file to the client
	if _, err := io.Copy(c.Writer, &buffer); err != nil {
		c.JSON(500, gin.H{"error": "Failed to stream file"})
		return
	}
}
