package core

import (
	"github.com/nspiguelman/zeus/domain"
	"gorm.io/gorm"
)

func CreateKahootCore (db *gorm.DB) {
	domain.CreateKahootDomain(db)
}
