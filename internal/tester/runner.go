package tester

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func runAttack(targetURL string, rps int, duration time.Duration, maxWorkers uint64, label string, csvFilename string) (*vegeta.Metrics, error) {
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    targetURL,
	})

	rate := vegeta.Rate{Freq: rps, Per: time.Second}
	attacker := vegeta.NewAttacker(vegeta.MaxWorkers(maxWorkers))

	file, err := os.Create(csvFilename)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать файл %s: %w", csvFilename, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"timestamp_sec", "latency_ms", "status_code"})

	var metrics vegeta.Metrics
	startTime := time.Now()

	for res := range attacker.Attack(targeter, rate, duration, label) {
		metrics.Add(res)

		relTime := res.Timestamp.Sub(startTime).Seconds()
		latencyMs := res.Latency.Milliseconds()

		writer.Write([]string{
			fmt.Sprintf("%.2f", relTime),
			strconv.FormatInt(latencyMs, 10),
			strconv.Itoa(int(res.Code)),
		})
	}
	
	metrics.Close()
	return &metrics, nil
}

func printMetrics(m *vegeta.Metrics) {
	fmt.Printf("  - Всего запросов: %d\n", m.Requests)
	fmt.Printf("  - Успешность: %.2f%%\n", m.Success*100)
	fmt.Printf("  - Среднее время отклика: %s\n", m.Latencies.Mean)
	fmt.Printf("  - 95-й процентиль: %s\n", m.Latencies.P95)
	fmt.Printf("  - Макс. время отклика: %s\n", m.Latencies.Max)
	fmt.Print("  - Статус-коды: ")
	for code, count := range m.StatusCodes {
		fmt.Printf("[%s: %d] ", code, count)
	}
	fmt.Println()
}