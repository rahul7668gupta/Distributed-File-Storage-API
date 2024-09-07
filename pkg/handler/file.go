package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetFiles(c *gin.Context) {
	rows, err := h.db.Raw("SELECT file_id, filename FROM file_metadata").Rows()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve files"})
		return
	}
	defer rows.Close()

	var files []gin.H
	for rows.Next() {
		var fileID, filename string
		if err := rows.Scan(&fileID, &filename); err != nil {
			c.JSON(500, gin.H{"error": "Failed to scan file data"})
			return
		}
		files = append(files, gin.H{"file_id": fileID, "filename": filename})
	}

	c.JSON(200, gin.H{"files": files})
}
