package validator

import (
	"net/mail"

	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/util"
)

func ValidateUser(user requestModel.User) (bool, *errs.XError) {
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return false, errs.NewXError(errs.VALIDATION, "Invalid Email", err)
	}

	if util.IsNilOrEmptyString(&user.FirstName) {
		return false, errs.NewXError(errs.VALIDATION, "Invalid Name", nil)
	}

	if util.IsNilOrEmptyString(&user.Extension) || util.IsNilOrEmptyString(&user.PhoneNumber) {
		return false, errs.NewXError(errs.VALIDATION, "Invalid PhoneNumber", nil)
	}
	return true, nil
}
