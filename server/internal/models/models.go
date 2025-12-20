package models

import (
	"errors"
	"time"
)

var(
    NotFound = errors.New("record not found")
)

type User struct {
    ID           int       `json:"id"`
    Email        string    `json:"email"`
    Username     string    `json:"username"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// Profiles table equivalent
type Profile struct {
    ID            int       `json:"id"`
    UserID        int       `json:"user_id"`
    FullName      string    `json:"full_name"`
    DateOfBirth   time.Time `json:"date_of_birth"`
    AadhaarNumber string    `json:"aadhaar_number"`
    UniqueID      string    `json:"unique_id"`    
    PhoneNumber   string    `json:"phone_number"`
    Address       string    `json:"address"`     
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}