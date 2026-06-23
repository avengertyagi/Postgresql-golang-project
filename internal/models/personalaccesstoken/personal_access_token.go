package personalaccesstoken

import (
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/models"
)

func init() {
	models.Register(&PersonalAccessToken{})
}

type PersonalAccessToken struct {
	ID        uint      `json:"id"         gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"user_id"    gorm:"not null;index"`
	TokenHash string    `json:"-"          gorm:"type:varchar(255);not null;uniqueIndex"`
	Name      string    `json:"name"       gorm:"type:varchar(100);default:null"`
	Revoked   bool      `json:"revoked"    gorm:"default:false"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
