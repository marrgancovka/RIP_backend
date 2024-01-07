package ds

import "time"

type Application struct {
	ID             uint `gorm:"primarykey"`
	Status         string
	Date_creation  time.Time
	Date_formation time.Time
	Date_end       time.Time
	Id_user        uint
	Id_admin       uint
}

type ApplicationReq struct {
	ID             uint `gorm:"primarykey"`
	Status         string
	Date_creation  time.Time
	Date_formation time.Time
	Date_end       time.Time
	User           string
	Admin          string
}
