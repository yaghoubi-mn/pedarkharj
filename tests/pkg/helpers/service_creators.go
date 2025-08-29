package helpers

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	app_device "github.com/yaghoubi-mn/pedarkharj/internal/application/device"
	app_user "github.com/yaghoubi-mn/pedarkharj/internal/application/user"
	domain_device "github.com/yaghoubi-mn/pedarkharj/internal/domain/device"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	repository "github.com/yaghoubi-mn/pedarkharj/internal/infrastructure/repository/gorm"
	"github.com/yaghoubi-mn/pedarkharj/pkg/cache"
	"github.com/yaghoubi-mn/pedarkharj/pkg/database"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/s3"
	"github.com/yaghoubi-mn/pedarkharj/pkg/validator"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error

	err = godotenv.Load("../../../.env")
	if err != nil {
		slog.Error("cannot load variables", "error", err)
		os.Exit(1)
	}

	// setup s3
	s3.Init()

	db, err = database.SetupGromForTest()
	if err != nil {
		slog.Error("database", "error", err)
		os.Exit(1)
	}

	err = database.MigrateTables(
		db,
		domain_user.User{},
		domain_device.Device{},
		cache.Cache{},
	)

	if err != nil {
		slog.Error("cannot migrate tables", "error", err)
		os.Exit(1)
	}

	db.Raw("delete from users")
	db.Raw("delete from devices")
	db.Raw("delete from caches")

}

func GetCacheRepository() datatypes.CacheRepository {
	return cache.New(db)
}

func GetUserDomainRepository() domain_user.UserDomainRepository {
	return repository.NewGormUserRepository(db)
}

func GetDeviceAppService() app_device.DeviceAppService {
	vld := validator.NewValidator()

	return app_device.NewDeviceAppService(
		repository.NewGormDeviceRepository(db),
		domain_device.NewDeviceService(vld),
	)
}

func GetUserAppService() app_user.UserAppService {
	vld := validator.NewValidator()

	return app_user.NewUserService(
		GetUserDomainRepository(),
		GetCacheRepository(),
		GetDeviceAppService(),
		domain_user.NewUserService(vld),
	)
}
