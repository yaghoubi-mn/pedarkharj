package cache

import (
	"log/slog"
	"time"

	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
	"gorm.io/gorm"
)

// database table
type Cache struct {
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
	return db.AutoMigrate(&Cache{})
}

func (g GormCacheRepository) Save(key string, value map[string]string, expireTime time.Duration) error {
	var c Cache
	c.Key = key
	c.Expire = time.Now().Add(expireTime)

	var err error
	c.Value, err = utils.ConvertMapToString(value)
	if err != nil {
		return err
	}

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
	if err := g.DB.Where("expire < ?", time.Now()).Delete(&Cache{}).Error; err != nil {
		slog.Error("cannot delete expired records", "error", err)
	}
}

func (g GormCacheRepository) Get(key string) (map[string]string, time.Time, error) {
	var c Cache
	if err := g.DB.First(&c, Cache{Key: key}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, c.Expire, database_errors.ErrRecordNotFound
		}
		return nil, c.Expire, err
	}

	// check expire
	if c.Expire.Sub(time.Now()).Seconds() < 0 {
		return nil, c.Expire, database_errors.ErrExpired
	}

	value, err := utils.ConvertStringToMap(c.Value)

	return value, c.Expire, err
}

func (g GormCacheRepository) Delete(key string) error {
	if err := g.DB.Delete(&Cache{}, &Cache{Key: key}).Error; err != nil {
		return err
	}

	return nil

}
