package repository

//отвечает за обращение к хранилищам данных
import (
	"awesomeProject/internal/app/ds"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

// получить все услуги
func (r *Repository) GetAllShip() (*[]ds.Ship, error) {
	var ships []ds.Ship
	res := r.db.Where("Is_delete = ?", "False").Find(&ships)
	return &ships, res.Error
}

// получить услугу по айди
func (r *Repository) GetShipByID(id int) (*ds.Ship, error) {
	ship := &ds.Ship{}
	err := r.db.First(ship, "id=?", "1").Error
	if err != nil {
		return nil, err
	}
	return ship, nil
}

func (r *Repository) DeleteShip(id int) error {
	var ship ds.Ship
	if res := r.db.First(&ship, id); res.Error != nil {
		return res.Error
	}
	if ship.ID == 0 {
		return fmt.Errorf("Ship not found")
	}

	ship.Is_delete = true
	res := r.db.Save(&ship)
	return res.Error

}

func (r *Repository) CreateShip(ship ds.Ship) error {
	return r.db.Create(ship).Error
}
