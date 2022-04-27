package common

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

func (g *Game) End(db *gorm.DB, reason Reason) error {
	return db.Model(g).
		Where(&Game{
			ID: g.ID,
		}).
		Updates(&Game{
			EndTime: sql.NullTime{Time: time.Now(), Valid: true},
			Reason:  reason,
		}).Error
}
