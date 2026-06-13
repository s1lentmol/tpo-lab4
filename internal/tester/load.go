package tester

import (
	"fmt"
	"time"

	"tpo-lab4/config"
)

func RunLoadTesting() {
	fmt.Println(" Нагрузочное тестирование ")

	const targetRPS = 7
	const testDuration = 1 * time.Minute
	const maxWorkers = 11

	for _, cfg := range config.HardwareConfigs {
		targetURL := fmt.Sprintf("%s?token=%s&user=%s&config=%d", config.BaseURL, config.Token, config.User, cfg.ID)
		fmt.Printf("\nТестирование Конфигурации #%d (Цена: $%d) на %d RPS...\n", cfg.ID, cfg.Price, targetRPS)

		rawFilename := fmt.Sprintf("scripts/load_cfg%d_raw.csv", cfg.ID)
		metrics, err := runAttack(targetURL, targetRPS, testDuration, maxWorkers, fmt.Sprintf("load-cfg-%d", cfg.ID), rawFilename)
		if err != nil {
			fmt.Printf("[Ошибка] Не удалось протестировать конфигурацию %d: %v\n", cfg.ID, err)
			continue
		}

		printMetrics(metrics)
		fmt.Printf("  [+] Сырые логи сохранены в '%s'\n", rawFilename)
	}
}