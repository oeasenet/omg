package omg

import (
	"bytes"
	"dario.cat/mergo"
	"github.com/oeasenet/go-premailer/premailer"
	"github.com/oeasenet/omg/pkg/html2text"
	"github.com/oeasenet/sprig"
	"github.com/russross/blackfriday/v2"
	"html/template"
)

// Theme is an interface to implement when creating a new theme
type Theme interface {
	Name() string              // The name of the theme
	HTMLTemplate() string      // The golang template for HTML emails
	PlainTextTemplate() string // The golang template for plain text emails (can be basic HTML)
}

// TextDirection of the text in HTML email
type TextDirection string

// TDLeftToRight is the text direction from left to right (default)
const TDLeftToRight TextDirection = "ltr"

// TDRightToLeft is the text direction from right to left
const TDRightToLeft TextDirection = "rtl"

func NewOEmail() *OEmailGenerationEngine {
	return &OEmailGenerationEngine{}
}

var templateFuncs = template.FuncMap{
	"url": func(s string) template.URL {
		return template.URL(s)
	},
}

// ToHTML converts Markdown to HTML
func (c Markdown) ToHTML() template.HTML {
	return template.HTML(blackfriday.Run([]byte(c)))
}

func setDefaultEmailValues(e *Email) error {
	// DefaultEmailTheme values of an email
	defaultEmail := &Email{
		Body: &Body{
			Intros:     []string{},
			Dictionary: []*Entry{},
			Outros:     []string{},
			Signature:  "Yours truly",
			Greeting:   "Hi",
		},
	}
	// Merge the given email with default one
	// DefaultEmailTheme one overrides all zero values
	return mergo.Merge(e, defaultEmail)
}

// default values of the engine
func setDefaultEngineValues(h *OEmailGenerationEngine) error {
	defaultTextDirection := TDLeftToRight
	defaultHermes := &OEmailGenerationEngine{
		Theme:         new(DefaultEmailTheme),
		TextDirection: defaultTextDirection,
		Product: &Product{
			Name:        "OEASE OMS",
			Copyright:   "Copyright © 2024 OEASE OMS. All rights reserved.",
			TroubleText: "If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser.",
		},
	}
	// Merge the given hermes engine configuration with default one
	// DefaultEmailTheme one overrides all zero values
	err := mergo.Merge(h, defaultHermes)
	if err != nil {
		return err
	}
	if h.TextDirection != TDLeftToRight && h.TextDirection != TDRightToLeft {
		h.TextDirection = defaultTextDirection
	}
	return nil
}

func (h *OEmailGenerationEngine) generateTemplate(email *Email, tplStr string) (string, error) {

	err := setDefaultEmailValues(email)
	if err != nil {
		return "", err
	}

	// Generate the email from Golang template
	// Allow usage of simple function from sprig : https://github.com/oeasenet/sprig
	t, err := template.New("omg").Funcs(sprig.FuncMap()).Funcs(templateFuncs).Funcs(template.FuncMap{
		"safe": func(s string) template.HTML { return template.HTML(s) }, // Used for keeping comments in generated template
	}).Parse(tplStr)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = t.Execute(&b, &Template{h, email})
	if err != nil {
		return "", err
	}

	res := b.String()
	if h.DisableCSSInlining {
		return res, nil
	}

	// Inlining CSS
	prem, err := premailer.NewPremailerFromString(res, premailer.NewOptions())
	if err != nil {
		return "", err
	}
	html, err := prem.Transform()
	if err != nil {
		return "", err
	}
	return html, nil
}

// GenerateHTML generates the email body from data to an HTML Reader
// This is for modern email clients
func (h *OEmailGenerationEngine) GenerateHTML(email *Email) (string, error) {
	err := setDefaultEngineValues(h)
	if err != nil {
		return "", err
	}
	return h.generateTemplate(email, h.Theme.HTMLTemplate())
}

// GeneratePlainText generates the email body from data
// This is for old email clients
func (h *OEmailGenerationEngine) GeneratePlainText(email *Email) (string, error) {
	err := setDefaultEngineValues(h)
	if err != nil {
		return "", err
	}
	generateTemplate, err := h.generateTemplate(email, h.Theme.PlainTextTemplate())
	if err != nil {
		return "", err
	}
	return html2text.FromString(generateTemplate, html2text.Options{PrettyTables: true})
}
