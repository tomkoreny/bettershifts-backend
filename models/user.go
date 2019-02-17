package models

import "github.com/jinzhu/gorm"

type User struct {
  ID     string
	FirstName string
	LastName string
	UserName string
	IsAdmin bool
	Password string
	Shifts []Shift
	Workplaces []Workplace
	Wage int
	Benefits []Benefit
  gorm.Model
}
