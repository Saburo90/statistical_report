package goods_model

import (
	"github.com/shopspring/decimal"
	"github.com/xormplus/xorm"
)

type Book struct {
	Id             uint64              `xorm:"not null autoincr pk INT(11) 'id'"`
	Buniacid       uint64              `xorm:"not null INT(11) 'uniacid'"`
	BRatingAverage decimal.NullDecimal `xorm:"not null default 0.00 DECIMAL(10,2) 'rating_average'"`
}

type EweiShopBookModel struct {
	session *xorm.Session
}

func NewEweiShopBookModel(session *xorm.Session) *EweiShopBookModel {
	return &EweiShopBookModel{session}
}

func (g *Book) TableName() string {
	return "ims_ewei_shop_book"
}

// 通过ISBN获取豆瓣评分
func (bModel *EweiShopBookModel) GetDoubanScoreBySN(buniacid uint64, bsn string) (decimal.NullDecimal, error) {
	book := Book{}

	var dScore decimal.NullDecimal
	exist, err := bModel.session.Select("id, rating_average").
		Where("uniacid = ?", buniacid).
		And("isbn13 = ?", bsn).
		Get(&book)

	if err != nil {
		return dScore, err
	}

	if !exist {
		return dScore, nil
	}

	return book.BRatingAverage, nil
}
