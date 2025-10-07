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
	projectList *ProjectList       // Backend project list for selected project
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
	projectThemeInfo      *widget.Label
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

	// Apply project theme if there's an active project
	if app.projectList != nil && app.currentProject != "" {
		app.applyProjectTheme()
	}

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

	// Project settings button
	projectSettingsBtn := widget.NewButton("🎨 Theme", func() {
		if app.currentProject == "" {
			dialog.ShowInformation("Thông báo", "Chọn project trước khi thay đổi theme", app.window)
			return
		}
		app.showProjectThemeDialog()
	})
	projectSettingsBtn.Importance = widget.MediumImportance

	projectButtons := container.NewHBox(projectSettingsBtn, addProjectBtn)
	projectSelector := container.NewBorder(nil, nil, nil, projectButtons, app.projectSelect)

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

	// Project theme info - will be updated when project is loaded
	themeInfo := widget.NewLabel("Chưa chọn project")
	themeInfo.TextStyle = fyne.TextStyle{Italic: true}
	
	// Store reference for updating later
	app.projectThemeInfo = themeInfo

	// Main container
	return container.NewBorder(
		container.NewVBox(
			widget.NewLabel("📁 Quản lý Projects"),
			themeInfo,
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

	// Migrate old project file format if needed
	app.migrateOldProjectFile(projectName)

	// Get project color, theme and background image
	projectColor := app.getProjectColor(projectName)
	backgroundImage := app.getProjectBackgroundImage(projectName)
	app.projectColor = projectColor

	// Create ProjectList with color, theme and background image
	app.projectList = NewProjectList(filename, projectName, projectColor, projectColor, backgroundImage)

	// Refresh project lists and apply project theme
	app.refreshAllLists()
	app.applyProjectTheme()

	fmt.Printf("📁 Loaded project: %s (%s) - Background: %s\n", projectName, app.projectList.GetColor(), backgroundImage)
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

// getProjectBackgroundImage returns the background image for a project
func (app *TodoApp) getProjectBackgroundImage(projectName string) string {
	filename := filepath.Join("data/project", projectName+".txt")
	content, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# BackgroundImage: ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# BackgroundImage: "))
		}
	}
	return ""
}

// showCreateProjectDialog shows the create project dialog
func (app *TodoApp) showCreateProjectDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nhập tên project...")

	colorSelect := widget.NewSelect([]string{"blue", "red", "green", "yellow", "orange", "purple", "brown", "black"}, nil)
	colorSelect.SetSelected("blue")

	// Background image selection
	var selectedImagePath string
	imageLabel := widget.NewLabel("Chưa chọn ảnh nền")

	selectImageBtn := widget.NewButton("📁 Chọn ảnh nền", func() {
		app.showImageSelectionDialog(func(imagePath string) {
			selectedImagePath = imagePath
			if imagePath != "" {
				imageLabel.SetText("✅ Đã chọn: " + filepath.Base(imagePath))
			} else {
				imageLabel.SetText("Chưa chọn ảnh nền")
			}
		})
	})

	clearImageBtn := widget.NewButton("🗑️ Xóa ảnh", func() {
		selectedImagePath = ""
		imageLabel.SetText("Chưa chọn ảnh nền")
	})

	imageContainer := container.NewHBox(selectImageBtn, clearImageBtn)

	form := container.NewVBox(
		widget.NewLabel("Tạo Project Mới"),
		widget.NewSeparator(),
		widget.NewFormItem("Tên:", nameEntry).Widget,
		widget.NewFormItem("Màu:", colorSelect).Widget,
		widget.NewFormItem("Ảnh nền:", imageContainer).Widget,
		imageLabel,
	)

	dialog.ShowCustomConfirm("Tạo Project", "Tạo", "Hủy", form, func(response bool) {
		if response && nameEntry.Text != "" && colorSelect.Selected != "" {
			app.createProject(nameEntry.Text, colorSelect.Selected, selectedImagePath)
		}
	}, app.window)
}

