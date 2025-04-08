package common

// Page 分页参数结构体
type Page struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}
