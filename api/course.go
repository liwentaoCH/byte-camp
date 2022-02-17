package api

import (
	"camp-course-selection/service"
	"camp-course-selection/vo"
	"github.com/gin-gonic/gin"
)

var courService service.CourService

// CreateCourse 创建课程
func CreateCourse(c *gin.Context) {
	var courseVo vo.CreateCourseRequest
	if err := c.ShouldBind(&courseVo); err == nil {
		res := courService.CreateCourse(&courseVo)
		c.JSON(200, res)
	} else {
		c.JSON(200, vo.CreateCourseResponse{
			Code: vo.UnknownError,
		})
	}
}

// GetCourse 获取课程
func GetCourse(c *gin.Context) {
	var courseVo vo.GetCourseRequest
	//courseVo.CourseID = c.Query("CourseID")
	if err := c.ShouldBind(&courseVo); err == nil {
		res := courService.GetCourse(&courseVo)
		c.JSON(200, res)
	} else {
		c.JSON(200, vo.GetCourseResponse{
			Code: vo.UnknownError,
			Data: vo.TCourse{},
		})
	}
}

// BindCourse 绑定课程
func BindCourse(c *gin.Context) {
	var bindCourse vo.BindCourseRequest
	var bindRes vo.BindCourseResponse
	if err := c.ShouldBind(&bindCourse); err == nil {
		bindRes = courService.BindCourseService(&bindCourse)
		c.JSON(200, bindRes)
	} else {
		bindRes.Code = vo.UnknownError
		c.JSON(200, bindRes)
	}
}

// UnBindCourse 解绑课程服务
func UnBindCourse(c *gin.Context) {
	var unBindCourse vo.UnbindCourseRequest
	var unbindRes vo.UnbindCourseResponse
	if err := c.ShouldBind(&unBindCourse); err == nil {
		unbindRes = courService.UnBindCourseService(&unBindCourse)
		c.JSON(200, unbindRes)
	} else {
		unbindRes.Code = vo.UnknownError
		c.JSON(200, unbindRes)
	}
}

// GetTeacherCourse 获取该老师的所有课程（返回课程名称列表）
func GetTeacherCourse(c *gin.Context) {
	var teacherCourse vo.GetTeacherCourseRequest
	var TeacherCourseRes vo.GetTeacherCourseResponse
	if err := c.ShouldBind(&teacherCourse); err == nil {
		TeacherCourseRes := courService.GetTeacherCourseService(&teacherCourse)
		c.JSON(200, TeacherCourseRes)
	} else {
		TeacherCourseRes.Code = vo.UnknownError
		c.JSON(200, TeacherCourseRes)
	}
}

// ScheduleCourse 排课求解器
func ScheduleCourse(c *gin.Context) {
	var schedule vo.ScheduleCourseRequest
	var res vo.ScheduleCourseResponse
	if err := c.ShouldBind(&schedule); err == nil {
		res = courService.ScheduleCourse(schedule)
		c.JSON(200, res)
	} else {
		res.Code = vo.UnknownError
		c.JSON(200, res)
	}
}
