package sq

import (
	"fmt"
	"io"
)

type part struct {
	pred interface{}
	args []interface{}
}

func newPart(pred interface{}, args ...interface{}) QueryBuilder {
	return &part{pred, args}
}

func (p part) ToSQL() (sql string, args []interface{}, err error) {
	switch pred := p.pred.(type) {
	case nil:
		// no-op
	case QueryBuilder:
		sql, args, err = pred.ToSQL()
	case string:
		sql = pred
		args = p.args
	default:
		err = fmt.Errorf("expected string or QueryBuilder, not %T", pred)
	}
	return
}

func appendToSQL(parts []QueryBuilder, w io.Writer, sep string, args []interface{}) ([]interface{}, error) {
	for i, p := range parts {
		partSQL, partArgs, err := p.ToSQL()
		if err != nil {
			return nil, err
		} else if len(partSQL) == 0 {
			continue
		}

		if i > 0 {
			_, err := io.WriteString(w, sep)
			if err != nil {
				return nil, err
			}
		}

		_, err = io.WriteString(w, partSQL)
		if err != nil {
			return nil, err
		}
		args = append(args, partArgs...)
	}
	return args, nil
}
