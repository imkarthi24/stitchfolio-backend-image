package db

import (
	"context"

	"github.com/imkarthi24/sf-backend/pkg/constants"
	"github.com/imkarthi24/sf-backend/pkg/util"
	"gorm.io/gorm"
)

type DBTransactionManager interface {
	WithTransaction(ctx *context.Context, opts ...TransactionOption) *gorm.DB
	Commit(ctx *context.Context)
	Rollback(ctx *context.Context)
	ExecuteStoredProc(ctx *context.Context, name string, params map[string]interface{}) ([]ResultSet, error)
}

type TransactionOption func(*gorm.DB)

type transactionManager struct {
	*StoredProcExecutor
	db *gorm.DB
}

func ProvideDBTransactionManager(db *gorm.DB) DBTransactionManager {
	return &transactionManager{
		db:                 db,
		StoredProcExecutor: &StoredProcExecutor{db: db},
	}
}

// WithTransaction returns the current transaction from context if exists, else creates a new transaction and stores it in context
func (txn *transactionManager) WithTransaction(ctx *context.Context, opts ...TransactionOption) *gorm.DB {

	var gormDB *gorm.DB
	if util.ReadValueFromContext(ctx, constants.TRANSACTION_KEY) == nil {
		gormDB = txn.createTransaction(ctx)
	} else {
		transactionObj := util.ReadValueFromContext(ctx, constants.TRANSACTION_KEY)
		gormDB = transactionObj.(*gorm.DB)
	}

	for _, opt := range opts {
		opt(gormDB)
	}

	return gormDB

}

func (txn *transactionManager) Commit(ctx *context.Context) {
	transaction := txn.WithTransaction(ctx)
	transaction.Commit()
}

func (txn *transactionManager) Rollback(ctx *context.Context) {
	transaction := txn.WithTransaction(ctx)
	transaction.Rollback()
}

func (txn *transactionManager) ExecuteStoredProc(ctx *context.Context, name string, params map[string]interface{}) ([]ResultSet, error) {
	return txn.StoredProcExecutor.CallStoredProcedure(ctx, name, params)
}

func (txn *transactionManager) createTransaction(ctx *context.Context) *gorm.DB {

	transaction := txn.db.Begin()
	newCtx := util.NewContextWithValue(ctx, constants.TRANSACTION_KEY, transaction)
	*ctx = newCtx

	return transaction

}
