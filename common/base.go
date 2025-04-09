package common

import "fmt"

// Page 分页参数结构体
type Page struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

func (p *Page) String() string {
	return fmt.Sprintf("skip: %d, limit: %d", p.Skip, p.Limit)
}
