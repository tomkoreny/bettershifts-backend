package models

import "github.com/jinzhu/gorm"

type User struct {
  ID     string `gorm:"primary_key"`
	FirstName string
	LastName string
	UserName string
	IsAdmin bool
	Password string
	Shifts []Shift
	Workplaces []*Workplace `gorm:"many2many:users_workplaces;"`
	Wage int
	Benefits []Benefit
  gorm.Model
}
