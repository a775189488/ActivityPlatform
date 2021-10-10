package metadata

import "entrytask/internal/model"

type UserRegisterReq struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Aliasname string `json:"aliasname" binding:"required"`
	Headpic   string `json:"headpic" binding:"required"`
}

func (u UserRegisterReq) ToUser() *model.User {
	return &model.User{
		Username:  u.Username,
		Password:  u.Password,
		Email:     u.Email,
		Aliasname: u.Aliasname,
		Headpic:   u.Headpic,
	}
}

type UserRegisterResp struct {
	Id        uint64 `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Aliasname string `json:"aliasname"`
	Headpic   string `json:"headpic"`
}

func FormatUserRegisterResp(user *model.User) *UserRegisterResp {
	return &UserRegisterResp{
		user.Id,
		user.Username,
		user.Email,
		user.Aliasname,
		user.Headpic,
	}
}

type UploadHeadpicResp struct {
	FilePath string `json:"file_path"`
}

type ListUserReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type GetUserResp struct {
	Id        uint64 `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Aliasname string `json:"aliasname"`
	Headpic   string `json:"headpic"`
}

type UpdateUserReq struct {
	Email     string `json:"email"`
	Aliasname string `json:"aliasname"`
	Headpic   string `json:"headpic"`
}

type ListUserActReq struct {
	Page int32 `form:"page" binding:"required"`
	Size int32 `form:"size" binding:"required"`
}

func (u *UpdateUserReq) ToUser() *model.User {
	return &model.User{
		Email:     u.Email,
		Aliasname: u.Aliasname,
		Headpic:   u.Headpic,
	}
}

func FormatGetUserResp(user *model.User) *GetUserResp {
	return &GetUserResp{
		user.Id,
		user.Username,
		user.Email,
		user.Aliasname,
		user.Headpic,
	}
}

type ListUserResp struct {
	ListBase
	Data []*GetUserResp `json:"data"`
}

func FormatListUserResp(users []*model.User, total int, page int, size int) *ListUserResp {
	data := make([]*GetUserResp, 0)
	for _, u := range users {
		data = append(data, &GetUserResp{
			u.Id,
			u.Username,
			u.Email,
			u.Aliasname,
			u.Headpic,
		})
	}
	return &ListUserResp{
		ListBase{
			page,
			size,
		},
		data,
	}
}

type GetUserActResp struct {
	ListBase
	Data []*GetActResp `json:"data"`
}

func FormateGetUserActResp(acts []*model.Activity, page int, size int) *GetUserActResp {
	data := make([]*GetActResp, 0)
	for _, a := range acts {
		data = append(data, &GetActResp{
			Id:          a.Id,
			Title:       a.Title,
			BeginAt:     a.BeginAt,
			EndAt:       a.EndAt,
			Description: a.Description,
			Type:        a.ActType,
			Address:     a.Address,
		})
	}
	return &GetUserActResp{
		ListBase{
			page,
			size,
		},
		data,
	}
}
