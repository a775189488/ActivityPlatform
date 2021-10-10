package metadata

import "entrytask/internal/model"

type ActTypeCreateReq struct {
	Name string `json:"name" binding:"required"`
}

func (a *ActTypeCreateReq) ToActType() *model.ActivityType {
	return &model.ActivityType{
		Name: a.Name,
	}
}

type ActTypeUpdateReq struct {
	ActTypeCreateReq
}

type GetActTypeResp struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type ListActTypeReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ListActTypeResp struct {
	ListBase
	Data []*GetActTypeResp `json:"data"`
}

func FormateActTypeCreateResp(at *model.ActivityType) *GetActTypeResp {
	return &GetActTypeResp{
		Id:   at.Id,
		Name: at.Name,
	}
}

func FormateActTypeUpdateResp(at *model.ActivityType) *GetActTypeResp {
	return &GetActTypeResp{
		Id:   at.Id,
		Name: at.Name,
	}
}

func FormatListActTypeResp(ats []*model.ActivityType, total int, page int, size int) *ListActTypeResp {
	data := make([]*GetActTypeResp, 0)
	for _, a := range ats {
		data = append(data, &GetActTypeResp{Id: a.Id, Name: a.Name})
	}
	return &ListActTypeResp{
		ListBase: ListBase{
			page,
			size,
		},
		Data: data,
	}
}
