package datatypes

import "time"

type Table interface {
	Table() string
}

type CacheRepository interface {
	Save(key string, value string, expireTime time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type Validator interface {
	Struct(st interface{}) map[string]string
}
