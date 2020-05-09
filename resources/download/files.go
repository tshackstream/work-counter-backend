// cf https://qiita.com/tebakane/items/2f2ed2558357c274c478

package download

import (
	"backend/resources"
	"backend/resources/holiday"
	"backend/resources/monthly_work_results"
	"backend/resources/projects"
	"backend/resources/work_info"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/loadoff/excl"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func WorkSheet(c echo.Context) error {
	year := c.Param("year")
	month := c.Param("month")
	projectId := c.Param("project_id")

	dateList, err := work_info.GetMonthInfo(year, month, projectId)

	if err != nil {
		return err
	}

	workResultModel := monthly_work_results.MonthlyWorkResult{}
	condition := []resources.WhereCondition{
		{Column: "project_id", Operation: "=", Value: projectId},
		{Column: "month", Operation: "=", Value: year + "-" + month + "-01"},
	}

	res, err := resources.FetchOne(&workResultModel, condition)

	if err != nil {
		return nil
	}

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, nil)
	} else if res.Error != nil {
		return res.Error
	}

	fileName, err := makeWorkSheet(year, month, dateList, workResultModel)

	if err != nil {
		return err
	}

	return c.Attachment("data/"+fileName, fileName)
}

func Invoice(c echo.Context) error {
	year := c.Param("year")
	month := c.Param("month")
	projectId := c.Param("project_id")

	workResultModel := monthly_work_results.MonthlyWorkResult{}
	condition := []resources.WhereCondition{
		{Column: "project_id", Operation: "=", Value: projectId},
		{Column: "month", Operation: "=", Value: year + "-" + month + "-01"},
	}

	workResultRes, err := resources.FetchOne(&workResultModel, condition)

	if err != nil {
		return nil
	}

	if errors.Is(workResultRes.Error, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, nil)
	} else if workResultRes.Error != nil {
		return workResultRes.Error
	}

	projectModel := projects.Project{}
	projectCondition := []resources.WhereCondition{
		{Column: "id", Operation: "=", Value: c.Param("project_id")},
	}

	projectRes, err := resources.FetchOne(&projectModel, projectCondition)
	if err != nil {
		return err
	}

	if errors.Is(projectRes.Error, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusOK, nil)
	} else if projectRes.Error != nil {
		return projectRes.Error
	}

	dateList, err := work_info.GetMonthInfo(year, month, projectId)

	if err != nil {
		return err
	}

	fileName, err := makeInvoice(year, month, workResultModel, projectModel, dateList)

	if err != nil {
		return err
	}

	return c.Attachment("data/"+fileName, fileName)
}

func makeWorkSheet(year string, month string, dateList map[int]work_info.DateInfo, workResult monthly_work_results.MonthlyWorkResult) (string, error) {
	fileName := year + month + "worksheet.xlsx"

	if f, err := os.Stat("data"); os.IsNotExist(err) || !f.IsDir() {
		return "", errors.New("dataディレクトリが存在しません")
	}

	if _, err := os.Stat("data/" + fileName); err != nil {
		os.Remove("data/" + fileName)
	}

	file, err := excl.Open("data/work_sheet.xlsx")
	if err != nil {
		return "", err
	}

	sheet, _ := file.OpenSheet("Sheet1")

	// 対象年月
	row := sheet.GetRow(4)
	row.SetString(year+"/"+month, 9)

	// 各データ
	rowNum := 10
	dateLayout := "2006-01-02"

	// mapは順番がバラバラなのでキーをソートしつつ処理
	// cf https://golang.hateblo.jp/entry/2019/10/07/202140
	indexes := make([]int, len(dateList), len(dateList))

	i := 0
	for key := range dateList {
		indexes[i] = key
		i++
	}

	sort.Ints(indexes)

	for i := 0; i < len(indexes); i++ {
		data := dateList[indexes[i]]

		row := sheet.GetRow(rowNum)
		dateObj, _ := time.Parse(dateLayout, data.Date)

		row.SetString(month+"/"+strconv.Itoa(dateObj.Day()), 2)
		row.SetString(data.DayOfWeek, 3)

		if data.Start != nil && data.Status != nil && *data.Status != 2 {
			row.SetString(*data.Start, 4)
		}
		if data.End != nil && data.Status != nil && *data.Status != 2 {
			row.SetString(*data.End, 5)
		}
		if data.Rest != nil && data.Status != nil && *data.Status != 2 {
			row.SetString(*data.Rest, 6)
		}
		if data.Total != nil && data.Status != nil && *data.Status != 2 {
			row.SetString(*data.Total, 7)
		}
		if data.Note != nil {
			row.SetString(*data.Note, 8)
		}

		if data.DayOfWeek == "土" {
			for i := 2; i <= 8; i++ {
				cell := row.GetCell(i)
				cell.SetBackgroundColor("A4C2F4")
			}
		}

		if data.DayOfWeek == "日" || data.IsHoliday {
			for i := 2; i <= 8; i++ {
				cell := row.GetCell(i)
				cell.SetBackgroundColor("EA9999")
			}

			if data.IsHoliday && len(data.HolidayNote) > 0 {
				row.SetString(data.HolidayNote, 8)
			}
		}

		rowNum++
	}

	inputDateRow := sheet.GetRow(42)
	inputDateRow.SetNumber(workResult.InputDay, 7)

	workTimeRow := sheet.GetRow(43)
	workTimeRow.SetString(workResult.WorkTime, 7)
	sheet.Close()

	file.Save("data/" + fileName)

	return fileName, nil
}

