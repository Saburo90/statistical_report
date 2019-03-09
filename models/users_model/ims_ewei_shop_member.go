package users_model

import (
	"fmt"
	"github.com/Saburo90/statistical_report/components"
	"github.com/xormplus/xorm"
	"time"
)

const (
	defaultExpires = time.Duration(1 * 24 * 60 * 60)
)

type Member struct {
	Id         uint64 `xorm:"not null autoincr pk INT(11) 'id' COMMENT('primarykey')"`
	Uid        uint64 `xorm:"not null INT(11) 'uid' COMMENT('mc member ID')"`
	Uniacid    uint64 `xorm:"not null unsigned TINYINT 'uniacid' COMMENT('applets ID')"`
	Openid     string `xorm:"not null VARCHAR(64) 'openid' COMMENT('pub flag')"`
	OpenidWa   string `xorm:"not null VARCHAR(64) 'openid_wa' COMMENT('miniprogram flag')"`
	Unionid    string `xorm:"not null VARCHAR(64) 'unionid' COMMENT('paltform unique flag')"`
	CreateTime int64  `xorm:"not null INT(11) 'createtime' COMMENT('register time')"`
}

type EweiShopMemberModel struct {
	session *xorm.Session
}

func (m *Member) TableName() string {
	return "ims_ewei_shop_member"
}
func NewEweiShopMemberModel(session *xorm.Session) *EweiShopMemberModel {
	return &EweiShopMemberModel{session}
}

func dedupUsers(sUnion string, rUnios []string) {
	for i := 0; i < len(rUnios); i++ {
		if sUnion == rUnios[i] {
			rSuNo, err := components.Redis.Get("saburo02").Int()
			if err != nil {
				fmt.Println(err)
			}
			// if visit roaming whale minus 1
			rSuNo--
		}
	}
}

func (eweiShopMemberModel *EweiShopMemberModel) GetPlatRUsersNum(roamUniacid uint64, createTime int64) (int64, error) {
	tPUsers, err := eweiShopMemberModel.session.Where("unionid != ?", "").
		And("createtime > ?", 0).
		And("createtime <= ?", createTime).
		Count(&Member{Uniacid: roamUniacid})

	if err != nil {
		return 0, err
	}

	return tPUsers, nil
}

func (eweiShopMemberModel *EweiShopMemberModel) GetRoamAppletsUsersNum(appletsOpenid, miniOpenid string, appletsUniacid uint64, createTime int64) (int64, error) {
	tRusers, err := eweiShopMemberModel.session.Where("openid != ?", appletsOpenid).
		And("openid_wa = ?", miniOpenid).
		And("uniacid = ?", appletsUniacid).
		And("createtime > ?", 0).
		And("createtime <= ?", createTime).
		Count(&Member{})
	if err != nil {
		return 0, err
	}
	return tRusers, nil
}

func (eweiShopMemberModel *EweiShopMemberModel) GetRoamMiniUsersNum(appletsOpenid, miniOpenid string, appletsUniacid uint64, createTime int64) (int64, error) {
	tRmUsers, err := eweiShopMemberModel.session.Where("openid_wa != ?", miniOpenid).
		And("openid = ?", appletsOpenid).
		And("createtime > ?", 0).
		And("createtime <= ?", createTime).
		Count(&Member{Uniacid: appletsUniacid})

	if err != nil {
		return 0, nil
	}

	return tRmUsers, nil
}

func (eweiShopMemberModel *EweiShopMemberModel) GetYesterdayPNewUsersNum(roamUniacid, shareUniacid uint64, emptyUnionid string, yesterdayStart, yesterdayEnd int64) (ryNewUsers, pyNewUsers, syNewUsers int64, err error) {
	ryNewUsers, err = eweiShopMemberModel.session.Where("createtime > ?", yesterdayStart).
		And("createtime <= ?", yesterdayEnd).
		Count(&Member{Uniacid: roamUniacid})

	if err != nil {
		return 0, 0, 0, err
	}

	syNewU := make([]*Member, 0)

	syNewUsers, err = eweiShopMemberModel.session.Distinct("unionid").
		Where("unionid != ?", emptyUnionid).
		And("createtime > ?", yesterdayStart).
		And("createtime <= ?", yesterdayEnd).
		FindAndCount(&syNewU, &Member{Uniacid: shareUniacid})

	if err != nil {
		return 0, 0, 0, err
	}

	for _, sNuser := range syNewU {
		rMNo, err := eweiShopMemberModel.session.Where("unionid = ?", sNuser.Unionid).
			Count(&Member{Uniacid: roamUniacid})

		if err != nil {
			continue
		}

		if rMNo > 0 {
			syNewUsers--
		}
	}

	pyNewUsers = ryNewUsers + syNewUsers

	return ryNewUsers, pyNewUsers, syNewUsers, nil
}

