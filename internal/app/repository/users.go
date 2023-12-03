package repository

import (
	"awesomeProject/internal/app/ds"
)

func (r *Repository) Register(user *ds.Users) error {
	return r.db.Create(user).Error
}
