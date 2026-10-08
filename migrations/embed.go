// Package migrations embeds the SQL migrations so the binary can apply them (see the -migrate flag).
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