func (eweiShopMemberModel *EweiShopMemberModel) GetRapltAndRminiNewUsersNum(roamUniacid uint64, emptyOpenid, emptyOpenidwa string, yesterdayStart, yesterdayEnd int64) (ryApltNewUsers, ryMiniNewUsers int64, err error) {
	ryApltNewUsers, err = eweiShopMemberModel.session.Where("openid != ?", emptyOpenid).
		And("openid_wa = ?", emptyOpenidwa).
		And("createtime > ?", yesterdayStart).
		And("createtime <= ?", yesterdayEnd).
		Count(&Member{Uniacid: roamUniacid})

	if err != nil {
		return 0, 0, err
	}

	ryMiniNewUsers, err = eweiShopMemberModel.session.Where("openid_wa != ?", emptyOpenidwa).
		And("openid = ?", emptyOpenid).
		And("createtime > ?", yesterdayStart).
		And("createtime <= ?", yesterdayEnd).
		Count(&Member{Uniacid: roamUniacid})

	if err != nil {
		return 0, 0, err
	}

	return ryApltNewUsers, ryMiniNewUsers, nil
}

func (eweiShopMemberModel *EweiShopMemberModel) GetTimeRangeUsers(apltsUniacid uint64, timeRangeStart, timeRangeEnd int64) ([]Member, error) {
	rangeMember := make([]Member, 0)

	err := eweiShopMemberModel.session.Select("id, openid, openid_wa, uniacid, unionid, createtime").
		Where("uniacid = ?", apltsUniacid).
		And("createtime > ?", timeRangeStart).
		And("createtime <= ?", timeRangeEnd).
		Find(&rangeMember)

	if err != nil {
		return nil, err
	}

	return rangeMember, nil
}

// 通过unionid获取openid
func (eweiShopMemberModel *EweiShopMemberModel) GetMemberOpenid(munionid string) (string, error) {
	mem := Member{}

	exist, err := eweiShopMemberModel.session.Select("openid").
		Where("unionid = ?", munionid).
		Get(&mem)

	if err != nil {
		return "", err
	}

	if !exist {
		return "", nil
	}
	return mem.Openid, nil
}

// 通过openid获取unionid
func (eweiShopMemberModel *EweiShopMemberModel) GetMemberUnionidByOpenid(mopenid string) (string, error) {
	mem := Member{}

	exist, err := eweiShopMemberModel.session.Select("unionid").
		Where("openid = ?", mopenid).
		Limit(1).
		Get(&mem)

	if err != nil {
		return "", err
	}

	if !exist {
		return "", nil
	}
	return mem.Unionid, nil
}

// 通过openid_wa获取unionid
func (eweiShopMemberModel *EweiShopMemberModel) GetMemberUnionidByOpenidWa(mopenidwa string) (string, error) {
	mem := Member{}

	exist, err := eweiShopMemberModel.session.Select("unionid").
		Where("openid_wa = ?", mopenidwa).
		Limit(1).
		Get(&mem)

	if err != nil {
		return "", err
	}

	if !exist {
		return "", nil
	}
	return mem.Unionid, nil
}

// 通过unionid获取openid_wa
func (eweiShopMemberModel *EweiShopMemberModel) GetMemberOpenidWa(munionid string) (string, error) {
	mem := Member{}

	exist, err := eweiShopMemberModel.session.Select("openid_wa").
		Where("unionid = ?", munionid).
		Get(&mem)

	if err != nil {
		return "", err
	}

	if !exist {
		return "", nil
	}
	return mem.OpenidWa, nil
}

func (eweiShopMemberModel *EweiShopMemberModel) GetPlatformUsersNum(roamUniacid, shareUniacid uint64, emptyUnionid string) (int64, error) {
	rUsers := make([]*Member, 0)
	sUsers := make([]*Member, 0)
	rUnionids := make([]string, 0)
	sUnionids := make([]string, 0)

	tRusers, err := eweiShopMemberModel.session.Distinct("unionid").
		And("unionid != ?", emptyUnionid).
		FindAndCount(&rUsers, &Member{Uniacid: roamUniacid})

	if err != nil {
		return 0, err
	}

	for _, runi := range rUsers {
		rUnionids = append(rUnionids, runi.Unionid)
	}
	//fmt.Println(rUsers[0].Unionid)
	//fmt.Println(rUnionids)
	//fmt.Println(len(rUnionids))
	rSusers, err := eweiShopMemberModel.session.Distinct("unionid").
		And("unionid != ?", emptyUnionid).
		FindAndCount(&sUsers, &Member{Uniacid: shareUniacid})

	if err != nil {
		return 0, err
	}
	components.Redis.Set("saburo01", tRusers, defaultExpires*time.Second)
	components.Redis.Set("saburo02", rSusers, defaultExpires*time.Second)
	for _, suni := range sUsers {
		sUnionids = append(sUnionids, suni.Unionid)
	}

	//for _, sUnionid := range sUnionids {
	//	go dedupUsers(sUnionid, rUnionids)
	//	//for i := 0; i < len(rUnionids); i++ {
	//	//	if sUnionid == rUnionids[i] {
	//	//		// if visit roaming whale minus 1
	//	//		rSusers--
	//	//	}
	//	//}
	//}
	return tRusers + rSusers, err
}
