package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/franciscoescher/gosimplerest/repository/local"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/testify/assert"
)

func TestSearchHandler(t *testing.T) {
	// Prepare the test
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	t1 := time.Now()
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, time.UTC)
	data := map[string]interface{}{
		"uuid":       "4c017ccf-0749-4744-a5a6-9c92725411b9",
		"first_name": "Fulano Search Test",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": t1,
		"created_at": t1.Add(-time.Hour * 24),
	}
	data2 := map[string]interface{}{
		"uuid":       "6b548c12-5cac-42e9-aaf1-465c31fafd63",
		"first_name": "Fulano Search Test",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": t1.Add(-time.Hour * 48),
		"created_at": t1.Add(-time.Hour * 72),
	}
	data3 := map[string]interface{}{
		"uuid":       "aed2737d-c105-4ee1-ab41-2341871ecd1a",
		"first_name": "John Search Test",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": t1.Add(-time.Hour * 80),
		"created_at": t1.Add(-time.Hour * 90),
	}
	_, _ = base.Repository.Insert(&testResource, data)
	_, _ = base.Repository.Insert(&testResource, data2)
	_, _ = base.Repository.Insert(&testResource, data3)

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
	request, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	q := request.URL.Query()
	q.Add("first_name", "Fulano Search Test")
	request.URL.RawQuery = q.Encode()
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(SearchHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusOK, response.Code)

	dataJson, err := json.Marshal([]map[string]interface{}{data, data2})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(dataJson), strings.TrimSpace(response.Body.String()))

	// Test head method
	request, err = http.NewRequest(http.MethodHead, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	q = request.URL.Query()
	q.Add("first_name", "Fulano Search Test")
	request.URL.RawQuery = q.Encode()
	response = httptest.NewRecorder()
	handler = http.HandlerFunc(SearchHandler(base))
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "", response.Body.String())
}

func TestSearchHandlerNoContent(t *testing.T) {
	// Prepare the test
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
	request, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	q := request.URL.Query()
	q.Add("first_name", "ABCDEF")
	request.URL.RawQuery = q.Encode()
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(SearchHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestSearchHandlerBadRequest(t *testing.T) {
	// Prepare the test
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
	request, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	q := request.URL.Query()
	q.Add("first_name", "A")
	request.URL.RawQuery = q.Encode()
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(SearchHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestSearchHandlerUnsearchable(t *testing.T) {
	// Prepare the test
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
	request, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	q := request.URL.Query()
	q.Add("phone", "+55 (11) 99999-9999")
	request.URL.RawQuery = q.Encode()
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(SearchHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusBadRequest, response.Code)
}
