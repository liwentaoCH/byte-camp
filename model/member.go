package model

import (
	"golang.org/x/crypto/bcrypt"
)

// Member 用户模型

type TMember struct {
	UserID   int64  `gorm:"primarykey;type:bigint(20);not null;comment:用户ID"`
	UserName string `gorm:"type:varchar(32);uniqueIndex;not null;comment:用户名"`
	Nickname string `gorm:"type:varchar(32);not null;comment:用户昵称"`
	Password string `gorm:"type:varchar(128);not null;comment:密码"` // 密码存储的是加密后的
	UserType int    `gorm:"type:tinyint(4);not null;comment:用户类型 1-管理员，2-学生，3-教师"`
	Status   int    `gorm:"type:tinyint(4);not null;comment:用户状态 0-删除，1-正常"`
}

func (TMember) TableName() string {
	return "t_member"
}

// GetUser 用ID获取用户
func GetUser(ID interface{}) (TMember, error) {
	var user TMember
	result := DB.First(&user, ID)
	return user, result.Error
}

// CheckPassword 校验密码
func (member *TMember) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(password))
	return err == nil
	//return true
}
