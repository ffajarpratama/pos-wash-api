package model

import (
	"fmt"
	"time"

	"github.com/ffajarpratama/pos-wash-api/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Media struct {
	MediaID   uuid.UUID      `json:"media_id" gorm:"primaryKey; default:gen_random_uuid()"`
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Size      int            `json:"size"`
	Mimetype  string         `json:"mimetype"`
	MediaURL  string         `json:"media_url"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`

	// json field
	ShortMediaURL string `json:"short_media_url" gorm:"-"`
}

func (Media) TableName() string {
	return "tr_media"
}

func (m *Media) AfterFind(db *gorm.DB) (err error) {
	m.ShortMediaURL = fmt.Sprintf("%s/bo_40px_solid_brown/%s", config.CLOUDINARY_BASE_URL, m.Path)
	return
}
