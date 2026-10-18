package main

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// TaskState описывает жизненный цикл задачи[cite: 1]
type TaskState string

const (
	StatePending    TaskState = "Pending"
	StateProcessing TaskState = "Processing"
	StateCompleted  TaskState = "Completed"
)

// Result хранит результат нашей работы (количество слов и символов)[cite: 1]
type Result struct {
	Words      int
	Characters int
}

// Task — это сама задача со всеми ее данными[cite: 1]
type Task struct {
	ID     int
	Input  string
	State  TaskState
	Result *Result
}

// TaskManager — интерфейс, который описывает, что должно уметь приложение[cite: 1, 3]
type TaskManager interface {
	CreateTask(input string) (int, error)
	GetTask(id int) (*Task, error)
	ProcessTask(id int) error
}

// InMemoryManager хранит задачи прямо в оперативной памяти (без внешних баз данных)[cite: 1, 3]
type InMemoryManager struct {
	tasks  map[int]*Task
	nextID int
}

// NewInMemoryManager создает новое хранилище
func NewInMemoryManager() *InMemoryManager {
	return &InMemoryManager{
		tasks:  make(map[int]*Task),
		nextID: 1,
	}
}

// CreateTask проверяет текст и создает новую задачу[cite: 1]
func (m *InMemoryManager) CreateTask(input string) (int, error) {
	// Валидация: если текст пустой или состоит только из пробелов — отклоняем и выдаем ошибку[cite: 1]
	if strings.TrimSpace(input) == "" {
		return 0, errors.New("invalid input: text cannot be empty")
	}

	task := &Task{
		ID:    m.nextID,
		Input: input,
		State: StatePending,
	}
	m.tasks[task.ID] = task
	m.nextID++

	return task.ID, nil
}

// GetTask ищет задачу по ее ID[cite: 1]
func (m *InMemoryManager) GetTask(id int) (*Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

// ProcessTask выполняет подсчет слов и символов[cite: 1]
func (m *InMemoryManager) ProcessTask(id int) error {
	task, err := m.GetTask(id)
	if err != nil {
		return err
	}

	// Защита от неправильных действий: если задача уже выполнена, её нельзя выполнить снова[cite: 1]
	if task.State == StateCompleted {
		return errors.New("invalid operation: task is already completed")
	}

	// Переводим в статус "В процессе"[cite: 1]
	task.State = StateProcessing

	// Считаем слова и буквы в тексте[cite: 1]
	words := len(strings.Fields(task.Input))
	chars := utf8.RuneCountInString(task.Input)

	task.Result = &Result{
		Words:      words,
		Characters: chars,
	}

	// Ставим статус "Завершено"[cite: 1]
	task.State = StateCompleted
	return nil
}

// функция main обязательна для запускаемой программы на Go
func main() {
}
