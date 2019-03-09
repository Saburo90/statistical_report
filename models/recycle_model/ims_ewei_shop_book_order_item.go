package recycle_model

import "github.com/xormplus/xorm"

type EewiShopBookOrderGds struct {
	session *xorm.Session
}

type BookOrderGoods struct {
	Id      uint64 `xoml:"not null autoincr pk INT(11) 'id' COMMENT('primary key')"`
	Unionid string `xorm:"not null VARCHAR(64) 'unionid'"`
}
