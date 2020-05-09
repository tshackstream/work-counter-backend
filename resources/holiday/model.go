package holiday

type Holiday struct {
	ID          int    `json:"id"`
	Date        string `json:"date" form:"date"`
	HolidayName string `json:"holiday_name"`
}

func (Holiday) TableName() string {
	return "holiday_master"
}