// createProject creates a new project
func (app *TodoApp) createProject(name, color string, backgroundImage ...string) {
	projectDir := "data/project"
	os.MkdirAll(projectDir, 0755)

	filename := filepath.Join(projectDir, name+".txt")

	// Build content with optional background image
	content := fmt.Sprintf("# Project: %s\n# Color: %s\n# Created: %s\n",
		name, color, time.Now().Format("2006-01-02 15:04:05"))

	if len(backgroundImage) > 0 && backgroundImage[0] != "" {
		content += fmt.Sprintf("# BackgroundImage: %s\n", backgroundImage[0])
	}

	content += "\n"

	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}

	app.refreshProjectList()
	app.projectSelect.SetSelected(name)

	imageInfo := ""
	if len(backgroundImage) > 0 && backgroundImage[0] != "" {
		imageInfo = " với ảnh nền"
	}

	dialog.ShowInformation("Thành công", fmt.Sprintf("Đã tạo project: %s%s", name, imageInfo), app.window)
}

// showImageSelectionDialog shows dialog to select background image
func (app *TodoApp) showImageSelectionDialog(callback func(string)) {
	// Create file dialog for image selection
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, app.window)
			callback("")
			return
		}
		if reader == nil {
			callback("")
			return
		}
		defer reader.Close()

		// Get file path
		uri := reader.URI()
		imagePath := uri.Path()

		// Validate image file extension
		ext := strings.ToLower(filepath.Ext(imagePath))
		validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}

		isValidImage := false
		for _, validExt := range validExts {
			if ext == validExt {
				isValidImage = true
				break
			}
		}

		if !isValidImage {
			dialog.ShowError(fmt.Errorf("định dạng file không hỗ trợ. Chỉ chấp nhận: %s", strings.Join(validExts, ", ")), app.window)
			callback("")
			return
		}

		// Copy image to themes directory
		themesDir := "data/themes/images"
		os.MkdirAll(themesDir, 0755)

		destPath := filepath.Join(themesDir, filepath.Base(imagePath))

		// Copy file
		sourceFile, err := os.Open(imagePath)
		if err != nil {
			dialog.ShowError(fmt.Errorf("không thể mở file: %v", err), app.window)
			callback("")
			return
		}
		defer sourceFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			dialog.ShowError(fmt.Errorf("không thể tạo file đích: %v", err), app.window)
			callback("")
			return
		}
		defer destFile.Close()

		// Copy content
		_, err = sourceFile.Seek(0, 0)
		if err == nil {
			_, err = destFile.ReadFrom(sourceFile)
		}

		if err != nil {
			dialog.ShowError(fmt.Errorf("không thể copy file: %v", err), app.window)
			callback("")
			return
		}

		callback(destPath)
	}, app.window)

	fileDialog.Show()
}

