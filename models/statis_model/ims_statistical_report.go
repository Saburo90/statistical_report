package statis_model

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/xormplus/xorm"
)

type StatisticalReport struct {
	Id                       uint64              `xorm:"not null autoincr pk INT(11) 'id'"`
	Daily                    string              `xorm:"not null VARCHAR(10) 'daily' COMMENT('日期字符串')"`
	TotalRoamU               int64               `xorm:"not null unsigned INT(11) 'total_roamu' COMMENT('漫游鲸截至昨日结束时间总用户数')"`
	TotalRapltsU             int64               `xorm:"not null unsigned INT(11) 'total_rapltsu' COMMENT('纯漫游鲸公众号截至昨日结束时间用户数')"`
	TotalRminiU              int64               `xorm:"not null unsigned INT(11) 'total_rminiu' COMMENT('纯漫游鲸小程序截至昨日结束时间用户数')"`
	TotalRplatNewU           int64               `xorm:"not null unsigned INT(11) 'total_rplat_newu' COMMENT('漫游平台昨日新增用户数')"`
	TotalRoamNewU            int64               `xorm:"not null unsigned INT(11) 'total_roam_newu' COMMENT('漫游鲸昨日新增用户数')"`
	TotalRapltsNewU          int64               `xorm:"not null unsigned INT(11) 'total_raplts_newu' COMMENT('漫游鲸公众号昨日新增用户数')"`
	TotalSapltsNewU          int64               `xorm:"not null unsigned INT(11) 'total_saplt_newu' COMMENT('漫游共享店昨日新增用户数,去重')"`
	TotalRminiNewU           int64               `xorm:"not null unsigned INT(11) 'total_rmini_newu' COMMENT('漫游鲸小程序昨日新增用户数')"`
	TotalRapltsFans          int64               `xorm:"not null unsigned INT(11) 'total_raplts_fans' COMMENT('漫游鲸公众号截至昨日结束时间总粉丝数')"`
	TotalRapltsNewF          int64               `xorm:"not null unsigned INT(11) 'total_raplts_newf' COMMENT('漫游鲸公众号昨日新增粉丝数')"`
	TotalRapltsNewCF         int64               `xorm:"not null unsigned INT(11) 'total_raplts_newcf' COMMENT('漫游鲸公众号昨日取关粉丝数')"`
	TotalSapltsFans          int64               `xorm:"not null unsigned INT(11) 'total_saplts_fans' COMMENT('漫游共享店公众号截至昨日结束时间总粉丝数')"`
	TotalSapltsNewF          int64               `xorm:"not null unsigned INT(11) 'total_saplts_newf' COMMENT('漫游共享店公众号昨日新增粉丝数')"`
	TotalSapltsNewCF         int64               `xorm:"not null unsigned INT(11) 'total_saplts_newcf' COMMENT('漫游共享店公众号昨日取关粉丝数')"`
	TotalRminiVisitsU        int64               `xorm:"not null unsigned INT(11) 'total_rmini_visitsu' COMMENT('漫游鲸小程序截至昨日结束时间总访问人数')"`
	TotalRminiNewVisitsUv    int64               `xorm:"not null unsigned INT(11) 'total_rmini_newvistsuv' COMMENT('漫游鲸小程序昨日新增访问人数')"`
	TotalRminiNewVisitsNo    int64               `xorm:"not null unsigned INT(11) 'total_rmini_newvistno' COMMENT('漫游鲸小程序昨日新增访问次数')"`
	TotalRminiNewVisitsUsr   int64               `xorm:"not null unsigned INT(11) 'total_rmini_newvistusr' COMMENT('漫游鲸小程序昨日新增用户数')"`
	TotalRminiNewRcyPageU    int64               `xorm:"not null unsigned INT(11) 'total_rmini_newrcypageu' COMMENT('漫游鲸小程序回收页昨日访问人数')"`
	TotalRminiNewRcyPageN    int64               `xorm:"not null unsigned INT(11) 'total_rmini_newrcypagen' COMMENT('漫游鲸小程序回收页昨日访问次数')"`
	TotalRminiNewBdePageU    int64               `xorm:"not null unsigned INT(11) 'total_rmini_newbdepageu' COMMENT('漫游鲸小程序商品详情页昨日访问人数')"`
	TotalRminiNewBdePageN    int64               `xorm:"not null unsigned INT(11) 'total_rmini_newbdepagen' COMMENT('漫游鲸小程序商品详情页昨日访问次数')"`
	TotalRminiNewVitStay     float64             `xorm:"not null unsigned INT(11) 'total_rmini_newvitstay' COMMENT('漫游鲸小程序昨日用户人均停留时长')"`
	TotalRminiNewSearchPv    int64               `xorm:"not null unsigned INT(11) 'total_rmini_newsearpv' COMMENT('漫游鲸小程序昨日搜索商品次数')"`
	TotalRminiNewSearchUv    int64               `xorm:"not null unsigned INT(11) 'total_rmini_newsearuv' COMMENT('漫游鲸小程序昨日搜索商品人数')"`
	Meansearch               decimal.NullDecimal `xorm:"not null unsigned INT(11) 'total_rmini_newmeansear' COMMENT('漫游鲸小程序昨日人均搜索次数')"`
	TotalRorderNewNo         int64               `xorm:"not null unsigned INT(11) 'total_rorder_newno' COMMENT('漫游鲸昨日销售订单数')"`
	TotalRorderGoodsNewNo    int64               `xorm:"not null unsigned INT(11) 'total_rordgoods_newno' COMMENT('漫游鲸昨日销售商品数')"`
	MeanOrdPrice             decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'mean_ordprice' COMMENT('漫游鲸昨日销售客单价')"`
	MeanOrdGoods             decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'meam_ordgoods' COMMENT('漫游鲸昨日单均销量')"`
	ConversionOrdRate        decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(5,2) 'conversion_ordrate' COMMENT('漫游鲸昨日销售转化率')"`
	OrderNumIn24             int64               `xorm:"not null unsigned INT(11) 'ordernum_in24' COMMENT('漫游鲸24小时下单数')"`
	ConversionOrdRate24      decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(5,2) 'conversion_ordrate24' COMMENT('漫游鲸24小时销售转化率')"`
	MeanOrdPrice24           decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'mean_ordprice24' COMMENT('漫游鲸24小时销售客单价')"`
	MeanOrdGoods24           decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'meam_ordgoods24' COMMENT('漫游鲸24小时单均销量')"`
	TotalRcyRorderNewNo      int64               `xorm:"not null unsigned INT(11) 'totalrcy_rorder_newno' COMMENT('漫游鲸昨日回收订单数')"`
	TotalRcyRorderGoodsNewNo int64               `xorm:"not null unsigned INT(11) 'totalrcy_rordgoods_newno' COMMENT('漫游鲸昨日回收商品数')"`
	MeanRcyOrdPrice          decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'meanrcy_ordprice' COMMENT('漫游鲸昨日回收客单码洋')"`
	MeanRcyOrdGoods          decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'meamrcy_ordgoods' COMMENT('漫游鲸昨日单均回收本数')"`
	ConversionRcyOrdRate     decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(5,2) 'conversion_rcyordrate' COMMENT('漫游鲸昨日回收转化率')"`
	RcyOrderNumIn24          int64               `xorm:"not null unsigned INT(11) 'rcyordernum_in24' COMMENT('漫游鲸24小时回收下单数')"`
	ConversionRcyOrdRate24   decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(5,2) 'conversion_rcyordrate24' COMMENT('漫游鲸24小时回收转化率')"`
	MeanRcyOrdPrice24        decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'mean_rcyordprice24' COMMENT('漫游鲸24小时回收客单码洋')"`
	MeanRcyOrdGoods24        decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'meam_rcyordgoods24' COMMENT('漫游鲸24小时单均回收本数')"`
	RecordTime               string              `xorm:"not null unsigned INT(11) 'record_time' COMMENT('记录时间')"`
	Operator                 string              `xorm:"not null VARCHAR(32) 'operator' COMMENT('operator')"`
}

