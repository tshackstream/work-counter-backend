package projects

type Project struct {
	ID                 int    `json:"id" form:"id"`
	ProjectName        string `json:"project_name" form:"project_name"`
	RewardType         int    `json:"reward_type" form:"reward_type"`
	LowerLimitTime     *int   `json:"lower_limit_time" form:"lower_limit_time"`
	LimitTime          *int   `json:"limit_time" form:"limit_time"`
	UnitPrice          *int   `json:"unit_price" form:"unit_price"`
	OverUnitPrice      *int   `json:"over_unit_price" form:"over_unit_price"`
	DeductionUnitPrice *int   `json:"deduction_unit_price" form:"deduction_unit_price"`
	HourlyWage         *int   `json:"hourly_wage" form:"hourly_wage"`
	WorkTimePerDay     *int   `json:"work_time_per_day"`
	RestTime           *int   `json:"rest_time" form:"rest_time"`
}
