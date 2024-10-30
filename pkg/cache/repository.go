package cache

import (
	"log/slog"
	"time"

	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"gorm.io/gorm"
)

// database table
type CacheTable struct {
	ID     uint64
	Key    string `gorm:"unique"`
	Value  string
	Expire time.Time
}

type GormCacheRepository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) datatypes.CacheRepository {

	return GormCacheRepository{
		DB: db,
	}
}

func MigrateTables(db *gorm.DB) error {
	return db.AutoMigrate(&CacheTable{})
}

func (g GormCacheRepository) Save(key string, value string, expireTime time.Duration) error {
	var c CacheTable
	c.Key = key
	c.Value = value
	c.Expire = time.Now().Add(expireTime)

	// delete if already saved
	if err := g.Delete(key); err != nil {
		return err
	}

	// save to database
	if err := g.DB.Save(&c).Error; err != nil {
		return err
	}

	go g.DeleteExpiredRecords()

	return nil
}

func (g GormCacheRepository) DeleteExpiredRecords() {
	if err := g.DB.Where("expire < ?", time.Now()).Delete(&CacheTable{}).Error; err != nil {
		slog.Error("cannot delete expired records", "error", err)
	}
}

func (g GormCacheRepository) Get(key string) (string, time.Time, error) {
	var c CacheTable
	if err := g.DB.First(&c, CacheTable{Key: key}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", c.Expire, database_errors.ErrRecordNotFound
		}
		return "", c.Expire, err
	}

	// check expire
	if c.Expire.Sub(time.Now()).Seconds() < 0 {
		return "", c.Expire, database_errors.ErrExpired
	}

	return c.Value, c.Expire, nil
}

func (g GormCacheRepository) Delete(key string) error {
	if err := g.DB.Delete(&CacheTable{}, &CacheTable{Key: key}).Error; err != nil {
		return err
	}

	return nil

}
