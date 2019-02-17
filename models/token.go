package models

import (
  "github.com/jinzhu/gorm"
)

type Token struct {
  ID     string
  UserID string
	Token   string
  gorm.Model
}
