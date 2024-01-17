package repository

import (
	"awesomeProject/internal/app/ds"
	"time"
)

// вывод списка всех заявок без услуг включенных в них + фильтрация по статусу и дате формирования
func (r *Repository) Select_applications(status string, date time.Time, date_end time.Time) (*[]ds.Application, error) {
	var applications []ds.Application
	if status != "" {
		res := r.db.Where("status = ? AND status != 'delete'", status).Where("date_creation BETWEEN ? AND ?", date, date_end).Find(&applications)
		return &applications, res.Error
	}
	res := r.db.Where("status <> ?", "delete").Where("date_creation BETWEEN ? AND ?", date, date_end).Find(&applications)
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
func (r *Repository) Update_application_admin(id uint, status string) error {
	app := ds.Application{}
	res1 := r.db.First(&app, "id = ? AND status = 'formated'", id)
	if res1.Error != nil {
		return res1.Error
	}
	app.Status = status
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
