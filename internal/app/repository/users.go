package repository

import (
	"awesomeProject/internal/app/ds"
)

func (r *Repository) Register(user *ds.Users) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByLogin(login string) (*ds.Users, error) {
	user := &ds.Users{}

	res := r.db.Where("username = ?", login).First(user)
	return user, res.Error
}
