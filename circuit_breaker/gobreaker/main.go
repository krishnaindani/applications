package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/sony/gobreaker/v2"
	"io"
	"net/http"
)

var cb *gobreaker.CircuitBreaker[[]byte]

func init() {

	var st gobreaker.Settings
	st.Name = "HTTP GET"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRation := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests > 3 && failureRation >= 0.6
	}

	cb = gobreaker.NewCircuitBreaker[[]byte](st)
}

func Get(url string) ([]byte, error) {
	body, err := cb.Execute(func() ([]byte, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	})

	if err != nil {
		return nil, err
	}
	return body, nil
}

func main() {
	body, err := Get("https://www.google.com/robots.txt")
	if err != nil {
		log.Fatal().
			Err(err)
	}

	fmt.Println(string(body))
}
