package main

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TodoApp struct {
	todoList       *TodoList
	window         fyne.Window
	allList        *widget.List
	activeList     *widget.List
	completedList  *widget.List
	addEntry       *widget.Entry
	allTodos       []Todo
	activeTodos    []Todo
	completedTodos []Todo
	tabs           *container.AppTabs
	inputContainer *fyne.Container
	addButton      *widget.Button
	showingInput   bool
}

func main() {
	fmt.Println("🚀 Starting Todo App...")

	// Force software rendering to fix input display issues
	os.Setenv("FYNE_DRIVER", "x11")
	os.Setenv("FYNE_SOFTWARE", "1")
	os.Setenv("FYNE_DISABLE_HARDWARE_ACCELERATION", "1")
	os.Setenv("GTK_IM_MODULE", "ibus")
	os.Setenv("QT_IM_MODULE", "ibus")
	os.Setenv("XMODIFIERS", "@im=ibus")

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
		todoList: NewTodoList("todos.txt"),
		window:   myWindow,
	}
	fmt.Println("📋 TodoApp struct created")

	todoApp.setupUI()
	fmt.Println("🎛️ UI setup complete")

	// Show window explicitly before running
	myWindow.Show()
	fmt.Println("👁️ Window shown, starting main loop...")
	myApp.Run()
}

func (app *TodoApp) setupUI() {
	// Tạo nút "Thêm" ban đầu với simple handler
	app.addButton = widget.NewButton("➕ Thêm công việc mới", func() {
		fmt.Println("📱 Add button clicked")
		app.showAddInput()
	})
	app.addButton.Importance = widget.HighImportance

	// Tạo widget.Entry thông thường
	app.addEntry = widget.NewEntry()
	app.addEntry.SetPlaceHolder("Nhập công việc mới...")

	// Set change handler
	app.addEntry.OnChanged = func(text string) {
		fmt.Printf("📝 Text changed: '%s'\n", text)
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

	// Tabs
	app.tabs = container.NewAppTabs(
		container.NewTabItem("Tất cả", app.allList),
		container.NewTabItem("Chưa hoàn thành", app.activeList),
		container.NewTabItem("Đã hoàn thành", app.completedList),
	)

	// Nội dung chính
	content := container.NewVBox(
		widget.NewCard("", "Todo List Desktop App", nil),
		paddedInput,
		widget.NewSeparator(),
		app.tabs,
	)

	app.window.SetContent(content)
	app.window.Resize(fyne.NewSize(600, 500))

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

// showAddInput hiển thị input field để nhập todo mới
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

// hideAddInput ẩn input field và hiện lại nút Thêm
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

// TodoItem represents a single todo item widget
type TodoItem struct {
	widget.BaseWidget
	todo        Todo
	label       *widget.Label
	completeBtn *widget.Button
	deleteBtn   *widget.Button
	onComplete  func(int)
	onDelete    func(int, string)
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

// TodoItemWidget tạo widget tùy chỉnh cho todo item
type TodoItemWidget struct {
	widget.Card
	todo        Todo
	completeBtn *widget.Button
	deleteBtn   *widget.Button
	app         *TodoApp
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
	// Tạo card trống sẽ được cập nhật sau
	return widget.NewCard("", "", widget.NewLabel(""))
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

	// Cập nhật card
	status := "📌"
	if todo.Completed {
		status = "✅"
	}

	card.SetTitle(fmt.Sprintf("%s %s", status, todo.Description))
	card.SetSubTitle(fmt.Sprintf("ID: %d", todo.ID))

	// Tạo buttons mới
	completeBtn := widget.NewButton("✅", func() {
		if todo.Completed {
			dialog.ShowInformation("Thông báo", "Công việc này đã hoàn thành", app.window)
		} else {
			app.markComplete(todo.ID)
		}
	})

	deleteBtn := widget.NewButton("🗑️", func() {
		app.confirmDelete(todo.ID, todo.Description)
	})

	// Style buttons
	if todo.Completed {
		completeBtn.SetText("✓")
		completeBtn.Importance = widget.MediumImportance
	} else {
		completeBtn.Importance = widget.SuccessImportance
	}
	deleteBtn.Importance = widget.DangerImportance

	// Đặt buttons ở bên phải
	buttonContainer := container.NewHBox(completeBtn, deleteBtn)
	card.SetContent(buttonContainer)
}

func (app *TodoApp) addTodo() {
	// Get text from entry widget
	description := strings.TrimSpace(app.addEntry.Text)
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
	err := app.todoList.MarkComplete(todoID)
	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}

	app.refreshAllLists()
	dialog.ShowInformation("Thành công",
		fmt.Sprintf("Đã đánh dấu hoàn thành công việc ID %d", todoID), app.window)
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

// Legacy methods to maintain compatibility (can be removed later)
func (app *TodoApp) refreshList() {
	app.refreshAllLists()
}

func (app *TodoApp) showCompleteDialog() {
	// This method is now replaced by individual buttons, but kept for compatibility
	activeTodos := app.todoList.GetActiveTodos()
	if len(activeTodos) == 0 {
		dialog.ShowInformation("Thông báo", "Không có công việc nào chưa hoàn thành", app.window)
		return
	}
	dialog.ShowInformation("Hướng dẫn", "Sử dụng nút ✅ bên cạnh mỗi công việc để đánh dấu hoàn thành", app.window)
}

func (app *TodoApp) showDeleteDialog() {
	// This method is now replaced by individual buttons, but kept for compatibility
	todos := app.todoList.GetTodos()
	if len(todos) == 0 {
		dialog.ShowInformation("Thông báo", "Không có công việc nào để xóa", app.window)
		return
	}
	dialog.ShowInformation("Hướng dẫn", "Sử dụng nút 🗑️ bên cạnh mỗi công việc để xóa", app.window)
}