type OverViewStatis struct {
	Daily           string `xorm:"not null VARCHAR(10) 'daily' COMMENT('日期字符串')"`
	TotalRoamU      int64  `xorm:"not null unsigned INT(11) 'total_roamu' COMMENT('漫游鲸截至昨日结束时间总用户数')"`
	TotalRapltsU    int64  `xorm:"not null unsigned INT(11) 'total_rapltsu' COMMENT('纯漫游鲸公众号截至昨日结束时间用户数')"`
	TotalRminiU     int64  `xorm:"not null unsigned INT(11) 'total_rminiu' COMMENT('纯漫游鲸小程序截至昨日结束时间用户数')"`
	TotalRplatNewU  int64  `xorm:"not null unsigned INT(11) 'total_rplat_newu' COMMENT('漫游平台昨日新增用户数')"`
	TotalRoamNewU   int64  `xorm:"not null unsigned INT(11) 'total_roam_newu' COMMENT('漫游鲸昨日新增用户数')"`
	TotalRapltsNewU int64  `xorm:"not null unsigned INT(11) 'total_raplts_newu' COMMENT('漫游鲸公众号昨日新增用户数')"`
	TotalSapltsNewU int64  `xorm:"not null unsigned INT(11) 'total_saplt_newu' COMMENT('漫游共享店昨日新增用户数,去重')"`
	TotalRminiNewU  int64  `xorm:"not null unsigned INT(11) 'total_rmini_newu' COMMENT('漫游鲸小程序昨日新增用户数')"`
}

