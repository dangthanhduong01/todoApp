package main

import (
	"fmt"
	"os"
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

// TodoApp represents the main application structure containing UI components and state
type TodoApp struct {
	todoList       *TodoList          // Backend todo list handler
	window         fyne.Window        // Main application window
	allList        *widget.List       // List widget for all todos
	activeList     *widget.List       // List widget for active todos
	completedList  *widget.List       // List widget for completed todos
	addEntry       *widget.Entry      // Entry widget for adding new todos
	allTodos       []Todo             // Cache of all todos
	activeTodos    []Todo             // Cache of active todos
	completedTodos []Todo             // Cache of completed todos
	tabs           *container.AppTabs // Tab container for different todo views
	inputContainer *fyne.Container    // Container for input elements
	addButton      *widget.Button     // Button to show add input
	settingsButton *widget.Button     // Button to show settings
	isDarkTheme    bool               // Current theme state
	myApp          fyne.App           // Reference to the Fyne application
	showingInput   bool               // State of input visibility

	// Navigation and Project management
	navbar          *container.AppTabs // Main navigation bar
	navbarButton    *widget.Button     // Button to toggle navbar mode
	isNavbarMode    bool               // Current view mode state
	projectsList    *widget.List       // List widget for projects
	currentProject  string             // Currently selected project name
	projectTodoList *TodoList          // Backend for current project todos
	projectTabs     *container.AppTabs // Tab container for project views
}

// main initializes and starts the Todo application with proper environment setup
func main() {
	fmt.Println("🚀 Starting Todo App...")

	// Force software rendering and fix input display issues
	os.Setenv("FYNE_DRIVER", "x11")
	os.Setenv("FYNE_SOFTWARE", "1")
	os.Setenv("FYNE_DISABLE_HARDWARE_ACCELERATION", "1")
	// Disable problematic input methods that cause text rendering issues
	os.Setenv("GTK_IM_MODULE", "")
	os.Setenv("QT_IM_MODULE", "")
	os.Setenv("XMODIFIERS", "")
	os.Setenv("SDL_IM_MODULE", "")
	// Force proper text rendering
	os.Setenv("FYNE_FONT", "")
	os.Setenv("FYNE_THEME", "light")

	fmt.Println("📱 Environment variables set for software rendering")

	myApp := app.New()
	fmt.Println("✅ App created")

	myApp.SetIcon(theme.DocumentIcon())
	fmt.Println("🎨 Icon set")

	myWindow := myApp.NewWindow("📝 Todo List Application")
	fmt.Println("🪟 Window created")

	myWindow.Resize(fyne.NewSize(700, 600))
	myWindow.CenterOnScreen()
	fmt.Println("📐 Window sized and centered")

	todoApp := &TodoApp{
		todoList:    NewTodoList("todos.txt"),
		window:      myWindow,
		myApp:       myApp,
		isDarkTheme: false, // Mặc định theme sáng
	}
	fmt.Println("📋 TodoApp struct created")

	todoApp.setupUI()
	fmt.Println("🎛️ UI setup complete")

	// Show window explicitly before running
	myWindow.Show()
	fmt.Println("👁️ Window shown, starting main loop...")
	myApp.Run()
}

// setupUI configures and initializes all user interface components
func (app *TodoApp) setupUI() {
	// Tạo nút "Thêm" ban đầu với simple handler
	app.addButton = widget.NewButton("➕ Thêm công việc mới", func() {
		fmt.Println("📱 Add button clicked")
		app.showAddInput()
	})
	app.addButton.Importance = widget.HighImportance

	// Tạo widget.Entry với multiline để tránh text rendering issues
	app.addEntry = widget.NewMultiLineEntry()
	app.addEntry.SetPlaceHolder("Nhập công việc mới...")
	app.addEntry.Wrapping = fyne.TextWrapWord
	app.addEntry.Resize(fyne.NewSize(400, 60))

	// Set change handler với refresh để force text hiển thị
	app.addEntry.OnChanged = func(text string) {
		fmt.Printf("📝 Text changed: '%s'\n", text)
		// Force refresh để text hiện ngay
		app.addEntry.Refresh()
	}

	// Submit handler
	handleSubmit := func() {
		app.addTodo()
	}
	app.addEntry.OnSubmitted = func(text string) {
		handleSubmit()
	}

	// Container chính sẽ switch giữa nút Thêm và input field
	app.inputContainer = container.NewVBox(app.addButton)
	paddedInput := container.NewPadded(app.inputContainer)

	// Tạo list cho từng tab
	app.allList = app.createList("all")
	app.activeList = app.createList("active")
	app.completedList = app.createList("completed")

	// Tabs cho todos (không bao gồm projects)
	app.tabs = container.NewAppTabs(
		container.NewTabItem("Tất cả", app.allList),
		container.NewTabItem("Chưa hoàn thành", app.activeList),
		container.NewTabItem("Đã hoàn thành", app.completedList),
	)

	// Tạo nút Settings
	app.settingsButton = widget.NewButton("⚙️ Cài đặt", func() {
		app.showSettingsDialog()
	})
	app.settingsButton.Importance = widget.MediumImportance

	// Tạo nút navbar toggle
	app.navbarButton = widget.NewButton("≡", func() {
		app.toggleNavbarMode()
	})
	app.navbarButton.Importance = widget.MediumImportance

	// Nội dung chính với expanded layout
	header := widget.NewCard("", "Todo List Desktop App", nil)
	header.Resize(fyne.NewSize(600, 60))

	// Header với nút navbar và settings
	headerWithButtons := container.NewBorder(
		nil, nil, 
		app.navbarButton, // Trái: nút navbar
		container.NewHBox(app.settingsButton), // Phải: nút settings
		header,
	)

	// Tạo view chính ban đầu (chế độ bình thường)
	mainView := container.NewBorder(
		container.NewVBox(headerWithButtons, paddedInput, widget.NewSeparator()),
		nil, nil, nil,
		app.tabs,
	)

	// Khởi tạo với chế độ bình thường
	app.isNavbarMode = false
	app.window.SetContent(mainView)
	app.window.Resize(fyne.NewSize(800, 700))

	// Load dữ liệu ban đầu
	app.refreshAllLists()
	fmt.Println("📊 Data loaded")

	// Simple keyboard handling - ESC to cancel input
	app.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyEscape && app.showingInput {
			// ESC để hủy input
			app.hideAddInput()
		}
	})

	// Input sẽ được focus khi user click nút Thêm
}

