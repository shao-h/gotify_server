package database

import (
	"gotify_server/model"

	"github.com/jinzhu/gorm"
)

//CreateClient creates a new client
func (d *GormDatabase) CreateClient(client *model.Client) error {
	return d.DB.Create(client).Error
}

//GetClientByToken returns client for given token
func (d *GormDatabase) GetClientByToken(token string) (*model.Client, error) {
	client := new(model.Client)
	err := d.DB.Where("token = ?", token).Find(client).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return client, err
}
