package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMainFunction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	envFile, err := os.CreateTemp("", ".env")
	assert.NoError(t, err)
	defer os.Remove(envFile.Name())

	err = os.WriteFile(envFile.Name(), []byte(""), 0644)
	assert.NoError(t, err)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	go func() {
		main()
	}()

	time.Sleep(2 * time.Second)

	resp, err := http.Get("http://localhost:8080/metrics")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	w.Close()
	os.Stdout = old

	var buf []byte
	r.Read(buf)
}