type AppletsStatis struct {
	Daily            string `xorm:"not null VARCHAR(10) 'daily' COMMENT('日期字符串')"`
	TotalRapltsFans  int64  `xorm:"not null unsigned INT(11) 'total_raplts_fans' COMMENT('漫游鲸公众号截至昨日结束时间总粉丝数')"`
	TotalRapltsNewF  int64  `xorm:"not null unsigned INT(11) 'total_raplts_newf' COMMENT('漫游鲸公众号昨日新增粉丝数')"`
	TotalRapltsNewCF int64  `xorm:"not null unsigned INT(11) 'total_raplts_newcf' COMMENT('漫游鲸公众号昨日取关粉丝数')"`
	TotalSapltsFans  int64  `xorm:"not null unsigned INT(11) 'total_saplts_fans' COMMENT('漫游共享店公众号截至昨日结束时间总粉丝数')"`
	TotalSapltsNewF  int64  `xorm:"not null unsigned INT(11) 'total_saplts_newf' COMMENT('漫游共享店公众号昨日新增粉丝数')"`
	TotalSapltsNewCF int64  `xorm:"not null unsigned INT(11) 'total_saplts_newcf' COMMENT('漫游共享店公众号昨日取关粉丝数')"`
}

type MiniStatis struct {
	Daily                  string              `xorm:"not null VARCHAR(10) 'daily' COMMENT('日期字符串')"`
	TotalRminiVisitsU      int64               `xorm:"not null unsigned INT(11) 'total_rmini_visitsu' COMMENT('漫游鲸小程序截至昨日结束时间总访问人数')"`
	TotalRminiNewVisitsUv  int64               `xorm:"not null unsigned INT(11) 'total_rmini_newvistsuv' COMMENT('漫游鲸小程序昨日新增访问人数')"`
	TotalRminiNewVisitsNo  int64               `xorm:"not null unsigned INT(11) 'total_rmini_newvistno' COMMENT('漫游鲸小程序昨日新增访问次数')"`
	TotalRminiNewVisitsUsr int64               `xorm:"not null unsigned INT(11) 'total_rmini_newvistusr' COMMENT('漫游鲸小程序昨日新增用户数')"`
	TotalRminiNewRcyPageU  int64               `xorm:"not null unsigned INT(11) 'total_rmini_newrcypageu' COMMENT('漫游鲸小程序回收页昨日访问人数')"`
	TotalRminiNewRcyPageN  int64               `xorm:"not null unsigned INT(11) 'total_rmini_newrcypagen' COMMENT('漫游鲸小程序回收页昨日访问次数')"`
	TotalRminiNewBdePageU  int64               `xorm:"not null unsigned INT(11) 'total_rmini_newbdepageu' COMMENT('漫游鲸小程序商品详情页昨日访问人数')"`
	TotalRminiNewBdePageN  int64               `xorm:"not null unsigned INT(11) 'total_rmini_newbdepagen' COMMENT('漫游鲸小程序商品详情页昨日访问次数')"`
	TotalRminiNewVitStay   float64             `xorm:"not null unsigned INT(11) 'total_rmini_newvitstay' COMMENT('漫游鲸小程序昨日用户人均停留时长')"`
	TotalRminiNewSearchPv  int64               `xorm:"not null unsigned INT(11) 'total_rmini_newsearpv' COMMENT('漫游鲸小程序昨日搜索商品次数')"`
	TotalRminiNewSearchUv  int64               `xorm:"not null unsigned INT(11) 'total_rmini_newsearuv' COMMENT('漫游鲸小程序昨日搜索商品人数')"`
	Meansearch             decimal.NullDecimal `xorm:"not null unsigned INT(11) 'total_rmini_newmeansear' COMMENT('漫游鲸小程序昨日人均搜索次数')"`
}

