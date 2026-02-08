package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/imkarthi24/sf-backend/internal/utils"
	"github.com/loop-kar/pixie/db"
	"github.com/loop-kar/pixie/errs"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	Create(*context.Context, *entities.User) *errs.XError
	Update(*context.Context, *entities.User) *errs.XError
	UpdateChannel(ctx *context.Context, userId uint, channelId uint) *errs.XError
	GetAllUsers(ctx *context.Context, search string) ([]entities.User, *errs.XError)
	GetUserByEmail(*context.Context, string) (*entities.User, *errs.XError)
	Get(*context.Context, uint) (*entities.User, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	SetPasswordReset(*context.Context, uint, string) *errs.XError
	ResetPassword(*context.Context, string, string) (*entities.User, *errs.XError)
	GetUsersForAutoComplete(ctx *context.Context, name string, role []string) ([]entities.User, *errs.XError)

	//User Config
	CreateUserConfig(*context.Context, *entities.UserConfig) *errs.XError
	UpdateUserConfig(*context.Context, *entities.UserConfig) *errs.XError
	GetUserConfig(*context.Context, uint) (*entities.UserConfig, *errs.XError)

	//UserChannelDetail
	CreateUserChannelDetail(ctx *context.Context, channelDetail *entities.UserChannelDetail) *errs.XError
	UpdateUserChannelDetail(ctx *context.Context, det *entities.UserChannelDetail) *errs.XError
	DeleteUserChannelDetail(ctx *context.Context, id uint) *errs.XError

	GetUserAccessibleChannels(ctx *context.Context, userId uint) ([]entities.UserChannelDetail, *errs.XError)
}

type userRepository struct {
	GormDAL
}

func ProvideUserRepository(customDB GormDAL) UserRepository {
	return &userRepository{GormDAL: customDB}
}

func (ur *userRepository) Create(ctx *context.Context, user *entities.User) *errs.XError {
	res := ur.WithDB(ctx).Create(&user)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save user", res.Error)
	}
	return nil

}

func (ur *userRepository) Update(ctx *context.Context, user *entities.User) *errs.XError {
	return ur.GormDAL.Update(ctx, *user)
}

func (ur *userRepository) GetUserByEmail(ctx *context.Context, email string) (*entities.User, *errs.XError) {
	user := entities.User{}
	res := ur.WithDB(ctx).
		Limit(1).
		Where("is_active = ? AND email = ?", true, email).
		Preload("UserChannelDetails", scopes.IsActive(), scopes.SelectFields("user_id", "user_channel_id")).
		Find(&user)
	if res.Error != nil || res.RowsAffected != 1 {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find user", res.Error)
	}
	return &user, nil
}

func (ur *userRepository) GetAllUsers(ctx *context.Context, search string) ([]entities.User, *errs.XError) {

	users := new([]entities.User)

	res := ur.WithDB(ctx).
		Scopes(scopes.Channel(), scopes.IsActive()).
		Scopes(scopes.ILike(search, "first_name", "last_name", "email")).
		Scopes(scopes.AccessibleChannels(utils.GetAccessibleLocationIds(ctx))).
		Scopes(db.Paginate(ctx)).
		Where("role != ?", entities.SYSTEM_ADMIN). //Skip SystemAdmin
		Find(users)

	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to fetch all users", res.Error)
	}

	return *users, nil
}

func (repo *userRepository) Get(ctx *context.Context, id uint) (*entities.User, *errs.XError) {
	user := entities.User{}
	//-- Channel scope removed to allow system admin to access other channel while trying to switch channels
	res := repo.WithDB(ctx).
		//Scopes(scopes.Channel()).
		Find(&user, id)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find user", res.Error)
	}
	return &user, nil
}

func (repo *userRepository) Delete(ctx *context.Context, id uint) *errs.XError {

	user := &entities.User{Model: &entities.Model{ID: id, IsActive: false}}

	err := repo.GormDAL.Delete(ctx, user)

	return err

}

func (repo *userRepository) SetPasswordReset(ctx *context.Context, userId uint, resetString string) *errs.XError {
	user := &entities.User{
		Model:               &entities.Model{ID: userId},
		ResetPasswordString: &resetString,
	}
	res := repo.WithDB(ctx).Updates(user)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to reset password", res.Error)
	}

	return nil
}

