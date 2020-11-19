package sq

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplacePlaceholders(t *testing.T) {
	sql := "x = ? AND y = ?"
	s, _ := replacePlaceholders(sql)
	assert.Equal(t, "x = $1 AND y = $2", s)
}

func TestPlaceholders(t *testing.T) {
	assert.Equal(t, placeholders(2), "?,?")
}

func TestEscape(t *testing.T) {
	sql := "SELECT uuid, \"data\" #> '{tags}' AS tags FROM nodes WHERE  \"data\" -> 'tags' ??| array['?'] AND enabled = ?"
	s, _ := replacePlaceholders(sql)
	assert.Equal(t, "SELECT uuid, \"data\" #> '{tags}' AS tags FROM nodes WHERE  \"data\" -> 'tags' ?| array['$1'] AND enabled = $2", s)
}

func BenchmarkPlaceholdersArray(b *testing.B) {
	var count = b.N
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}
	var _ = strings.Join(placeholders, ",")
}

func BenchmarkPlaceholdersStrings(b *testing.B) {
	placeholders(b.N)
}
