package hanging

import (
	"github.com/Saburo90/statistical_report/components"
	"github.com/Saburo90/statistical_report/constant"
	"github.com/Saburo90/statistical_report/models/goods_model"
	"github.com/Luxurioust/excelize"
	"go.uber.org/zap"
	"strconv"
)

func CompleteBnameAndBpriceAndDScore() {
	zap.L().Info("开始读取Excel文件")
	xlsx, err := excelize.OpenFile("/home/gopath/excel/statistical_report/purchase.xlsx")

	if err != nil {
		zap.L().Error("读取失败", zap.Error(err))
		return
	}

	session := components.NewDBSession()

	defer session.Close()

	gModel := goods_model.NewEweiShopGoodsModel(session)

	bModel := goods_model.NewEweiShopBookModel(session)

	rows := xlsx.GetRows("Sheet1")
	for key, row := range rows {
		// 获取当前书籍书名、原价、豆瓣评分
		if key > 0 && row[0] != "" {
			title, productPrice, err := gModel.GetTitleAndProductPriceBySN(constant.RoamingApplets, row[0])

			if err != nil || title == "" {

			} else {
				xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(key+1), title)
				xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(key+1), productPrice.Decimal)
			}

			dScore, err := bModel.GetDoubanScoreBySN(constant.RoamingApplets, row[0])

			if err == nil {
				xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(key+1), dScore.Decimal)
			}
		}
	}

	xlsx.Save()
}
