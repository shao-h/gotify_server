package model

//Application is an app which can send messages.
type Application struct {
	ID          uint   `gorm:"primary_key"`
	Token       string `gorm:"type:varchar(180);unique_index"`
	UserID      uint   `gorm:"index"`
	Name        string `gorm:"type:text" form:"name" binding:"required"`
	Description string `gorm:"type:text" form:"description"`
	Image       string `gorm:"type:text"`
}
