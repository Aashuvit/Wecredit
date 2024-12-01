package main

import "time"

type UserDetails struct {
	Id          int    `gorm:"primaryKey" json:"id,omitempty"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Mobile      string `gorm:"type:varchar(15);unique;not null" json:"mobile"`
	Address     string `gorm:"type:varchar(255);not null" json:"address"`
	Fingerprint string `gorm:"type:varchar(255);not null"`
}

type UserOtp struct {
	ID        int       `gorm:"primaryKey"`
	Mobile    string    `gorm:"type:varchar(15);unique;not null"`
	OTP       string    `gorm:"type:varchar(6);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type LoginRequest struct {
	Mobile      string `json:"mobile"`
	OTP         string `json:"otp"`
	Fingerprint string `json:"fingerprint"`
}
