package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveHandler(t *testing.T) {
	// Prepare the test
	base := &resource.Base{Resource: &testResource, Logger: logrus.New(), DB: testDB, Validate: validator.New()}

	t1 := time.Now()
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, time.UTC)
	var data = map[string]interface{}{
		"uuid":       "644159a8-0b21-4250-8184-9f06457435c8",
		"first_name": "Fulano",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": t1,
		"created_at": t1.Add(-time.Hour * 24),
	}

	_, err := testDB.Exec(fmt.Sprintf("INSERT INTO %s (uuid, first_name, phone, created_at, deleted_at) VALUES (?,?,?,?,?)", testResource.Table),
		data["uuid"], data["first_name"], data["phone"], data["created_at"], data["deleted_at"],
	)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := fmt.Sprintf("/%s", strcase.KebabCase(testResource.Table))
	request, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	request = GetRequestWithParams(request, map[string]string{"id": data["uuid"].(string)})
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(RetrieveHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusOK, response.Code)

	dataJson, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(dataJson), strings.TrimSpace(response.Body.String()))

	// Test head method
	request, err = http.NewRequest(http.MethodHead, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	request = GetRequestWithParams(request, map[string]string{"id": data["uuid"].(string)})
	response = httptest.NewRecorder()
	handler = http.HandlerFunc(RetrieveHandler(base))
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "", response.Body.String())
}

func TestRetrieveHandlerBadRequest(t *testing.T) {
	// Prepare the test
	base := &resource.Base{Resource: &testResource, Logger: logrus.New(), DB: testDB, Validate: validator.New()}

	// Make the request
	route := fmt.Sprintf("/%s", strcase.KebabCase(testResource.Table))
	request, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	request = GetRequestWithParams(request, map[string]string{})
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(RetrieveHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestRetrieveHandlerNotFound(t *testing.T) {
	// Prepare the test
	base := &resource.Base{Resource: &testResource, Logger: logrus.New(), DB: testDB, Validate: validator.New()}

	// Make the request
	route := fmt.Sprintf("/%s", strcase.KebabCase(testResource.Table))
	request, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	request = GetRequestWithParams(request, map[string]string{"id": "644159a8-0b21-4250-8184-9f06457435c9"})
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(RetrieveHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusNotFound, response.Code)
}
