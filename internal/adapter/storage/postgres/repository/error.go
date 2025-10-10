package repository

import (
	"errors"
	"log/slog"

	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// 400
	ErrPgInvalidTextRepresentation = "22P02" // 400 - Formato de texto inválido
	ErrPgInvalidDatetimeFormat     = "22007" // 400 - Formato de data/hora inválido
	ErrPgDatetimeFieldOverflow     = "22008" // 400 - Data/hora fora do range
	ErrPgStringDataRightTruncation = "22001" // 400 - Dados muito longos
	ErrPgNumericValueOutOfRange    = "22003" // 400 - Valor numérico inválido
	ErrPgUniqueConstraint          = "23505" // 400* ou 409 - Violação de unicidade
	ErrPgNotNullViolation          = "23502" // 400 - Campo obrigatório não preenchido
	ErrPgForeignKeyViolation       = "23503" // 400 - Referência inválida
	ErrPgCheckViolation            = "23514" // 400 - Violação de regra de negócio

	// 409
	ErrPgRestrictViolation    = "23001" // 409 - Operação restrita por regras
	ErrPgSerializationFailure = "40001" // 409 - Conflito de serialização
	ErrPgDeadLockDetected     = "40P01" // 409 - Deadlock detectado
	ErrPgLockNotAvailable     = "55P03" // 409 - Recurso bloqueado
	ErrPgObjectInUse          = "55006" // 409 - Objeto em uso

	// 404
	ErrPgUndefinedTable  = "42P01" // 404* ou 500 - Tabela não existe
	ErrPgUndefinidColumn = "42703" // 404* ou 500 - Coluna não existe

	// 500
	ErrPgInFailedSqlTransaction = "25P02" // 500 - Transação em estado inconsistente
	ErrPgSyntaxError            = "42601" // 500 - Erro de sintaxe SQL (bug na aplicação)

	// 503
	ErrPgTooManyConnections = "53300" // 503 - Muitas conexões
	ErrPgQueryCanceled      = "57014" // 503 - Query cancelada (timeout/overload)
)

func handleDatabaseError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		// 400 - Erros de entrada/dados inválidos
		case ErrPgInvalidTextRepresentation:
			slog.Error("Formato de texto inválido", "datatype", pgErr.DataTypeName, "msg", pgErr.Message)
			return domain.ErrBadRequest
		case ErrPgInvalidDatetimeFormat:
			slog.Error("Formato de data/hora inválido", "msg", pgErr.Message)
			return domain.ErrBadRequest
		case ErrPgDatetimeFieldOverflow:
			slog.Error("Data/hora fora do intervalo permitido", "msg", pgErr.Message)
			return domain.ErrBadRequest
		case ErrPgStringDataRightTruncation:
			slog.Error("Dados muito longos", "column", pgErr.ColumnName)
			return domain.ErrBadRequest
		case ErrPgNumericValueOutOfRange:
			slog.Error("Valor numérico fora do intervalo", "column", pgErr.ColumnName)
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

		// 404 - Entidades ou colunas não encontradas
		case ErrPgUndefinedTable:
			slog.Error("Tabela não encontrada", "table", pgErr.TableName)
			return domain.ErrDataNotFound
		case ErrPgUndefinidColumn:
			slog.Error("Coluna não encontrada", "column", pgErr.ColumnName)
			return domain.ErrDataNotFound

		// 409 - Conflitos ou bloqueios
		case ErrPgRestrictViolation:
			slog.Error("Violação de restrição (RESTRICT)", "constraint", pgErr.ConstraintName)
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
		case ErrPgObjectInUse:
			slog.Error("Objeto em uso", "msg", pgErr.Message)
			return domain.ErrConflictingData

		// 500 - Erros internos da aplicação
		case ErrPgInFailedSqlTransaction:
			slog.Error("Transação em estado inconsistente", "msg", pgErr.Message)
			return domain.ErrInternal
		case ErrPgSyntaxError:
			slog.Error("Erro de sintaxe SQL (provável bug da aplicação)", "msg", pgErr.Message)
			return domain.ErrInternal

		// 503 - Sobrecarga ou falhas temporárias
		case ErrPgTooManyConnections:
			slog.Error("Muitas conexões com o banco de dados", "msg", pgErr.Message)
			return domain.ErrServiceUnavailable
		case ErrPgQueryCanceled:
			slog.Error("Query cancelada (timeout/overload)", "msg", pgErr.Message)
			return domain.ErrServiceUnavailable

		// Caso não mapeado
		default:
			slog.Error("Erro PostgreSQL não tratado",
				"code", pgErr.Code,
				"message", pgErr.Message,
				"detail", pgErr.Detail,
				"where", pgErr.Where)
			return domain.ErrInternal
		}
	}

	// Se não for erro do PostgreSQL
	return err
}
