package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Response struct {
	Query        string `json:"query"`
	OfficialName string `json:"officialName"`
	ISOCode      string `json:"isoCode"`
}

func main() {
	// Flags
	inputFile := flag.String("input", "countries.txt", "Input file with country names")
	outputFile := flag.String("output", "results.log", "Output log file")
	rps := flag.Int("rps", 10, "Requests per second (1-1000)")
	flag.Parse()

	// Files
	in, err := os.Open(*inputFile)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	// Read all countries into memory
	var countries []string
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		country := strings.TrimSpace(scanner.Text())
		if country != "" {
			countries = append(countries, country)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// --- Key Changes Start Here ---

	// 1. Create a channel to pass log messages from workers to the writer.
	// A buffered channel helps prevent workers from blocking if the writer is busy.
	results := make(chan string, len(countries))

	out, err := os.Create(*outputFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// 2. Start a single, dedicated goroutine for writing to the file.
	var writerWg sync.WaitGroup
	writerWg.Add(1)
	go func() {
		defer writerWg.Done()
		writer := bufio.NewWriter(out)
		for result := range results {
			writer.WriteString(result + "\n")
		}
		writer.Flush() // Flush any remaining buffered data at the very end.
	}()

	// --- End Key Changes ---

	client := &http.Client{Timeout: 10 * time.Second}

	// Concurrency control for API requests
	var workerWg sync.WaitGroup
	limiter := time.NewTicker(time.Second / time.Duration(*rps))
	defer limiter.Stop()

	for _, country := range countries {
		<-limiter.C // pace requests
		workerWg.Add(1)
		go func(country string) {
			defer workerWg.Done()
			// 3. Workers now send results to the channel instead of writing to a file.
			handleRequest(client, country, results)
		}(country)
	}

	// Wait for all the API requests to finish.
	workerWg.Wait()

	// 4. Close the channel to signal the writer goroutine that there's no more work.
	close(results)

	// Wait for the writer goroutine to finish writing everything.
	writerWg.Wait()
}

// handleRequest now sends its result string to a channel.
func handleRequest(client *http.Client, country string, results chan<- string) {
	encoded := url.QueryEscape(country)
	apiURL := fmt.Sprintf("https://api.country2iso.0conv.com/api/convert?country=%s", encoded)

	resp, err := client.Get(apiURL)
	if err != nil {
		results <- fmt.Sprintf("%s\tERROR\t%v", country, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		results <- fmt.Sprintf("%s\tHTTP_%d\t%s", country, resp.StatusCode, string(body)[:min(200, len(body))])
		return
	}

	var r Response
	if err := json.Unmarshal(body, &r); err != nil {
		snippet := string(body)
		if len(snippet) > 200 {
			snippet = snippet[:200]
		}
		results <- fmt.Sprintf("%s\tDECODE_ERROR\t%s", country, snippet)
		return
	}

	if r.ISOCode == "" {
		results <- fmt.Sprintf("%s\tMISSING", country)
	} else {
		results <- fmt.Sprintf("%s\tISO=%s\tOfficial=%s", country, r.ISOCode, r.OfficialName)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
