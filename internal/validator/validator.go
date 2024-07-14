package validator

import (
	"strings"
	"time"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.NonFieldErrors) == 0 && len(v.FieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddFieldError(key, message string) {

	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func UnavailabilityTimeRange(start, end time.Time) bool {
	return start.Before(end)
}

func UnavailabilityNotPassed(start time.Time) bool {
	currentDate := time.Now()

	return currentDate.Before(start)
}
