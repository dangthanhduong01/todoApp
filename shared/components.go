package shared

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// SharedComponents contains common UI components used across the application
type SharedComponents struct{}

// CreateStyledButton creates a button with consistent styling
func (sc *SharedComponents) CreateStyledButton(text string, callback func(), importance widget.Importance) *widget.Button {
	btn := widget.NewButton(text, callback)
	btn.Importance = importance
	return btn
}

// CreateFormContainer creates a standardized form container
func (sc *SharedComponents) CreateFormContainer(items ...fyne.CanvasObject) *fyne.Container {
	return container.NewVBox(items...)
}

// CreateHeaderCard creates a standardized header card
func (sc *SharedComponents) CreateHeaderCard(title, subtitle string) *widget.Card {
	return widget.NewCard(title, subtitle, nil)
}
