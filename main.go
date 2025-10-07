package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"todoapp/fireworks"
)

// TodoApp represents the main application structure
type TodoApp struct {
	todoList    *TodoList          // Backend todo list for main todos
	projectList *TodoList          // Backend todo list for selected project
	window      fyne.Window        // Main application window
	tabs        *container.AppTabs // Tab container
	myApp       fyne.App           // Reference to the Fyne application
	isDarkTheme bool               // Current theme state

	// Todo tab widgets
	allList        *widget.List
	activeList     *widget.List
	completedList  *widget.List
	allTodos       []Todo
	activeTodos    []Todo
	completedTodos []Todo

	// Project tab widgets
	projectAllList        *widget.List
	projectActiveList     *widget.List
	projectCompletedList  *widget.List
	projectAllTodos       []Todo
	projectActiveTodos    []Todo
	projectCompletedTodos []Todo
	projectSelect         *widget.Select
	projectTodoEntry      *widget.Entry
	currentProject        string
	projectColor          string
}

// main initializes and starts the application
func main() {
	fmt.Println("🚀 Starting Todo App...")

	// Force software rendering for better compatibility
	os.Setenv("FYNE_DRIVER", "x11")
	os.Setenv("FYNE_SOFTWARE", "1")
	os.Setenv("FYNE_DISABLE_HARDWARE_ACCELERATION", "1")
	os.Setenv("GTK_IM_MODULE", "")
	os.Setenv("QT_IM_MODULE", "")
	os.Setenv("XMODIFIERS", "")
	os.Setenv("SDL_IM_MODULE", "")
	os.Setenv("FYNE_FONT", "")

	fmt.Println("📱 Environment variables set")

	myApp := app.New()
	myApp.SetIcon(theme.DocumentIcon())

	myWindow := myApp.NewWindow("📝 Todo List Application")
	myWindow.Resize(fyne.NewSize(900, 700))
	myWindow.CenterOnScreen()

	todoApp := &TodoApp{
		todoList:    NewTodoList("todos.txt"),
		window:      myWindow,
		myApp:       myApp,
		isDarkTheme: false,
	}

	todoApp.setupUI()
	myWindow.Show()
	myApp.Run()
}

// setupUI configures the main interface
func (app *TodoApp) setupUI() {
	// Settings button
	settingsButton := widget.NewButton("⚙️ Cài đặt", func() {
		app.showSettingsDialog()
	})
	settingsButton.Importance = widget.MediumImportance

	// Setup tabs
	todoTabContent := app.setupTodoTab()
	projectTabContent := app.setupProjectTab()

	// Main tabs
	app.tabs = container.NewAppTabs(
		container.NewTabItem("📋 Todos", todoTabContent),
		container.NewTabItem("📁 Projects", projectTabContent),
	)

	// Header
	header := widget.NewCard("", "Todo List Desktop App", nil)
	headerWithButtons := container.NewBorder(
		nil, nil, nil,
		container.NewHBox(settingsButton),
		header,
	)

	// Main view
	mainView := container.NewBorder(
		container.NewVBox(headerWithButtons, widget.NewSeparator()),
		nil, nil, nil,
		app.tabs,
	)

	app.window.SetContent(mainView)

	// Load initial data
	app.refreshAllLists()
	app.applyTheme()

	fmt.Println("🎛️ UI setup complete")
}

// setupTodoTab creates the todos tab content
func (app *TodoApp) setupTodoTab() *fyne.Container {
	// Create lists
	app.allList = app.createList("all", false)
	app.activeList = app.createList("active", false)
	app.completedList = app.createList("completed", false)

	// Input for adding todos
	todoEntry := widget.NewEntry()
	todoEntry.SetPlaceHolder("Nhập công việc mới...")

	addTodoBtn := widget.NewButton("+ Thêm Todo", func() {
		app.addTodo(todoEntry.Text, false)
		todoEntry.SetText("")
	})
	addTodoBtn.Importance = widget.HighImportance

	// Enter key support
	todoEntry.OnSubmitted = func(text string) {
		app.addTodo(text, false)
		todoEntry.SetText("")
	}

	todoInputContainer := container.NewBorder(nil, nil, nil, addTodoBtn, todoEntry)

	// Todo sub-tabs
	todoSubTabs := container.NewAppTabs(
		container.NewTabItem("Tất cả", container.NewScroll(app.allList)),
		container.NewTabItem("Chưa hoàn thành", container.NewScroll(app.activeList)),
		container.NewTabItem("Đã hoàn thành", container.NewScroll(app.completedList)),
	)

	// Main container
	return container.NewBorder(
		container.NewVBox(
			widget.NewLabel("📋 Quản lý Todos"),
			widget.NewSeparator(),
			todoInputContainer,
			widget.NewSeparator(),
		),
		nil, nil, nil,
		todoSubTabs,
	)
}

