package main

import (
	"log"

	"xenon007/xenesis/internal/app"
)

// main запускает CLI-приложение и завершает выполнение при ошибках.
func main() {
	logger := log.Default()
	application := app.New(logger)

	if err := application.Run(); err != nil {
		logger.Fatal(err)
	}
}
