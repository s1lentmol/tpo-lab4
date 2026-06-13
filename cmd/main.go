package main

import (
	"flag"
	"fmt"
	"os"

	"tpo-lab4/internal/tester"
)

func main() {
	mode := flag.String("mode", "load", "Режим работы: 'load' (нагрузочный тест) или 'stress' (стресс-тест)")
	configID := flag.Int("config", 1, "ID конфигурации для стресс-теста (1, 2 или 3)")
	flag.Parse()

	switch *mode {
	case "load":
		tester.RunLoadTesting()
	case "stress":
		tester.RunStressTesting(*configID)
	default:
		fmt.Println("Неизвестный режим. Используйте -mode=load или -mode=stress")
		os.Exit(1)
	}
}