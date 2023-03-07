package validator

type Validator interface {
	// Var validates a single variable against a single validation rule
	Var(field interface{}, tag string) error
	// ValidateMap validates a map of variables against a map of validation rules
	ValidateMap(data map[string]interface{}, rules map[string]interface{}) map[string]interface{}
}

type BlankValidator struct{}

func (b *BlankValidator) Var(field interface{}, tag string) error {
	return nil
}

func (b *BlankValidator) ValidateMap(data map[string]interface{}, rules map[string]interface{}) map[string]interface{} {
	return nil
}
