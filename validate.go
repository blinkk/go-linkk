package main

import (
	"fmt"
	"reflect"
	"strings"
)

// ValidationError Errors from validation of models.
type ValidationError struct {
	Type   string            `json:"type"`
	Errors map[string]string `json:"errors"`
}

func (f ValidationError) Error() string {
	return fmt.Sprintf("Invalid values found during validation of %s", f.Type)
}

// NewValidationError Create a validation error object for an object.
func NewValidationError(m interface{}) ValidationError {
	return ValidationError{
		Type:   reflect.TypeOf(m).Name(),
		Errors: make(map[string]string),
	}
}

// Validate Account type.
func (msg *Linkk) Validate() error {
	validationErrors := NewValidationError(*msg)

	// Paths start with a /.
	if !strings.HasPrefix(msg.Path, "/") {
		validationErrors.Errors["Path"] = "Path must start with a /."
	}

	// Paths cannot end with a /.
	if strings.HasSuffix(msg.Path, "/") {
		validationErrors.Errors["Path"] = "Path cannot end with a /."
	}

	// Paths cannot be root.
	if msg.Path == "/" {
		validationErrors.Errors["Path"] = "Path cannot be root."
	}

	// Paths cannot start with a reserved prefix.
	reservedPrefixes := []string{
		"/_/", "/_ah/",
		"/css/", "/js/", "/static/",
		"/~/",
		"/favicon.ico", "/robots.txt", "/sitemap.xml",
	}
	prefix := stringPrefixInSlice(msg.Path, reservedPrefixes)
	if prefix != "" {
		validationErrors.Errors["Path"] = fmt.Sprintf("Path cannot start with a reserved prefix: %s.", prefix)
	}

	// URL needs to start with the right protocols.
	validProtocols := []string{"http://", "https://"}
	if stringPrefixInSlice(msg.URL, validProtocols) == "" {
		validationErrors.Errors["URL"] = fmt.Sprintf("Url needs to be a valid protocol: %v.", validProtocols)
	}

	if len(validationErrors.Errors) > 0 {
		return validationErrors
	}

	return nil
}
