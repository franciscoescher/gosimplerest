package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/testify/assert"
)

func TestDeleteHandler(t *testing.T) {
	// Prepare the test
	base := &resource.Base{Resource: &testResource, Logger: logrus.New(), DB: testDB}

	t1 := time.Now()
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, time.UTC)
	var data = map[string]interface{}{
		"uuid":       "1f79de62-97b4-48cd-b89d-18628bf50395",
		"first_name": "Fulano",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": nil,
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

	sqlResult := base.DB.QueryRow(fmt.Sprintf(`SELECT deleted_at FROM %s WHERE uuid = ? LIMIT 1`,
		testResource.Table), data["uuid"])
	dataDB := make([]string, 1)
	err = sqlResult.Scan(&dataDB[0])
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, dataDB[0])
}
