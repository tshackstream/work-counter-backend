package work_info

type WorkInfo struct {
	ID          int     `json:"id" form:"id"`
	ProjectId   int     `json:"project_id" form:"project_id"`
	Date        string  `json:"date" form:"date"`
	Status      *int    `json:"status" form:"status"`
	StartHour   *string `json:"start_hour" form:"start_hour"`
	StartMinute *string `json:"start_minute" form:"start_minute"`
	EndHour     *string `json:"end_hour" form:"end_hour"`
	EndMinute   *string `json:"end_minute" form:"end_minute"`
	RestHour    *string `json:"rest_hour" form:"rest_hour"`
	RestMinute  *string `json:"rest_minute" form:"rest_minute"`
	Total       *string `json:"total" form:"total"`
	Note        *string `json:"note" form:"note"`
}

func (WorkInfo) TableName() string {
	return "work_info"
}
