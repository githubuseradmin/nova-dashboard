// Package migrations embeds the SQL migration files, applied at startup in
// lexical filename order (0001_*, 0002_*, ...).
package migrations

import "embed"

// Files holds all *.sql migrations.
//
//go:embed *.sql
var Files embed.FS
