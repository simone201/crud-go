package db

import (
	"fmt"
	"strings"
)

type NoRowsAffectedError struct{}

func (e NoRowsAffectedError) Error() string {
	return "sql: no rows affected"
}

type ParamsNotValidError struct {
	params []string
}

func (e ParamsNotValidError) Error() string {
	return fmt.Sprintf("sql: params not valid [%s]", strings.Join(e.params, ","))
}
