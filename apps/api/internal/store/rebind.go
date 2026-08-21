package store

import (
	"strconv"
	"strings"
)

// rebindPostgres replaces each standalone '?' placeholder with $1..$n in
// positional order so module/migration SQL written with the unified '?'
// convention runs on PostgreSQL (R1 v1.4 §3). The rebind is positional and
// does not parse SQL: string literals containing '?' as data are a contract
// violation caught by R3 migration review.
func rebindPostgres(query string) string {
	if !strings.ContainsRune(query, '?') {
		return query
	}
	var b strings.Builder
	b.Grow(len(query) + 16)
	n := 0
	for _, r := range query {
		if r == '?' {
			n++
			b.WriteByte('$')
			b.WriteString(strconv.Itoa(n))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
