package database

import (
	"gotify_server/model"

	"github.com/jinzhu/gorm"
)

//CreateApplication creates an application.
func (d *GormDatabase) CreateApplication(app *model.Application) error {
	return d.DB.Create(app).Error
}

//GetApplicationsByUser returns all applications from a user
func (d *GormDatabase) GetApplicationsByUser(uid uint) ([]*model.Application, error) {
	apps := []*model.Application{}
	err := d.DB.Where("user_id = ?", uid).Find(&apps).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return apps, err
}

//GetApplicationByID returns application by the given id.
func (d *GormDatabase) GetApplicationByID(id uint) (*model.Application, error) {
	app := new(model.Application)
	err := d.DB.First(app, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return app, err
}

//GetApplicationByToken returns application by the given token.
func (d *GormDatabase) GetApplicationByToken(token string) (*model.Application, error) {
	app := new(model.Application)
	err := d.DB.Where("token = ?", token).First(app).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return app, err
}

//UpdateApplicationField updates application for the given field
func (d *GormDatabase) UpdateApplicationField(app *model.Application, a map[string]interface{}) error {
	return d.DB.Model(app).Updates(a).Error
}
