package arrival_model

import (
	"github.com/xormplus/xorm"
)

type ArrivalNotice struct {
	Id           uint64 `xoml:"not null autoincr pk INT(11) 'id'"`
	ArriUniacid  uint64 `xorm:"not null unsigned INT(11) 'uniacid' COMMENT('用户所属公众号')"`
	ArriGoodsid  uint64 `xorm:"not null unsigned INT(11) 'goodsid' COMMENT('用户点击到货提醒的商品ID')"`
	ArriOpenid   string `xorm:"not null VARCHAR(50) 'openid' COMMENT('用户公众号openid')"`
	ArriDeleted  uint64 `xorm:"not null unsigned TINYINT(1) 'deleted' COMMENT('是否删除')"`
	ArriCreate   uint64 `xorm:"not null unsigned INT(11) 'createtime' COMMENT('创建时间')"`
	ArriNotice   uint64 `xorm:"not null unsigned TINYINT(4) 'arrival_notice' COMMENT('是否为到货通知1是0否')"`
	GoodsSN      string `xorm:"not null VARCHAR(50) 'goods_isbn' COMMENT('用户点击到货提醒商品的ISBN编码')"`
	ArriOpenidwa string `xorm:"not null VARCHAR(100) 'openid_wa' COMMENT('用户小程序openid')"`
	ArriUnionid  string `xorm:"not null VARCHAR(100) 'unionid' COMMENT('用户unionid')"`
}

type EweiShopMemFavoriteModel struct {
	session *xorm.Session
}

func (arri *ArrivalNotice) TableName() string {
	return "ims_ewei_shop_member_favorite"
}

func NewArrivalNoticeModel(session *xorm.Session) *EweiShopMemFavoriteModel {
	return &EweiShopMemFavoriteModel{session}
}

// 获取openid为空的到货提醒数据
func (arriModel *EweiShopMemFavoriteModel) GetEmptyOpenidArrival(appletsUniacid uint64) ([]ArrivalNotice, error) {
	emptyOpenidArri := make([]ArrivalNotice, 0)

	err := arriModel.session.Select("id, openid_wa, unionid").
		Where("uniacid = ?", appletsUniacid).
		And("openid = ?", "").
		And("unionid != ?", "").
		OrderBy("id DESC").
		Limit(20000).
		Find(&emptyOpenidArri)

	if err != nil {
		return nil, err
	}

	return emptyOpenidArri, nil
}

// 更新到货通知openid
func (arriModel *EweiShopMemFavoriteModel) UpdateOpenid(bean map[string]interface{}, arriId uint64) bool {
	_, err := arriModel.session.Table("ims_ewei_shop_member_favorite").
		Where("id = ?", arriId).
		Update(bean)

	if err != nil {
		return false
	}
	return true
}

// 获取unionid为空的到货提醒数据
func (arriModel *EweiShopMemFavoriteModel) GetEmptyUnionidArrival(appletsUniacid uint64) ([]ArrivalNotice, error) {
	emptyUnionidArri := make([]ArrivalNotice, 0)

	err := arriModel.session.Select("id, openid, openid_wa").
		Where("uniacid = ?", appletsUniacid).
		And("unionid = ?", "").
		OrderBy("id DESC").
		Limit(5000).
		Find(&emptyUnionidArri)

	if err != nil {
		return nil, err
	}

	return emptyUnionidArri, nil
}

// 更新到货通知unionid
func (arriModel *EweiShopMemFavoriteModel) UpdateUnionid(bean map[string]interface{}, arriId uint64) bool {
	_, err := arriModel.session.Table("ims_ewei_shop_member_favorite").
		Where("id = ?", arriId).
		Update(bean)

	if err != nil {
		return false
	}

	return true
}

// 删除到货通知记录
func (arriModel *EweiShopMemFavoriteModel) DeleteRecordById(id, uniacid uint64) error {
	_, err := arriModel.session.Delete(&ArrivalNotice{Id: id, ArriUniacid: uniacid})

	return err
}

// 获取openid_wa为空的到货提醒数据
func (arriModel *EweiShopMemFavoriteModel) GetEmptyOpenidWaArrival(appletsUniacid uint64) ([]ArrivalNotice, error) {
	emptyOpenidWaArri := make([]ArrivalNotice, 0)

	err := arriModel.session.Select("id, openid_wa, unionid").
		Where("uniacid = ?", appletsUniacid).
		And("openid_wa = ?", "").
		And("unionid != ?", "").
		OrderBy("id DESC").
		Limit(20000).
		Find(&emptyOpenidWaArri)

	if err != nil {
		return nil, err
	}

	return emptyOpenidWaArri, nil
}

// 更新到货通知openid_wa
func (arriModel *EweiShopMemFavoriteModel) UpdateOpenidWa(bean map[string]interface{}, arriId uint64) bool {
	_, err := arriModel.session.Table("ims_ewei_shop_member_favorite").
		Where("id = ?", arriId).
		Update(bean)

	if err != nil {
		return false
	}
	return true
}
