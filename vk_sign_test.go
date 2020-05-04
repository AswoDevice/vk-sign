package vksign

import (
	"testing"
)

func TestParse(t *testing.T) {
	url := "your-url"
	secret := "your-secret"

	_, err := Parse(url, secret)
	if err != nil {
		t.Error(err)
	}
}
