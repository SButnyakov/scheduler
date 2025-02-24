package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"scheduler/internal/generator"
	"scheduler/internal/scheduler"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
}

func main() {
	s := scheduler.New()

	s.Run()

	go func() {
		for i := 0; i < 10; i++ {
			t, err := generator.GenerateTask()
			if err != nil {
				log.Printf("failed to create task: %v\n", err)
				continue
			}
			s.AddNewTask(t)
		}
	}()

	<-s.StopChan
	fmt.Println("Симуляция завершена.")

}
