package monthly_work_results

type MonthlyWorkResult struct {
	ID                        int     `json:"id"`
	ProjectId                 int     `json:"project_id"`
	Month                     string  `json:"month"`
	BusinessDay               int     `json:"business_day"`
	InputDay                  int     `json:"input_day"`
	WorkTime                  string  `json:"work_time"`
	ProspectedWorkTime        string  `json:"prospected_work_time"`
	ProspectedDecimalWorkTime float64 `json:"prospected_decimal_work_time"`
	ProspectedReward          int     `json:"prospected_reward"`
}
