package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "time"
)

func main() {
	// create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

	// create request with context
    req, err := http.NewRequestWithContext(ctx, "GET", "https://httpbin.org/get", nil)
    if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }

    req.Header.Set("X-Custom-Header", "hello-from-ayo")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error making request:", err)
        return
    }
    defer resp.Body.Close()

    fmt.Println("Status code:", resp.StatusCode)

    body, _ := io.ReadAll(resp.Body)
    fmt.Println("Response:", string(body))
}