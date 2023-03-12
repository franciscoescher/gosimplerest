package handlers

import (
	"os"
	"testing"

	"github.com/franciscoescher/gosimplerest/resource"
	null "gopkg.in/guregu/null.v3"
)

var testResource = resource.Resource{
	Name:       "users_test",
	PrimaryKey: "uuid",
	Fields: map[string]resource.Field{
		"uuid":       {Validator: "uuid4"},
		"first_name": {Validator: "required,min=4,max=20"},
		"phone":      {Unsearchable: true},
		"created_at": {Immutable: true},
		"deleted_at": {},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
