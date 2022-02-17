package api

import (
	"camp-course-selection/service"
	"camp-course-selection/vo"
	"github.com/gin-gonic/gin"
)

var studentService service.StudentService

//抢课
func BookCourse(c *gin.Context) {
	var courseVo vo.BookCourseRequest
	var res vo.BookCourseResponse

	if err := c.ShouldBind(&courseVo); err == nil {
		res = studentService.BookCourse(&courseVo)
		c.JSON(200, res)
	} else {
		res.Code = vo.UnknownError
		c.JSON(200, res)
	}
}

//查询学生课表
func GetStudentCourse(c *gin.Context) {
	var courseVo vo.GetStudentCourseRequest
	res := vo.GetStudentCourseResponse{}
	if err := c.ShouldBind(&courseVo); err == nil {
		res = studentService.GetStudentCourse(&courseVo)
		c.JSON(200, res)
	} else {
		res.Code = vo.UnknownError
		c.JSON(200, res)
	}
}