// setupProjectTab creates the projects tab content
func (app *TodoApp) setupProjectTab() *fyne.Container {
	// Project selection
	app.projectSelect = widget.NewSelect([]string{}, func(selected string) {
		if selected != "" && selected != "Chưa có project nào" {
			app.loadProject(selected)
		}
	})

	// Create project button
	addProjectBtn := widget.NewButton("+ Tạo Project", func() {
		app.showCreateProjectDialog()
	})
	addProjectBtn.Importance = widget.HighImportance

	projectSelector := container.NewBorder(nil, nil, nil, addProjectBtn, app.projectSelect)

	// Project todo lists
	app.projectAllList = app.createList("all", true)
	app.projectActiveList = app.createList("active", true)
	app.projectCompletedList = app.createList("completed", true)

	// Project todo input
	app.projectTodoEntry = widget.NewEntry()
	app.projectTodoEntry.SetPlaceHolder("Nhập công việc cho project...")

	addProjectTodoBtn := widget.NewButton("+ Thêm", func() {
		if app.currentProject == "" {
			dialog.ShowInformation("Thông báo", "Chọn project trước khi thêm todo", app.window)
			return
		}
		app.addTodo(app.projectTodoEntry.Text, true)
		app.projectTodoEntry.SetText("")
	})
	addProjectTodoBtn.Importance = widget.HighImportance

	// Enter key support
	app.projectTodoEntry.OnSubmitted = func(text string) {
		if app.currentProject == "" {
			dialog.ShowInformation("Thông báo", "Chọn project trước khi thêm todo", app.window)
			return
		}
		app.addTodo(text, true)
		app.projectTodoEntry.SetText("")
	}

	projectTodoInputContainer := container.NewBorder(nil, nil, nil, addProjectTodoBtn, app.projectTodoEntry)

	// Project sub-tabs
	projectSubTabs := container.NewAppTabs(
		container.NewTabItem("Tất cả", container.NewScroll(app.projectAllList)),
		container.NewTabItem("Chưa hoàn thành", container.NewScroll(app.projectActiveList)),
		container.NewTabItem("Đã hoàn thành", container.NewScroll(app.projectCompletedList)),
	)

	// Load available projects
	app.refreshProjectList()

	// Main container
	return container.NewBorder(
		container.NewVBox(
			widget.NewLabel("📁 Quản lý Projects"),
			widget.NewSeparator(),
			projectSelector,
			widget.NewSeparator(),
			projectTodoInputContainer,
			widget.NewSeparator(),
		),
		nil, nil, nil,
		projectSubTabs,
	)
}

