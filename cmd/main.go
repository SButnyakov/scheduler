package main

import (
	"fmt"
	"log/slog"
	"os"
	"scheduler/internal/generator"
	"scheduler/internal/scheduler"
	"time"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
}

func main() {
	s := scheduler.Scheduler{}

	// Генерация задач каждые 200 мс
	go func() {
		for i := 0; i < 10; i++ {
			tsk := generator.GenerateTask()
			if tsk != nil {
				s.AddNewTask(tsk)
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	// Работаем 3 секунды, затем останавливаем симуляцию
	s.Run()
	time.Sleep(1000 * time.Second)
	fmt.Println("Симуляция завершена.")
}
