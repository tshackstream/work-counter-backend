package status

type Status struct {
	ID         int    `json:"value"`
	StatusName string `json:"label"`
}

func (Status) TableName() string {
	return "status_master"
}
