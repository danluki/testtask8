package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danluki/test-task-8/internal/config"
	handler "github.com/danluki/test-task-8/internal/delivery/http"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	h := handler.NewHandler(nil)

	require.IsType(t, &handler.Handler{}, h)
}

func TestNewHandler_Init(t *testing.T) {
	h := handler.NewHandler(nil)

	router := h.Init(&config.Config{})

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
}
