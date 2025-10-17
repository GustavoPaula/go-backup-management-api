package repository

import (
	"errors"
	"log/slog"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ErrPgConnectionDoesNotExist    = "08003"
	ErrPgConnectionFailure         = "08006"
	ErrPgInvalidTextRepresentation = "22P02"
	ErrPgNotNullViolation          = "23502"
	ErrPgForeignKeyViolation       = "23503"
	ErrPgUniqueConstraint          = "23505"
	ErrPgCheckViolation            = "23514"
	ErrPgInsufficientPrivilege     = "42501"
	ErrPgSerializationFailure      = "40001"
	ErrPgDeadLockDetected          = "40P01"
	ErrPgLockNotAvailable          = "55P03"
)

func handleDatabaseError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case ErrPgInvalidTextRepresentation:
			slog.Error("Formato de texto inválido", "datatype", pgErr.DataTypeName, "msg", pgErr.Message)
			return domain.ErrBadRequest
		case ErrPgNotNullViolation:
			slog.Error("Campo obrigatório não preenchido", "column", pgErr.ColumnName)
			return domain.ErrBadRequest
		case ErrPgCheckViolation:
			slog.Error("Violação de regra de negócio (CHECK)", "constraint", pgErr.ConstraintName)
			return domain.ErrBadRequest
		case ErrPgForeignKeyViolation:
			slog.Error("Violação de chave estrangeira", "constraint", pgErr.ConstraintName)
			return domain.ErrBadRequest
		case ErrPgUniqueConstraint:
			slog.Error("Violação de unicidade (chave duplicada)", "constraint", pgErr.ConstraintName)
			return domain.ErrConflictingData
		case ErrPgSerializationFailure:
			slog.Error("Falha de serialização (conflito concorrente)")
			return domain.ErrConflictingData
		case ErrPgDeadLockDetected:
			slog.Error("Deadlock detectado", "msg", pgErr.Message)
			return domain.ErrConflictingData
		case ErrPgLockNotAvailable:
			slog.Error("Recurso bloqueado", "msg", pgErr.Message)
			return domain.ErrConflictingData
		case ErrPgInsufficientPrivilege:
			slog.Error("Privilégios insuficientes", "msg", pgErr.Message)
			return domain.ErrPermissionDenied
		case ErrPgConnectionFailure, ErrPgConnectionDoesNotExist:
			slog.Error("Banco de dados indisponível", "msg", pgErr.Message)
			return domain.ErrDatabaseUnavailable
		default:
			slog.Error("Erro PostgreSQL não tratado", "code", pgErr.Code, "message", pgErr.Message, "detail", pgErr.Detail, "where", pgErr.Where)
			return domain.ErrInternal
		}
	}

	return err
}
