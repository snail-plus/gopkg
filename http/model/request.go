package model

type PageRequest struct {
	PageSize int `form:"pageSize" example:"50" binding:"required"` //pageSize
	PageNum  int `form:"pageNum" example:"1" binding:"required"`   //pageNum
	Params   map[string]interface{}
}

func (r *PageRequest) AddParam(key string, value interface{}) {
	if r.Params == nil {
		r.Params = make(map[string]interface{})
	}
	r.Params[key] = value
}

func (r *PageRequest) Offset() int {
	return (r.PageNum - 1) * r.PageSize
}
