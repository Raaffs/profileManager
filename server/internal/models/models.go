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
    PasswordHash string    `json:"password_hash"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// Profiles table equivalent
type Profile struct {
    ID           int       `json:"id"`
    UserID       int       `json:"user_id"`
    FullName     string    `json:"full_name"`
    DateOfBirth  time.Time `json:"date_of_birth"`
    AadhaarNumber string   `json:"aadhaar_number"`
    UniqueID     string    `json:"unique_id"`
    Address      string    `json:"address,omitempty"` // omitempty if null
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}