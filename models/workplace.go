package models

import (
  "github.com/jinzhu/gorm"
)

type Workplace struct {
  ID     string `gorm:"primary_key"`
	Name   string
  gorm.Model
  users []*User  `gorm:"many2many:users_workplaces;"`
}
