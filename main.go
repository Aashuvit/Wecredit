package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	// Replace with your MySQL DSN
	dsn := "root:1234@tcp(127.0.0.1:3306)/wecredit?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automatically migrate schema
	err = db.AutoMigrate(&UserDetails{}, &UserOtp{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully!")

	// Initialize Gin router
	r := gin.Default()

	// Define the Register API endpoint
	r.POST("/register", func(c *gin.Context) {
		RegisterUser(db, c)
	})

	r.GET("/otp", func(c *gin.Context) {
		GetOtp(db, c)
	})

	r.POST("/login", func(c *gin.Context) {
		LoginUser(db, c)
	})

	// Start the server
	r.Run(":8080") // Runs on localhost:8080
	fmt.Println("listening")
}
