package repository

import (
	"awesomeProject/internal/app/ds"
	"fmt"
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

func (r *Repository) SelectApplication(id int) (*ds.Application, []ApplicationService, error) {
	var application ds.Application
	var flights []ds.Flights

	// Ищем заявку
	if err := r.db.First(&application, "id = ?", id).Error; err != nil {
		return nil, nil, fmt.Errorf("ошибка при поиске заявки: %v", err)
	}

	// Ищем связанные полеты
	if err := r.db.Where("id_application = ?", id).Order("id_ship").Find(&flights).Error; err != nil {
		return nil, nil, fmt.Errorf("ошибка при поиске полетов: %v", err)
	}

	response := make([]ApplicationService, len(flights))

	for i, fl := range flights {
		ship, err := r.getShipByID(int(fl.Id_Ship))
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка при получении корабля: %v", err)
		}

		cosmodromBegin, err := r.getCosmodromTitleByID(int(fl.Id_Cosmodrom_Begin))
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка при получении космодрома (начало): %v", err)
		}

		cosmodromEnd, err := r.getCosmodromTitleByID(int(fl.Id_cosmodrom_End))
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка при получении космодрома (конец): %v", err)
		}

		response[i] = ApplicationService{
			ID:             int(ship.ID),
			Title:          ship.Title,
			CosmodromBegin: cosmodromBegin,
			CosmodromEnd:   cosmodromEnd,
			Date:           fl.Date_Flight,
		}
	}

	return &application, response, nil
}

func (r *Repository) getShipByID(id int) (*ds.Ship, error) {
	var ship ds.Ship
	if err := r.db.Select("id, title").Where("id = ?", id).First(&ship).Error; err != nil {
		return nil, err
	}
	return &ship, nil
}

func (r *Repository) getCosmodromTitleByID(id int) (string, error) {
	var cosmodrom ds.Cosmodroms
	if err := r.db.Table("cosmodroms").Select("title").Where("id = ?", id).First(&cosmodrom).Error; err != nil {
		return "", err
	}
	return cosmodrom.Title, nil
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
