package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/franciscoescher/gosimplerest/repository/local"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/testify/assert"
)

func TestDeleteHandler(t *testing.T) {
	// Prepare the test
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	t1 := time.Now()
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, time.UTC)
	var data = map[string]interface{}{
		"uuid":       "1f79de62-97b4-48cd-b89d-18628bf50395",
		"first_name": "Fulano",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": nil,
		"created_at": t1.Add(-time.Hour * 24),
	}
	_, _ = base.Repository.Insert(&testResource, data)

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
	request, err := http.NewRequest(http.MethodDelete, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	request = GetRequestWithParams(request, map[string]string{"id": data["uuid"].(string)})
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusOK, response.Code)

	dataDB, _ := base.Repository.Find(&testResource, data["uuid"])
	assert.Len(t, dataDB, 0)
}

func TestDeleteHandlerNotFound(t *testing.T) {
	// Prepare the test
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
	request, err := http.NewRequest(http.MethodDelete, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	request = GetRequestWithParams(request, map[string]string{"id": "1f79de62-97b4-48cd-b89d-18628bf50396"})
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestDeleteHandlerBadRequest(t *testing.T) {
	// Prepare the test
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
	request, err := http.NewRequest(http.MethodDelete, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	request = GetRequestWithParams(request, map[string]string{})
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusBadRequest, response.Code)
}
