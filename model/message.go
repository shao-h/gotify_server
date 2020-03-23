package model

import "time"

//Message holds information about a message
type Message struct {
	ID            uint `gorm:"primary_key"`
	ApplicationID uint
	Message       string `gorm:"type:text"`
	Title         string `gorm:"type:text"`
	Priority      int
	Extras        []byte
	Created       time.Time
}

//MessageExternal holds information about a message sent by application
type MessageExternal struct {
	ID            uint
	ApplicationID uint
	Message       string `form:"message" binding:"required"`
	Title         string `form:"title"`
	Priority      int    `form:"priority"`
	Extras        map[string]interface{}
	Created       time.Time
}
