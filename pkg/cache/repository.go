package cache

import (
	"time"

	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"gorm.io/gorm"
)

// database table
type CacheTable struct {
	ID     uint64
	Key    string
	Value  string
	Expire time.Time
}

type GormCacheRepository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) datatypes.CacheRepository {

	db.AutoMigrate(&CacheTable{})

	return GormCacheRepository{
		DB: db,
	}
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
	if err := g.DB.Create(&c).Error; err != nil {
		return err
	}

	return nil
}

func (g GormCacheRepository) Get(key string) (string, error) {
	var c CacheTable
	if err := g.DB.First(&c, CacheTable{Key: key}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", database_errors.ErrRecordNotFound
		}
		return "", err
	}

	// check expire
	if c.Expire.Sub(time.Now()).Seconds() < 0 {
		return "", database_errors.ErrExpired
	}

	return c.Value, nil
}

func (g GormCacheRepository) Delete(key string) error {
	if err := g.DB.Delete(&CacheTable{}, &CacheTable{Key: key}).Error; err != nil {
		return err
	}

	return nil

}
