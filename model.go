package omg

import "html/template"

// OEmailGenerationEngine is an instance of the OEmailGenerationEngine email generator
type OEmailGenerationEngine struct {
	Theme              Theme
	TextDirection      TextDirection
	Product            *Product
	DisableCSSInlining bool
}

// Product represents your company product (brand)
// Appears in header & footer of e-mails
type Product struct {
	Name        string
	Link        string // e.g. https://matcornic.github.io
	Logo        string // e.g. https://matcornic.github.io/img/logo.png
	Copyright   string // Copyright © 2024 OEmailGenerationEngine. All rights reserved.
	TroubleText string // TroubleText is the sentence at the end of the email for users having trouble with the button (default to `If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser.`)
}

// Email is the email containing a body
type Email struct {
	Body *Body
}

// Markdown is a HTML template (a string) representing Markdown content
// https://en.wikipedia.org/wiki/Markdown
type Markdown template.HTML

// Body is the body of the email, containing all interesting data
type Body struct {
	Name         string    // The name of the contacted person
	Intros       []string  // Intro sentences, first displayed in the email
	Dictionary   []*Entry  // A list of key+value (useful for displaying parameters/settings/personal info)
	Table        *Table    // Table is an table where you can put data (pricing grid, a bill, and so on)
	Actions      []*Action // Actions are a list of actions that the user will be able to execute via a button click
	Outros       []string  // Outro sentences, last displayed in the email
	Greeting     string    // Greeting for the contacted person (default to 'Hi')
	Signature    string    // Signature for the contacted person (default to 'Yours truly')
	Title        string    // Title replaces the greeting+name when set
	FreeMarkdown Markdown  // Free markdown content that replaces all content other than header and footer
}

// Entry is a simple entry of a map
// Allows using a slice of entries instead of a map
// Because Golang maps are not ordered
type Entry struct {
	Key   string
	Value string
}

// Table is an table where you can put data (pricing grid, a bill, and so on)
type Table struct {
	Data    [][]*Entry // Contains data
	Columns *Columns   // Contains meta-data for display purpose (width, alignement)
}

// Columns contains meta-data for the different columns
type Columns struct {
	CustomWidth     map[string]string
	CustomAlignment map[string]string
}

// Action is anything the user can act on (i.e., click on a button, view an invite code)
type Action struct {
	Instructions string
	Button       *Button
	InviteCode   string
}

// Button defines an action to launch
type Button struct {
	Color     string
	TextColor string
	Text      string
	Link      string
}

// Template is the struct given to Golang templating
// Root object in a template is this struct
type Template struct {
	Engine *OEmailGenerationEngine
	Email  *Email
}
