package sq

import "github.com/Masterminds/squirrel"

// PgSb postgres arguments squirrel builder
func PgSb() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
