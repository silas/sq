package sq

import (
	"bytes"
	"strconv"
	"strings"
)

func replacePlaceholders(sql string) (string, error) {
	return replacePlaceholdersIter(sql, func(buf *bytes.Buffer, i int) error {
		buf.WriteString("$")
		buf.WriteString(strconv.Itoa(i))
		return nil
	})
}

func placeholders(count int) string {
	if count < 1 {
		return ""
	}

	return strings.Repeat(",?", count)[1:]
}

func replacePlaceholdersIter(sql string, replace func(buf *bytes.Buffer, i int) error) (string, error) {
	buf := &bytes.Buffer{}
	i := 0
	for {
		p := strings.Index(sql, "?")
		if p == -1 {
			break
		}

		if len(sql[p:]) > 1 && sql[p:p+2] == "??" { // escape ?? => ?
			buf.WriteString(sql[:p])
			buf.WriteString("?")
			if len(sql[p:]) == 1 {
				break
			}
			sql = sql[p+2:]
		} else {
			i++
			buf.WriteString(sql[:p])
			if err := replace(buf, i); err != nil {
				return "", err
			}
			sql = sql[p+1:]
		}
	}

	buf.WriteString(sql)
	return buf.String(), nil
}