type OrderStatis struct {
	Daily                 string              `xorm:"not null VARCHAR(10) 'daily' COMMENT('日期字符串')"`
	TotalRorderNewNo      int64               `xorm:"not null unsigned INT(11) 'total_rorder_newno' COMMENT('漫游鲸昨日销售订单数')"`
	TotalRorderGoodsNewNo int64               `xorm:"not null unsigned INT(11) 'total_rordgoods_newno' COMMENT('漫游鲸昨日销售商品数')"`
	MeanOrdPrice          decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'mean_ordprice' COMMENT('漫游鲸昨日销售客单价')"`
	MeanOrdGoods          decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'meam_ordgoods' COMMENT('漫游鲸昨日单均销量')"`
	ConversionOrdRate     decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(5,2) 'conversion_ordrate' COMMENT('漫游鲸昨日销售转化率')"`
	OrderNumIn24          int64               `xorm:"not null unsigned INT(11) 'ordernum_in24' COMMENT('漫游鲸24小时下单数')"`
	ConversionOrdRate24   decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(5,2) 'conversion_ordrate24' COMMENT('漫游鲸24小时销售转化率')"`
	MeanOrdPrice24        decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'mean_ordprice24' COMMENT('漫游鲸24小时销售客单价')"`
	MeanOrdGoods24        decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'meam_ordgoods24' COMMENT('漫游鲸24小时单均销量')"`
}

type RcyOrderStatis struct {
	Daily                    string              `xorm:"not null VARCHAR(10) 'daily' COMMENT('日期字符串')"`
	TotalRcyRorderNewNo      int64               `xorm:"not null unsigned INT(11) 'totalrcy_rorder_newno' COMMENT('漫游鲸昨日回收订单数')"`
	TotalRcyRorderGoodsNewNo int64               `xorm:"not null unsigned INT(11) 'totalrcy_rordgoods_newno' COMMENT('漫游鲸昨日回收商品数')"`
	MeanRcyOrdPrice          decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'meanrcy_ordprice' COMMENT('漫游鲸昨日回收客单码洋')"`
	MeanRcyOrdGoods          decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'meamrcy_ordgoods' COMMENT('漫游鲸昨日单均回收本数')"`
	ConversionRcyOrdRate     decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(5,2) 'conversion_rcyordrate' COMMENT('漫游鲸昨日回收转化率')"`
	RcyOrderNumIn24          int64               `xorm:"not null unsigned INT(11) 'rcyordernum_in24' COMMENT('漫游鲸24小时回收下单数')"`
	ConversionRcyOrdRate24   decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(5,2) 'conversion_rcyordrate24' COMMENT('漫游鲸24小时回收转化率')"`
	MeanRcyOrdPrice24        decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'mean_rcyordprice24' COMMENT('漫游鲸24小时回收客单码洋')"`
	MeanRcyOrdGoods24        decimal.NullDecimal `xorm:"not null unsigned default 0.00 DECIMAL(10,2) 'meam_rcyordgoods24' COMMENT('漫游鲸24小时单均回收本数')"`
	RecordTime               int64               `xorm:"not null unsigned INT(11) 'record_time' COMMENT('记录时间')"`
	Operator                 string              `xorm:"not null VARCHAR(32) 'operator' COMMENT('operator')"`
}

type CheckDailyStatis struct {
	Daily string `xorm:"not null VARCHAR(10) 'daily' COMMENT('日期字符串')"`
}

type GetYesterdayRecord struct {
	Daily                 string `xorm:"not null VARCHAR(10) 'daily' COMMENT('日期字符串')"`
	TotalRminiNewVisitsUv int64  `xorm:"not null unsigned INT(11) 'total_rmini_newvistsuv' COMMENT('漫游鲸小程序昨日新增访问人数')"`
}

type GetYesterdayRcyRecord struct {
	Daily                 string `xorm:"not null VARCHAR(10) 'daily' COMMENT('日期字符串')"`
	TotalRminiNewRcyPageU int64  `xorm:"not null unsigned INT(11) 'total_rmini_newrcypageu' COMMENT('漫游鲸小程序回收页昨日访问人数')"`
}

type GetAftStatisRecord struct {
	Daily                  string `xorm:"not null VARCHAR(10) 'daily' COMMENT('日期字符串')"`
	TotalRminiNewVisitsUsr int64  `xorm:"not null unsigned INT(11) 'total_rmini_newvistusr' COMMENT('漫游鲸小程序昨日新增用户数')"`
}

type StatisticalReportModel struct {
	session *xorm.Session
}

func NewStatisticalReportModel(session *xorm.Session) *StatisticalReportModel {
	return &StatisticalReportModel{session}
}

func (statis *CheckDailyStatis) TableName() string {
	return "ims_statistical_report"
}

func (sModel *StatisticalReportModel) CheckDaily(dailyStr string) (int64, error) {
	dailyReport := CheckDailyStatis{Daily: dailyStr}
	dailyCount, err := sModel.session.Table(dailyReport.TableName()).Count(&dailyReport)
	fmt.Println("err = ", err)
	if err != nil {
		return 0, err
	}
	return dailyCount, nil
}