func (repo *userRepository) ResetPassword(ctx *context.Context, resetString, password string) (*entities.User, *errs.XError) {
	user := &entities.User{}
	updateMap := map[string]interface{}{
		"password":              password,
		"reset_password_string": nil,
	}

	res := repo.WithDB(ctx).Model(&user).
		Clauses(clause.Returning{}).
		Where("reset_password_string = ?", resetString).
		Updates(updateMap)

	if res.Error != nil || res.RowsAffected == 0 {
		return nil, errs.NewXError(errs.DATABASE, "Unable to reset password", res.Error)
	}

	return user, nil
}

func (repo *userRepository) GetUsersForAutoComplete(ctx *context.Context, name string, role []string) ([]entities.User, *errs.XError) {

	users := new([]entities.User)

	query := repo.WithDB(ctx).
		Scopes(scopes.SearchNameOrEmailOrPhone_Filter(name), scopes.AccessibleChannels(utils.GetAccessibleLocationIds(ctx)), scopes.IsActive()).
		Scopes(db.Paginate(ctx)).
		Where("role != ?", entities.SYSTEM_ADMIN).
		Select("id", "first_name", "last_name")

	if len(role) > 0 {
		query = query.Where("role IN (?)", role)
	}

	res := query.Find(users)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to fetch users for autocomplete", res.Error)
	}

	return *users, nil
}

func (ur *userRepository) UpdateChannel(ctx *context.Context, userId uint, channelId uint) *errs.XError {

	var chanId interface{}
	if channelId == 0 {
		chanId = nil
	} else {
		chanId = channelId
	}

	//Using Exec with Raw query since Model is declared in such a way that
	//we wont be able to update channel using normal ORM operations
	//(gorm:"<-:create)

	res := ur.WithDB(ctx).Exec("UPDATE \"Users\" SET channel_id = ? WHERE id = ?", chanId, userId)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to update User channel", res.Error)
	}

	return nil
}

func (ur *userRepository) CreateUserConfig(ctx *context.Context, config *entities.UserConfig) *errs.XError {
	res := ur.WithDB(ctx).Create(&config)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save user config", res.Error)
	}
	return nil
}

func (ur *userRepository) GetUserConfig(ctx *context.Context, userId uint) (*entities.UserConfig, *errs.XError) {

	userConfig := entities.UserConfig{}
	res := ur.WithDB(ctx).Limit(1).Where("is_active = ? AND user_id = ?", true, userId).Find(&userConfig)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find user config", res.Error)
	}
	return &userConfig, nil
}

func (ur *userRepository) UpdateUserConfig(ctx *context.Context, config *entities.UserConfig) *errs.XError {
	return ur.GormDAL.Update(ctx, *config)
}

func (ur *userRepository) CreateUserChannelDetail(ctx *context.Context, channelDetail *entities.UserChannelDetail) *errs.XError {
	res := ur.WithDB(ctx).Create(&channelDetail)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save user channel details", res.Error)
	}
	return nil
}

func (ur *userRepository) UpdateUserChannelDetail(ctx *context.Context, det *entities.UserChannelDetail) *errs.XError {
	return ur.GormDAL.Update(ctx, *det)
}

func (repo *userRepository) DeleteUserChannelDetail(ctx *context.Context, id uint) *errs.XError {
	user := &entities.UserChannelDetail{Model: &entities.Model{ID: id, IsActive: false}}
	err := repo.GormDAL.Delete(ctx, user)
	return err
}

func (ur *userRepository) GetUserAccessibleChannels(ctx *context.Context, userId uint) ([]entities.UserChannelDetail, *errs.XError) {

	channelDetails := make([]entities.UserChannelDetail, 0)
	res := ur.WithDB(ctx).
		Preload("UserChannel", scopes.IsActive(), scopes.SelectFields("name")).
		Where("is_active = ? AND user_id = ?", true, userId).
		Find(&channelDetails)
	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find user config", res.Error)
	}
	return channelDetails, nil
}
