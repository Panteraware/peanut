package peanut

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type CustomHeaders struct {
	Label string
	Value string
}

var client = http.Client{Timeout: 10 * time.Second}

func getRequest(url string, headers []CustomHeaders) *http.Response {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("error creating request")
	}

	for _, value := range headers {
		req.Header.Set(value.Label, value.Value)
	}

	res, err := client.Do(req)

	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("error executing request")
	}

	fmt.Printf(res.Status)

	return res
}
