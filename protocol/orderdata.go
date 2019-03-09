package protocol

import "github.com/shopspring/decimal"

type (
	GetOrderDataResp struct {
		YNorder           int64           `json:"yesterday_new_order"`
		YNorderPrice      float64         `json:"yesterday_neworder_price"`
		MeanOrdPrice      decimal.Decimal `json:"yesterday_mean_ordprice"`
		MeanOrdGoos       decimal.Decimal `json:"yesterday_mean_ordgoods"`
		ConversionOrdRate decimal.Decimal `json:"yesterday_conv_ordrate"`
		TNorderGoodsT     int64           `json:"yesterday_totalorder_goods"`
		OrderNumIn24      int64           `json:"order_num_24hours"`
		ConverOrdRate24   decimal.Decimal `json:"conv_ordrate_in24"`
		MeanOrdPrice24    decimal.Decimal `json:"mean_ordprice_in24"`
		MeanOrdGoodsIn24  decimal.Decimal `json:"mean_ordgoods_in24"`
	}

	GetRcyOrderDataResp struct {
		RcyYNorder           int64           `json:"Rcyyesterday_new_order"`
		RcyYNorderPrice      float64         `json:"Rcyyesterday_neworder_price"`
		RcyMeanOrdPrice      decimal.Decimal `json:"Rcyyesterday_mean_ordprice"`
		RcyMeanOrdGoos       decimal.Decimal `json:"Rcyyesterday_mean_ordgoods"`
		RcyConversionOrdRate decimal.Decimal `json:"Rcyyesterday_conv_ordrate"`
		RcyTNorderGoodsT     int64           `json:"Rcyyesterday_totalorder_goods"`
		RcyOrderNumIn24      int64           `json:"Rcyorder_num_24hours"`
		RcyConverOrdRate24   decimal.Decimal `json:"Rcyconv_ordrate_in24"`
		RcyMeanOrdPrice24    decimal.Decimal `json:"Rcymean_ordprice_in24"`
		RcyMeanOrdGoodsIn24  decimal.Decimal `json:"Rcymean_ordgoods_in24"`
	}
)
