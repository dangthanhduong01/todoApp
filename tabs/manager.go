package tabs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// TabManager manages tab creation and handling
type TabManager struct {
	tabs *container.AppTabs
}

// NewTabManager creates a new tab manager
func NewTabManager() *TabManager {
	return &TabManager{
		tabs: container.NewAppTabs(),
	}
}

// AddTab adds a new tab to the container
func (tm *TabManager) AddTab(title string, content fyne.CanvasObject) {
	tab := container.NewTabItem(title, content)
	tm.tabs.Append(tab)
}

// GetTabs returns the tabs container
func (tm *TabManager) GetTabs() *container.AppTabs {
	return tm.tabs
}

// SetSelectedTab sets the active tab by index
func (tm *TabManager) SetSelectedTab(index int) {
	if index >= 0 && index < len(tm.tabs.Items) {
		tm.tabs.SelectTab(tm.tabs.Items[index])
	}
}
