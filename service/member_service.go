package service

import (
	"camp-course-selection/common/constants"
	"camp-course-selection/common/util"
	"camp-course-selection/model"
	"camp-course-selection/vo"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type MemberService struct {
}

// CreateMember 用户注册
func (m *MemberService) CreateMember(memberVo *vo.CreateMemberRequest, c *gin.Context) (res vo.CreateMemberResponse) {
	//获取session中的用户
	user, _ := c.Get("user")
	if user == nil {
		res.Code = vo.LoginRequired
		return
	}
	u, _ := user.(*model.TMember)
	//检查权限
	if vo.UserType(u.UserType) != vo.Admin {
		res.Code = vo.PermDenied
		return
	}
	// 表单验证
	if code := CreateMemberValid(memberVo); code != vo.OK {
		res.Code = code
		return
	}

	// 雪花ID
	node, err := snowflake.NewNode(1)
	if err != nil {
		util.Log().Error("generate SnowFlakeID Error: %v", err)
		res.Code = vo.UnknownError
		return
	}
	id := node.Generate()
	member := model.TMember{
		Nickname: memberVo.Nickname,
		UserName: memberVo.Username,
		UserType: int(memberVo.UserType),
		Status:   constants.Active,
	}
	member.UserID = int64(id)
	// 加密密码
	bytes, _ := bcrypt.GenerateFromPassword([]byte(memberVo.Password), bcrypt.DefaultCost)
	member.Password = string(bytes)

	// 创建用户
	if err := model.DB.Create(&member).Error; err != nil {
		util.Log().Error("create member Error: %v", err)
		res.Code = vo.UnknownError
		return
	}
	res.Code = vo.OK
	res.Data.UserID = strconv.FormatInt(member.UserID, 10)
	return
}

// CreateMemberValid 用户注册验证表单
func CreateMemberValid(memberVo *vo.CreateMemberRequest) (code vo.ErrNo) {
	count := int64(0)
	model.DB.Model(&model.TMember{}).Where("user_name = ?", memberVo.Username).Count(&count)
	if count > 0 {
		return vo.UserHasExisted
	}
	//检查参数是否正确
	if nick_size := len(memberVo.Nickname); nick_size > 20 || nick_size < 4 {
		return vo.ParamInvalid
	}
	if name_size := len(memberVo.Username); name_size > 20 || name_size < 8 {
		return vo.ParamInvalid
	}
	if pass_size := len(memberVo.Password); pass_size > 20 || pass_size < 8 {
		return vo.ParamInvalid
	}

	pw := memberVo.Password
	var CapitalLetter, LowercaseLetter, Number bool
	for i := 0; i < len(pw); i++ {
		switch {
		case 64 < pw[i] && pw[i] < 91:
			CapitalLetter = true
		case 96 < pw[i] && pw[i] < 123:
			LowercaseLetter = true
		case 47 < pw[i] && pw[i] < 58:
			Number = true
		default:
			return vo.ParamInvalid
		}
	}

	if CapitalLetter && LowercaseLetter && Number {
		return vo.OK
	} else {
		return vo.ParamInvalid
	}
}

// GetMember 获取用户信息
func (m *MemberService) GetMember(memberVo *vo.GetMemberRequest) (res vo.GetMemberResponse) {
	var member model.TMember
	sid, _ := strconv.ParseInt(memberVo.UserID, 10, 64)
	if err := model.DB.First(&member, sid).Error; err != nil {
		res.Code = vo.UserNotExisted
		return
	}
	if member.Status == constants.InActive {
		res.Code = vo.UserHasDeleted
		return
	}
	res.Code = vo.OK
	res.Data.UserID = strconv.FormatInt(member.UserID, 10)
	res.Data.Username = member.UserName
	res.Data.Nickname = member.Nickname
	res.Data.UserType = vo.UserType(member.UserType)
	return
}

// GetMemberList 批量获取用户
func (m *MemberService) GetMemberList(memberVo *vo.GetMemberListRequest) (res vo.GetMemberListResponse) {
	memberList := make([]model.TMember, 0)
	if err := model.DB.Limit(memberVo.Limit).Offset(memberVo.Offset).Find(&memberList).Error; err != nil {
		res.Code = vo.ParamInvalid
		res.Data.MemberList = nil
		return
	}
	resMemberList := make([]vo.TMember, len(memberList))
	for i := 0; i < len(memberList); i++ {
		resMemberList[i].Username = memberList[i].UserName
		resMemberList[i].Nickname = memberList[i].Nickname
		resMemberList[i].UserType = vo.UserType(memberList[i].UserType)
		resMemberList[i].UserID = strconv.FormatInt(memberList[i].UserID, 10)
	}
	res.Data.MemberList = resMemberList
	res.Code = vo.OK
	return
}

// UpdateMember 更新用户信息
func (m *MemberService) UpdateMember(memberVo *vo.UpdateMemberRequest) (res vo.UpdateMemberResponse) {
	var member model.TMember
	sid, _ := strconv.ParseInt(memberVo.UserID, 10, 64)
	if err := model.DB.First(&member, sid).Error; err != nil {
		res.Code = vo.UserNotExisted
		return
	}

	if member.Status == constants.InActive {
		res.Code = vo.UserHasDeleted
		return
	}

	if err := model.DB.Model(&member).Update("nickname", memberVo.Nickname).Error; err != nil {
		res.Code = vo.UnknownError
		return
	}
	res.Code = vo.OK
	return
}

// DeleteMember 软删除
func (m *MemberService) DeleteMember(memberVo *vo.DeleteMemberRequest) (res vo.DeleteMemberResponse) {
	var member model.TMember
	sid, _ := strconv.ParseInt(memberVo.UserID, 10, 64)
	if err := model.DB.First(&member, sid).Error; err != nil {
		res.Code = vo.UserNotExisted
		return
	}

	if member.Status == constants.InActive {
		res.Code = vo.UserHasDeleted
		return
	}

	if err := model.DB.Model(&member).Update("status", constants.InActive).Error; err != nil {
		res.Code = vo.UnknownError
		return
	}
	res.Code = vo.OK
	return
}
