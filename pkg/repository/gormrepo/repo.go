package gormrepo

import (
	"context"
	"fmt"

	"ptcg_trader/internal/errors"
	"ptcg_trader/pkg/repository"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgconn"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// some default constant
const (
	DefaultPerPage int = 50
)

// define some useful gorm clauses
var (
	selectForUpdate clause.Expression = clause.Locking{Strength: "UPDATE"}
)

// RepoParams define params for create repository
type RepoParams struct {
	fx.In

	DB *gorm.DB
}

type _repository struct {
	db              *gorm.DB
	isInTransaction bool
}

// NewRepository support DI tool to create a new gorm repository instance
func NewRepository(param RepoParams) (repository.Repositorier, error) {
	return &_repository{
		db: param.DB,
	}, nil
}

// DB get gorm.DB with context
func (repo *_repository) DB(ctx context.Context) *gorm.DB {
	return repo.db.WithContext(ctx)
}

// IsInTransaction check that is currently repository handling a transaction
func (repo *_repository) IsInTransaction() bool {
	return repo.isInTransaction
}

// Begin begins a transaction
func (repo *_repository) Begin(ctx context.Context) repository.Repositorier {
	return &_repository{
		db:              repo.db.WithContext(ctx).Begin(),
		isInTransaction: true,
	}
}

// Commit commit a transaction
func (repo *_repository) Commit() error {
	if repo.db == nil {
		return errors.Wrap(errors.ErrInternalError, "repo.db is nil")
	}
	defer func() {
		repo.isInTransaction = true
	}()
	err := repo.db.Commit().Error
	if err != nil {
		return errors.Wrapf(errors.ErrInternalError, "%+v", err)
	}
	return nil
}

// Rollback rollback a transaction
func (repo *_repository) Rollback() error {
	if repo.db == nil {
		return errors.Wrap(errors.ErrInternalError, "repo.db is nil")
	}
	defer func() {
		repo.isInTransaction = true
	}()
	err := repo.db.Rollback().Error
	if err != nil {
		return errors.Wrapf(errors.ErrInternalError, "%+v", err)
	}
	return nil
}

// Transaction handle a transaction in a callback function, from begin to commmit
func (repo *_repository) Transaction(ctx context.Context, f func(context.Context, repository.Repositorier) error) (txErr error) {
	txRepo := repo.Begin(ctx)
	defer func() {
		r := recover()
		if r != nil {
			txErr = errors.Wrap(errors.ErrInternalError, fmt.Sprint(r))
		}
		if txErr != nil {
			_ = txRepo.Rollback()
		} else {
			_ = txRepo.Commit()
		}
	}()

	txErr = f(ctx, txRepo)
	if txErr != nil {
		return txErr
	}
	return nil
}

// Close closes the database and prevents new queries from starting.
func (repo *_repository) Close() error {
	defer func() {
		repo.isInTransaction = true
	}()

	sqlDB, err := repo.db.DB()
	if err != nil {
		return errors.Wrapf(errors.ErrInternalError, "%+v", err)
	}

	err = sqlDB.Close()
	if err != nil {
		return errors.Wrapf(errors.ErrInternalError, "%+v", err)
	}
	return nil
}

func notFoundOrInternalError(err error) error {
	if err == nil {
		return nil
	}
	if err == gorm.ErrRecordNotFound {
		return errors.ErrResourceNotFound
	}
	return errors.ErrInternalError
}

func duplicateOrInternalError(err error) error {
	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case *mysql.MySQLError:
		if e.Number == 1062 {
			return errors.ErrResourceAlreadyExists
		}
	case *pgconn.PgError:
		if e.Code == "23505" {
			return errors.ErrResourceAlreadyExists
		}
	}
	return errors.ErrInternalError
}
