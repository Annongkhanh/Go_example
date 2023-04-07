package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidSecretCode = regexp.MustCompile(`^[a-z]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {

	if length := len(value); length < minLength || length > maxLength {
		return fmt.Errorf("must contain from %d to %d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("username must only contains lowercase letters, digits and underscore")
	}
	return nil
}

func ValidatePassword(value string) error {
	if err := ValidateString(value, 6, 100); err != nil {
		return err
	}
	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 6, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

func ValidateFullname(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidFullname(value) {
		return fmt.Errorf("fullname must only contains letters and spaces")
	}
	return nil
}

func ValidateEmailId(value int64) error{
	if value <= 0 {
		return fmt.Errorf("email id must be a positive integer: %d", value)
	}
	return nil
}

func ValidateSecretCode(value string) error {
	if len(value) != 32{
		return fmt.Errorf("secret code must be 32 characters")
	}
	if !isValidSecretCode(value) {
		return fmt.Errorf("secret code must only contain lowercase letters")
	}
	return nil
}