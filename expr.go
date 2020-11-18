package sq

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type expr struct {
	sql  string
	args []interface{}
}

// Expr builds value expressions for InsertBuilder and UpdateBuilder.
//
// Ex:
//     .data(Expr("FROM_UNIXTIME(?)", t))
func Expr(sql string, args ...interface{}) expr {
	return expr{sql: sql, args: args}
}

func (e expr) ToSQL() (string, []interface{}, error) {
	if !hasQueryBuilder(e.args) {
		return e.sql, e.args, nil
	}

	args := make([]interface{}, 0, len(e.args))
	sql, err := replacePlaceholders(e.sql, func(buf *bytes.Buffer, i int) error {
		if i > len(e.args) {
			buf.WriteRune('?')
			return nil
		}
		switch arg := e.args[i-1].(type) {
		case QueryBuilder:
			sql, vs, err := arg.ToSQL()
			if err != nil {
				return err
			}
			args = append(args, vs...)
			fmt.Fprint(buf, sql)
		default:
			args = append(args, arg)
			buf.WriteRune('?')
		}
		return nil
	})
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}

type exprs []expr

func (es exprs) AppendToSQL(w io.Writer, sep string, args []interface{}) ([]interface{}, error) {
	for i, e := range es {
		if i > 0 {
			_, err := io.WriteString(w, sep)
			if err != nil {
				return nil, err
			}
		}
		_, err := io.WriteString(w, e.sql)
		if err != nil {
			return nil, err
		}
		args = append(args, e.args...)
	}
	return args, nil
}

// aliasExpr helps to alias part of SQL query generated with underlying "expr"
type aliasExpr struct {
	expr  QueryBuilder
	alias string
}

// Alias allows to define alias for column in SelectBuilder. Useful when column is
// defined as complex expression like IF or CASE
// Ex:
//		.Column(Alias(caseStmt, "case_column"))
func Alias(expr QueryBuilder, alias string) aliasExpr {
	return aliasExpr{expr, alias}
}

func (e aliasExpr) ToSQL() (sql string, args []interface{}, err error) {
	sql, args, err = e.expr.ToSQL()
	if err == nil {
		sql = fmt.Sprintf("(%s) AS %s", sql, e.alias)
	}
	return
}

// Eq is syntactic sugar for use with Where/Having/Set methods.
// Ex:
//     .Where(Eq{"id": 1})
type Eq map[string]interface{}

func (eq Eq) toSQL(useNotOpr bool) (sql string, args []interface{}, err error) {
	var (
		exprs    []string
		equalOpr = "="
		inOpr    = "IN"
		nullOpr  = "IS"
	)

	if useNotOpr {
		equalOpr = "<>"
		inOpr = "NOT IN"
		nullOpr = "IS NOT"
	}

	for key, val := range eq {
		expr := ""

		switch v := val.(type) {
		case driver.Valuer:
			if val, err = v.Value(); err != nil {
				return
			}
		}

		if val == nil {
			expr = fmt.Sprintf("%s %s NULL", key, nullOpr)
		} else {
			valVal := reflect.ValueOf(val)
			if valVal.Kind() == reflect.Array || valVal.Kind() == reflect.Slice {
				if valVal.Len() == 0 {
					err = fmt.Errorf("equality condition must contain at least one paramater")
					return
				}
				for i := 0; i < valVal.Len(); i++ {
					args = append(args, valVal.Index(i).Interface())
				}
				expr = fmt.Sprintf("%s %s (%s)", key, inOpr, Placeholders(valVal.Len()))
			} else {
				expr = fmt.Sprintf("%s %s ?", key, equalOpr)
				args = append(args, val)
			}
		}
		exprs = append(exprs, expr)
	}
	sql = strings.Join(exprs, " AND ")
	return
}

// ToSQL builds the query into a SQL string and bound args.
func (eq Eq) ToSQL() (sql string, args []interface{}, err error) {
	return eq.toSQL(false)
}

// NotEq is syntactic sugar for use with Where/Having/Set methods.
// Ex:
//     .Where(NotEq{"id": 1}) == "id <> 1"
type NotEq Eq

