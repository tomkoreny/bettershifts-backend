package models

import (
  //"github.com/jinzhu/gorm"
  "time"
)

type Token struct {
  ID     string
  UserID string
	Token   string
  CreatedAt time.Time
  UpdatedAt time.Time
 // gorm.Model
}
