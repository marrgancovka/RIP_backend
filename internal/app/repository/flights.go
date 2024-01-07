package repository

import (
	"awesomeProject/internal/app/ds"
	"time"
)

// выводит список космодромов
func (r *Repository) Select_cosmodroms() (*[]ds.Cosmodroms, error) {
	var cosmodroms []ds.Cosmodroms
	res := r.db.Find(&cosmodroms)
	return &cosmodroms, res.Error
}

// изменение даты полета
func (r *Repository) Update_flight_date(app uint, ship uint, date time.Time, cosm_begin uint, cosm uint) (*ds.Flights, error) {
	var flight ds.Flights
	res := r.db.First(&flight, "id_ship = ? AND id_application = ?", ship, app)
	if res.Error != nil {
		return nil, res.Error
	}
	flight.Date_Flight = date
	flight.Id_Cosmodrom_Begin = cosm_begin
	flight.Id_cosmodrom_End = cosm

	result := r.db.Where("id_ship = ? AND id_application = ?", ship, app).Save(&flight)
	return &flight, result.Error
}

// // изменение космодрома вылета
// func (r *Repository) Update_cosmodrom_begin(app uint, ship uint, cosm_begin uint) (*ds.Flights, error) {
// 	var flight ds.Flights
// 	res := r.db.Where("id_ship = ? AND id_application = ?", ship, app).First(&flight)
// 	if res.Error != nil {
// 		return nil, res.Error
// 	}
// 	flight.Id_Cosmodrom_Begin = cosm_begin
// 	result := r.db.Where("id_ship = ? AND id_application = ?", ship, app).Save(&flight)
// 	return &flight, result.Error
// }

// // изменение космодрома прилета
// func (r *Repository) Update_cosmodrom_end(app uint, ship uint, cosm uint) (*ds.Flights, error) {
// 	var flight ds.Flights
// 	res := r.db.Where("id_ship = ? AND id_application = ?", ship, app).First(&flight)
// 	if res.Error != nil {
// 		return nil, res.Error
// 	}
// 	flight.Id_cosmodrom_End = cosm
// 	result := r.db.Where("id_ship = ? AND id_application = ?", ship, app).Save(&flight)
// 	return &flight, result.Error
// }

// удаление полета из заявки
func (r *Repository) Delete_flight(app int, ship int) error {
	var flight ds.Flights
	res := r.db.Where("id_ship = ? AND id_application = ?", ship, app).Delete(&flight)

	return res.Error
}
