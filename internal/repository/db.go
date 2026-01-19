package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/utils"
	"github.com/imkarthi24/sf-backend/pkg/constants"
	"github.com/imkarthi24/sf-backend/pkg/db"
	"github.com/imkarthi24/sf-backend/pkg/errs"
	"gorm.io/gorm"
)

// GormDAL provides data access layer functionalities using GORM and DBTransactionManager.
// It encapsulates and acts as a wrapper for common database operations like Create, Update, Delete, and executing stored procedures.
type GormDAL struct {
	tm db.DBTransactionManager
}

func ProvideGormDAL(txn db.DBTransactionManager) GormDAL {
	return GormDAL{tm: txn}
}

// WithDB returns the current transaction from context if exists,
// else creates a new one and stores it in the context.
// It also prepares the transaction with user and channel info (if present in session).
func (customDB *GormDAL) WithDB(ctx *context.Context, opts ...db.TransactionOption) *gorm.DB {
	preparedOpts := []db.TransactionOption{withSessionInfo(ctx)}
	preparedOpts = append(preparedOpts, opts...)
	return customDB.tm.WithTransaction(ctx, preparedOpts...)
}

// ExecuteStoredProc executes a stored procedure using the underlying DBTransactionManager.
func (customDB *GormDAL) ExecuteStoredProc(ctx *context.Context, name string, params map[string]interface{}) ([]db.ResultSet, error) {
	return customDB.tm.ExecuteStoredProc(ctx, name, params)
}

func (customDB *GormDAL) Delete(ctx *context.Context, model interface{}) *errs.XError {

	session := utils.GetSession(ctx)
	if session == nil {
		return errs.NewXError(errs.SMTP_ERROR, "Unable to get user session", nil)
	}

	res := customDB.WithDB(ctx).Model(model).
		Updates(map[string]interface{}{
			"is_active":     false,
			"updated_by_id": session.UserId,
		})

	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to delete", res.Error)
	}
	return nil
}

func (customDB *GormDAL) Update(ctx *context.Context, model interface{}) *errs.XError {

	// updateMap := util.PrepareEntityForUpdate(model)

	// val := util.ReadValueFromContext(ctx, constants.SESSION)
	// if val != nil {
	// 	updateMap["updated_by_id"] = val.(*models.Session).UserId
	// }

	res := customDB.WithDB(ctx).Model(model).Save(model)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to update entity", res.Error)
	}

	return nil
}

func (customDB *GormDAL) Create(ctx *context.Context, value interface{}) *errs.XError {

	res := customDB.WithDB(ctx).Create(value)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to create entity", res.Error)
	}

	return nil
}

func (customDB *GormDAL) BatchCreate(ctx *context.Context, value []interface{}) *errs.XError {

	res := customDB.WithDB(ctx).CreateInBatches(value, 10)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to batch create entity", res.Error)
	}

	return nil
}

// withSessionInfo prepares the transaction with user and channel info from the session in context.
func withSessionInfo(ctx *context.Context) db.TransactionOption {

	return func(db *gorm.DB) {
		session := utils.GetSession(ctx)
		if session == nil {
			return
		}

		db.Set(constants.USER_ID, session.UserId).
			Set(constants.CHANNEL_ID, session.ChannelId)
	}

}
