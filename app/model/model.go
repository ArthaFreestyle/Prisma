package model

type WebResponse[T any] struct {
	Data   T             `json:"data"`
	Paging *PageMetaData `json:"paging"`
	Errors string        `json:"errors"`
}

type PageMetaData struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int   `json:"total_page"`
}