// ToSQL builds the query into a SQL string and bound args.
func (neq NotEq) ToSQL() (sql string, args []interface{}, err error) {
	return Eq(neq).toSQL(true)
}

// Lt is syntactic sugar for use with Where/Having/Set methods.
// Ex:
//     .Where(Lt{"id": 1})
type Lt map[string]interface{}

func (lt Lt) toSQL(opposite, orEq bool) (sql string, args []interface{}, err error) {
	var (
		exprs []string
		opr   = "<"
	)

	if opposite {
		opr = ">"
	}

	if orEq {
		opr = fmt.Sprintf("%s%s", opr, "=")
	}

	for key, val := range lt {
		expr := ""

		switch v := val.(type) {
		case driver.Valuer:
			if val, err = v.Value(); err != nil {
				return
			}
		}

		if val == nil {
			err = fmt.Errorf("cannot use null with less than or greater than operators")
			return
		} else if v, ok := val.(QueryBuilder); ok {
			var s string
			var a []interface{}
			s, a, err = v.ToSQL()
			if err != nil {
				return
			}

			expr = fmt.Sprintf("%s %s %s", key, opr, s)
			args = append(args, a...)
		} else {
			valVal := reflect.ValueOf(val)
			if valVal.Kind() == reflect.Array || valVal.Kind() == reflect.Slice {
				err = fmt.Errorf("cannot use array or slice with less than or greater than operators")
				return
			} else {
				expr = fmt.Sprintf("%s %s ?", key, opr)
				args = append(args, val)
			}
		}
		exprs = append(exprs, expr)
	}
	sql = strings.Join(exprs, " AND ")
	return
}

func (lt Lt) ToSQL() (sql string, args []interface{}, err error) {
	return lt.toSQL(false, false)
}

// LtOrEq is syntactic sugar for use with Where/Having/Set methods.
// Ex:
//     .Where(LtOrEq{"id": 1}) == "id <= 1"
type LtOrEq Lt

func (ltOrEq LtOrEq) ToSQL() (sql string, args []interface{}, err error) {
	return Lt(ltOrEq).toSQL(false, true)
}

// Gt is syntactic sugar for use with Where/Having/Set methods.
// Ex:
//     .Where(Gt{"id": 1}) == "id > 1"
type Gt Lt

func (gt Gt) ToSQL() (sql string, args []interface{}, err error) {
	return Lt(gt).toSQL(true, false)
}

// GtOrEq is syntactic sugar for use with Where/Having/Set methods.
// Ex:
//     .Where(GtOrEq{"id": 1}) == "id >= 1"
type GtOrEq Lt

func (gtOrEq GtOrEq) ToSQL() (sql string, args []interface{}, err error) {
	return Lt(gtOrEq).toSQL(true, true)
}

type conj []QueryBuilder

func (c conj) join(sep string) (sql string, args []interface{}, err error) {
	var sqlParts []string
	for _, queryBuilder := range c {
		partSQL, partArgs, err := queryBuilder.ToSQL()
		if err != nil {
			return "", nil, err
		}
		if partSQL != "" {
			sqlParts = append(sqlParts, partSQL)
			args = append(args, partArgs...)
		}
	}
	if len(sqlParts) > 0 {
		sql = fmt.Sprintf("(%s)", strings.Join(sqlParts, sep))
	}
	return
}

// And is syntactic sugar that glues where/having parts with AND clause
// Ex:
//     .Where(And{Expr("a > ?", 15), Expr("b < ?", 20), Expr("c is TRUE")})
type And conj

// ToSQL builds the query into a SQL string and bound args.
func (a And) ToSQL() (string, []interface{}, error) {
	return conj(a).join(" AND ")
}

// Or is syntactic sugar that glues where/having parts with OR clause
// Ex:
//     .Where(And{Expr("a > ?", 15), Expr("b < ?", 20), Expr("c is TRUE")})
type Or conj

// ToSQL builds the query into a SQL string and bound args.
func (o Or) ToSQL() (string, []interface{}, error) {
	return conj(o).join(" OR ")
}

func hasQueryBuilder(args []interface{}) bool {
	for _, arg := range args {
		_, ok := arg.(QueryBuilder)
		if ok {
			return true
		}
	}
	return false
}
