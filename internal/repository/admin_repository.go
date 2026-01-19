package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/pkg/errs"
)

type AdminRepository interface {
	SwitchBranch(ctx *context.Context, params map[string]interface{}) *errs.XError
}

type adminRepository struct {
	GormDAL
}

func ProvideAdminRepository(customDB GormDAL) AdminRepository {
	return &adminRepository{GormDAL: customDB}
}

func (ur *adminRepository) SwitchBranch(ctx *context.Context, params map[string]interface{}) *errs.XError {
	_, err := ur.ExecuteStoredProc(ctx, "SwitchChannelForRecord", params)
	if err != nil {
		return errs.NewXError(errs.DATABASE, "Unable to execute func SwitchChannelForRecord", err)
	}
	return nil
}
