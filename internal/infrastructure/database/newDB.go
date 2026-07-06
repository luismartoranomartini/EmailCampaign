package database

import (
	"projeto-golang/internal/domain/campaign"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	dsn := "user:luis1407@tcp(127.0.0.1:3306)/email_campaign?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("fail to connect to database")
	}
	// Para criar referência de chave estrangeira
	db.AutoMigrate(&campaign.Campaign{}, &campaign.Contact{})
	return db
}
