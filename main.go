package main

import (
	"bufio" // Добавляем для чтения целых строк с пробелами
	"errors"
	"fmt"
	"os" // Добавляем для работы со стандартным вводом (клавиатурой)
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

// Result хранит результат нашей работы
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

// InMemoryManager хранит задачи в памяти
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

// Новая интерактивная функция main
func main() {
	manager := NewInMemoryManager()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== SIS1: Анализатор текста ===")
	fmt.Println("Введите текст для анализа (или напишите 'exit' для выхода):")
	fmt.Print("\n> ")

	// Считываем ввод пользователя строка за строкой
	for scanner.Scan() {
		input := scanner.Text()

		// Завершаем работу, если пользователь ввёл exit
		if strings.ToLower(strings.TrimSpace(input)) == "exit" {
			fmt.Println("Завершение программы...")
			break
		}

		// 1. Создаём задачу
		id, err := manager.CreateTask(input)
		if err != nil {
			fmt.Printf("Ошибка создания задачи: %v\n", err)
			fmt.Print("\n> ")
			continue
		}

		fmt.Printf("Задача создана! ID: %d. Статус: Pending\n", id)

		// 2. Обрабатываем задачу
		err = manager.ProcessTask(id)
		if err != nil {
			fmt.Printf("Ошибка обработки задачи: %v\n", err)
			fmt.Print("\n> ")
			continue
		}

		// 3. Получаем и выводим результат
		task, _ := manager.GetTask(id)
		fmt.Printf("Задача %d выполнена! Статус: %s\n", task.ID, task.State)
		fmt.Printf("Результат: Слов - %d, Символов - %d\n", task.Result.Words, task.Result.Characters)

		fmt.Print("\n> ")
	}

	// Финальная проверка ошибок считывания (убирает предупреждение линтера)
	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка чтения ввода: %v\n", err)
	}
}