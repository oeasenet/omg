package test

import (
	"github.com/oeasenet/omg"
	"testing"
)

func TestDefaultEmailTheme_HTMLTemplate(t *testing.T) {
	oe := omg.NewOEmail()
	emailStr, err := oe.GenerateHTML(&omg.Email{})
	if err != nil {
		t.Error(err)
	}
	if emailStr == "" {
		t.Error("Expected HTML content to be non-empty")
	}
	t.Log(emailStr)
}

func TestDefaultEmailTheme_PlainTextTemplate(t *testing.T) {
	oe := omg.NewOEmail()
	emailStr, err := oe.GeneratePlainText(&omg.Email{})
	if err != nil {
		t.Error(err)
	}
	if emailStr == "" {
		t.Error("Expected TEXT content to be non-empty")
	}
	t.Log(emailStr)
}
