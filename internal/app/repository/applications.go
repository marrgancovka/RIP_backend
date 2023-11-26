package repository

import (
	"awesomeProject/internal/app/ds"
	"fmt"
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

// вывод одной заявки со списком её услуг
func (r *Repository) Select_application(id int) (*ds.Application, *[]struct {
	Title           string
	Cosmodrom_begin string
	Cosmodrom_end   string
	Date            time.Time
}, error) {
	var applications ds.Application
	var flights []ds.Flights
	var ship ds.Ship
	var cosmodrom ds.Cosmodroms

	//ищем такую заявку
	result := r.db.First(&applications, "id =?", id)
	if result.Error != nil {
		return nil, nil, result.Error
	}
	fmt.Println("1111", applications)
	//ищем м-м заявки
	res := r.db.Where("Id_Application = ?", id).Find(&flights)
	if res.Error != nil {
		return nil, nil, res.Error
	}
	fmt.Println("22222", applications)

	var response []struct {
		Title           string
		Cosmodrom_begin string
		Cosmodrom_end   string
		Date            time.Time
	}
	for i, fl := range flights {
		var entry struct {
			Title           string
			Cosmodrom_begin string
			Cosmodrom_end   string
			Date            time.Time
		}
		response = append(response, entry)
		r.db.Table("ships").Select("title").Where("id = ?", fl.Id_Ship).First(&ship)
		response[i].Title = ship.Title
		r.db.Table("cosmodroms").Select("title").Where("id = ?", fl.Id_Cosmodrom_Begin).First(&cosmodrom)
		response[i].Cosmodrom_begin = cosmodrom.Title
		r.db.Table("cosmodroms").Select("title").Where("id = ?", fl.Id_cosmodrom_End).First(&cosmodrom)
		response[i].Cosmodrom_end = cosmodrom.Title
		response[i].Date = fl.Date_Flight
	}
	fmt.Println("33333", applications)
	fmt.Println(response)

	return &applications, &response, nil
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
