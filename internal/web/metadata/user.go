package metadata

import "entrytask/internal/model"

type UserRegisterReq struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Aliasname string `json:"aliasname"`
	Headpic   string `json:"headpic"`
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
	Count int            `json:"count"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	Data  []*GetUserResp `json:"data"`
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
		total,
		page,
		size,
		data,
	}
}
