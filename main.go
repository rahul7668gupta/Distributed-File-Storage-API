package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rahul7668gupta/dfsa/pkg/constants"
	"github.com/rahul7668gupta/dfsa/pkg/handler"
	"github.com/rahul7668gupta/dfsa/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize database connection
	db := initDB()

	// Migrate the database
	err := db.AutoMigrate(&model.FileChunk{}, model.FileMetadata{})
	if err != nil {
		log.Fatal("Failed to migrate the database:", err)
	}

	// Set up Gin router
	router := gin.Default()

	handlerSrv := handler.NewHandler(db)

	// Define API routes
	router.POST("/upload", handlerSrv.UploadFile)
	router.GET("/files", handlerSrv.GetFiles)
	router.GET("/download/:id", handlerSrv.DownloadFile)

	// Start the server
	port := os.Getenv(constants.PORT)
	if port == "" {
		port = "8080"
	}
	log.Fatal(router.Run(":" + port))

	log.Println("Server started on port " + port)
}

func initDB() *gorm.DB {
	var err error
	dbURL := os.Getenv(constants.DATABASE_URL)
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	if err = sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	fmt.Println("Successfully connected to the database")
	return db
}
