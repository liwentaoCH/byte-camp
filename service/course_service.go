package service

import (
	"camp-course-selection/cache"
	"camp-course-selection/common/util"
	"camp-course-selection/model"
	"camp-course-selection/vo"
	"github.com/bwmarrin/snowflake"
	"strconv"
)

type CourService struct{}

// CreateCourse 创建课程
func (m *CourService) CreateCourse(courseVo *vo.CreateCourseRequest) (res vo.CreateCourseResponse) {
	course := model.TCourse{
		Name:        courseVo.Name,
		CourseStock: courseVo.Cap,
	}
	// 雪花ID
	node, err := snowflake.NewNode(1)
	if err != nil {
		util.Log().Error("generate SnowFlakeID Error: %v", err)
		res.Code = vo.UnknownError
		return
	}

	id := node.Generate()
	course.CourseID = int64(id)
	//course.TeacherID = -1 // -1代表该课程未被绑定

	// 创建课程
	if err := model.DB.Create(&course).Error; err != nil {
		res.Code = vo.UnknownError
	} else {
		res.Code = vo.OK
		res.Data.CourseID = strconv.FormatInt(course.CourseID, 10)
		//课程容量存入redis中 redis中的key格式 CourseCap:课程ID
		cache.RedisClient.Set("CourseCap:"+res.Data.CourseID, courseVo.Cap, 0)
	}
	return
}

// GetCourse 获取课程
func (m *CourService) GetCourse(v *vo.GetCourseRequest) (res vo.GetCourseResponse) {
	course := &model.TCourse{}
	if err := model.DB.Where("course_id = ?", v.CourseID).First(&course).Error; err != nil {
		res.Code = vo.CourseNotExisted
		return
	}

	res.Code = vo.OK
	res.Data = vo.TCourse{
		CourseID:  strconv.FormatInt(course.CourseID, 10),
		Name:      course.Name,
		TeacherID: strconv.FormatInt(course.TeacherID, 10),
	}
	return
}

// BindCourseService 绑定课程
func (m *CourService) BindCourseService(v *vo.BindCourseRequest) (res vo.BindCourseResponse) {
	teacher := &model.TMember{}
	course := &model.TCourse{}

	// 检测老师是否正确
	if err := model.DB.Where("user_id = ?", v.TeacherID).First(&teacher).Error; err != nil {
		res.Code = vo.UserNotExisted
		return
	}

	if vo.UserType(teacher.UserType) != vo.Teacher {
		res.Code = vo.NotTeacher
		return
	}

	// 检查课程是否正确
	if err := model.DB.Where("course_id = ?", v.CourseID).First(&course).Error; err != nil {
		res.Code = vo.CourseNotExisted
		return
	}

	if course.TeacherID != 0 {
		res.Code = vo.CourseHasBound
		return
	}

	// 绑定课程
	course.TeacherID = teacher.UserID
	if err := model.DB.Model(&course).Update("teacher_id", course.TeacherID).Error; err == nil {
		res.Code = vo.OK
	} else {
		res.Code = vo.UnknownError
	}
	return
}

// UnBindCourseService 解绑课程
func (m *CourService) UnBindCourseService(v *vo.UnbindCourseRequest) (res vo.UnbindCourseResponse) {
	teacher := &model.TMember{}
	course := &model.TCourse{}

	// 检测老师是否正确
	if err := model.DB.Where("user_id = ?", v.TeacherID).First(&teacher).Error; err != nil {
		res.Code = vo.UserNotExisted
		return
	}

	if vo.UserType(teacher.UserType) != vo.Teacher {
		res.Code = vo.NotTeacher
		return
	}

	// 检查课程是否正确
	if err := model.DB.Where("course_id = ?", v.CourseID).First(&course).Error; err != nil {
		res.Code = vo.CourseNotExisted
		return
	}

	if course.TeacherID == 0 {
		// 课程已经解绑
		res.Code = vo.OK
		return
	}

	// 解绑课程
	course.TeacherID = 0
	if err := model.DB.Model(&course).Update("teacher_id", course.TeacherID).Error; err != nil {
		res.Code = vo.UnknownError
	} else {
		res.Code = vo.OK
	}

	return
}

// GetTeacherCourseService 获取老师所有课程
func (m *CourService) GetTeacherCourseService(v *vo.GetTeacherCourseRequest) (res vo.GetTeacherCourseResponse) {
	teacher := &model.TMember{}

	// 检测老师是否正确
	if err := model.DB.Where("user_id = ?", v.TeacherID).First(&teacher).Error; err != nil {
		res.Code = vo.UserNotExisted
		return
	}
	if vo.UserType(teacher.UserType) != vo.Teacher {
		res.Code = vo.NotTeacher
		return
	}

	// 获取结果
	courses := make([]model.TCourse, 0)
	result := make([]*vo.TCourse, 0)

	tid, _ := strconv.ParseInt(v.TeacherID, 10, 64)
	if err := model.DB.Where("teacher_id = ?", tid).Find(&courses).Error; err != nil {
		res.Code = vo.UnknownError
		return
	}
	for _, v := range courses {
		course := vo.TCourse{
			strconv.FormatInt(v.CourseID, 10),
			v.Name,
			strconv.FormatInt(v.TeacherID, 10),
		}
		result = append(result, &course)
	}
	res.Code = vo.OK
	res.Data.CourseList = result
	return
}

// ScheduleCourse 排课求解器
func (m *CourService) ScheduleCourse(schedule vo.ScheduleCourseRequest) (res vo.ScheduleCourseResponse) {
	// result----key:课程, val：对应的老师
	result := make(map[string]string, 0)
	teachers := make([]string, 0)
	courses := make([]string, 0)

	// set----key:课程, val:是否被被遍历 find函数需要
	set := make(map[string]bool, 0)
	relation := &schedule.TeacherCourseRelationShip

	// 提取teacher与course集合
	for k, v := range *relation {
		teachers = append(teachers, k)
		for _, val := range v {
			if _, ok := set[val]; ok == false {
				set[val] = false
				courses = append(courses, val)
			}
		}
	}

	// 二分图匹配
	for _, v := range teachers {
		clear(&set)
		find(v, relation, &set, &result)
	}
	// result中key是课程， ans中key是老师
	ans := make(map[string]string, 0)
	for k, v := range result {
		ans[v] = k
	}
	res.Code = vo.OK
	res.Data = ans
	return
}

func find(v string, relation *map[string][]string, set *map[string]bool, result *map[string]string) bool {
	courses, _ := (*relation)[v]
	for _, course := range courses {
		if ok, _ := (*set)[course]; ok == false {
			(*set)[course] = true
			if val, ok := (*result)[course]; ok == false || find(val, relation, set, result) == true {
				(*result)[course] = v
				return true
			}
		}
	}
	return false
}

func clear(m *map[string]bool) {
	for k, _ := range *m {
		(*m)[k] = false
	}
}
