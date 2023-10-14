package repository

import "awesomeProject/internal/app/ds"

// возвращает список космических кораблей
func (r *Repository) Select_ships(search string) (*[]ds.Ship, error) {
	var ships []ds.Ship
	if search != "" {
		res := r.db.Where("Is_delete = ?", "False").Where("Title LIKE ?", "%"+search+"%").Find(&ships)
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
	res := r.db.Create(&ship)
	return res.Error
}

// создание новой заявки
func (r *Repository) Insert_application(application *ds.Application) error {
	res := r.db.Create(&application)
	return res.Error
}

// изменение информации о космическом корабле
func (r *Repository) Update_ship(updateShip *ds.Ship) error {
	var ship ds.Ship
	res := r.db.First(&ship, "id =?", updateShip.ID)
	if res.Error != nil {
		return res.Error
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
		ship.Image_url = updateShip.Image_url
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
