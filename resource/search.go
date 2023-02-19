package resource

import (
	"strings"
)

// Search searches for rows in the database with where clauses
func (b *Resource) Search(base *Base, query map[string][]string) ([]map[string]any, error) {
	fields := b.GetFieldNames()

	// build query
	where := make([]string, len(query))
	values := make([]any, 0)
	i := 0
	for field, value := range query {
		if len(value) == 1 {
			where[i] = ConcatStr(field, " = ?")
			values = append(values, value[0])
		} else {
			where[i] = ConcatStr(field, " IN (", strings.Repeat("?,", len(value)-1)+"?", ")")
			for _, v := range value {
				values = append(values, v)
			}
		}
		i++
	}
	whereStr := ""
	if len(where) > 0 {
		whereStr = "WHERE " + strings.Join(where, " AND ")
	}
	sqlStr := ConcatStr(`SELECT `, strings.Join(fields, ","), ` FROM `, b.Table, ` `, whereStr, ` ORDER BY `, b.PrimaryKey)
	response, err := base.DB.Query(sqlStr, values...)
	if err != nil {
		return nil, err
	}
	defer response.Close()
	return b.parseRows(response)
}