func makeInvoice(year string, month string,
	workResult monthly_work_results.MonthlyWorkResult, project projects.Project, dateList map[int]work_info.DateInfo) (string, error) {

	fileName := year + month + "invoice.xlsx"

	if f, err := os.Stat("data"); os.IsNotExist(err) || !f.IsDir() {
		return "", errors.New("dataディレクトリが存在しません")
	}

	if _, err := os.Stat("data/" + fileName); err != nil {
		os.Remove("data/" + fileName)
	}

	file, err := excl.Open("data/invoice.xlsx")
	if err != nil {
		return "", err
	}
	sheet, _ := file.OpenSheet("Sheet1")

	// 最終日
	yearInt, _ := strconv.Atoi(year)
	monthInt, _ := strconv.Atoi(month)
	t := time.Date(yearInt, time.Month(monthInt+1), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	dateLayout := "2006-01-02"

	//請求日・請求書番号
	row := sheet.GetRow(12)
	row.SetString(t.Format(dateLayout), 13)
	row = sheet.GetRow(13)
	row.SetString(t.Format("200601"), 13)

	// 合計金額
	row = sheet.GetRow(15)
	row.SetNumber(workResult.ProspectedReward, 5)

	// 請求金額
	row = sheet.GetRow(18)
	row.SetNumber(workResult.ProspectedReward, 13)

	// 振込期日

	// 祝日を取得
	var holidayModel []holiday.Holiday
	condition := []resources.WhereCondition{
		{Column: "date", Operation: "BETWEEN", Value: []string{year + "-01-01", year + "-12-31"}},
	}
	holiday, err := resources.Fetch(&holidayModel, condition, nil)
	if err != nil {
		return "", err
	} else if holiday.Error != nil {
		return "", holiday.Error
	}
	holidayList := map[string]string{}
	if holiday.RowsAffected != 0 {
		for _, holidays := range holidayModel {
			holidayList[holidays.Date] = holidays.HolidayName
		}
	}

	var deadLine string
	// 二十日以降一週間のうち
	for i := 20; i < 27; i++ {
		date := time.Date(yearInt, time.Month(monthInt+1), i, 0, 0, 0, 0, time.Local)
		formattedDate := date.Format(dateLayout)
		// 土日ではない
		if date.Weekday() != 0 && date.Weekday() != 6 {
			// 祝日ではない
			if _, exists := holidayList[formattedDate]; !exists {
				deadLine = formattedDate
				break
			}
		}
	}
	row = sheet.GetRow(23)
	row.SetString(deadLine, 4)

	// 備考
	decimalWorkTimeStr := strconv.FormatFloat(workResult.ProspectedDecimalWorkTime, 'f', 2, 64)
	row = sheet.GetRow(30)
	row.SetString(decimalWorkTimeStr+"h", 4)
	if project.LowerLimitTime != nil && project.LimitTime != nil &&
		(workResult.ProspectedDecimalWorkTime < float64(*project.LowerLimitTime) ||
			float64(*project.LimitTime) > workResult.ProspectedDecimalWorkTime) {

		var labelPrefix string
		var timeFormula string
		var deductOrOverAmountFormula string
		var invoiceFomula string

		lowerLimitTimeFloat := float64(*project.LowerLimitTime)
		limitTimeFloat := float64(*project.LimitTime)
		lowerLimitTimeStr := strconv.Itoa(*project.LowerLimitTime)
		limitTimeStr := strconv.Itoa(*project.LimitTime)

		if workResult.ProspectedDecimalWorkTime < lowerLimitTimeFloat {
			labelPrefix = "控除"

			deductedTime := lowerLimitTimeFloat - workResult.ProspectedDecimalWorkTime
			deductedTimeStr := strconv.FormatFloat(deductedTime, 'f', 2, 64)
			timeFormula = lowerLimitTimeStr + "h - " + decimalWorkTimeStr + "h = " + deductedTimeStr + "h"

			deductOrOverAmount := int(float64(*project.DeductionUnitPrice) * deductedTime)
			deductOrOverAmountStr := numFormat(deductOrOverAmount)
			deductOrOverAmountFormula = "\\" + numFormat(*project.DeductionUnitPrice) + " * " + deductedTimeStr + "h = \\" + deductOrOverAmountStr

			invoice := *project.UnitPrice - deductOrOverAmount
			invoiceFomula = "\\" + numFormat(*project.UnitPrice) + " - " + "\\" + deductOrOverAmountStr + " = \\" + numFormat(invoice)

			// 控除/超過時間
			row = sheet.GetRow(31)
			row.SetString(labelPrefix+"時間", 3)
			row.SetString(timeFormula, 4)

			// 控除/超過金額
			row = sheet.GetRow(32)
			row.SetString(labelPrefix+"金額", 3)
			row.SetString(deductOrOverAmountFormula, 4)
			row.SetString("(小数点以下切り捨て)", 6)

			//請求金額
			row = sheet.GetRow(33)
			row.SetString("請求金額", 3)
			row.SetString(invoiceFomula, 4)
		} else if workResult.ProspectedDecimalWorkTime > limitTimeFloat {
			labelPrefix = "超過"

			overTime := workResult.ProspectedDecimalWorkTime - limitTimeFloat
			overTimeStr := strconv.FormatFloat(overTime, 'f', 2, 64)
			timeFormula = limitTimeStr + "h - " + decimalWorkTimeStr + "h = " + overTimeStr + "h"

			deductOrOverAmount := int(float64(*project.OverUnitPrice) * overTime)
			deductOrOverAmountStr := numFormat(deductOrOverAmount)
			deductOrOverAmountFormula = "\\" + strconv.Itoa(*project.OverUnitPrice) + " * " + overTimeStr + "h = \\" + deductOrOverAmountStr

			invoice := *project.UnitPrice + deductOrOverAmount
			invoiceFomula = "\\" + numFormat(*project.UnitPrice) + " + " + "\\" + deductOrOverAmountStr + " = " + "\\" + numFormat(invoice)

			// 控除/超過時間
			row = sheet.GetRow(31)
			row.SetString(labelPrefix+"時間", 3)
			row.SetString(timeFormula, 4)

			// 控除/超過金額
			row = sheet.GetRow(32)
			row.SetString(labelPrefix+"金額", 3)
			row.SetString(deductOrOverAmountFormula, 4)
			row.SetString("(小数点以下切り捨て)", 6)

			//請求金額
			row = sheet.GetRow(33)
			row.SetString("請求金額", 3)
			row.SetString(invoiceFomula, 4)
		}
	}

	sheet.Close()

	file.Save("data/" + fileName)

	return fileName, nil
}

// cf http://psychedelicnekopunch.com/archives/1520
func numFormat(i int) string {
	arr := strings.Split(fmt.Sprintf("%d", i), "")
	cnt := len(arr) - 1
	res := ""
	i2 := 0
	for i := cnt; i >= 0; i-- {
		if i2 > 2 && i2%3 == 0 {
			res = fmt.Sprintf(",%s", res)
		}
		res = fmt.Sprintf("%s%s", arr[i], res)
		i2++
	}
	return res
}
