package repository

import "awesomeProject/internal/app/ds"

// выводит список космодромов
func (r *Repository) Select_cosmodroms() (*[]ds.Cosmodroms, error) {
	var cosmodroms []ds.Cosmodroms
	res := r.db.Find(&cosmodroms)
	return &cosmodroms, res.Error
}

// изменение даты полета
func (r *Repository) Update_flight_date() {

}

// изменение космодрома вылета
func (r *Repository) Update_cosmodrom_begin() {

}

// изменение космодрома прилета
func (r *Repository) Update_cosmodrom_end() {

}

// удаление полета из заявки
func (r *Repository) Delete_flight() {

}