// showAddInput displays the input field for adding new todos
func (app *TodoApp) showAddInput() {
	// Guard: nếu đã đang hiện input thì không làm gì
	if app.showingInput {
		fmt.Println("⚠️ Input already showing, ignoring duplicate call")
		return
	}

	fmt.Println("🎯 Showing add input")
	app.showingInput = true

	// Tạo input layout với buttons
	confirmButton := widget.NewButton("✅ Xác nhận", func() {
		app.addTodo()
	})
	confirmButton.Importance = widget.SuccessImportance

	cancelButton := widget.NewButton("❌ Hủy", app.hideAddInput)
	cancelButton.Importance = widget.LowImportance

	inputWithButtons := container.NewBorder(
		nil, nil, nil,
		container.NewHBox(confirmButton, cancelButton),
		app.addEntry,
	)

	// Clear và thay thế content
	app.inputContainer.Objects = []fyne.CanvasObject{inputWithButtons}
	app.inputContainer.Refresh()

	// Focus vào input field sau khi refresh
	app.window.Canvas().Focus(app.addEntry)
	fmt.Println("🎯 Input field focused")
}

// hideAddInput hides the input field and shows the Add button again
func (app *TodoApp) hideAddInput() {
	// Guard: nếu đã đang ẩn thì không làm gì
	if !app.showingInput {
		fmt.Println("⚠️ Input already hidden, ignoring duplicate call")
		return
	}

	fmt.Println("🔙 Hiding add input")
	app.showingInput = false

	// Clear input text
	app.addEntry.SetText("")

	// Thay thế lại bằng nút Thêm
	app.inputContainer.Objects = []fyne.CanvasObject{app.addButton}
	app.inputContainer.Refresh()
}

func (app *TodoApp) createList(listType string) *widget.List {
	list := widget.NewList(
		func() int {
			switch listType {
			case "all":
				return len(app.allTodos)
			case "active":
				return len(app.activeTodos)
			case "completed":
				return len(app.completedTodos)
			default:
				return 0
			}
		},
		func() fyne.CanvasObject {
			return app.createTodoItem()
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			app.updateTodoItem(id, item, listType)
		},
	)

	// Add double-click handler for completing todos
	list.OnSelected = func(id widget.ListItemID) {
		app.handleTodoSelection(id, listType)
	}

	return list
}

// TodoItem represents a single todo item widget with custom rendering
type TodoItem struct {
	widget.BaseWidget
	todo        Todo              // The todo data
	label       *widget.Label     // Label for displaying todo text
	completeBtn *widget.Button    // Button to mark as complete
	deleteBtn   *widget.Button    // Button to delete todo
	onComplete  func(int)         // Callback when todo is completed
	onDelete    func(int, string) // Callback when todo is deleted
}

func NewTodoItem(todo Todo, onComplete func(int), onDelete func(int, string)) *TodoItem {
	item := &TodoItem{
		todo:       todo,
		onComplete: onComplete,
		onDelete:   onDelete,
	}
	item.ExtendBaseWidget(item)
	return item
}

