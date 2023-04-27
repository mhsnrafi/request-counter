package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mhsnrafi/request-counter/counter"
)

const testPersistenceFile = "test_counts.json"

func TestHandleRequest(t *testing.T) {
	config := &counter.Config{
		TimeWindow: 60 * time.Second,
		PersistenceFile: testPersistenceFile,
	}

	cnt := counter.NewCounter(config)

	t.Run("IncrementAndCount", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		handler := handleRequest(cnt)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
		}

		responseBody := strings.TrimSpace(rr.Body.String())
		countStr := strings.TrimPrefix(responseBody, "Requests in the last 60 seconds: ")
		count, err := strconv.Atoi(countStr)

		if err != nil {
			t.Errorf("Error parsing count from response: %v", err)
		}

		if count != 1 {
			t.Errorf("Expected count to be 1, got %d", count)
		}
	})
	os.Remove(config.PersistenceFile)
}
