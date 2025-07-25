package models

type DutySchedule struct {
	TenantId string     `json:"tenantId"`
	DutyId   string     `json:"dutyId"`
	Time     string     `json:"time"`
	UserId   string     `json:"userid"`
	Username string     `json:"username"`
	Status   string     `json:"status"`
	Users    []DutyUser `json:"users" gorm:"users;serializer:json"`
}

type DutyUser struct {
	UserId   string `json:"userid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
}
