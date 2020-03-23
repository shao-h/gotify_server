package model

//User holds info about the credentials of a user and other stuff
type User struct {
	ID    uint   `gorm:"primary_key"`
	Name  string `gorm:"type:varchar(180);unique_index"`
	Pass  []byte
	Admin bool
}
