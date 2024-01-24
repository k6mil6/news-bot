package summary

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

type LocalSummariser struct {
	url string
}

func NewLocalSummariser(url string) *LocalSummariser {
	s := &LocalSummariser{
		url: url,
	}

	log.Printf("local summariser is enabled: %v", url != "")

	return s
}

func (s *LocalSummariser) Summarise(text string) (string, error) {
	url := s.url + "prompt"
	requestBody, err := json.Marshal(map[string]string{"prompt": strconv.Quote(text)})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json") // установка заголовка Content-Type

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
