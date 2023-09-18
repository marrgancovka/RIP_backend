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

func (r *Repository) GetShipByID(id int) (*ds.Ship, error) {
	ship := &ds.Ship{}
	err := r.db.First(ship, "id=?", "1").Error
	if err != nil {
		return nil, err
	}
	return ship, nil
}

func (r *Repository) CreateShip(ship ds.Ship) error {
	return r.db.Create(ship).Error
}
