package utils

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"cms-backend/models"
)

var (
	emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	validStatus = map[string]bool{
		"planning": true, "active": true, "on-hold": true, "completed": true,
	}
)

type ValidationError struct{ msg string }

func (e *ValidationError) Error() string { return e.msg }

func ValidateProject(p *models.Project) error {
	if utf8.RuneCountInString(strings.TrimSpace(p.Title)) < 2 {
		return &ValidationError{"title must be at least 2 characters"}
	}
	if utf8.RuneCountInString(p.Title) > 120 {
		return &ValidationError{"title must be 120 characters or fewer"}
	}
	if utf8.RuneCountInString(strings.TrimSpace(p.ClientName)) < 2 {
		return &ValidationError{"client name must be at least 2 characters"}
	}
	if p.ClientEmail != "" && !emailRegexp.MatchString(p.ClientEmail) {
		return &ValidationError{"invalid client email address"}
	}
	if p.Status != "" && !validStatus[p.Status] {
		return &ValidationError{"status must be one of: planning, active, on-hold, completed"}
	}
	if p.Progress < 0 || p.Progress > 100 {
		return &ValidationError{"progress must be between 0 and 100"}
	}
	return nil
}

func ValidateUpdate(u *models.Update) error {
	if utf8.RuneCountInString(strings.TrimSpace(u.Title)) < 2 {
		return &ValidationError{"update title must be at least 2 characters"}
	}
	if u.Progress < 0 || u.Progress > 100 {
		return &ValidationError{"progress must be between 0 and 100"}
	}
	if len(u.Images) > 20 {
		return &ValidationError{"maximum 20 images per update"}
	}
	return nil
}
