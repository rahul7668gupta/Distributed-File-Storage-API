package model

import "time"

// FileChunk represents a chunk of a file in the database
type FileChunk struct {
	FileID     string `gorm:"column:file_id;type:varchar(255);not null"`
	ChunkIndex int    `gorm:"column:chunk_index;type:int;not null"`
	ChunkData  []byte `gorm:"column:chunk_data;type:bytea;not null"`
}

// TableName sets the insert table name for this struct type
func (FileChunk) TableName() string {
	return "file_chunks"
}

type FileMetadata struct {
	FileID    string `gorm:"primaryKey"`
	Filename  string
	MimeType  string
	CreatedAt time.Time
}
