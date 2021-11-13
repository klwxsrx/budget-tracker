package budget

import "embed"

// nolint:gochecknoglobals
//go:embed *.sql
var MysqlMigrations embed.FS
