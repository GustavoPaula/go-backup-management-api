package repository

import (
	"errors"
	"log/slog"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ErrPgNotNullViolation    = "23502"
	ErrPgForeignKeyViolation = "23503"
	ErrPgUniqueConstraint    = "23505"
)

func handlePgDatabaseError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case ErrPgNotNullViolation:
			slog.Error("Campo obrigatório não preenchido", "column", pgErr.ColumnName)
			return domain.ErrBadRequest
		case ErrPgForeignKeyViolation:
			slog.Error("Violação de chave estrangeira", "constraint", pgErr.ConstraintName)
			return domain.ErrBadRequest
		case ErrPgUniqueConstraint:
			slog.Error("Violação de unicidade (chave duplicada)", "constraint", pgErr.ConstraintName)
			return domain.ErrConflictingData
		default:
			slog.Error("Erro PostgreSQL não tratado", "code", pgErr.Code, "message", pgErr.Message, "detail", pgErr.Detail, "where", pgErr.Where)
			return domain.ErrInternal
		}
	}

	return err
}
