package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var FileCache *FileCacheUtility

type FileCacheUtility struct{}

// DoJsonRequest Method for execute http Request to another service
func (helper *FileCacheUtility) DoJsonRequest(Method string, URL string, RequestData []byte) []byte {
	req, err := http.NewRequest(Method, URL, bytes.NewBuffer(RequestData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)

	return body
}

func (helper *FileCacheUtility) FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (helper *FileCacheUtility) IsFileOlderThanOneDay(t time.Time) bool {
	return time.Now().Sub(t) > time.Hour
}

// CacheData Caching response to file for faster future executions
func (helper *FileCacheUtility) CacheData(filename string, actualData []byte) {
	if helper.FileExists(filename) {
		err := os.Remove(filename)
		if err != nil {
		}
	}

	f, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
		}
	}(f)

	_, _ = f.Write(actualData)
}
