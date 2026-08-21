// Command dbgdump is a throwaway read-only diagnostic for the users list
// 500 (workspace-014 R3 follow-up). Delete after use.
package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

const baseSelect = "SELECT u.id, u.username, u.name, u.roles, u.password_hash, u.token_version, u.failed_login_count, u.locked_until, u.enabled, u.avatar_url, u.must_change_password, u.created_at, u.updated_at, EXISTS(SELECT 1 FROM user_mfa um WHERE um.user_id = u.id AND um.status = 'active') AS mfa_enabled FROM users u"

func describe(db *sql.DB, label, query string) {
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("[%s] query err: %v\n", label, err)
		return
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		fmt.Printf("[%s] columntypes err: %v\n", label, err)
		return
	}
	fmt.Printf("[%s]\n", label)
	last := len(colTypes) - 1
	for i, ct := range colTypes {
		scan := ct.ScanType()
		scanName := "nil"
		if scan != nil {
			scanName = scan.String()
		}
		if i == last || strings.Contains(strings.ToUpper(ct.DatabaseTypeName()), "BOOL") {
			fmt.Printf("  col %d name=%s decl=%q scanType=%s\n", i, ct.Name(), ct.DatabaseTypeName(), scanName)
		}
	}

	n := 0
	for rows.Next() {
		vals := make([]any, len(colTypes))
		ptrs := make([]any, len(vals))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			fmt.Printf("  row %d scan err: %v\n", n+1, err)
			continue
		}
		n++
		lastVal := vals[len(vals)-1]
		fmt.Printf("  row %d lastcol goType=%T val=%v\n", n, lastVal, lastVal)
	}
}

func main() {
	db, err := sql.Open("sqlite", os.Args[1])
	if err != nil {
		fmt.Println("open:", err)
		os.Exit(1)
	}
	defer db.Close()

	describe(db, "plain (no order/limit)", baseSelect)
	describe(db, "server-shaped (ORDER BY username + LIMIT/OFFSET)", baseSelect+" ORDER BY u.username ASC, u.id ASC LIMIT 20 OFFSET 0")
	describe(db, "server-shaped q variant (ORDER BY created_at DESC)", baseSelect+" ORDER BY u.created_at DESC, u.id ASC LIMIT 20 OFFSET 0")
}
