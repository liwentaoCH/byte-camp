package api

import (
	"camp-course-selection/model"
	"camp-course-selection/service"
	"camp-course-selection/vo"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

var authService service.AuthService

// Login 用户登录接口
func Login(c *gin.Context) {
	var res = vo.LoginResponse{}
	var loginVo vo.LoginRequest
	if err := c.ShouldBind(&loginVo); err == nil {
		res = authService.Login(&loginVo, c)
		c.JSON(200, res)
	} else {
		res.Code = vo.UnknownError
		c.JSON(200, res)
	}
}

// Logout 用户登出
func Logout(c *gin.Context) {
	res := vo.LogoutResponse{}
	s := sessions.Default(c)
	s.Clear()
	s.Save()
	res.Code = vo.OK
	c.JSON(200, res)
}

// Whoami 获取当前用户
func Whoami(c *gin.Context) {
	res := vo.WhoAmIResponse{}
	if user, _ := c.Get("user"); user != nil {
		if u, ok := user.(*model.TMember); ok {
			res.Code = vo.OK
			res.Data.UserID = strconv.FormatInt(u.UserID, 10)
			res.Data.Username = u.UserName
			res.Data.Nickname = u.Nickname
			res.Data.UserType = vo.UserType(u.UserType)
			c.JSON(200, res)
			return
		}
	}
	res.Code = vo.LoginRequired
	c.JSON(200, res)
}
