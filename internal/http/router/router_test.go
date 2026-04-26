package router

import (
	"encoding/json"
	employeeapp "github.com/jluisv16/hcm-go/internal/employees/application"
	"github.com/jluisv16/hcm-go/internal/employees/infrastructure/memory"
	employeehttp "github.com/jluisv16/hcm-go/internal/employees/interfaces/http"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jluisv16/hcm-go/internal/config"
	"github.com/jluisv16/hcm-go/internal/http/handlers"
)

func newTestRouter() http.Handler {
	cfg := config.Config{
		AppName: "hcm-go-test",
		AppEnv:  "test",
	}
	healthHandler := handlers.NewHealthHandler(cfg.AppName, "test", time.Now().UTC())
	employeeRepository := memory.NewRepository(memory.SeedEmployees())
	employeeService := employeeapp.NewService(employeeRepository)
	employeeHandler := employeehttp.NewHandler(employeeService)

	return New(cfg, healthHandler, employeeHandler)
}

func TestHealthz(t *testing.T) {
	t.Parallel()
	engine := newTestRouter()

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status field to be ok, got %v", body["status"])
	}
}

func TestPing(t *testing.T) {
	t.Parallel()
	engine := newTestRouter()

	request := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["message"] != "pong" {
		t.Fatalf("expected message field to be pong, got %v", body["message"])
	}
}

func TestListEmployeesHasSeedData(t *testing.T) {
	t.Parallel()
	engine := newTestRouter()

	request := httptest.NewRequest(http.MethodGet, "/api/v1/employees", nil)
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	rawEmployees, ok := body["employees"].([]any)
	if !ok {
		t.Fatalf("employees field is missing or invalid")
	}

	if len(rawEmployees) != 10 {
		t.Fatalf("expected 10 seeded employees, got %d", len(rawEmployees))
	}
}
