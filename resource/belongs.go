package resource

import (
	"database/sql"
	"fmt"
	"strings"
)

// FindFromBelongsTo finds all rows of a model with the belongsTo relationship
func (b *Resource) FindFromBelongsTo(base *Base, id any, belongsTo BelongsTo) ([]map[string]any, error) {
	fields := b.GetFieldNames()
	sqlStatement := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ?`, strings.Join(fields, ","), b.GetTable(), belongsTo.Field)
	response, err := base.DB.Query(sqlStatement, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]map[string]any, 0), nil
		}
		return nil, err
	}
	defer response.Close()
	return b.parseRows(response)
}
