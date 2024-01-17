package repository

import (
	"awesomeProject/internal/app/ds"
	"time"

	"gorm.io/gorm"
)

// вывод списка всех заявок без услуг включенных в них + фильтрация по статусу и дате формирования
func (r *Repository) Select_applications(status string, date time.Time, date_end time.Time) (*[]ds.ApplicationReq, error) {
	var applications []ds.Application
	var user ds.Users
	var admin = ""
	var res *gorm.DB
	if status != "" {
		res = r.db.Where("status = ? AND status != 'delete' AND status <> 'created'", status).Where("date_formation BETWEEN ? AND ?", date, date_end).Order("date_formation DESC").Find(&applications)

	} else {
		res = r.db.Where("status <> ? AND status <> ?", "delete", "created").Where("date_formation BETWEEN ? AND ?", date, date_end).Order("date_formation DESC").Find(&applications)
	}
	response := make([]ds.ApplicationReq, len(applications))
	for i, app := range applications {
		r.db.Table("users").Select("username").Where("id = ?", app.Id_user).First(&user)
		if app.Id_admin != 0 {
			var admin_ds ds.Users
			adminRecord := r.db.Table("users").Select("username").Where("id = ?", app.Id_admin).First(&admin_ds)
			if adminRecord.Error != nil {
				return nil, adminRecord.Error
			}
			admin = admin_ds.UserName
		}

		response[i] = ds.ApplicationReq{
			ID:             app.ID,
			Status:         app.Status,
			Date_creation:  app.Date_creation,
			Date_formation: app.Date_formation,
			Date_end:       app.Date_end,
			User:           user.UserName,
			Admin:          admin,
		}
	}
	return &response, res.Error
}

func (r *Repository) Select_applications_buyer(status string, date time.Time, date_end time.Time, id_user uint) (*[]ds.Application, error) {
	var applications []ds.Application
	if status != "" {
		res := r.db.Where("id_user = ?", id_user).Where("status = ? AND status != 'delete'", status).Where("date_creation BETWEEN ? AND ?", date, date_end).Order("date_creation DESC").Find(&applications)
		return &applications, res.Error
	}
	res := r.db.Where("id_user = ?", id_user).Where("status <> ?", "delete").Where("date_creation BETWEEN ? AND ?", date, date_end).Order("date_creation DESC").Find(&applications)
	return &applications, res.Error
}

type ApplicationService struct {
	ID             int
	Title          string
	CosmodromBegin string
	CosmodromEnd   string
	Date           time.Time
}

// вывод одной заявки со списком её услуг
func (r *Repository) Select_application(id int) (*ds.Application, []ApplicationService, error) {
	var applications ds.Application
	var applicationService []ApplicationService

	//ищем такую заявку
	result := r.db.First(&applications, "id =?", id)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	if err := r.db.Table("applications").
		Select("Ships.id, Ships.title, CosmodromsBegin.title as cosmBegin, CosmodromsEnd.title as cosmEnd, Flights.date_flight").
		Joins("JOIN Flights ON Flights.id_application = Applications.id").
		Joins("JOIN Ships ON Flights.id_ship = Ships.id").
		Joins("JOIN Cosmodroms as CosmodromsBegin ON Flights.id_cosmodrom_begin = CosmodromsBegin.id").
		Joins("JOIN Cosmodroms as CosmodromsEnd ON Flights.id_cosmodrom_end = CosmodromsEnd.id").
		Where("Applications.id = ?", id).
		Scan(&applicationService).Error; err != nil {
		return nil, nil, err
	}

	return &applications, applicationService, nil
}

// изменение статуса модератора
func (r *Repository) Update_application_admin(id uint, status string, id_admin uint) error {
	app := ds.Application{}
	res1 := r.db.First(&app, "id = ? AND status = 'formated'", id)
	if res1.Error != nil {
		return res1.Error
	}
	app.Status = status
	app.Date_end = time.Now()
	app.Id_admin = id_admin
	res := r.db.Save(&app)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// изменение статуса клиента
func (r *Repository) Update_application_client(id uint, status string) error {
	app := ds.Application{}
	res1 := r.db.First(&app, "id = ? and status = 'created'", id)
	if res1.Error != nil {
		return res1.Error
	}
	app.Status = status
	app.Date_formation = time.Now()
	res := r.db.Save(&app)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// присваивает статус удалено
func (r *Repository) Delete_application(id int) error {
	var app ds.Application
	res := r.db.First(&app, "id =? AND status = 'created'", id)
	if res.Error != nil {
		return res.Error
	}
	app.Status = "delete"
	result := r.db.Save(app)
	return result.Error
}
