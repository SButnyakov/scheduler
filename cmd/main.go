package main

import (
	"fmt"
	"log/slog"
	"os"
	"scheduler/internal/scheduler"
	"scheduler/internal/task"
	"time"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
}

func main() {
	s := scheduler.New()

	// Генерация задач каждые 200 мс
	//go generator.Generate(s)

	// Работаем 3 секунды, затем останавливаем симуляцию

	s.Run()
	time.Sleep(1 * time.Second)
	t, _ := task.New(task.Basic, task.P2, task.Suspended)
	s.AddNewTask(t)
	time.Sleep(2 * time.Second)
	t, _ = task.New(task.Extended, task.P2, task.Ready)
	s.AddNewTask(t)
	time.Sleep(15 * time.Second)
	t, _ = task.New(task.Basic, task.P1, task.Ready)
	s.AddNewTask(t)
	time.Sleep(1000 * time.Second)
	fmt.Println("Симуляция завершена.")
}
