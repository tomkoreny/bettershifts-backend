package models

import (
	"time"
  "github.com/jinzhu/gorm"
)

type Todo struct {
  ID     string
	Name   string
	Done   bool
	UserID string
	WorkplaceID string
	Date time.Time
	Benefit int
  gorm.Model
}
