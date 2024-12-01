package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math/big"
	"net/http"
	"time"
)

func RegisterUser(db *gorm.DB, c *gin.Context) {
	var userDetails UserDetails

	// Bind JSON input to the UserDetails struct
	if err := c.ShouldBindJSON(&userDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		fmt.Println(err)
		return
	}

	// Save the user details in the database
	result := db.Create(&userDetails)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user details"})
		fmt.Println(result.Error)
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user":    userDetails,
	})
}

func GetOtp(db *gorm.DB, c *gin.Context) {
	// Retrieve the phone number from query parameters
	phoneNumber := c.Query("number")
	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'number' query parameter"})
		return
	}
	var user UserDetails

	// Perform the select query with the WHERE condition on mobile number
	result := db.Where("mobile = ?", phoneNumber).First(&user)

	// Check if the query was successful
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Return JSON response if the user is not found
			c.JSON(404, gin.H{"message": "User doesn't exist"})
		} else {
			// Handle other errors
			c.JSON(500, gin.H{"error": "Internal server error"})
		}
		return
	}
	// Generate a 6-digit OTP
	otp := generateOTP(6)
	userOtp := UserOtp{
		Mobile:    phoneNumber,
		OTP:       otp,
		CreatedAt: time.Now(),
	}
	result = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "mobile"}},         // Specify the unique column (Mobile)
		DoUpdates: clause.AssignmentColumns([]string{"otp"}), // Update the OTP if conflict occurs
	}).Create(&userOtp)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user otp"})
		fmt.Println(result.Error)
		return
	}
	c.JSON(http.StatusOK, gin.H{"otp": otp})
}

// generateOTP generates a random OTP of the given length.
func generateOTP(length int) string {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)
	otpNum, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%0*d", length, otpNum.Int64())
}

func LoginUser(db *gorm.DB, c *gin.Context) {
	var loginRequest LoginRequest

	// Bind JSON input to the UserDetails struct
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		fmt.Println(err)
		return
	}
	// Query the database for the matching mobile number and OTP
	var userOtp UserOtp
	result := db.Where("mobile = ? AND otp = ?", loginRequest.Mobile, loginRequest.OTP).First(&userOtp)

	// Check if the query was successful
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// If no record found, return error
			c.JSON(404, gin.H{"message": "Invalid mobile number or OTP"})
		} else {
			// Handle other errors (e.g., database connection issues)
			c.JSON(500, gin.H{"error": "Internal server error"})
		}
		return
	} else {

		// Check if the OTP is expired (more than 5 minutes)
		if time.Since(userOtp.CreatedAt).Minutes() > 5 {
			c.JSON(400, gin.H{"message": "OTP has expired"})
			return
		}
		// If OTP is correct, fetch the user details from the UserDetails table
		var userDetails UserDetails

		userResult := db.Where("mobile = ?", loginRequest.Mobile).First(&userDetails)

		// Check if user details were found
		if userResult.Error != nil {
			if errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
				c.JSON(404, gin.H{"message": "User not found"})

			} else {
				c.JSON(500, gin.H{"error": "Error fetching user details"})

			}
			return
		} else {
			// check if fingerprint matches
			if userDetails.Fingerprint != loginRequest.Fingerprint {
				c.JSON(400, gin.H{"message": "Device not recognized. Please verify your device."})
				return
			}

			// Return user details along with OTP success message
			c.JSON(200, gin.H{
				"message": "OTP verified successfully",
				"user": gin.H{
					"id":      userDetails.Id,
					"name":    userDetails.Name,
					"mobile":  userDetails.Mobile,
					"address": userDetails.Address,
				},
			})
		}
	}

}
