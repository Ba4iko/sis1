package main

import (
	"testing"
)

// Table-driven тест для проверки валидации и создания задачи[cite: 1]
func TestCreateTask(t *testing.T) {
	manager := NewInMemoryManager()

	// Набор тестовых данных (таблица)
	tests := []struct {
		name         string
		input        string
		expectingErr bool
	}{
		{"Valid Input", "Hello Go world", false}, // Нормальный текст - ошибки быть не должно
		{"Empty Input", "", true},                // Пустой текст - ждем ошибку[cite: 1]
		{"Whitespace Input", "   ", true},        // Только пробелы - ждем ошибку
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := manager.CreateTask(tt.input)

			if tt.expectingErr && err == nil {
				t.Errorf("ожидалась ошибка для ввода %q, но ее нет", tt.input)
			}
			if !tt.expectingErr && err != nil {
				t.Errorf("неожиданная ошибка для ввода %q: %v", tt.input, err)
			}
			if !tt.expectingErr && id <= 0 {
				t.Errorf("ожидался корректный ID задачи, получено %d", id)
			}
		})
	}
}

// Тестирование обработки, сохранения результата и изменения состояний[cite: 1]
func TestTaskLifecycleAndProcessing(t *testing.T) {
	manager := NewInMemoryManager()

	// Успешное создание[cite: 1]
	id, err := manager.CreateTask("Hello Go world")
	if err != nil {
		t.Fatalf("ошибка при создании задачи: %v", err)
	}

	// Успешное получение задачи по ID[cite: 1]
	task, _ := manager.GetTask(id)
	if task.State != StatePending {
		t.Errorf("ожидался статус %s, получено %s", StatePending, task.State)
	}

	// Успешная обработка задачи[cite: 1]
	err = manager.ProcessTask(id)
	if err != nil {
		t.Fatalf("ошибка при обработке задачи: %v", err)
	}

	// Корректное состояние и результат после обработки[cite: 1]
	if task.State != StateCompleted {
		t.Errorf("ожидался статус %s, получено %s", StateCompleted, task.State)
	}
	if task.Result.Words != 3 {
		t.Errorf("ожидалось 3 слова, получено %d", task.Result.Words)
	}
	if task.Result.Characters != 14 {
		t.Errorf("ожидалось 14 символов, получено %d", task.Result.Characters)
	}

	// Невалидная операция: попытка обработать уже завершенную задачу[cite: 1]
	err = manager.ProcessTask(id)
	if err == nil {
		t.Error("ожидалась ошибка при попытке обработать уже завершенную задачу")
	}
}

// Попытка получить несуществующую задачу[cite: 1]
func TestGetNonExistentTask(t *testing.T) {
	manager := NewInMemoryManager()
	_, err := manager.GetTask(999)
	if err == nil {
		t.Error("ожидалась ошибка при запросе несуществующей задачи")
	}
}

// Проверка независимости задач друг от друга в памяти[cite: 1]
func TestMultipleIndependentTasks(t *testing.T) {
	manager := NewInMemoryManager()
	id1, _ := manager.CreateTask("First task")
	id2, _ := manager.CreateTask("Second task longer")

	// Обрабатываем только первую задачу
	manager.ProcessTask(id1)

	task1, _ := manager.GetTask(id1)
	task2, _ := manager.GetTask(id2)

	if task1.State != StateCompleted {
		t.Errorf("Задача 1 должна быть завершена")
	}
	if task2.State != StatePending {
		t.Errorf("Задача 2 должна оставаться в статусе Pending")
	}
}
