package models

import (
  "github.com/jinzhu/gorm"
)

type Workplace struct {
  ID     string
	Name   string
  gorm.Model
}
