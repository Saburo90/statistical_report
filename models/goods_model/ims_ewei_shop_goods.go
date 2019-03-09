package goods_model

import (
	"github.com/shopspring/decimal"
	"github.com/xormplus/xorm"
)

type Good struct {
	Id            uint64              `xorm:"not null autoincr pk INT(11) 'id'"`
	Guniacid      uint64              `xorm:"not null INT(11) 'uniacid'"`
	GTitle        string              `xorm:"not null VARCHAR(100) 'title'"`
	GProductSn    string              `xorm:"not null VARCHAR(50) 'productsn'"`
	GProductPrice decimal.NullDecimal `xorm:"not null default 0.00 DECIMAL(10,2) 'productprice'"`
}

type EweiShopGoodsModel struct {
	session *xorm.Session
}

func NewEweiShopGoodsModel(session *xorm.Session) *EweiShopGoodsModel {
	return &EweiShopGoodsModel{session}
}

func (g *Good) TableName() string {
	return "ims_ewei_shop_goods"
}

// 通过ISBN获取书名及原价
func (gModel *EweiShopGoodsModel) GetTitleAndProductPriceBySN(guniacid uint64, gsn string) (string, decimal.NullDecimal, error) {
	good := Good{}

	var proPrice decimal.NullDecimal
	exist, err := gModel.session.Select("id, title, productsn, productprice").
		Where("uniacid = ?", guniacid).
		And("productsn = ?", gsn).
		Get(&good)

	if err != nil {
		return "", proPrice, err
	}

	if !exist {
		return "", proPrice, nil
	}

	return good.GTitle, good.GProductPrice, nil
}
