package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	StatusCode int
	Error      error
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 0, "Número total de requests")
	concurrency := flag.Int("concurrency", 0, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" || *requests <= 0 || *concurrency <= 0 {
		fmt.Println("Uso: --url=<URL> --requests=<NUM> --concurrency=<NUM>")
		return
	}

	fmt.Printf("Iniciando teste de carga...\n")
	fmt.Printf("URL: %s\n", *url)
	fmt.Printf("Requests: %d\n", *requests)
	fmt.Printf("Concorrência: %d\n\n", *concurrency)

	startTime := time.Now()
	results := runLoadTest(*url, *requests, *concurrency)
	totalTime := time.Since(startTime)

	generateReport(results, *requests, totalTime)
}

func runLoadTest(url string, totalRequests, concurrency int) []Result {
	results := make([]Result, totalRequests)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	client := &http.Client{Timeout: 30 * time.Second}

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			resp, err := client.Get(url)
			if err != nil {
				results[index] = Result{Error: err}
				return
			}
			defer resp.Body.Close()

			results[index] = Result{StatusCode: resp.StatusCode}
		}(i)
	}

	wg.Wait()
	return results
}

func generateReport(results []Result, totalRequests int, totalTime time.Duration) {
	statusCounts := make(map[int]int)
	successCount := 0
	errorCount := 0

	for _, result := range results {
		if result.Error != nil {
			errorCount++
			continue
		}

		statusCounts[result.StatusCode]++
		if result.StatusCode == 200 {
			successCount++
		}
	}

	fmt.Println("### RELATÓRIO DE TESTE DE CARGA ###")
	fmt.Printf("Tempo total de execução: %.2fs\n", totalTime.Seconds())
	fmt.Printf("Total de requests realizados: %d\n", totalRequests)
	fmt.Printf("Requests com status 200: %d\n", successCount)

	fmt.Println("\nDistribuição de códigos de status:")
	for status, count := range statusCounts {
		fmt.Printf("Status %d: %d requests\n", status, count)
	}

	if errorCount > 0 {
		fmt.Printf("Erros de conexão: %d\n", errorCount)
	}

	fmt.Println("### FIM DO RELATÓRIO ###")
}
