package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Pagination struct {
	Page     int
	PageSize int
}

type UriParse struct {
	ID []string `json:"id" uri:"id"`
}

type BaseModel struct {
	ID        string          `gorm:"primary_key;type:char(36);not null" json:"id"`
	CreatorID *string         `gorm:"type:char(36);" json:"creator_id,omitempty"`
	UpdaterID *string         `gorm:"type:char(36);" json:"updater_id,omitempty"`
	CreatedAt time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type BaseParam struct {
	Search   string `json:"search" form:"search"`
	Sort     string `json:"sort" form:"sort"`
	Filter   string `json:"filter" form:"filter"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
}

type Object struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// This functions are called before creating Base
func (u *BaseModel) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
