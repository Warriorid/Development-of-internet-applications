package repository

import "gorm.io/gorm"

type PitsPostgres struct {
    db *gorm.DB
}

func NewPitsPostgres(db *gorm.DB) *PitsPostgres {
    return &PitsPostgres{db: db}
}