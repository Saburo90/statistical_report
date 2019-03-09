package recycle_model

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/xormplus/xorm"
)

type BookOrder struct {
	Id            uint64              `xorm:"not null autoincr pk INT(11) 'id' COMMENT('order id')"`
	UniacidOrd    uint64              `xorm:"not null TINYINT 'uniacid' COMMENT('applets flag')"`
	OrderSN       string              `xorm:"not null VARCHAR(30) 'ordersn' COMMENT('order sn')"`
	Price         decimal.NullDecimal `xorm:"not null default 0.00 DECIMAL(10,2) 'price' COMMENT('order price')"`
	StatusOrd     int64               `xorm:"not null TINYINT 'status' COMMENT('order status')"`
	OpenidOrd     string              `xorm:"not null VARCHAR(50) 'openid' COMMENT('order user openid')"`
	OpenidWaOrd   string              `xorm:"not null VARCHAR(64) 'openid_wa' COMMENT('order user openid_wa')"`
	UnionidOrd    string              `xorm:"not null VARCHAR(64) 'unionid' COMMENT('order user unionid')"`
	CreateTimeOrd int64               `xorm:"not null INT(11) 'createtime' COMMENT('order create time')"`
}

type compareCreateTime struct {
	BookOrderCreateTime int64 `xorm:"not null INT(11) 'rcy_ordcreatetime'"`
	MemberCreateTime    int64 `xorm:"not null INT(11) 'm_createtime'"`
}

type calculateBookGoodsPrice struct {
	BordCreateTime int64               `xorm:"not null INT(11) 'rcy_ordcreatetime'"`
	MembCreateTime int64               `xorm:"not null INT(11) 'm_createtime'"`
	BookOrdPrice   decimal.NullDecimal `xorm:"not null default 0.00 DECIMAL(10,2) 'rcy_ordprice' COMMENT('price')"`
}

type EweiShopBookOrderModel struct {
	session *xorm.Session
}

func NewEweiShopBookOrderModel(session *xorm.Session) *EweiShopBookOrderModel {
	return &EweiShopBookOrderModel{session}
}

func (o *BookOrder) TableName() string {
	return "ims_ewei_shop_book_order"
}

func (eweiShopBookOrderModel *EweiShopBookOrderModel) GetOrderStatisData(apltUniacid uint64, yesterdayStart, yesterdayEnd, orderStatus int64) (yNewOrder int64, yNewOrderTPrice float64, err error) {
	yNewOrder, err = eweiShopBookOrderModel.session.Where("createtime > ?", yesterdayStart).
		And("createtime <= ?", yesterdayEnd).
		And("status != ?", orderStatus).
		And("status != ?", 6).
		Count(&BookOrder{UniacidOrd: apltUniacid})

	if err != nil {
		return 0, 0, err
	}

	yNewOrderTPrice, err = eweiShopBookOrderModel.session.Where("createtime > ?", yesterdayStart).
		And("createtime <= ?", yesterdayEnd).
		And("status != ?", orderStatus).
		And("status != ?", 6).
		Sum(&BookOrder{UniacidOrd: apltUniacid}, "price")

	if err != nil {
		return 0, 0, err
	}
	return
}

func (eweiShopBookOrderModel *EweiShopBookOrderModel) CountRangeOrder(apltsUniacid uint64, timeRangeSrt, timeRangeEnd int64, inCondition ...interface{}) (int64, error) {
	compareTime := make([]compareCreateTime, 0)

	err := eweiShopBookOrderModel.session.Table([]string{"ims_ewei_shop_book_order", "b"}).
		Select("b.createtime rcy_ordcreatetime, m.createtime m_createtime").
		Join("LEFT", []string{"ims_ewei_shop_member", "m"}, "b.unionid = m.unionid").
		Where("b.createtime > ?", timeRangeSrt).
		And("b.createtime <= ?", timeRangeEnd).
		And("b.status != ?", -1).
		And("b.status != ?", 6).
		In("b.unionid", inCondition[0]).
		Find(&compareTime)

	if err != nil {
		fmt.Println("Error24 = ", err)
		return 0, err
	}

	rangeOrdNum := len(compareTime)

	for _, cmpt := range compareTime {

		if cmpt.BookOrderCreateTime-cmpt.MemberCreateTime > 86400 {
			rangeOrdNum--
		}
	}

	return int64(rangeOrdNum), nil
}

func (eweiShopBookOrderModel *EweiShopBookOrderModel) CountRangeOrdersPrice(apltsUniacid uint64, timeRangeSrt, timeRangeEnd int64, inCondition ...interface{}) (rangeOrdPrice decimal.Decimal, err error) {
	compareItems := make([]calculateBookGoodsPrice, 0)

	err = eweiShopBookOrderModel.session.Table([]string{"ims_ewei_shop_book_order", "b"}).
		Select("b.createtime rcy_ordcreatetime, m.createtime m_createtime, b.price rcy_ordprice").
		Join("LEFT", []string{"ims_ewei_shop_member", "m"}, "b.unionid = m.unionid").
		Where("b.createtime > ?", timeRangeSrt).
		And("b.createtime <= ?", timeRangeEnd).
		And("b.status != ?", -1).
		And("b.status != ?", 6).
		In("b.unionid", inCondition[0]).
		Find(&compareItems)

	if err != nil {
		return rangeOrdPrice, err
	}

	for _, cmpt := range compareItems {

		if cmpt.BordCreateTime-cmpt.MembCreateTime > 86400 {
			continue
		}
		rangeOrdPrice = rangeOrdPrice.Add(cmpt.BookOrdPrice.Decimal)
	}

	return rangeOrdPrice, nil
}
