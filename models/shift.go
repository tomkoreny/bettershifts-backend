package models

import (
	"time"
  "github.com/jinzhu/gorm"
)

type Shift struct {
  ID string
	Start time.Time
	End *time.Time
	UserID string
	WorkplaceID string
  gorm.Model
}

