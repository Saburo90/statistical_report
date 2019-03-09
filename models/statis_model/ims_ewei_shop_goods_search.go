package statis_model

type SearchGoods struct {
	Id          uint64 `xorm:"not null autoincr pk 'id'"`
	UniacidSear uint64 `xorm:"not null unsigned INT(11) 'uniacid'"`
	OpenidSear  string `xorm:"not null VARCHAR(50) 'openid'"`
	Searchkey   string `xorm:"not null VARCHAR(100) 'search'"`
	SearchTime  uint64 `xorm:"not null unsigned INT(11) 'search_time'"`
	UnionidSear string `xorm:"not null VARCHAR(64) 'unionid'"`
	UidSear     uint64 `xorm:"not null unsigned INT(11) 'uid'"`
}

type CountSearchPv struct {
	Id         uint64 `xorm:"not null autoincr pk 'id'"`
	SearchTime uint64 `xorm:"not null unsigned INT(11) 'search_time'"`
}

func (s *SearchGoods) TableName() string {
	return "ims_ewei_shop_goods_search"
}
