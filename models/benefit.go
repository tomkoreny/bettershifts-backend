package models

import (
	"time"
  "github.com/jinzhu/gorm"
)

type Benefit struct {
  ID string
	Date time.Time
	Reason string
	Amount int
	UserID string
  gorm.Model
}
