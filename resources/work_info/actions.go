package work_info

import (
	"backend/resources"
	"backend/resources/holiday"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

type DateInfo struct {
	ID          int     `json:"id"`
	ProjectID   int     `json:"project_id"`
	Date        string  `json:"date"`
	DayOfWeek   string  `json:"day_of_week"`
	IsHoliday   bool    `json:"is_holiday"`
	HolidayNote string  `json:"holiday_note"`
	Status      *int    `json:"status"`
	StartHour   *string `json:"start_hour" form:"start_hour"`
	StartMinute *string `json:"start_minute" form:"start_minute"`
	EndHour     *string `json:"end_hour" form:"end_hour"`
	EndMinute   *string `json:"end_minute" form:"end_minute"`
	RestHour    *string `json:"rest_hour" form:"rest_hour"`
	RestMinute  *string `json:"rest_minute" form:"rest_minute"`
	Total       *string `json:"total"`
	Note        *string `json:"note"`
}

func List(c echo.Context) error {
	year := c.Param("year")
	month := c.Param("month")
	projectId := c.Param("project_id")

	dateList, err := GetMonthInfo(year, month, projectId)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dateList)
}

func GetMonthInfo(year string, month string, projectId string) (map[int]DateInfo, error) {
	var model []WorkInfo

	// 最終日
	yearInt, _ := strconv.Atoi(year)
	monthInt, _ := strconv.Atoi(month)
	t := time.Date(yearInt, time.Month(monthInt+1), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	firstDate := year + "-" + month + "-" + "01"
	lastDate := year + "-" + month + "-" + fmt.Sprintf("%02d", t.Day())
	condition := []resources.WhereCondition{
		{Column: "project_id", Operation: "=", Value: projectId},
		{Column: "date", Operation: "BETWEEN", Value: []string{firstDate, lastDate}},
	}
	order := []resources.Order{
		{Column: "date"},
	}

	workInfo, err := resources.Fetch(&model, condition, order)
	if err != nil {
		return nil, err
	} else if workInfo.Error != nil {
		return nil, workInfo.Error
	}

	var workInfoList []WorkInfo
	projectIdInt, _ := strconv.Atoi(projectId)
	if workInfo.RowsAffected == 0 {
		// その月の情報がない場合は作成
		for i := 1; i <= t.Day(); i++ {
			day := fmt.Sprintf("%02d", i)
			date := year + "-" + month + "-" + day
			data := WorkInfo{}
			data.ID = 0
			data.ProjectId = projectIdInt
			data.Date = date
			_, err := resources.Save(&data)
			if err != nil {
				return nil, err
			}
			workInfoList = append(workInfoList, data)
		}
	} else {
		workInfoList = model
	}

	var holidayModel []holiday.Holiday
	holidayCondition := []resources.WhereCondition{
		{Column: "date", Operation: "BETWEEN", Value: []string{firstDate, lastDate}},
	}

	holidayMst, err := resources.Fetch(&holidayModel, holidayCondition, nil)
	if err != nil {
		return nil, err
	} else if holidayMst.Error != nil {
		return nil, holidayMst.Error
	}

	holidayList := map[string]string{}
	if holidayMst.RowsAffected != 0 {
		for _, holidays := range holidayModel {
			holidayList[holidays.Date] = holidays.HolidayName
		}
	}

	dateList := map[int]DateInfo{}
	wDays := [...]string{"日", "月", "火", "水", "木", "金", "土"}
	dateLayout := "2006-01-02"
	for _, datum := range workInfoList {
		dateObj, _ := time.Parse(dateLayout, datum.Date)
		dayOfWeek := dateObj.Weekday()

		dateInfo := DateInfo{}
		dateInfo.ID = datum.ID
		dateInfo.Date = datum.Date
		dateInfo.ProjectID = datum.ProjectId
		dateInfo.DayOfWeek = wDays[dayOfWeek]

		if len(holidayList) > 0 {
			_, isHoliday := holidayList[datum.Date]
			dateInfo.IsHoliday = isHoliday
			if isHoliday {
				dateInfo.HolidayNote = holidayList[datum.Date]
			}
		}

		dateInfo.Status = datum.Status
		dateInfo.StartHour = datum.StartHour
		dateInfo.StartMinute = datum.StartMinute
		dateInfo.EndHour = datum.EndHour
		dateInfo.EndMinute = datum.EndMinute
		dateInfo.RestHour = datum.RestHour
		dateInfo.RestMinute = datum.RestMinute
		dateInfo.Total = datum.Total
		dateInfo.Note = datum.Note

		dateList[datum.ID] = dateInfo
	}

	return dateList, nil
}

func BulkUpdate(c echo.Context) error {
	model := new(map[int]WorkInfo)
	if err := c.Bind(model); err != nil {
		return err
	}

	var response []WorkInfo
	db, err := resources.DbConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	tx := db.Begin()
	for _, params := range *model {
		res := tx.Save(&params)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
		response = append(response, params)
	}
	tx.Commit()

	return c.JSON(http.StatusOK, response)
}
