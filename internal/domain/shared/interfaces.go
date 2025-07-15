package domain_shared

import "time"

type Validator interface {
	ValidateField(v any, tag string) error
	ValidateFieldByFieldName(fieldName string, fieldValue any, model any) error
}

type CacheRepository interface {
	Save(key string, value map[string]string, expireTime time.Duration) error
	Get(key string) (map[string]string, time.Time, error)
	Delete(key string) error
}
