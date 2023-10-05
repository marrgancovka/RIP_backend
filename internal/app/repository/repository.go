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
func (r *Repository) GetAllShip(search string) (*[]ds.Ship, error) {
	var ships []ds.Ship
	if search != "" {
		res := r.db.Where("Is_delete = ?", "False").Where("Title ILIKE ?", "%"+search+"%").Find(&ships)
		fmt.Println(&ships)
		return &ships, res.Error
	}
	res := r.db.Where("Is_delete = ?", "False").Find(&ships)
	fmt.Println(&ships)
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

// удалить услугу по айди
func (r *Repository) DeleteShip(id int) error {
	err := r.db.Exec("UPDATE ships SET is_delete=true WHERE id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}
