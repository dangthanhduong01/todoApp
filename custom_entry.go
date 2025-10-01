package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CustomEntry - Manual text input widget để fix render issues
type CustomEntry struct {
	widget.BaseWidget
	text        string
	placeholder string
	label       *widget.Label
	onSubmit    func(string)
	onChange    func(string)
	cursor      int
	focused     bool
}

func NewCustomEntry() *CustomEntry {
	entry := &CustomEntry{
		text:        "",
		placeholder: "",
		cursor:      0,
		focused:     false,
	}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *CustomEntry) SetPlaceHolder(text string) {
	e.placeholder = text
	e.updateDisplay()
}

func (e *CustomEntry) SetOnSubmitted(callback func(string)) {
	e.onSubmit = callback
}

func (e *CustomEntry) SetOnChanged(callback func(string)) {
	e.onChange = callback
}

func (e *CustomEntry) SetText(text string) {
	e.text = text
	e.cursor = len(text)
	e.updateDisplay()
	if e.onChange != nil {
		e.onChange(text)
	}
}

func (e *CustomEntry) GetText() string {
	return e.text
}

func (e *CustomEntry) updateDisplay() {
	displayText := e.text
	if displayText == "" && !e.focused {
		displayText = e.placeholder
	}

	// Add cursor indicator
	if e.focused && e.cursor <= len(e.text) {
		if e.cursor == len(e.text) {
			displayText = e.text + "|"
		} else {
			displayText = e.text[:e.cursor] + "|" + e.text[e.cursor:]
		}
	}

	// Simple direct update - let Fyne handle threading
	if e.label != nil {
		e.label.SetText(displayText)
		e.Refresh() // Refresh the widget itself, not the label
	}
}

func (e *CustomEntry) TypedRune(r rune) {
	if r == '\r' || r == '\n' {
		// Enter pressed
		if e.onSubmit != nil {
			e.onSubmit(e.text)
		}
		return
	}

	// Insert character at cursor
	if e.cursor <= len(e.text) {
		e.text = e.text[:e.cursor] + string(r) + e.text[e.cursor:]
		e.cursor++
		e.updateDisplay()

		if e.onChange != nil {
			e.onChange(e.text)
		}

		fmt.Printf("✍️ Typed '%c' - text: '%s'\n", r, e.text)
	}
}

func (e *CustomEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyBackspace:
		if e.cursor > 0 && len(e.text) > 0 {
			e.text = e.text[:e.cursor-1] + e.text[e.cursor:]
			e.cursor--
			e.updateDisplay()

			if e.onChange != nil {
				e.onChange(e.text)
			}
		}
	case fyne.KeyDelete:
		if e.cursor < len(e.text) {
			e.text = e.text[:e.cursor] + e.text[e.cursor+1:]
			e.updateDisplay()

			if e.onChange != nil {
				e.onChange(e.text)
			}
		}
	case fyne.KeyLeft:
		if e.cursor > 0 {
			e.cursor--
			e.updateDisplay()
		}
	case fyne.KeyRight:
		if e.cursor < len(e.text) {
			e.cursor++
			e.updateDisplay()
		}
	case fyne.KeyHome:
		e.cursor = 0
		e.updateDisplay()
	case fyne.KeyEnd:
		e.cursor = len(e.text)
		e.updateDisplay()
	}
}

func (e *CustomEntry) FocusGained() {
	e.focused = true
	e.updateDisplay()
	fmt.Println("🎯 CustomEntry gained focus")
}

func (e *CustomEntry) FocusLost() {
	e.focused = false
	e.updateDisplay()
	fmt.Println("😴 CustomEntry lost focus")
}

func (e *CustomEntry) CreateRenderer() fyne.WidgetRenderer {
	e.label = widget.NewLabel("")
	e.label.Alignment = fyne.TextAlignLeading

	// Style như Entry widget
	bg := widget.NewCard("", "", nil)
	bg.SetContent(container.NewPadded(e.label))

	e.updateDisplay()

	return widget.NewSimpleRenderer(bg)
}

func (e *CustomEntry) Tapped(*fyne.PointEvent) {
	// Handle tap to position cursor (simplified)
	e.cursor = len(e.text)
	e.updateDisplay()
}

func (e *CustomEntry) DoubleTapped(*fyne.PointEvent) {
	// Select all text (simplified)
	fmt.Println("Double tapped - selecting all")
}
