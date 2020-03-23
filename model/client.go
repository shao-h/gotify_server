package model

//Client holds information about a device which can receive notifications.
type Client struct {
	ID     uint   `gorm:"primary_key"`
	Token  string `gorm:"type:varchar(180);unique_index"`
	UserID uint   `gorm:"index"`
	Name   string `gorm:"type:text" form:"name"`
}
