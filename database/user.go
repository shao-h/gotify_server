package database

import (
	"gotify_server/model"

	"github.com/jinzhu/gorm"
)

//GetUserByName returns user by the given name.
func (d *GormDatabase) GetUserByName(name string) (*model.User, error) {
	user := new(model.User)
	err := d.DB.Where("name = ?", name).Find(user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}

//GetUserByID returns user by the given id
func (d *GormDatabase) GetUserByID(id uint) (*model.User, error) {
	user := new(model.User)
	err := d.DB.First(user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, err
}
