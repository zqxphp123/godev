package data

import "time"

type Goods struct {
	ID       uint64    `json:"id"`
	Mobile   string    `json:"mobile"`
	NickName string    `json:"nick_name"`
	Birthday time.Time `gorm:"type:datetime"`
	Gender   string    `json:"gender"`
	Role     int32     `json:"role"`
	PassWord string    `json:"password"`
}

type GoodsListDO struct {
	TotalCount int64   `json:"totalCount,omitempty"`
	Items      []*User `json:"items"`
}
