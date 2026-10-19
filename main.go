package main

import (
	"errors"
	"fmt" // Добавили пакет fmt для вывода текста в терминал
	"strings"
	"unicode/utf8"
)

// TaskState описывает жизненный цикл задачи
type TaskState string

const (
	StatePending    TaskState = "Pending"
	StateProcessing TaskState = "Processing"
	StateCompleted  TaskState = "Completed"
)

// Result хранит результат нашей работы (количество слов и символов)
type Result struct {
	Words      int
	Characters int
}

// Task — это сама задача со всеми ее данными
type Task struct {
	ID     int
	Input  string
	State  TaskState
	Result *Result
}

// TaskManager — интерфейс
type TaskManager interface {
	CreateTask(input string) (int, error)
	GetTask(id int) (*Task, error)
	ProcessTask(id int) error
}

// InMemoryManager хранит задачи прямо в оперативной памяти
type InMemoryManager struct {
	tasks  map[int]*Task
	nextID int
}

func NewInMemoryManager() *InMemoryManager {
	return &InMemoryManager{
		tasks:  make(map[int]*Task),
		nextID: 1,
	}
}

func (m *InMemoryManager) CreateTask(input string) (int, error) {
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

func (m *InMemoryManager) GetTask(id int) (*Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (m *InMemoryManager) ProcessTask(id int) error {
	task, err := m.GetTask(id)
	if err != nil {
		return err
	}

	if task.State == StateCompleted {
		return errors.New("invalid operation: task is already completed")
	}

	task.State = StateProcessing

	words := len(strings.Fields(task.Input))
	chars := utf8.RuneCountInString(task.Input)

	task.Result = &Result{
		Words:      words,
		Characters: chars,
	}

	task.State = StateCompleted
	return nil
}

// Вот здесь мы "оживили" запуск программы!
func main() {
	fmt.Println("=== Запуск SIS1: Анализатор текста ===")
	manager := NewInMemoryManager()

	// Сценарий 1: Нормальный текст
	fmt.Println("\n--- Сценарий 1: Успешная задача ---")
	text := "Привет KBTU! Изучаем язык Golang."
	fmt.Printf("Ввод: %q\n", text)
	
	id1, err := manager.CreateTask(text)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Задача создана! Ей присвоен ID: %d\n", id1)
		
		// Обрабатываем задачу
		manager.ProcessTask(id1)
		
		// Достаем результат
		task1, _ := manager.GetTask(id1)
		fmt.Printf("СТАТУС: %s\n", task1.State)
		fmt.Printf("РЕЗУЛЬТАТ: Слов - %d, Символов - %d\n", task1.Result.Words, task1.Result.Characters)
	}

	// Сценарий 2: Передаем пустой текст (ошибка валидации)
	fmt.Println("\n--- Сценарий 2: Передача пустого текста ---")
	_, err = manager.CreateTask("    ")
	if err != nil {
		fmt.Printf("Получена ожидаемая ошибка: %v\n", err)
	}

	// Сценарий 3: Ищем задачу, которой не существует
	fmt.Println("\n--- Сценарий 3: Поиск задачи с ID 999 ---")
	_, err = manager.GetTask(999)
	if err != nil {
		fmt.Printf("Получена ожидаемая ошибка: %v\n", err)
	}
	
	fmt.Println("\n======================================")
}