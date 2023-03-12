package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franciscoescher/gosimplerest/repository/local"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandlerOK(t *testing.T) {
	// Prepare the test
	var data = map[string]interface{}{
		"first_name": "Fulano",
		"phone":      "+55 (11) 99999-9999",
	}

	params := &GetHandlerFuncParams{
		Resource:   &testResource,
		Logger:     logrus.New(),
		Validate:   validator.New(),
		Repository: local.NewRepository(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
	request, err := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateHandler(params))
	handler.ServeHTTP(response, request)

	// Make assertions
	var bodyJson map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &bodyJson)
	if err != nil {
		t.Fatal(err)
	}
	data["uuid"] = bodyJson["uuid"]

	assert.Equal(t, http.StatusOK, response.Code)

	dataInDB, _ := params.Repository.Find(params.Resource, bodyJson["uuid"])
	dataOnlyInsertedFields := map[string]interface{}{
		"first_name": dataInDB["first_name"],
		"phone":      dataInDB["phone"],
		"uuid":       dataInDB["uuid"],
	}

	assert.Equal(t, data, dataOnlyInsertedFields)
}

func TestCreateHandlerBadRequest(t *testing.T) {
	// Prepare the test
	params := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: testRepo, Validate: validator.New()}

	var data = map[string]interface{}{
		"first_name": "",
		"phone":      "+55 (11) 99999-9999",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
	request, err := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateHandler(params))
	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}