func (t *TodoItem) CreateRenderer() fyne.WidgetRenderer {
	// Create label
	t.label = widget.NewLabel("")
	t.label.Wrapping = fyne.TextWrapWord

	// Create buttons
	t.completeBtn = widget.NewButton("✅", func() {
		if t.onComplete != nil {
			t.onComplete(t.todo.ID)
		}
	})
	t.deleteBtn = widget.NewButton("🗑️", func() {
		if t.onDelete != nil {
			t.onDelete(t.todo.ID, t.todo.Description)
		}
	})

	// Style buttons
	t.completeBtn.Importance = widget.SuccessImportance
	t.deleteBtn.Importance = widget.DangerImportance

	t.refresh()

	// Create layout
	content := container.NewVBox(
		t.label,
		container.NewHBox(t.completeBtn, t.deleteBtn),
	)

	return widget.NewSimpleRenderer(content)
}

func (t *TodoItem) refresh() {
	// Update label
	status := "📌"
	if t.todo.Completed {
		status = "✅"
	}
	t.label.SetText(fmt.Sprintf("%s %s", status, t.todo.Description))

	// Update complete button
	if t.todo.Completed {
		t.completeBtn.SetText("↩️")
	} else {
		t.completeBtn.SetText("✅")
	}
}

func (t *TodoItem) SetTodo(todo Todo) {
	t.todo = todo
	t.refresh()
}

// TodoItemWidget creates a custom widget for todo items using Card layout
type TodoItemWidget struct {
	widget.Card
	todo        Todo           // The todo data
	completeBtn *widget.Button // Button to complete todo
	deleteBtn   *widget.Button // Button to delete todo
	app         *TodoApp       // Reference to main app
}

func NewTodoItemWidget(todo Todo, app *TodoApp) *TodoItemWidget {
	item := &TodoItemWidget{
		todo: todo,
		app:  app,
	}

	// Tạo buttons
	item.completeBtn = widget.NewButton("✅", func() {
		if todo.Completed {
			dialog.ShowInformation("Thông báo", "Công việc này đã hoàn thành", app.window)
		} else {
			app.markComplete(todo.ID)
		}
	})

	item.deleteBtn = widget.NewButton("🗑️", func() {
		app.confirmDelete(todo.ID, todo.Description)
	})

	// Style buttons
	if todo.Completed {
		item.completeBtn.SetText("✓")
		item.completeBtn.Importance = widget.MediumImportance
	} else {
		item.completeBtn.Importance = widget.SuccessImportance
	}
	item.deleteBtn.Importance = widget.DangerImportance

	// Set up card
	status := "📌"
	if todo.Completed {
		status = "✅"
	}

	item.SetTitle(fmt.Sprintf("%s %s", status, todo.Description))
	item.SetSubTitle(fmt.Sprintf("ID: %d", todo.ID))

	// Nút ở bên phải
	buttonContainer := container.NewHBox(item.completeBtn, item.deleteBtn)
	item.SetContent(buttonContainer)

	return item
}

func (app *TodoApp) createTodoItem() fyne.CanvasObject {
	// Tạo card đơn giản với height cố định
	card := widget.NewCard("", "", widget.NewLabel(""))
	card.Resize(fyne.NewSize(750, 120)) // Height cố định để tránh overlap
	return card
}

func (app *TodoApp) updateTodoItem(id widget.ListItemID, item fyne.CanvasObject, listType string) {
	var todos []Todo
	switch listType {
	case "all":
		todos = app.allTodos
	case "active":
		todos = app.activeTodos
	case "completed":
		todos = app.completedTodos
	default:
		return
	}

	if id >= len(todos) {
		return
	}

	todo := todos[id]
	card := item.(*widget.Card)

	// Reset card title/subtitle
	card.SetTitle("")
	card.SetSubTitle("")

	// Tạo label ngày với font nhỏ, mờ
	dateLabel := widget.NewLabel(todo.CreatedAt.Format("02/01/2006 15:04"))
	dateLabel.TextStyle = fyne.TextStyle{Italic: true}
	dateLabel.Resize(fyne.NewSize(120, 30)) // Fixed width cho date

	// Tạo label nội dung với font lớn, đậm
	contentLabel := widget.NewLabel(todo.Description)
	contentLabel.TextStyle = fyne.TextStyle{Bold: true}
	contentLabel.Wrapping = fyne.TextWrapWord

	// Tạo checkbox cho trạng thái hoàn thành
	var completeCheck *widget.Check
	completeCheck = widget.NewCheck("", func(checked bool) {
		if !todo.Completed && checked {
			// Chỉ cho phép đánh dấu hoàn thành, không cho phép bỏ tích
			app.markComplete(todo.ID)
		} else if todo.Completed && !checked {
			dialog.ShowInformation("Thông báo", "Công việc đã hoàn thành không thể bỏ tích", app.window)
			// Reset lại trạng thái checkbox
			completeCheck.SetChecked(true)
		}
	})
	completeCheck.SetChecked(todo.Completed)
	completeCheck.Resize(fyne.NewSize(30, 30))

	// Nút xóa
	deleteBtn := widget.NewButton("🗑️", func() {
		app.confirmDelete(todo.ID, todo.Description)
	})
	deleteBtn.Resize(fyne.NewSize(40, 30))

	// Buttons container với checkbox và nút xóa
	buttonsContainer := container.NewHBox(completeCheck, deleteBtn)
	buttonsContainer.Resize(fyne.NewSize(80, 35))

	// Layout ngang: ngày bên trái, nội dung ở giữa (expand), buttons bên phải
	horizontalLayout := container.NewBorder(
		nil, nil,
		dateLabel,        // Trái: ngày tạo
		buttonsContainer, // Phải: buttons
		contentLabel,     // Giữa: nội dung (sẽ expand)
	)

	card.SetContent(container.NewPadded(horizontalLayout))
}

