package database

import (
	"gotify_server/auth/password"
	"gotify_server/config"
	"gotify_server/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //enable the mysql dialect
)

//New creates a new wrapper for the gorm database framework.
func New(conf *config.Configuration) (*GormDatabase, error) {
	db, err := gorm.Open(conf.Database.Dialect, conf.Database.Connection)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxOpenConns(10)
	if err := db.AutoMigrate(&model.User{}, new(model.Application), new(model.Client), new(model.Message)).Error; err != nil {
		return nil, err
	}

	userCount := 0
	db.Find(new(model.User)).Count(&userCount)
	if userCount == 0 {
		db.Create(&model.User{
			Name:  conf.DefaultUser.Name,
			Pass:  password.CreatePassword(conf.DefaultUser.Pass, conf.PassStrength),
			Admin: true,
		})
	}

	return &GormDatabase{DB: db}, nil
}

//GormDatabase is a wrapper for the gorm framework
type GormDatabase struct {
	DB *gorm.DB
}

//Close closes the gorm database connection.
func (d *GormDatabase) Close() {
	d.DB.Close()
}
