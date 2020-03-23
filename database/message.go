package database

import "gotify_server/model"

//CreateMessage creates a message
func (d *GormDatabase) CreateMessage(msg *model.Message) error {
	return d.DB.Create(msg).Error
}
