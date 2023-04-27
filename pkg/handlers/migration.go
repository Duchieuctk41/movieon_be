package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/praslar/cloud0/logger"
	"gorm.io/gorm"
	"movieon_be/pkg/model"
)

type MigrationHandler struct {
	db *gorm.DB
}

func NewMigrationHandler(db *gorm.DB) *MigrationHandler {
	return &MigrationHandler{db: db}
}

func (h *MigrationHandler) BaseMigrate(ctx *gin.Context, tx *gorm.DB) error {
	//log := logger.WithCtx(ctx, "BaseMigrate")
	//if err := tx.Exec(`
	//		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	//	`).Error; err != nil {
	//	log.Errorf(err.Error())
	//}

	models := []interface{}{
		&model.Movie{},
	}

	for _, m := range models {
		err := h.db.AutoMigrate(m)
		if err != nil {
			_ = ctx.Error(err)
		}
	}

	//if err := tx.Exec(`
	//	ALTER TABLE answer ADD CONSTRAINT uc_answer_key UNIQUE (question_id, content, type);
	//`).Error; err != nil {
	//	log.Warn(err)
	//}

	return nil
}

func (h *MigrationHandler) Migrate(ctx *gin.Context) {
	log := logger.WithCtx(ctx, "Migrate")
	// put your migrations at the end of the list
	migrate := gormigrate.New(h.db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			// add column view_count
			ID: "20230426191358",
			Migrate: func(tx *gorm.DB) error {
				if err := h.db.AutoMigrate(&model.Movie{}, &model.ViewMovie{}); err != nil {
					return err
				}
				return nil
			},
		},
		{
			// add column view_count
			ID: "20230422194716",
			Migrate: func(tx *gorm.DB) error {
				if err := h.db.AutoMigrate(&model.Rating{}, &model.User{}); err != nil {
					return err
				}
				return nil
			},
		},
		{
			// add column user_id_old in table rating
			// add column id_old in table users
			ID: "20230422192407",
			Migrate: func(tx *gorm.DB) error {
				if err := h.db.AutoMigrate(&model.Rating{}, &model.User{}); err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "20230404210428",
			Migrate: func(tx *gorm.DB) error {
				if err := h.db.AutoMigrate(&model.Rating{}); err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "20221128213038",
			Migrate: func(tx *gorm.DB) error {
				if err := h.db.AutoMigrate(&model.Movie{}); err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "20230404101636",
			Migrate: func(tx *gorm.DB) error {
				log.Info("Migrate 20230404101636 - create table mo_movie")
				if err := h.BaseMigrate(ctx, tx); err != nil {
					return err
				}
				return nil
			},
		},
	})
	err := migrate.Migrate()
	if err != nil {
		log.Errorf(err.Error())
	}
}