func (app *TodoApp) addTodo() {
	// Get text from entry widget và loại bỏ newlines
	description := strings.ReplaceAll(strings.TrimSpace(app.addEntry.Text), "\n", " ")
	description = strings.ReplaceAll(description, "\r", " ")
	fmt.Printf("addTodo called - Entry text: '%s'\n", app.addEntry.Text)
	fmt.Printf("After trim: '%s'\n", description)

	if description == "" {
		fmt.Println("Description is empty, showing error")
		dialog.ShowError(fmt.Errorf("vui lòng nhập mô tả công việc"), app.window)
		return
	}

	err := app.todoList.AddTodo(description)
	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}

	// Hide input field và về trạng thái ban đầu
	app.hideAddInput()

	// Refresh lists
	app.refreshAllLists()

	// Show success message
	dialog.ShowInformation("Thành công",
		fmt.Sprintf("Đã thêm công việc: %s", description), app.window)
}

func (app *TodoApp) markComplete(todoID int) {
	// Tìm todo để lấy description trước khi mark complete
	var todoDescription string
	for _, todo := range app.allTodos {
		if todo.ID == todoID {
			todoDescription = todo.Description
			break
		}
	}

	err := app.todoList.MarkComplete(todoID)
	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}

	app.refreshAllLists()

	// Hiển thị animation pháo hoa với physics thực tế
	fireworks.ShowFireworksDialog(todoDescription, app.window)
}

func (app *TodoApp) confirmDelete(todoID int, description string) {
	dialog.ShowConfirm("Xác nhận xóa",
		fmt.Sprintf("Bạn có chắc chắn muốn xóa công việc:\n'%s'?", description),
		func(confirmed bool) {
			if confirmed {
				err := app.todoList.DeleteTodo(todoID)
				if err != nil {
					dialog.ShowError(err, app.window)
					return
				}
				app.refreshAllLists()
				dialog.ShowInformation("Thành công",
					fmt.Sprintf("Đã xóa công việc: %s", description), app.window)
			}
		}, app.window)
}

func (app *TodoApp) refreshAllLists() {
	// Update data arrays
	app.allTodos = app.todoList.GetTodos()
	app.activeTodos = app.todoList.GetActiveTodos()
	app.completedTodos = app.todoList.GetCompletedTodos()

	// Refresh all lists
	app.allList.Refresh()
	app.activeList.Refresh()
	app.completedList.Refresh()
}

