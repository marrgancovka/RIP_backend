package repository

//отвечает за обращение к хранилищам данных
import (
	"awesomeProject/internal/app/ds"

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
	err := r.db.First(ship, "id=?", id).Error
	if err != nil {
		return nil, err
	}
	return ship, nil
}
