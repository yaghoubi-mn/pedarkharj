package datatypes

import "time"

type Table interface {
	Table() string
}

type CacheRepository interface {
	Save(key string, value interface{}, expireTime time.Duration) error
	Get(key string) (interface{}, error)
}

type Validator interface {
	Struct(st interface{}) map[string]string
}
