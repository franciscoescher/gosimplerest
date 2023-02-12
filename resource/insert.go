package resource

import (
	"fmt"
	"strings"
)

// Insert inserts a new row into the database
func (b *Resource) Insert(base *Base, data map[string]any) error {
	in := strings.TrimSuffix(strings.Repeat("?,", len(data)), ",")
	fields := make([]string, len(data))
	values := make([]any, len(data))
	i := 0
	for key, element := range data {
		fields[i] = key
		values[i] = element
		i++
	}

	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, b.Table, strings.Join(fields, ","), in)
	_, err := base.DB.Exec(sql, values...)
	return err
}
