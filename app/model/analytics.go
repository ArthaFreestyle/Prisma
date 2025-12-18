package model

type Statistics struct {
	Tahun string   `json:"tahun"`
	Data  Regional `json:"data"`
}

type Regional struct {
	International int `json:"international"`
	National      int `json:"national"`
	Regional      int `json:"regional"`
	Local         int `json:"local"`
}
