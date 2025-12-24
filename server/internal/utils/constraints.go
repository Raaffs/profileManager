package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type ValidationError struct {
	Key     string
	Message string
}	

var(
	MIN_NAME_LENGTH=3
	MAX_NAME_LENGTH=40
	AADHAR_LENGTH=12
	PHONE_NUMBER_LENGTH=10
)

var (
	ErrNameOutofRange      	= ValidationError{"name", fmt.Sprintf("name should be between %d - %d characters")}
	ErrFieldRequired       	= ValidationError{"field", "this field cannot be empty"}
	ErrInvalidEmail        	= ValidationError{"email", "invalid email address"}
	ErrPasswordTooWeak     	= ValidationError{"password", "password is too weak, must include letters, numbers, and special characters"}
	ErrInvalidPhone        	= ValidationError{"phone", "invalid phone number"}
	ErrInvalidAadharNumber 	= ValidationError{"aadhar", "invalid aadhar number"}
	ErrInvalidDate         	= ValidationError{"date", "invalid date format"}
)

func (v *Validator) NameLength(name string,min,max int) {
	v.Check(
     	len(strings.TrimSpace(name)) >= min && len(name) < max,
		ErrNameOutofRange.Key,
		fmt.Sprintf(ErrNameOutofRange.Message,min,max),
	)
}

func (v *Validator)Aadhar(aadhar string){
	regex := regexp.MustCompile(`^[2-9]{1}[0-9]{11}$`)
	v.Check(
	 	regex.MatchString(aadhar) && len(aadhar) == AADHAR_LENGTH,
		ErrInvalidPhone.Key,
		ErrInvalidPhone.Message,
	)
}

func (v *Validator) Phone(phone string)  {
	re := regexp.MustCompile(`^\+?(\d{1,3})?[-.\s]?\(?\d{1,4}?\)?[-.\s]?\d{1,4}[-.\s]?\d{1,9}$`)
	v.Check(
		re.MatchString(phone),
		ErrInvalidAadharNumber.Key,
		ErrInvalidAadharNumber.Message,
	)
}

func (v *Validator) Date(date string) {
    // 1. Check format with your regex
    re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
    formatMatch := re.MatchString(date)
    
    if !formatMatch {
        v.Check(false, ErrInvalidDate.Key, ErrInvalidDate.Message)
        return
    }

    // 2. Check if it's a real calendar date (e.g., handles Feb 29 on non-leap years)
    _, err := time.Parse("2006-01-02", date)
    v.Check(err == nil, ErrInvalidDate.Key, "Date must be a real calendar date")
}

func (v *Validator) Mail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}
