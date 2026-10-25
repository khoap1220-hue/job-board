package handler

import (
	"bytes"
	"net/http"
	"testing"
)

func TestValidateOptionalLinkedinURL(t *testing.T) {
	if !ValidateOptionalLinkedinURL("") {
		t.Fatalf("expected empty LinkedIn URL to be valid (optional)")
	}
	if !ValidateOptionalLinkedinURL("https://www.linkedin.com/in/johndoe") {
		t.Fatalf("expected valid www.linkedin.com URL to pass")
	}
	if !ValidateOptionalLinkedinURL("https://linkedin.com/in/johndoe") {
		t.Fatalf("expected valid linkedin.com URL to pass")
	}
	if ValidateOptionalLinkedinURL("https://github.com/johndoe") {
		t.Fatalf("expected non-linkedin URL to fail")
	}
}

func TestValidateDeveloperCVBytes(t *testing.T) {
	// 1. Empty payload
	if _, status := ValidateDeveloperCVBytes(nil, MaxDeveloperCVSize); status != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty CV, got %d", status)
	}

	// 2. Non-PDF payload
	if _, status := ValidateDeveloperCVBytes([]byte("not a pdf file"), MaxDeveloperCVSize); status != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415 for non-PDF CV, got %d", status)
	}

	// 3. Valid PDF header payload
	validPDF := append([]byte("%PDF-1.7\n"), bytes.Repeat([]byte("A"), 512)...)
	if errMsg, status := ValidateDeveloperCVBytes(validPDF, MaxDeveloperCVSize); status != http.StatusOK || errMsg != "" {
		t.Fatalf("expected 200 OK for valid PDF, got %d (%s)", status, errMsg)
	}

	// 4. Oversized PDF payload
	if _, status := ValidateDeveloperCVBytes(validPDF, 100); status != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413 for oversized PDF, got %d", status)
	}
}
