package sales_model

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/xormplus/xorm"
)

type Order struct {
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

type compareOrdCreateTime struct {
	OrderCreateTime int64 `xorm:"not null INT(11) 'ordcreatetime'"`
	MemberCreateT   int64 `xorm:"not null INT(11) 'mcreatetime'"`
}

type calculateOrdGoodsPrice struct {
	OrdCreateTime int64               `xorm:"not null INT(11) 'ordcreatetime'"`
	MembCreateT   int64               `xorm:"not null INT(11) 'mcreatetime'"`
	OrdPrice      decimal.NullDecimal `xorm:"not null default 0.00 DECIMAL(10,2) 'ordprice' COMMENT('price')"`
}

type EweiShopOrderModel struct {
	session *xorm.Session
}

func NewEweiShopOrderModel(session *xorm.Session) *EweiShopOrderModel {
	return &EweiShopOrderModel{session}
}

func (o *Order) TableName() string {
	return "ims_ewei_shop_order"
}

func (eweiShopOrderModel *EweiShopOrderModel) GetOrderStatisData(apltUniacid uint64, yesterdayStart, yesterdayEnd, orderStatus int64) (yNewOrder int64, yNewOrderTPrice float64, err error) {
	yNewOrder, err = eweiShopOrderModel.session.Where("createtime > ?", yesterdayStart).
		And("createtime <= ?", yesterdayEnd).
		And("status != ?", orderStatus).
		Count(&Order{UniacidOrd: apltUniacid})

	if err != nil {
		return 0, 0, err
	}

	yNewOrderTPrice, err = eweiShopOrderModel.session.Where("createtime > ?", yesterdayStart).
		And("createtime <= ?", yesterdayEnd).
		And("status != ?", orderStatus).
		Sum(&Order{UniacidOrd: apltUniacid}, "price")

	if err != nil {
		return 0, 0, err
	}
	return
}

func (eweiShopOrderModel *EweiShopOrderModel) CountRangeOrder(apltsUniacid uint64, timeRangeSrt, timeRangeEnd int64, inCondition ...interface{}) (int64, error) {
	compareTime := make([]compareOrdCreateTime, 0)

	err := eweiShopOrderModel.session.Table([]string{"ims_ewei_shop_order", "o"}).
		Select("o.createtime ordcreatetime, m.createtime mcreatetime").
		Join("LEFT", []string{"ims_ewei_shop_member", "m"}, "o.unionid = m.unionid").
		Where("o.createtime > ?", timeRangeSrt).
		And("o.createtime <= ?", timeRangeEnd).
		And("o.status != ?", -1).
		And("o.status != ?", 6).
		In("o.unionid", inCondition[0]).
		Find(&compareTime)

	if err != nil {
		fmt.Println("Error24 = ", err)
		return 0, err
	}

	rangeOrdNum := len(compareTime)

	for _, cmpt := range compareTime {

		if cmpt.OrderCreateTime-cmpt.MemberCreateT > 86400 {
			rangeOrdNum--
		}
	}

	return int64(rangeOrdNum), nil
}

func (eweiShopOrderModel *EweiShopOrderModel) CountRangeOrdersPrice(apltsUniacid uint64, timeRangeSrt, timeRangeEnd int64, inCondition ...interface{}) (rangeOrdPrice decimal.Decimal, err error) {
	compareItems := make([]calculateOrdGoodsPrice, 0)

	err = eweiShopOrderModel.session.Table([]string{"ims_ewei_shop_order", "o"}).
		Select("o.createtime ordcreatetime, m.createtime mcreatetime, o.price ordprice").
		Join("LEFT", []string{"ims_ewei_shop_member", "m"}, "o.unionid = m.unionid").
		Where("o.createtime > ?", timeRangeSrt).
		And("o.createtime <= ?", timeRangeEnd).
		And("o.status != ?", -1).
		And("o.status != ?", 6).
		In("o.unionid", inCondition[0]).
		Find(&compareItems)

	if err != nil {
		return rangeOrdPrice, err
	}

	for _, cmpt := range compareItems {

		if cmpt.OrdCreateTime-cmpt.MembCreateT > 86400 {
			continue
		}
		rangeOrdPrice = rangeOrdPrice.Add(cmpt.OrdPrice.Decimal)
	}

	return rangeOrdPrice, nil
}
