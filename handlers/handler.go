package handlers

import (
	"gorm.io/gorm"

	"go_demo_api/config"
)

type ApiHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewApiHandler(db *gorm.DB, cfg *config.Config) *ApiHandler {
	return &ApiHandler{db: db, cfg: cfg}
}
