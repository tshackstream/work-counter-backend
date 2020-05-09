package work_info

type WorkInfo struct {
	ID        int     `json:"id" form:"id"`
	ProjectId int     `json:"project_id" form:"project_id"`
	Date      string  `json:"date" form:"date"`
	Status    *int    `json:"status" form:"status"`
	Start     *string `json:"start" form:"start"`
	End       *string `json:"end" form:"end"`
	Rest      *string `json:"rest" form:"rest"`
	Total     *string `json:"total" form:"total"`
	Note      *string `json:"note" form:"note"`
}

func (WorkInfo) TableName() string {
	return "work_info"
}