// createList creates a todo list widget
func (app *TodoApp) createList(listType string, isProject bool) *widget.List {
	list := widget.NewList(
		func() int {
			if isProject {
				switch listType {
				case "all":
					return len(app.projectAllTodos)
				case "active":
					return len(app.projectActiveTodos)
				case "completed":
					return len(app.projectCompletedTodos)
				}
			} else {
				switch listType {
				case "all":
					return len(app.allTodos)
				case "active":
					return len(app.activeTodos)
				case "completed":
					return len(app.completedTodos)
				}
			}
			return 0
		},
		func() fyne.CanvasObject {
			return widget.NewCard("", "", widget.NewLabel(""))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			app.updateTodoItem(id, item, listType, isProject)
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		app.handleTodoSelection(id, listType, isProject)
	}

	return list
}

// updateTodoItem updates a todo item in the list
func (app *TodoApp) updateTodoItem(id widget.ListItemID, item fyne.CanvasObject, listType string, isProject bool) {
	var todos []Todo

	if isProject {
		switch listType {
		case "all":
			todos = app.projectAllTodos
		case "active":
			todos = app.projectActiveTodos
		case "completed":
			todos = app.projectCompletedTodos
		}
	} else {
		switch listType {
		case "all":
			todos = app.allTodos
		case "active":
			todos = app.activeTodos
		case "completed":
			todos = app.completedTodos
		}
	}

	if id >= len(todos) {
		return
	}

	// Newest first
	reversedIndex := len(todos) - 1 - id
	todo := todos[reversedIndex]
	card := item.(*widget.Card)

	// Date label
	dateLabel := widget.NewLabel(todo.CreatedAt.Format("02/01 15:04"))
	dateLabel.TextStyle = fyne.TextStyle{Italic: true}

	// Content label
	contentLabel := widget.NewLabel(todo.Description)
	contentLabel.TextStyle = fyne.TextStyle{Bold: true}
	contentLabel.Wrapping = fyne.TextWrapWord

	// Complete checkbox
	var completeCheck *widget.Check
	completeCheck = widget.NewCheck("", func(checked bool) {
		if !todo.Completed && checked {
			app.markComplete(todo.ID, isProject)
		} else if todo.Completed && !checked {
			dialog.ShowInformation("Thông báo", "Không thể bỏ tích công việc đã hoàn thành", app.window)
			completeCheck.SetChecked(true)
		}
	})
	completeCheck.SetChecked(todo.Completed)

	// Delete button
	deleteBtn := widget.NewButton("🗑️", func() {
		app.confirmDelete(todo.ID, todo.Description, isProject)
	})

	buttonsContainer := container.NewHBox(completeCheck, deleteBtn)

	// Layout
	horizontalLayout := container.NewBorder(
		nil, nil,
		dateLabel,
		buttonsContainer,
		contentLabel,
	)

	card.SetContent(container.NewPadded(horizontalLayout))
}

// addTodo adds a new todo item
func (app *TodoApp) addTodo(description string, isProject bool) {
	description = strings.TrimSpace(description)
	if description == "" {
		dialog.ShowError(fmt.Errorf("vui lòng nhập mô tả công việc"), app.window)
		return
	}

	var err error
	if isProject && app.projectList != nil {
		err = app.projectList.AddTodo(description)
	} else if !isProject {
		err = app.todoList.AddTodo(description)
	}

	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}

	app.refreshAllLists()
	dialog.ShowInformation("Thành công", fmt.Sprintf("Đã thêm: %s", description), app.window)
}

// markComplete marks a todo as completed
func (app *TodoApp) markComplete(todoID int, isProject bool) {
	var err error
	var todoDescription string

	// Find todo description
	var todos []Todo
	if isProject && app.projectList != nil {
		todos = app.projectList.GetTodos()
	} else {
		todos = app.todoList.GetTodos()
	}

	for _, todo := range todos {
		if todo.ID == todoID {
			todoDescription = todo.Description
			break
		}
	}

	// Mark complete
	if isProject && app.projectList != nil {
		err = app.projectList.MarkComplete(todoID)
	} else {
		err = app.todoList.MarkComplete(todoID)
	}

	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}

	app.refreshAllLists()
	fireworks.ShowFireworksDialog(todoDescription, app.window)
}

// confirmDelete shows confirmation dialog for deleting todo
func (app *TodoApp) confirmDelete(todoID int, description string, isProject bool) {
	dialog.ShowConfirm("Xác nhận xóa",
		fmt.Sprintf("Bạn có chắc chắn muốn xóa:\n'%s'?", description),
		func(confirmed bool) {
			if confirmed {
				var err error
				if isProject && app.projectList != nil {
					err = app.projectList.DeleteTodo(todoID)
				} else {
					err = app.todoList.DeleteTodo(todoID)
				}

				if err != nil {
					dialog.ShowError(err, app.window)
					return
				}

				app.refreshAllLists()
				dialog.ShowInformation("Thành công", fmt.Sprintf("Đã xóa: %s", description), app.window)
			}
		}, app.window)
}

