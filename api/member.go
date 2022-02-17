package api

import (
	"camp-course-selection/service"
	"camp-course-selection/vo"
	"github.com/gin-gonic/gin"
	"strconv"
)

var memberSerivce service.MemberService

// CreateMember 用户注册接口
func CreateMember(c *gin.Context) {
	res := vo.CreateMemberResponse{}
	var memberVo vo.CreateMemberRequest
	if err := c.ShouldBind(&memberVo); err == nil {
		res = memberSerivce.CreateMember(&memberVo, c)
		c.JSON(200, res)
	} else {
		res.Code = vo.UnknownError
		c.JSON(200, res)
	}
}

//GetMember 获取用户信息接口
func GetMember(c *gin.Context) {
	res := vo.GetMemberResponse{}
	var memberVo vo.GetMemberRequest
	if err := c.ShouldBind(&memberVo); err == nil {
		res = memberSerivce.GetMember(&memberVo)
		c.JSON(200, res)
	} else {
		res.Code = vo.UnknownError
		c.JSON(200, res)
	}
}

//GetMemberList 批量获取成员接口
func GetMemberList(c *gin.Context) {
	var memberVo vo.GetMemberListRequest
	var err error

	if memberVo.Offset, err = strconv.Atoi(c.Query("Offset")); err != nil {
		c.JSON(200, vo.GetMemberListResponse{
			Code: vo.UnknownError,
		})
		return
	}

	if memberVo.Limit, err = strconv.Atoi(c.Query("Limit")); err != nil {
		c.JSON(200, vo.GetMemberListResponse{
			Code: vo.UnknownError,
		})
		return
	}

	res := memberSerivce.GetMemberList(&memberVo)
	c.JSON(200, res)
}

//UpdateMember 更新成员数据(只允许更新昵称)接口
func UpdateMember(c *gin.Context) {
	var memberVo vo.UpdateMemberRequest
	if err := c.ShouldBind(&memberVo); err == nil {
		res := memberSerivce.UpdateMember(&memberVo)
		c.JSON(200, res)
	} else {
		c.JSON(200, vo.UpdateMemberResponse{
			Code: vo.UnknownError,
		})
	}
}

//DeleteMember 删除成员（软删除）接口
func DeleteMember(c *gin.Context) {
	var memberVo vo.DeleteMemberRequest
	if err := c.ShouldBind(&memberVo); err == nil {
		res := memberSerivce.DeleteMember(&memberVo)
		c.JSON(200, res)
	} else {
		c.JSON(200, vo.DeleteMemberResponse{
			Code: vo.UnknownError,
		})
	}
}
