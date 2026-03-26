package users

import (
	"fmt"
	"regexp"
	"strings"
)

var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

func normalizeUsername(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func validateCreateUserInput(input CreateUserInput) error {
	input.Username = strings.TrimSpace(input.Username)
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	input.Email = strings.TrimSpace(input.Email)

	if input.Username == "" {
		return fmt.Errorf("username is required")
	}

	if input.DisplayName == "" {
		return fmt.Errorf("display_name is required")
	}

	if input.Email == "" {
		return fmt.Errorf("email is required")
	}

	if input.Password == "" {
		return fmt.Errorf("password is required")
	}

	if len(input.Username) < 3 || len(input.Username) > 32 {
		return fmt.Errorf("username must be between 3 and 32 characters")
	}

	if !usernamePattern.MatchString(input.Username) {
		return fmt.Errorf("username contains invalid characters")
	}

	if !strings.Contains(input.Email, "@") {
		return fmt.Errorf("email must be valid")
	}

	if len(input.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	return nil
}
