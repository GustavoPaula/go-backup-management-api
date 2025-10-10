package repository

import (
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrPgProtocolViolation         = "08P01"
	ErrPgInvalidTextRepresentation = "22P02"
	ErrPgInvalidDatetimeFormat     = "22007"
	ErrPgDatetimeFieldOverflow     = "22008"
	ErrPgStringDataRightTruncation = "22001"
	ErrPgNumericValueOutOfRange    = "22003"
	ErrPgUniqueConstraint          = "23505"
	ErrPgNotNullViolation          = "23502"
	ErrPgForeignKeyViolation       = "23503"
	ErrPgRestrictViolation         = "23001"
	ErrPgCheckViolation            = "23514"
	ErrPgInFailedSqlTransaction    = "25P02"
	ErrPgSerializationFailure      = "40001"
	ErrPgDeadLockDetected          = "40P01"
	ErrPgUndefinedTable            = "42P01"
	ErrPgUndefinidColumn           = "42703"
	ErrPgSyntaxError               = "42601"
	ErrPgLockNotAvailable          = "55P03"
	ErrPgObjectInUse               = "55006"
	ErrPgTooManyConnections        = "53300"
	ErrPgQueryCanceled             = "57014"
)

func handleCreateError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case ErrPgUniqueConstraint:
			slog.Error("Violação de chave única", "constraint", pgErr.ConstraintName)
			return err
		case ErrPgForeignKeyViolation:
			slog.Error("Violação de chave estrangeira", "constraint", pgErr.ConstraintName)
			return err
		case ErrPgNotNullViolation:
			slog.Error("Violação de NOT NULL", "column", pgErr.ColumnName)
			return err
		case ErrPgCheckViolation:
			slog.Error("Violação de CHECK constraint", "constraint", pgErr.ConstraintName)
			return err
		case ErrPgStringDataRightTruncation:
			slog.Error("Dados truncados", "column", pgErr.ColumnName)
			return err
		case ErrPgInvalidTextRepresentation:
			slog.Error("Tipo de dado inválido", "datatype", pgErr.DataTypeName)
			return err
		case ErrPgRestrictViolation:
			slog.Error("Violação de restrição")
			return err
		default:
			slog.Error("Erro PostgreSQL em INSERT", "code", pgErr.Code, "message", pgErr.Message)
			return err
		}
	}
	return nil
}

func handleUpdateError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case ErrPgUniqueConstraint:
			slog.Error("Violação de chave única no UPDATE", "constraint", pgErr.ConstraintName)
			return err
		case ErrPgForeignKeyViolation:
			slog.Error("Chave estrangeira violada no UPDATE", "constraint", pgErr.ConstraintName)
			return err
		case ErrPgCheckViolation:
			slog.Error("CHECK constraint violada no UPDATE", "constraint", pgErr.ConstraintName)
			return err
		case ErrPgStringDataRightTruncation:
			slog.Error("Dados truncados no UPDATE", "column", pgErr.ColumnName)
			return err
		case ErrPgSerializationFailure:
			slog.Error("Falha de serialização")
			return err
		case ErrPgDeadLockDetected:
			slog.Error("Deadlock detectado")
			return err
		case ErrPgObjectInUse:
			slog.Error("Objeto em uso")
			return err
		default:
			slog.Error("Erro PostgreSQL em UPDATE", "code", pgErr.Code, "message", pgErr.Message)
			return err
		}
	}
	return nil
}

func handleDeleteError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case ErrPgForeignKeyViolation: // foreign_key_violation
			slog.Error("Violação de chave estrangeira no DELETE", "constraint", pgErr.ConstraintName)
			return err

		case ErrPgRestrictViolation: // restrict_violation
			slog.Error("Violação de restrição no DELETE")
			return err

		case ErrPgSerializationFailure: // serialization_failure
			slog.Error("Falha de serialização no DELETE")
			return err

		case ErrPgDeadLockDetected: // deadlock_detected
			slog.Error("Deadlock detectado no DELETE")
			return err

		case ErrPgObjectInUse: // object_in_use
			slog.Error("Objeto em uso durante DELETE")
			return err

		case ErrPgInFailedSqlTransaction: // in_failed_sql_transaction
			slog.Error("Transação falhou durante DELETE")
			return err

		default:
			slog.Error("Erro PostgreSQL em DELETE", "code", pgErr.Code, "message", pgErr.Message)
			return err
		}
	}
	return nil
}
