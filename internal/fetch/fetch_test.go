package fetch

import "testing"

func TestValidateURL(t *testing.T) {
	if _, err := ValidateURL("https://example.com"); err != nil {
		t.Fatalf("valid URL returned error: %v", err)
	}
	if _, err := ValidateURL("example.com"); err == nil {
		t.Fatal("expected invalid URL error")
	}
	if _, err := ValidateURL("ftp://example.com"); err == nil {
		t.Fatal("expected unsupported scheme error")
	}
}
