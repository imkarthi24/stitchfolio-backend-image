package service

import (
	"context"

	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type AdminService interface {
	SwitchBranch(ctx *context.Context, switc *requestModel.SwitchBranch) *errs.XError
}

type adminService struct {
	adminRepo repository.AdminRepository
}

func ProvideAdminService(dashboardRepo repository.AdminRepository) AdminService {
	return &adminService{
		adminRepo: dashboardRepo,
	}
}

func (svc *adminService) SwitchBranch(ctx *context.Context, switc *requestModel.SwitchBranch) (err *errs.XError) {
	params := make(map[string]interface{})
	params["_toChannel_id"] = switc.ToChannel
	params["_student_id"] = switc.StudentId
	params["_enquiry_id"] = switc.EnquiryId

	return svc.adminRepo.SwitchBranch(ctx, params)
}
