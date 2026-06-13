package tester

import (
	"strings"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"tpo-lab4/config"
)

func RunStressTesting(configID int) {
	fmt.Printf(" Стресс тестирование конфигурации #%d \n", configID)

	targetURL := fmt.Sprintf("%s?token=%s&user=%s&config=%d", config.BaseURL, config.Token, config.User, configID)
	
	steps := []int{5, 10, 20, 35, 50, 75, 100, 150}
	stepDuration := 30 * time.Second
	const maxWorkers = 11

	file, err := os.Create("scripts/stress_results.csv")
	if err != nil {
		fmt.Printf("[Ошибка] Не удалось создать CSV файл: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"rps", "p95_ms"})

	fmt.Println("RPS\tСреднее (мс)\t95-й % (мс)\tУспешность (%)\tКоды ответов")

	for _, rps := range steps {
		metrics, err := runAttack(targetURL, rps, stepDuration, maxWorkers, fmt.Sprintf("stress-%d-rps", rps), "scripts/stress_raw.csv")
		if err != nil {
			fmt.Printf("[Ошибка] Сбой на шаге %d RPS: %v\n", rps, err)
			break
		}

		meanMs := metrics.Latencies.Mean.Milliseconds()
		p95Ms := metrics.Latencies.P95.Milliseconds()
		successRate := metrics.Success * 100

		var statusCodesStr strings.Builder
		for code, count := range metrics.StatusCodes {
			statusCodesStr.WriteString(fmt.Sprintf("%s:%d ", code, count))
		}

		fmt.Printf("%d\t%d ms\t\t%d ms\t\t%.1f%%\t\t%s\n", rps, meanMs, p95Ms, successRate, statusCodesStr.String())

		writer.Write([]string{
			strconv.Itoa(rps),
			strconv.FormatInt(p95Ms, 10),
		})

		if p95Ms > 820 || successRate < 95.0 {
			fmt.Printf("\n[!] Система превысила лимиты ТЗ (820 мс или ошибки) при %d RPS.\n", rps)
			// break
		}
	}
	fmt.Println("\n[+] Результаты стресс-теста успешно сохранены в 'stress_results.csv'")
}