// handleTodoSelection handles when a todo is selected
func (app *TodoApp) handleTodoSelection(id widget.ListItemID, listType string, isProject bool) {
	var todos []Todo

	if isProject {
		switch listType {
		case "all":
			todos = app.projectAllTodos
		case "active":
			todos = app.projectActiveTodos
		case "completed":
			todos = app.projectCompletedTodos
		}
	} else {
		switch listType {
		case "all":
			todos = app.allTodos
		case "active":
			todos = app.activeTodos
		case "completed":
			todos = app.completedTodos
		}
	}

	if id >= len(todos) {
		return
	}

	reversedIndex := len(todos) - 1 - id
	todo := todos[reversedIndex]

	if todo.Completed {
		dialog.ShowConfirm("Công việc đã hoàn thành",
			fmt.Sprintf("Công việc: %s\nBạn muốn xóa?", todo.Description),
			func(confirmed bool) {
				if confirmed {
					app.confirmDelete(todo.ID, todo.Description, isProject)
				}
			}, app.window)
	} else {
		// Show options
		completeBtn := widget.NewButton("✅ Đánh dấu hoàn thành", func() {
			app.markComplete(todo.ID, isProject)
		})
		completeBtn.Importance = widget.SuccessImportance

		deleteBtn := widget.NewButton("🗑️ Xóa", func() {
			app.confirmDelete(todo.ID, todo.Description, isProject)
		})
		deleteBtn.Importance = widget.DangerImportance

		content := container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Công việc: %s", todo.Description)),
			widget.NewSeparator(),
			completeBtn,
			deleteBtn,
		)

		dialog.ShowCustom("Chọn hành động", "Hủy", content, app.window)
	}
}

// refreshAllLists refreshes all todo lists
func (app *TodoApp) refreshAllLists() {
	// Main todos
	app.allTodos = app.todoList.GetTodos()
	app.activeTodos = app.todoList.GetActiveTodos()
	app.completedTodos = app.todoList.GetCompletedTodos()

	if app.allList != nil {
		app.allList.Refresh()
	}
	if app.activeList != nil {
		app.activeList.Refresh()
	}
	if app.completedList != nil {
		app.completedList.Refresh()
	}

	// Project todos
	if app.projectList != nil {
		app.projectAllTodos = app.projectList.GetTodos()
		app.projectActiveTodos = app.projectList.GetActiveTodos()
		app.projectCompletedTodos = app.projectList.GetCompletedTodos()

		if app.projectAllList != nil {
			app.projectAllList.Refresh()
		}
		if app.projectActiveList != nil {
			app.projectActiveList.Refresh()
		}
		if app.projectCompletedList != nil {
			app.projectCompletedList.Refresh()
		}
	}
}

// refreshProjectList updates the project dropdown
func (app *TodoApp) refreshProjectList() {
	projectDir := "data/project"
	files, err := os.ReadDir(projectDir)
	if err != nil {
		app.projectSelect.Options = []string{"Chưa có project nào"}
		app.projectSelect.Refresh()
		return
	}

	var projects []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".txt") {
			projectName := strings.TrimSuffix(file.Name(), ".txt")
			projects = append(projects, projectName)
		}
	}

	if len(projects) == 0 {
		projects = []string{"Chưa có project nào"}
	} else {
		// Sort by modification time (newest first)
		sort.Slice(projects, func(i, j int) bool {
			pathI := filepath.Join(projectDir, projects[i]+".txt")
			pathJ := filepath.Join(projectDir, projects[j]+".txt")
			statI, _ := os.Stat(pathI)
			statJ, _ := os.Stat(pathJ)
			return statI.ModTime().After(statJ.ModTime())
		})
	}

	app.projectSelect.Options = projects
	app.projectSelect.Refresh()

	// Auto-select first project if available
	if len(projects) > 0 && projects[0] != "Chưa có project nào" {
		app.projectSelect.SetSelected(projects[0])
	}
}