func (app *TodoApp) handleTodoSelection(id widget.ListItemID, listType string) {
	var todos []Todo
	switch listType {
	case "all":
		todos = app.allTodos
	case "active":
		todos = app.activeTodos
	case "completed":
		todos = app.completedTodos
	default:
		return
	}

	if id >= len(todos) {
		return
	}

	todo := todos[id]

	// Show action dialog for the selected todo
	if todo.Completed {
		dialog.ShowConfirm("Công việc đã hoàn thành",
			fmt.Sprintf("Công việc: %s\nBạn muốn xóa công việc này?", todo.Description),
			func(confirmed bool) {
				if confirmed {
					app.confirmDelete(todo.ID, todo.Description)
				}
			}, app.window)
	} else {
		// For incomplete todos, show options
		completeBtn := widget.NewButton("✅ Đánh dấu hoàn thành", func() {
			app.markComplete(todo.ID)
		})
		completeBtn.Importance = widget.SuccessImportance

		deleteBtn := widget.NewButton("🗑️ Xóa công việc", func() {
			app.confirmDelete(todo.ID, todo.Description)
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

// showSettingsDialog hiển thị dialog cài đặt theme
func (app *TodoApp) showSettingsDialog() {
	// Tạo label để mô tả switch
	switchLabel := widget.NewLabel("Chế độ theme:")
	switchLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Label hiển thị trạng thái theme hiện tại
	themeLabel := widget.NewLabel(app.getThemeLabelText())

	// Tạo switch để chuyển đổi giữa theme sáng và tối
	var themeSwitch *widget.Button
	themeSwitch = widget.NewButton("", func() {
		app.isDarkTheme = !app.isDarkTheme
		app.applyTheme()
		app.updateSwitchAppearance(themeSwitch)
		themeLabel.SetText(app.getThemeLabelText())
	})

	// Khởi tạo appearance ban đầu
	app.updateSwitchAppearance(themeSwitch)

	// Thông tin hướng dẫn
	infoLabel := widget.NewLabel("Chọn giao diện sáng hoặc tối cho ứng dụng")
	infoLabel.TextStyle = fyne.TextStyle{Italic: true}

	// Layout cho dialog
	content := container.NewVBox(
		infoLabel,
		widget.NewSeparator(),
		themeLabel,
		themeSwitch,
	)

	// Hiển thị dialog
	dialog.ShowCustom("⚙️ Cài đặt Theme", "Đóng", content, app.window)
}

// applyTheme áp dụng theme sáng hoặc tối cho giao diện
func (app *TodoApp) applyTheme() {
	if app.isDarkTheme {
		// Áp dụng theme tối - sử dụng default theme với dark variant
		os.Setenv("FYNE_THEME", "dark")
		fmt.Println("🌙 Switched to dark theme")
	} else {
		// Áp dụng theme sáng - sử dụng default theme với light variant
		os.Setenv("FYNE_THEME", "light")
		fmt.Println("☀️ Switched to light theme")
	}

	// Refresh toàn bộ UI để áp dụng thay đổi
	app.window.Content().Refresh()
	app.refreshAllLists()
}

// getThemeLabelText trả về text cho theme label
func (app *TodoApp) getThemeLabelText() string {
	if app.isDarkTheme {
		return "🌙 Theme hiện tại: Tối"
	} else {
		return "☀️ Theme hiện tại: Sáng"
	}
}

// updateSwitchAppearance cập nhật appearance của switch button
func (app *TodoApp) updateSwitchAppearance(btn *widget.Button) {
	if app.isDarkTheme {
		btn.SetText("🌙 TỐI")
		btn.Importance = widget.HighImportance
	} else {
		btn.SetText("☀️ SÁNG")
		btn.Importance = widget.LowImportance
	}
}

// Project management methods

// createProjectsView creates the projects management view for navbar
func (app *TodoApp) createProjectsView() *fyne.Container {
	// Header với thông tin projects và settings
	projectHeader := widget.NewCard("", "📁 Projects Manager",
		widget.NewLabel("Tạo và quản lý các dự án todo riêng biệt"))

	// Tạo container cho header với settings button
	headerWithSettings := container.NewBorder(
		nil, nil, nil, app.settingsButton,
		projectHeader,
	)

	// Nút tạo project mới
	createProjectBtn := widget.NewButton("➕ Tạo Project Mới", func() {
		app.showCreateProjectDialog()
	})
	createProjectBtn.Importance = widget.HighImportance

	// Tạo list để hiển thị các projects
	app.projectsList = widget.NewList(
		func() int {
			projects := app.getProjectList()
			return len(projects)
		},
		func() fyne.CanvasObject {
			return app.createProjectItem()
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			projects := app.getProjectList()
			if id < len(projects) {
				app.updateProjectItem(id, item, projects[id])
			}
		},
	)

	// Layout cho projects view với header đẹp hơn
	content := container.NewBorder(
		container.NewVBox(headerWithSettings, createProjectBtn, widget.NewSeparator()),
		nil, nil, nil,
		container.NewPadded(app.projectsList),
	)

	return content
}

// showCreateProjectDialog hiển thị dialog tạo project mới
func (app *TodoApp) showCreateProjectDialog() {
	projectNameEntry := widget.NewEntry()
	projectNameEntry.SetPlaceHolder("Nhập tên project...")

	projectDescEntry := widget.NewMultiLineEntry()
	projectDescEntry.SetPlaceHolder("Mô tả project (tùy chọn)...")
	projectDescEntry.Resize(fyne.NewSize(300, 60))

	form := container.NewVBox(
		widget.NewLabel("Tên Project:"),
		projectNameEntry,
		widget.NewLabel("Mô tả:"),
		projectDescEntry,
	)

	dialog.ShowCustomConfirm("Tạo Project Mới", "Tạo", "Hủy", form, func(confirmed bool) {
		if confirmed {
			projectName := strings.TrimSpace(projectNameEntry.Text)
			if projectName == "" {
				dialog.ShowError(fmt.Errorf("tên project không được để trống"), app.window)
				return
			}

			// Kiểm tra project đã tồn tại
			if app.projectExists(projectName) {
				dialog.ShowError(fmt.Errorf("project '%s' đã tồn tại", projectName), app.window)
				return
			}

			// Tạo project mới
			err := app.createProject(projectName, projectDescEntry.Text)
			if err != nil {
				dialog.ShowError(err, app.window)
				return
			}

			// Refresh projects list
			app.projectsList.Refresh()

			dialog.ShowInformation("Thành công",
				fmt.Sprintf("Project '%s' đã được tạo thành công!", projectName), app.window)
		}
	}, app.window)
}

// getProjectList trả về danh sách tất cả projects
func (app *TodoApp) getProjectList() []string {
	files, err := os.ReadDir("data/project")
	if err != nil {
		return []string{}
	}

	var projects []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			projectName := strings.TrimSuffix(file.Name(), ".txt")
			projects = append(projects, projectName)
		}
	}
	return projects
}

// projectExists kiểm tra project có tồn tại không
func (app *TodoApp) projectExists(projectName string) bool {
	filename := fmt.Sprintf("data/project/%s.txt", projectName)
	_, err := os.Stat(filename)
	return err == nil
}

// createProject tạo project mới
func (app *TodoApp) createProject(projectName, description string) error {
	filename := fmt.Sprintf("data/project/%s.txt", projectName)

	// Tạo file project với header comment
	header := fmt.Sprintf("# Project: %s\n# Description: %s\n# Created: %s\n\n",
		projectName, description, time.Now().Format("02/01/2006 15:04"))

	err := os.WriteFile(filename, []byte(header), 0644)
	if err != nil {
		return fmt.Errorf("không thể tạo project file: %v", err)
	}

	return nil
}

// createProjectItem tạo widget cho project item
func (app *TodoApp) createProjectItem() fyne.CanvasObject {
	card := widget.NewCard("", "", widget.NewLabel(""))
	card.Resize(fyne.NewSize(700, 80))
	return card
}

// updateProjectItem cập nhật project item
func (app *TodoApp) updateProjectItem(id widget.ListItemID, item fyne.CanvasObject, projectName string) {
	card := item.(*widget.Card)

	// Đếm số todos trong project
	todoCount := app.getProjectTodoCount(projectName)

	// Tạo label info
	infoLabel := widget.NewLabel(fmt.Sprintf("📋 %d todos", todoCount))
	infoLabel.TextStyle = fyne.TextStyle{Italic: true}

	// Nút mở project
	openBtn := widget.NewButton("📂 Mở", func() {
		app.openProject(projectName)
	})
	openBtn.Importance = widget.SuccessImportance

	// Nút xóa project
	deleteBtn := widget.NewButton("🗑️", func() {
		app.confirmDeleteProject(projectName)
	})
	deleteBtn.Importance = widget.DangerImportance

	// Layout ngang
	layout := container.NewBorder(
		nil, nil,
		widget.NewLabel("📁 "+projectName),     // Trái: tên project
		container.NewHBox(openBtn, deleteBtn), // Phải: buttons
		infoLabel,                             // Giữa: info
	)

	card.SetContent(container.NewPadded(layout))
}

// getProjectTodoCount đếm số todos trong project
func (app *TodoApp) getProjectTodoCount(projectName string) int {
	filename := fmt.Sprintf("data/project/%s.txt", projectName)
	todoList := NewTodoList(filename)
	todos := todoList.GetTodos()
	return len(todos)
}

// openProject mở project trong cửa sổ mới hoặc tab mới
func (app *TodoApp) openProject(projectName string) {
	app.currentProject = projectName
	filename := fmt.Sprintf("data/project/%s.txt", projectName)
	app.projectTodoList = NewTodoList(filename)

	// Tạo cửa sổ mới cho project
	projectWindow := app.myApp.NewWindow(fmt.Sprintf("📁 Project: %s", projectName))
	projectWindow.Resize(fyne.NewSize(800, 600))
	projectWindow.CenterOnScreen()

	// Tạo UI cho project window
	app.setupProjectWindow(projectWindow, projectName)

	projectWindow.Show()
}

// confirmDeleteProject xác nhận xóa project
func (app *TodoApp) confirmDeleteProject(projectName string) {
	dialog.ShowConfirm("Xóa Project",
		fmt.Sprintf("Bạn có chắc chắn muốn xóa project '%s'?\nTất cả dữ liệu sẽ bị mất vĩnh viễn!", projectName),
		func(confirmed bool) {
			if confirmed {
				err := app.deleteProject(projectName)
				if err != nil {
					dialog.ShowError(err, app.window)
					return
				}

				app.projectsList.Refresh()
				dialog.ShowInformation("Thành công",
					fmt.Sprintf("Project '%s' đã được xóa!", projectName), app.window)
			}
		}, app.window)
}

// deleteProject xóa project
func (app *TodoApp) deleteProject(projectName string) error {
	filename := fmt.Sprintf("data/project/%s.txt", projectName)
	return os.Remove(filename)
}

// setupProjectWindow thiết lập UI cho cửa sổ project
func (app *TodoApp) setupProjectWindow(projectWindow fyne.Window, projectName string) {
	// Tạo todo lists cho project
	allProjectTodos := app.createProjectList("all")
	activeProjectTodos := app.createProjectList("active")
	completedProjectTodos := app.createProjectList("completed")

	// Tạo tabs cho project
	app.projectTabs = container.NewAppTabs(
		container.NewTabItem("Tất cả", allProjectTodos),
		container.NewTabItem("Chưa hoàn thành", activeProjectTodos),
		container.NewTabItem("Đã hoàn thành", completedProjectTodos),
	)

	// Input để thêm todo mới cho project
	projectAddEntry := widget.NewEntry()
	projectAddEntry.SetPlaceHolder("Nhập todo cho project " + projectName + "...")

	addProjectTodoBtn := widget.NewButton("➕ Thêm Todo", func() {
		app.addProjectTodo(projectAddEntry, projectName)
	})
	addProjectTodoBtn.Importance = widget.HighImportance

	// Header project
	header := widget.NewCard("", fmt.Sprintf("📁 Project: %s", projectName), nil)

	// Layout cho project window
	content := container.NewBorder(
		container.NewVBox(
			header,
			container.NewBorder(nil, nil, nil, addProjectTodoBtn, projectAddEntry),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		app.projectTabs,
	)

	projectWindow.SetContent(content)
}

// createProjectList tạo list cho project todos
func (app *TodoApp) createProjectList(listType string) *widget.List {
	list := widget.NewList(
		func() int {
			if app.projectTodoList == nil {
				return 0
			}
			switch listType {
			case "all":
				return len(app.projectTodoList.GetTodos())
			case "active":
				return len(app.projectTodoList.GetActiveTodos())
			case "completed":
				return len(app.projectTodoList.GetCompletedTodos())
			default:
				return 0
			}
		},
		func() fyne.CanvasObject {
			return app.createTodoItem()
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			app.updateProjectTodoItem(id, item, listType)
		},
	)
	return list
}

// updateProjectTodoItem cập nhật project todo item
func (app *TodoApp) updateProjectTodoItem(id widget.ListItemID, item fyne.CanvasObject, listType string) {
	if app.projectTodoList == nil {
		return
	}

	var todos []Todo
	switch listType {
	case "all":
		todos = app.projectTodoList.GetTodos()
	case "active":
		todos = app.projectTodoList.GetActiveTodos()
	case "completed":
		todos = app.projectTodoList.GetCompletedTodos()
	default:
		return
	}

	if id >= len(todos) {
		return
	}

	todo := todos[id]
	card := item.(*widget.Card)

	// Reset card title/subtitle
	card.SetTitle("")
	card.SetSubTitle("")

	// Tạo label ngày với font nhỏ, mờ
	dateLabel := widget.NewLabel(todo.CreatedAt.Format("02/01/2006 15:04"))
	dateLabel.TextStyle = fyne.TextStyle{Italic: true}
	dateLabel.Resize(fyne.NewSize(120, 30))

	// Tạo label nội dung với font lớn, đậm
	contentLabel := widget.NewLabel(todo.Description)
	contentLabel.TextStyle = fyne.TextStyle{Bold: true}
	contentLabel.Wrapping = fyne.TextWrapWord

	// Tạo checkbox cho trạng thái hoàn thành
	var completeCheck *widget.Check
	completeCheck = widget.NewCheck("", func(checked bool) {
		if !todo.Completed && checked {
			app.markProjectTodoComplete(todo.ID)
		} else if todo.Completed && !checked {
			dialog.ShowInformation("Thông báo", "Công việc đã hoàn thành không thể bỏ tích", app.window)
			completeCheck.SetChecked(true)
		}
	})
	completeCheck.SetChecked(todo.Completed)
	completeCheck.Resize(fyne.NewSize(30, 30))

	// Nút xóa
	deleteBtn := widget.NewButton("🗑️", func() {
		app.confirmDeleteProjectTodo(todo.ID, todo.Description)
	})
	deleteBtn.Resize(fyne.NewSize(40, 30))

	// Buttons container
	buttonsContainer := container.NewHBox(completeCheck, deleteBtn)
	buttonsContainer.Resize(fyne.NewSize(80, 35))

	// Layout ngang
	horizontalLayout := container.NewBorder(
		nil, nil,
		dateLabel,
		buttonsContainer,
		contentLabel,
	)

	card.SetContent(container.NewPadded(horizontalLayout))
}

// addProjectTodo thêm todo mới cho project
func (app *TodoApp) addProjectTodo(entry *widget.Entry, projectName string) {
	description := strings.TrimSpace(entry.Text)
	if description == "" {
		dialog.ShowError(fmt.Errorf("vui lòng nhập mô tả công việc"), app.window)
		return
	}

	if app.projectTodoList == nil {
		dialog.ShowError(fmt.Errorf("project chưa được khởi tạo"), app.window)
		return
	}

	err := app.projectTodoList.AddTodo(description)
	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}

	// Clear input
	entry.SetText("")

	// Refresh project tabs
	if app.projectTabs != nil {
		app.projectTabs.Refresh()
	}

	dialog.ShowInformation("Thành công",
		fmt.Sprintf("Đã thêm công việc vào project %s: %s", projectName, description), app.window)
}

// markProjectTodoComplete đánh dấu hoàn thành todo trong project
func (app *TodoApp) markProjectTodoComplete(todoID int) {
	if app.projectTodoList == nil {
		return
	}

	// Tìm todo để lấy description
	var todoDescription string
	for _, todo := range app.projectTodoList.GetTodos() {
		if todo.ID == todoID {
			todoDescription = todo.Description
			break
		}
	}

	err := app.projectTodoList.MarkComplete(todoID)
	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}

	// Refresh project tabs
	if app.projectTabs != nil {
		app.projectTabs.Refresh()
	}

	// Hiển thị pháo hoa
	fireworks.ShowFireworksDialog(todoDescription, app.window)
}

// confirmDeleteProjectTodo xác nhận xóa todo trong project
func (app *TodoApp) confirmDeleteProjectTodo(todoID int, description string) {
	dialog.ShowConfirm("Xác nhận xóa",
		fmt.Sprintf("Bạn có chắc chắn muốn xóa công việc:\n'%s'?", description),
		func(confirmed bool) {
			if confirmed {
				if app.projectTodoList == nil {
					return
				}

				err := app.projectTodoList.DeleteTodo(todoID)
				if err != nil {
					dialog.ShowError(err, app.window)
					return
				}

				// Refresh project tabs
				if app.projectTabs != nil {
					app.projectTabs.Refresh()
				}

				dialog.ShowInformation("Thành công",
					fmt.Sprintf("Đã xóa công việc: %s", description), app.window)
			}
		}, app.window)
}

// toggleNavbarMode chuyển đổi giữa chế độ bình thường và navbar dọc
func (app *TodoApp) toggleNavbarMode() {
	app.isNavbarMode = !app.isNavbarMode
	
	if app.isNavbarMode {
		// Chuyển sang chế độ navbar dọc
		app.showNavbarMode()
	} else {
		// Quay về chế độ bình thường
		app.showNormalMode()
	}
}

// showNavbarMode hiển thị giao diện navbar dọc
func (app *TodoApp) showNavbarMode() {
	// Tạo header cho navbar mode
	header := widget.NewCard("", "📱 Navigation Mode", nil)
	
	// Header với nút quay về và settings
	headerWithButtons := container.NewBorder(
		nil, nil, 
		widget.NewButton("←", func() { app.toggleNavbarMode() }), // Nút quay về
		app.settingsButton, // Nút settings
		header,
	)
	
	// Tạo view cho todos (không có input ở đây)
	todosView := container.NewBorder(
		container.NewVBox(widget.NewLabel("📋 Todos Management"), widget.NewSeparator()),
		nil, nil, nil,
		app.tabs,
	)
	
	// Tạo view cho projects
	projectsView := app.createProjectsView()
	
	// Tạo nút + cho projects
	addProjectBtn := widget.NewButton("+", func() {
		app.showCreateProjectDialog()
	})
	addProjectBtn.Importance = widget.HighImportance

	// Tạo dropdown để chọn project
	projectOptions := app.getProjectList()
	if len(projectOptions) == 0 {
		projectOptions = []string{"Chưa có project nào"}
	}
	
	projectSelect := widget.NewSelect(projectOptions, func(selected string) {
		if selected != "" && selected != "Chưa có project nào" {
			app.openProject(selected)
		}
	})
	
	// Set placeholder cho select
	if len(app.getProjectList()) > 0 {
		projectSelect.SetSelected("") // Clear selection để hiển thị placeholder
	}

	// Layout cho tab Projects với nút + và dropdown
	projectsTabContent := container.NewVBox(
		container.NewBorder(nil, nil, nil, addProjectBtn, 
			widget.NewLabel("📁 Projects")),
		projectSelect,
		widget.NewSeparator(),
		projectsView,
	)

	// Tạo navbar dọc với 2 tabs
	app.navbar = container.NewAppTabs(
		container.NewTabItem("📋 Todos", todosView),
		container.NewTabItem("📁 Projects", projectsTabContent),
	)
	
	// Layout chính cho navbar mode
	navbarContent := container.NewBorder(
		headerWithButtons,
		nil, nil, nil,
		app.navbar,
	)
	
	app.window.SetContent(navbarContent)
	app.navbarButton.SetText("←") // Đổi icon thành mũi tên quay về
}

// showNormalMode hiển thị giao diện bình thường
func (app *TodoApp) showNormalMode() {
	// Tạo header bình thường
	header := widget.NewCard("", "Todo List Desktop App", nil)
	header.Resize(fyne.NewSize(600, 60))
	
	// Container chính sẽ switch giữa nút Thêm và input field
	app.inputContainer = container.NewVBox(app.addButton)
	paddedInput := container.NewPadded(app.inputContainer)
	
	// Header với nút navbar và settings
	headerWithButtons := container.NewBorder(
		nil, nil, 
		app.navbarButton, // Trái: nút navbar
		container.NewHBox(app.settingsButton), // Phải: nút settings
		header,
	)
	
	// Tạo view chính
	mainView := container.NewBorder(
		container.NewVBox(headerWithButtons, paddedInput, widget.NewSeparator()),
		nil, nil, nil,
		app.tabs,
	)
	
	app.window.SetContent(mainView)
	app.navbarButton.SetText("≡") // Đổi icon về dấu "≡"
}
