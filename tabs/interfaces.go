package tabs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// TodoTabInterface defines the interface for todo tab components
type TodoTabInterface interface {
	GetContainer() *fyne.Container
	RefreshLists()
	AddTodo(description string) error
}

// ProjectTabInterface defines the interface for project tab components
type ProjectTabInterface interface {
	GetContainer() *fyne.Container
	RefreshLists()
	LoadProject(projectName string) error
	CreateProject(name, color string) error
}

// BaseTab provides common functionality for all tabs
type BaseTab struct {
	container *fyne.Container
	title     string
}

// NewBaseTab creates a new base tab
func NewBaseTab(title string) *BaseTab {
	return &BaseTab{
		title:     title,
		container: container.NewVBox(),
	}
}

// GetTitle returns the tab title
func (bt *BaseTab) GetTitle() string {
	return bt.title
}

// GetContainer returns the tab container
func (bt *BaseTab) GetContainer() *fyne.Container {
	return bt.container
}

// AddComponent adds a component to the tab
func (bt *BaseTab) AddComponent(component fyne.CanvasObject) {
	bt.container.Add(component)
}
