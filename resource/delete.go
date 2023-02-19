package resource

import (
	"database/sql"
	"fmt"
	"strings"
)

// Delete deletes a row with the given primary key from the database
func (b *Resource) Delete(base *Base, id string) error {
	var result sql.Result
	err := error(nil)
	if b.SoftDeleteField.Valid {
		sqlStr := ConcatStr(`UPDATE `, b.Table(), ` SET `, b.SoftDeleteField.String, ` = NOW() WHERE `, b.PrimaryKey, `=?`)
		result, err = base.DB.Exec(sqlStr, id)
		if err != nil {
			return err
		}
	} else {
		result, err = base.DB.Exec(ConcatStr(`DELETE FROM `, b.Table(), ` WHERE id=?`), id)
		if err != nil {
			return err
		}
	}
	affect, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affect == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func ConcatStr(strs ...string) string {
	var sb strings.Builder
	for _, s := range strs {
		sb.WriteString(s)
	}
	return sb.String()
}
