package server

import (
	"camp-course-selection/api"
	"camp-course-selection/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 中间件, 顺序不能改
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	r.Use(middleware.CurrentUser())

	g := r.Group("/api/v1")

	g.POST("/auth/login", api.Login) // 登录后设置了session   user_id   与CurrentUser 中的user_id相呼应

	// 成员管理
	g.GET("/member", api.GetMember)

	g.GET("/member/list", api.GetMemberList)
	g.POST("/member/update", api.UpdateMember)
	g.POST("/member/delete", api.DeleteMember)
	// 排课
	g.POST("/course/create", api.CreateCourse)
	g.GET("/course/get", api.GetCourse)

	g.POST("/teacher/bind_course", api.BindCourse)
	g.POST("/teacher/unbind_course", api.UnBindCourse)
	g.GET("/teacher/get_course", api.GetTeacherCourse)
	g.POST("/course/schedule", api.ScheduleCourse)

	// 抢课
	g.POST("/student/book_course", api.BookCourse)
	g.GET("/student/course", api.GetStudentCourse)
	auth := g.Group("")
	auth.Use(middleware.AuthRequired())
	{
		// 登录
		auth.POST("/auth/logout", api.Logout)
		auth.GET("/auth/whoami", api.Whoami)

		auth.POST("/member/create", api.CreateMember)
	}
	return r
}
