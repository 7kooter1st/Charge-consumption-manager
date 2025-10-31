package ds

import "sync"

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Login       string `gorm:"type:varchar(255);unique"`
	Password    string `gorm:"type:varchar(255)"`
	IsModerator bool
}

type singletonUser struct {
	id uint
}

var (
	instance *singletonUser
	once     sync.Once
)

func GetUser() *singletonUser {
	once.Do(func() {
		instance = &singletonUser{1}
	})
	return instance
}

func (s *singletonUser) GetId() uint {
	return s.id
}