// loadProject loads the selected project
func (app *TodoApp) loadProject(projectName string) {
	app.currentProject = projectName
	filename := fmt.Sprintf("data/project/%s.txt", projectName)
	app.projectList = NewTodoList(filename)

	// Get project color
	app.projectColor = app.getProjectColor(projectName)

	// Refresh project lists
	app.refreshAllLists()

	fmt.Printf("📁 Loaded project: %s (%s)\n", projectName, app.projectColor)
}

// getProjectColor returns the color for a project
func (app *TodoApp) getProjectColor(projectName string) string {
	filename := filepath.Join("data/project", projectName+".txt")
	content, err := os.ReadFile(filename)
	if err != nil {
		return "blue"
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# Color: ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# Color: "))
		}
	}
	return "blue"
}

// showCreateProjectDialog shows the create project dialog
func (app *TodoApp) showCreateProjectDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nhập tên project...")

	colorSelect := widget.NewSelect([]string{"blue", "red", "green", "yellow", "orange", "purple", "brown", "black"}, nil)
	colorSelect.SetSelected("blue")

	form := container.NewVBox(
		widget.NewLabel("Tạo Project Mới"),
		widget.NewSeparator(),
		widget.NewFormItem("Tên:", nameEntry).Widget,
		widget.NewFormItem("Màu:", colorSelect).Widget,
	)

	dialog.ShowCustomConfirm("Tạo Project", "Tạo", "Hủy", form, func(response bool) {
		if response && nameEntry.Text != "" && colorSelect.Selected != "" {
			app.createProject(nameEntry.Text, colorSelect.Selected)
		}
	}, app.window)
}

// createProject creates a new project
func (app *TodoApp) createProject(name, color string) {
	projectDir := "data/project"
	os.MkdirAll(projectDir, 0755)

	filename := filepath.Join(projectDir, name+".txt")
	content := fmt.Sprintf("# Project: %s\n# Color: %s\n# Created: %s\n\n",
		name, color, time.Now().Format("2006-01-02 15:04:05"))

	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}

	app.refreshProjectList()
	app.projectSelect.SetSelected(name)

	dialog.ShowInformation("Thành công", fmt.Sprintf("Đã tạo project: %s", name), app.window)
}

// showSettingsDialog shows the settings dialog
func (app *TodoApp) showSettingsDialog() {
	themeLabel := widget.NewLabel(app.getThemeLabelText())

	var themeSwitch *widget.Button
	themeSwitch = widget.NewButton("", func() {
		app.isDarkTheme = !app.isDarkTheme
		app.applyTheme()
		app.updateSwitchAppearance(themeSwitch)
		themeLabel.SetText(app.getThemeLabelText())
	})

	app.updateSwitchAppearance(themeSwitch)

	content := container.NewVBox(
		widget.NewLabel("Chọn giao diện sáng hoặc tối"),
		widget.NewSeparator(),
		themeLabel,
		themeSwitch,
	)

	dialog.ShowCustom("⚙️ Cài đặt", "Đóng", content, app.window)
}

// applyTheme applies the selected theme
func (app *TodoApp) applyTheme() {
	var customTheme fyne.Theme
	if app.isDarkTheme {
		customTheme = &customDarkTheme{}
	} else {
		customTheme = &customLightTheme{}
	}

	app.myApp.Settings().SetTheme(customTheme)
	app.window.Content().Refresh()
}

// getThemeLabelText returns the theme label text
func (app *TodoApp) getThemeLabelText() string {
	if app.isDarkTheme {
		return "🌙 Theme hiện tại: Tối"
	} else {
		return "☀️ Theme hiện tại: Sáng"
	}
}

// updateSwitchAppearance updates the theme switch appearance
func (app *TodoApp) updateSwitchAppearance(btn *widget.Button) {
	if app.isDarkTheme {
		btn.SetText("🌙 TỐI")
		btn.Importance = widget.HighImportance
	} else {
		btn.SetText("☀️ SÁNG")
		btn.Importance = widget.LowImportance
	}
}

// Theme implementations
type customDarkTheme struct{}

func (t *customDarkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, theme.VariantDark)
}

func (t *customDarkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *customDarkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *customDarkTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

type customLightTheme struct{}

func (t *customLightTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

func (t *customLightTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *customLightTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *customLightTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
