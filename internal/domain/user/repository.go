package domain_user

type UserDomainRepository interface {
	GetByID(id uint64) (User, error)
	GetByNumber(number string) (User, error)
	Create(user *User) error
	Update(user User) error
	UpdateColumns(user User) error
	Delete(id uint64) error
}
