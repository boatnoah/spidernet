package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	defaultBaseURL = "http://localhost:8080"
	defaultTimeout = 15 * time.Second
)

type crawlRequest struct {
	StartURL string `json:"start_url"`
	Depth    int    `json:"depth"`
}

type crawlResponse struct {
	JobID uuid.UUID `json:"job_id"`
}

type jobStatusResponse struct {
	ID        uuid.UUID `json:"id"`
	StartURL  string    `json:"start_url"`
	MaxDepth  int       `json:"max_depth"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	if len(os.Args) < 2 {
		printTopLevelUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "crawl":
		exitOnErr(runCrawl(os.Args[2:]))
	case "status":
		exitOnErr(runStatus(os.Args[2:]))
	case "graph":
		exitOnErr(runGraph(os.Args[2:]))
	case "health":
		exitOnErr(runHealth(os.Args[2:]))
	case "-h", "--help", "help":
		printTopLevelUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
		printTopLevelUsage()
		os.Exit(1)
	}
}

func runCrawl(args []string) error {
	fs := flag.NewFlagSet("crawl", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	baseURL := fs.String("api-url", defaultBaseURL, "spidernet API base URL")
	depth := fs.Int("depth", 1, "max crawl depth")
	timeout := fs.Duration("timeout", defaultTimeout, "request timeout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("usage: spidernet crawl [--api-url URL] [--depth N] [--timeout 15s] <start-url>")
	}

	startURL := fs.Arg(0)
	if _, err := url.ParseRequestURI(startURL); err != nil {
		return fmt.Errorf("invalid start URL: %w", err)
	}

	payload := crawlRequest{
		StartURL: startURL,
		Depth:    *depth,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	endpoint := buildEndpoint(*baseURL, "/v1/crawl")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return httpStatusError(resp)
	}

	var response crawlResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	fmt.Printf("job submitted: %s\n", response.JobID)
	return nil
}

func runStatus(args []string) error {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	baseURL := fs.String("api-url", defaultBaseURL, "spidernet API base URL")
	timeout := fs.Duration("timeout", defaultTimeout, "request timeout")
	asJSON := fs.Bool("json", false, "print full JSON response")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: spidernet status [--api-url URL] [--timeout 15s] [--json] <job-id>")
	}

	jobID, err := uuid.Parse(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid job id: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	endpoint := buildEndpoint(*baseURL, "/v1/jobs/"+jobID.String()+"/status")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return httpStatusError(resp)
	}

	if *asJSON {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	}

	var response jobStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	fmt.Printf("job: %s\n", response.ID)
	fmt.Printf("status: %s\n", response.Status)
	fmt.Printf("start-url: %s\n", response.StartURL)
	fmt.Printf("max-depth: %d\n", response.MaxDepth)
	fmt.Printf("created-at: %s\n", response.CreatedAt.Format(time.RFC3339))
	return nil
}

func runGraph(args []string) error {
	fs := flag.NewFlagSet("graph", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	baseURL := fs.String("api-url", defaultBaseURL, "spidernet API base URL")
	timeout := fs.Duration("timeout", defaultTimeout, "request timeout")
	outPath := fs.String("out", "", "output .png path (default: <job-id>.png)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: spidernet graph [--api-url URL] [--timeout 15s] [--out FILE] <job-id>")
	}

	jobID, err := uuid.Parse(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid job id: %w", err)
	}

	targetPath := *outPath
	if strings.TrimSpace(targetPath) == "" {
		targetPath = jobID.String() + ".png"
	}
	if filepath.Ext(targetPath) == "" {
		targetPath += ".png"
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	endpoint := buildEndpoint(*baseURL, "/v1/jobs/"+jobID.String()+"/graph")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return httpStatusError(resp)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := os.WriteFile(targetPath, b, 0o644); err != nil {
		return err
	}

	fmt.Printf("graph saved to %s\n", targetPath)
	return nil
}

func runHealth(args []string) error {
	fs := flag.NewFlagSet("health", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	baseURL := fs.String("api-url", defaultBaseURL, "spidernet API base URL")
	timeout := fs.Duration("timeout", defaultTimeout, "request timeout")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("usage: spidernet health [--api-url URL] [--timeout 15s]")
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	endpoint := buildEndpoint(*baseURL, "/v1/health")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return httpStatusError(resp)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func httpStatusError(resp *http.Response) error {
	b, _ := io.ReadAll(resp.Body)
	if len(b) == 0 {
		return fmt.Errorf("request failed with status %s", resp.Status)
	}
	return fmt.Errorf("request failed with status %s: %s", resp.Status, strings.TrimSpace(string(b)))
}

func buildEndpoint(baseURL string, path string) string {
	return strings.TrimRight(baseURL, "/") + path
}

func exitOnErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func printTopLevelUsage() {
	fmt.Fprintln(os.Stderr, "spidernet CLI")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  spidernet <command> [options]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  crawl   Submit a crawl job")
	fmt.Fprintln(os.Stderr, "  status  Show crawl job status")
	fmt.Fprintln(os.Stderr, "  graph   Download crawl graph as PNG")
	fmt.Fprintln(os.Stderr, "  health  Check API health endpoint")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Use `spidernet <command> --help` for command-specific options.")
}
