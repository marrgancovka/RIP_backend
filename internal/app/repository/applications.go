package repository

import (
	"awesomeProject/internal/app/ds"
	"time"
)

// вывод списка всех заявок без услуг включенных в них + фильтрация по статусу и дате формирования
func (r *Repository) Select_applications(status string, date time.Time) (*[]ds.Application, error) {
	var applications []ds.Application
	if status != "" && !date.IsZero() {
		res := r.db.Where("status = ? AND status != 'delete'", status).Where("date_creation = ?", date).Find(&applications)
		return &applications, res.Error
	}
	if status != "" && date.IsZero() {
		res := r.db.Where("status = ? AND status != 'delete'", status).Find(&applications)
		return &applications, res.Error
	}
	if status == "" && !date.IsZero() {
		res := r.db.Where("date_creation = ?", date).Find(&applications)
		return &applications, res.Error
	}

	res := r.db.Where("status <> ?", "delete").Find(&applications)
	return &applications, res.Error
}

// вывод одной заявки со списком её услуг
func (r *Repository) Select_application(id int) (*ds.Application, *[]ds.Flights, error) {
	var applications ds.Application
	var flights []ds.Flights

	//ищем такую заявку
	result := r.db.First(&applications, "id =?", id)
	if result.Error != nil {
		return nil, nil, result.Error
	}
	//ищем услуги в заявке
	res := r.db.Where("Id_Application = ?", id).Find(&flights)
	if res.Error != nil {
		return nil, nil, res.Error
	}

	return &applications, &flights, nil
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
	res := r.db.First(&app, "id =?", id)
	if res.Error != nil {
		return res.Error
	}
	app.Status = "delete"
	result := r.db.Save(app)
	return result.Error
}
