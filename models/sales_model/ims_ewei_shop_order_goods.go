package sales_model

import (
	"github.com/xormplus/xorm"
)

type EewiShopOrderGds struct {
	session *xorm.Session
}

type OrderGoods struct {
	Id uint64 `xoml:"not null autoincr pk INT(11) 'id' COMMENT('primary key')"`
}
