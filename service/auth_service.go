package service

import (
	"camp-course-selection/model"
	"camp-course-selection/vo"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

type AuthService struct{}

type Data struct {
	UserID string // int64 范围
}

// setSession 设置session
func setSession(c *gin.Context, member model.TMember) {
	s := sessions.Default(c)
	s.Clear()
	s.Set("user_id", member.UserID)
	s.Save()
}

// Login 用户登录函数
func (m *AuthService) Login(loginVo *vo.LoginRequest, c *gin.Context) (res vo.LoginResponse) {

	var member model.TMember

	if err := model.DB.Where("user_name = ?", loginVo.Username).First(&member).Error; err != nil {
		res.Code = vo.UserNotExisted
		return
	}

	if member.Status != 1 {
		res.Code = vo.UserHasDeleted
		return
	}

	if ok := member.CheckPassword(loginVo.Password); ok == false {
		res.Code = vo.WrongPassword
		return
	}

	setSession(c, member)

	res.Code = vo.OK
	res.Data.UserID = strconv.FormatInt(member.UserID, 10)
	return
}
