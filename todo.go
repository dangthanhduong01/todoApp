package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Todo represents a single todo item
type Todo struct {
	ID          int
	Description string
	Completed   bool
	CreatedAt   time.Time
}

// TodoList manages the list of todos and file operations
type TodoList struct {
	todos    []Todo
	filename string
	nextID   int
}

// NewTodoList creates a new TodoList instance
func NewTodoList(filename string) *TodoList {
	tl := &TodoList{
		todos:    []Todo{},
		filename: filename,
		nextID:   1,
	}
	tl.LoadFromFile()
	return tl
}

// LoadFromFile loads todos from the text file
func (tl *TodoList) LoadFromFile() {
	file, err := os.Open(tl.filename)
	if err != nil {
		// File doesn't exist yet, which is fine for a new todolist
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			// Skip empty lines and header metadata
			continue
		}

		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		completed := parts[2] == "true"
		createdAt, err := time.Parse(time.RFC3339, parts[3])
		if err != nil {
			createdAt = time.Now()
		}

		todo := Todo{
			ID:          id,
			Description: parts[1],
			Completed:   completed,
			CreatedAt:   createdAt,
		}

		tl.todos = append(tl.todos, todo)
		if id >= tl.nextID {
			tl.nextID = id + 1
		}
	}
}

// SaveToFile saves todos to the text file while preserving header metadata
func (tl *TodoList) SaveToFile() error {
	// Read existing header metadata
	var headerLines []string
	if existingFile, err := os.Open(tl.filename); err == nil {
		defer existingFile.Close()
		scanner := bufio.NewScanner(existingFile)
		for scanner.Scan() {
			line := scanner.Text()
			// Keep header lines (starting with #) and empty lines before todos
			if strings.HasPrefix(line, "#") || (strings.TrimSpace(line) == "" && len(headerLines) > 0) {
				headerLines = append(headerLines, line)
			} else if strings.Contains(line, "|") {
				// Stop when we hit todo data
				break
			}
		}
	}

	// Create new file
	file, err := os.Create(tl.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write header metadata first
	for _, headerLine := range headerLines {
		writer.WriteString(headerLine + "\n")
	}

	// Write todo data
	for _, todo := range tl.todos {
		line := fmt.Sprintf("%d|%s|%t|%s\n",
			todo.ID,
			todo.Description,
			todo.Completed,
			todo.CreatedAt.Format(time.RFC3339),
		)
		writer.WriteString(line)
	}

	return nil
}

// AddTodo adds a new todo item
func (tl *TodoList) AddTodo(description string) error {
	if strings.TrimSpace(description) == "" {
		return fmt.Errorf("mô tả không được để trống")
	}

	todo := Todo{
		ID:          tl.nextID,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}

	tl.todos = append(tl.todos, todo)
	tl.nextID++
	return tl.SaveToFile()
}

// MarkComplete marks a todo as completed
func (tl *TodoList) MarkComplete(id int) error {
	for i := range tl.todos {
		if tl.todos[i].ID == id {
			tl.todos[i].Completed = true
			return tl.SaveToFile()
		}
	}
	return fmt.Errorf("không tìm thấy công việc với ID %d", id)
}

// DeleteTodo removes a todo item
func (tl *TodoList) DeleteTodo(id int) error {
	for i, todo := range tl.todos {
		if todo.ID == id {
			tl.todos = append(tl.todos[:i], tl.todos[i+1:]...)
			return tl.SaveToFile()
		}
	}
	return fmt.Errorf("không tìm thấy công việc với ID %d", id)
}

// GetTodos returns all todos
func (tl *TodoList) GetTodos() []Todo {
	return tl.todos
}

// GetActiveTodos returns only incomplete todos
func (tl *TodoList) GetActiveTodos() []Todo {
	var active []Todo
	for _, todo := range tl.todos {
		if !todo.Completed {
			active = append(active, todo)
		}
	}
	return active
}

// GetCompletedTodos returns only completed todos
func (tl *TodoList) GetCompletedTodos() []Todo {
	var completed []Todo
	for _, todo := range tl.todos {
		if todo.Completed {
			completed = append(completed, todo)
		}
	}
	return completed
}

// ProjectList extends TodoList with project-specific features
type ProjectList struct {
	*TodoList              // Embedded TodoList for inheritance
	Color           string // Project color theme
	Theme           string // Theme name (e.g., "blue", "red", "green")
	Name            string // Project name
	BackgroundImage string // Background image path for project theme
}

// NewProjectList creates a new ProjectList instance
func NewProjectList(filename, name, color, theme, backgroundImage string) *ProjectList {
	pl := &ProjectList{
		TodoList:        NewTodoList(filename),
		Color:           color,
		Theme:           theme,
		Name:            name,
		BackgroundImage: backgroundImage,
	}
	return pl
}

// GetColor returns the project color
func (pl *ProjectList) GetColor() string {
	return pl.Color
}

// SetColor sets the project color
func (pl *ProjectList) SetColor(color string) {
	pl.Color = color
}

// GetTheme returns the project theme
func (pl *ProjectList) GetTheme() string {
	return pl.Theme
}

// SetTheme sets the project theme
func (pl *ProjectList) SetTheme(theme string) {
	pl.Theme = theme
}

// GetName returns the project name
func (pl *ProjectList) GetName() string {
	return pl.Name
}

// GetBackgroundImage returns the project background image path
func (pl *ProjectList) GetBackgroundImage() string {
	return pl.BackgroundImage
}

// SetBackgroundImage sets the project background image path
func (pl *ProjectList) SetBackgroundImage(imagePath string) {
	pl.BackgroundImage = imagePath
}

// HasBackgroundImage checks if project has a background image
func (pl *ProjectList) HasBackgroundImage() bool {
	return pl.BackgroundImage != ""
}