// showProjectThemeDialog shows dialog to change current project theme
func (app *TodoApp) showProjectThemeDialog() {
	if app.projectList == nil || app.currentProject == "" {
		return
	}

	currentColor := app.projectList.GetColor()
	currentImage := app.projectList.GetBackgroundImage()

	// Color selection
	colorSelect := widget.NewSelect([]string{"blue", "red", "green", "yellow", "orange", "purple", "brown", "black"}, nil)
	colorSelect.SetSelected(currentColor)

	// Current background image info
	var imageLabel *widget.Label
	if currentImage != "" {
		imageLabel = widget.NewLabel("✅ Hiện tại: " + filepath.Base(currentImage))
	} else {
		imageLabel = widget.NewLabel("Chưa có ảnh nền")
	}

	var selectedImagePath string = currentImage

	// Background image selection
	selectImageBtn := widget.NewButton("📁 Chọn ảnh mới", func() {
		app.showImageSelectionDialog(func(imagePath string) {
			selectedImagePath = imagePath
			if imagePath != "" {
				imageLabel.SetText("✅ Mới: " + filepath.Base(imagePath))
			} else {
				imageLabel.SetText("Chưa có ảnh nền")
			}
		})
	})

	clearImageBtn := widget.NewButton("🗑️ Xóa ảnh", func() {
		selectedImagePath = ""
		imageLabel.SetText("Chưa có ảnh nền")
	})

	// Preview button
	previewBtn := widget.NewButton("👁️ Xem trước", func() {
		if selectedImagePath != "" && selectedImagePath != currentImage {
			dialog.ShowInformation("Xem trước",
				fmt.Sprintf("Ảnh nền mới: %s\nMàu: %s",
					filepath.Base(selectedImagePath),
					colorSelect.Selected),
				app.window)
		} else {
			dialog.ShowInformation("Xem trước",
				fmt.Sprintf("Màu hiện tại: %s\nẢnh nền: %s",
					colorSelect.Selected,
					func() string {
						if currentImage != "" {
							return filepath.Base(currentImage)
						}
						return "Không có"
					}()),
				app.window)
		}
	})

	imageContainer := container.NewHBox(selectImageBtn, clearImageBtn, previewBtn)

	form := container.NewVBox(
		widget.NewCard("", "🎨 Cài đặt Theme Project",
			widget.NewLabel(fmt.Sprintf("Project: %s", app.currentProject))),
		widget.NewSeparator(),
		widget.NewFormItem("Màu chủ đề:", colorSelect).Widget,
		widget.NewFormItem("Ảnh nền:", imageContainer).Widget,
		imageLabel,
	)

	dialog.ShowCustomConfirm("Cài đặt Theme", "✅ Áp dụng", "❌ Hủy", form, func(response bool) {
		if response {
			// Update project theme
			newColor := colorSelect.Selected
			if newColor == "" {
				newColor = currentColor
			}

			app.projectList.SetColor(newColor)
			app.projectList.SetTheme(newColor)
			app.projectList.SetBackgroundImage(selectedImagePath)
			app.projectColor = newColor

			// Update project file
			app.updateProjectFile(app.currentProject, newColor, selectedImagePath)

			// Apply new theme
			app.applyProjectTheme()

			dialog.ShowInformation("Thành công",
				fmt.Sprintf("Đã cập nhật theme cho project %s", app.currentProject),
				app.window)
		}
	}, app.window)
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

// updateProjectFile updates the project file with new color and background image information
func (app *TodoApp) updateProjectFile(projectName, color, backgroundImage string) {
	filename := filepath.Join("data/project", projectName+".txt")
	content, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	colorUpdated := false
	backgroundUpdated := false

	for i, line := range lines {
		if strings.HasPrefix(line, "# Color: ") {
			lines[i] = fmt.Sprintf("# Color: %s", color)
			colorUpdated = true
		} else if strings.HasPrefix(line, "# BackgroundImage: ") {
			if backgroundImage != "" {
				lines[i] = fmt.Sprintf("# BackgroundImage: %s", backgroundImage)
			} else {
				// Remove background image line if empty
				lines = append(lines[:i], lines[i+1:]...)
			}
			backgroundUpdated = true
		}
	}

	// Add missing fields if not found
	if !colorUpdated {
		// Insert after project name line
		for i, line := range lines {
			if strings.HasPrefix(line, "# Project: ") {
				newLines := make([]string, 0, len(lines)+1)
				newLines = append(newLines, lines[:i+1]...)
				newLines = append(newLines, fmt.Sprintf("# Color: %s", color))
				newLines = append(newLines, lines[i+1:]...)
				lines = newLines
				break
			}
		}
	}

	if !backgroundUpdated && backgroundImage != "" {
		// Insert after color line
		for i, line := range lines {
			if strings.HasPrefix(line, "# Color: ") {
				newLines := make([]string, 0, len(lines)+1)
				newLines = append(newLines, lines[:i+1]...)
				newLines = append(newLines, fmt.Sprintf("# BackgroundImage: %s", backgroundImage))
				newLines = append(newLines, lines[i+1:]...)
				lines = newLines
				break
			}
		}
	}

	updatedContent := strings.Join(lines, "\n")
	os.WriteFile(filename, []byte(updatedContent), 0644)
}

// applyProjectTheme applies theme specific to the project tab
func (app *TodoApp) applyProjectTheme() {
	if app.projectList == nil || app.tabs == nil {
		return
	}

	projectTabIndex := 1 // Project tab is the second tab (index 1)
	if projectTabIndex >= len(app.tabs.Items) {
		return
	}

	projectTab := app.tabs.Items[projectTabIndex]
	projectColor := app.projectList.GetColor()
	projectName := app.projectList.GetName()
	backgroundImage := app.projectList.GetBackgroundImage()

	// Update tab title with color indicator
	colorEmoji := app.getColorEmoji(projectColor)
	projectTab.Text = fmt.Sprintf("📁 Projects %s", colorEmoji)

	// Update theme info label if available
	if app.projectThemeInfo != nil {
		themeMessage := fmt.Sprintf("🎨 %s • %s %s", projectName, colorEmoji, strings.ToUpper(projectColor))
		// Only show background image indicator if file actually exists
		if backgroundImage != "" {
			if _, err := os.Stat(backgroundImage); err == nil {
				themeMessage += " • 🖼️ Background Image"
			}
		}
		app.projectThemeInfo.SetText(themeMessage)
		app.projectThemeInfo.TextStyle = fyne.TextStyle{Italic: true, Bold: true}
		app.projectThemeInfo.Refresh()
	}

	// Apply background image safely if available and file exists
	if backgroundImage != "" {
		if _, err := os.Stat(backgroundImage); err == nil {
			themedContent := app.createProjectThemedContainer(projectTab.Content)
			projectTab.Content = themedContent
			app.tabs.Refresh()
		}
	}
	
	fmt.Printf("🎨 Applied project theme: %s (color: %s, image: %s)\n",
		projectName, projectColor, backgroundImage)
}

// createProjectThemedContainer creates a themed container for project tab
func (app *TodoApp) createProjectThemedContainer(originalContent fyne.CanvasObject) fyne.CanvasObject {
	if app.projectList == nil {
		return originalContent
	}

	// Only add background image if it exists and is valid
	backgroundImagePath := app.projectList.GetBackgroundImage()
	if backgroundImagePath != "" {
		if _, err := os.Stat(backgroundImagePath); err == nil {
			imageResource, err := fyne.LoadResourceFromPath(backgroundImagePath)
			if err == nil {
				// Create background image
				backgroundImage := widget.NewIcon(imageResource)
				
				// Stack background image behind content
				return container.NewStack(
					backgroundImage,
					container.NewPadded(originalContent),
				)
			} else {
				fmt.Printf("❌ Error loading background image: %v\n", err)
			}
		}
	}

	// If no background image or error, return original content
	return originalContent
}

// createColorBackground creates a colored background based on project color
func (app *TodoApp) createColorBackground(projectColor string) fyne.CanvasObject {
	// Create a simple colored card as background
	colorEmoji := app.getColorEmoji(projectColor)
	colorName := strings.ToUpper(projectColor)
	
	// Create a subtle indicator card
	colorCard := widget.NewCard("", fmt.Sprintf("%s %s Theme", colorEmoji, colorName), nil)
	
	// Create a container with very light background color indication
	colorIndicator := widget.NewLabel("")
	colorIndicator.Resize(fyne.NewSize(2000, 2000))
	
	return container.NewStack(colorIndicator, colorCard)
}

// getColorEmoji returns emoji for project color
func (app *TodoApp) getColorEmoji(color string) string {
	switch color {
	case "red":
		return "🔴"
	case "green":
		return "🟢"
	case "blue":
		return "🔵"
	case "yellow":
		return "🟡"
	case "orange":
		return "🟠"
	case "purple":
		return "🟣"
	case "brown":
		return "🟤"
	case "black":
		return "⚫"
	default:
		return "🔵" // default blue
	}
}

// migrateOldProjectFile adds header metadata to old project files that don't have it
func (app *TodoApp) migrateOldProjectFile(projectName string) {
	filename := filepath.Join("data/project", projectName+".txt")
	content, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	// Check if file already has header
	lines := strings.Split(string(content), "\n")
	hasHeader := false
	for _, line := range lines {
		if strings.HasPrefix(line, "# Project:") {
			hasHeader = true
			break
		}
	}

	// If file doesn't have header, add it
	if !hasHeader {
		header := fmt.Sprintf("# Project: %s\n# Color: blue\n# Created: %s\n\n",
			projectName, time.Now().Format("2006-01-02 15:04:05"))
		newContent := header + string(content)

		os.WriteFile(filename, []byte(newContent), 0644)
		fmt.Printf("🔄 Migrated old project file: %s\n", projectName)
	}
}

// createThemedProjectTabContent recreates the project tab content with theme applied
func (app *TodoApp) createThemedProjectTabContent(themeMessage string) *fyne.Container {
	// Don't recreate the entire tab content, just return the existing one with theme info
	// This approach is simpler and avoids widget conflicts
	
	// Theme info label
	themeInfoLabel := widget.NewLabel(themeMessage)
	themeInfoLabel.TextStyle = fyne.TextStyle{Italic: true, Bold: true}

	// Get existing project tab content and wrap with theme info
	existingContent := app.setupProjectTab()
	
	// Wrap with theme background if available
	if app.projectList != nil && (app.projectList.HasBackgroundImage() || app.projectList.GetColor() != "blue") {
		themedContent := app.createProjectThemedContainer(existingContent)
		return themedContent.(*fyne.Container)
	}

	return existingContent
}
