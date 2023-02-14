package resource

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// Insert inserts a new row into the database
func (b *Resource) Insert(base *Base, data map[string]any) (int64, error) {
	l := len(data)
	if b.AutoIncrementalPK {
		l--
	}
	in := strings.TrimSuffix(strings.Repeat("?,", l), ",")
	fields := make([]string, l)
	values := make([]any, l)
	i := 0
	for key, element := range data {
		if key == b.PrimaryKey() && b.AutoIncrementalPK {
			continue
		}
		fields[i] = key
		values[i] = element
		i++
	}
	logrus.Info("AQUIIII")
	logrus.Info(values)

	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, b.GetTable(), strings.Join(fields, ","), in)
	result, err := base.DB.Exec(sql, values...)
	if err != nil {
		return 0, err
	}
	if b.AutoIncrementalPK {
		return result.LastInsertId()
	}
	return 0, nil
}
