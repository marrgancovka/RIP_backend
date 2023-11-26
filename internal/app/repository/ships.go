package repository

import (
	"awesomeProject/internal/app/ds"
	"errors"
	"strings"
	"time"
)

const default_image_url = "http://localhost:9000/spacey/1.png"

// возвращает список космических кораблей
func (r *Repository) Select_ships(search string) (*[]ds.Ship, error) {
	var ships []ds.Ship
	if search != "" {
		res := r.db.Where("Is_delete = ?", "False").Where("LOWER(Title) LIKE ?", "%"+strings.ToLower(search)+"%").Find(&ships)
		return &ships, res.Error
	}
	res := r.db.Where("Is_delete = ?", "False").Find(&ships)
	return &ships, res.Error
}

// возвращает информацию о космическом корабле по айди
func (r *Repository) Select_ship(id int) (*ds.Ship, error) {
	ship := &ds.Ship{}
	res := r.db.First(&ship, "id = ?", id)
	return ship, res.Error
}

// добавление нового космического корабля
func (r *Repository) Insert_ship(ship *ds.Ship) error {
	ship.Image_url = default_image_url

	res := r.db.Create(&ship)
	return res.Error
}

// добавление космического корабля в заявку и создание заявки если ее не было
func (r *Repository) Insert_application(request *struct {
	Id_Ship uint
	Id_user uint
}) error {
	var app ds.Application
	r.db.Where("id_user = ?", request.Id_user).Where("status = ?", "created").First(&app)

	if app.ID == 0 {
		newApp := ds.Application{
			Id_user:       request.Id_user,
			Id_admin:      1,
			Status:        "created",
			Date_creation: time.Now(),
		}
		res := r.db.Create(&newApp)
		if res.Error != nil {
			return res.Error
		}
		app = newApp
	}
	flight := ds.Flights{
		Id_Ship:            request.Id_Ship,
		Id_Application:     app.ID,
		Id_Cosmodrom_Begin: 1,
		Id_cosmodrom_End:   1,
	}
	result := r.db.Create(&flight)
	return result.Error
}

// изменение информации о космическом корабле
func (r *Repository) Update_ship(updateShip *ds.Ship) error {
	var ship ds.Ship
	res := r.db.First(&ship, "id =?", updateShip.ID)
	if res.Error != nil {
		return res.Error
	}
	if updateShip.Is_delete != false {
		return errors.New("Нельзя менять статус услуги")
	}
	if updateShip.Title != "" {
		ship.Title = updateShip.Title
	}
	if updateShip.Type != "" {
		ship.Type = updateShip.Type
	}
	if updateShip.Description != "" {
		ship.Description = updateShip.Description
	}
	if updateShip.Image_url != "" {
		return errors.New("Поле для изображения должно быть пустое")
	}
	if updateShip.Rocket != "" {
		ship.Rocket = updateShip.Rocket
	}

	*updateShip = ship
	result := r.db.Save(updateShip)
	return result.Error
}

// логически удаляет космический корабль
func (r *Repository) Delete_ship(id int) error {
	var ship ds.Ship
	res := r.db.First(&ship, "id =?", id)
	if res.Error != nil {
		return res.Error
	}
	ship.Is_delete = true
	result := r.db.Save(ship)
	return result.Error
}

// фото
func (r *Repository) Update_image(id string, newUrl string) error {
	ship := ds.Ship{}
	res := r.db.First(&ship, id)
	if res.Error != nil {
		return res.Error
	}
	ship.Image_url = newUrl
	result := r.db.Save(ship)
	return result.Error
}
