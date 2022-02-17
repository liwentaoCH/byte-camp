package service

import (
	"camp-course-selection/cache"
	"camp-course-selection/common/constants"
	"camp-course-selection/common/util"
	"camp-course-selection/model"
	"camp-course-selection/vo"
	"encoding/json"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

type StudentService struct {
}

//课程容量缓存
var courseCapCache = make(map[string]bool)

//学生抢课缓存
var studentCourseCache = make(map[string]bool)

func (m *StudentService) BookCourse(v *vo.BookCourseRequest) (res vo.BookCourseResponse) {
	//学生-课程 两级缓存
	scStr := v.StudentID + "_" + v.CourseID
	//在内存中查学生-课程是否已经抢成功
	if _, ok := studentCourseCache[scStr]; ok {
		res.Code = vo.StudentHasCourse
		return
	}
	//在redis中查学生-课程是否已经抢成功
	a, _ := cache.RedisClient.Get("StuCourse:" + scStr).Result()
	if a != "" {
		studentCourseCache[scStr] = false
		res.Code = vo.StudentHasCourse
		return
	}
	//在内存中查课程是否有容量
	if _, ok := courseCapCache[v.CourseID]; ok {
		res.Code = vo.CourseNotAvailable
		return
	}
	//查询学生信息
	var member model.TMember
	sid, _ := strconv.ParseInt(v.StudentID, 10, 64)
	if err := model.DB.First(&member, sid).Error; err != nil {
		res.Code = vo.StudentNotExisted
		return
	}
	if vo.UserType(member.UserType) != vo.Student {
		res.Code = vo.StudentNotExisted
		return
	}
	if member.Status == constants.InActive {
		res.Code = vo.UserHasDeleted
		return
	}
	//查询绑课信息
	cid, _ := strconv.ParseInt(v.CourseID, 10, 64)
	count := int64(0)
	model.DB.Model(&model.StudentCourse{}).Where("STUDENT_ID = ? AND COURSE_ID = ?", sid, cid).Count(&count)
	if count > 0 {
		CacheStudentCourse(scStr)
		res.Code = vo.StudentHasCourse
		return
	}
	//分布式锁防止重复请求 redis中的key格式 BookCourseLock:StudentID_CourseID
	redisKey := "BookCourseLock:" + v.StudentID + "_" + v.CourseID
	b, err := cache.RedisClient.SetNX(redisKey, 1, time.Minute).Result()
	if err != nil {
		util.Log().Error("BookCourse SetNX Error : %v", err)
		res.Code = vo.UnknownError
		return
	}
	if b == false {
		//重复请求
		res.Code = vo.RepeatRequest
		return
	}
	// 此代码会造成线程不安全，引发超卖。 需要人工保证redis中课程容量缓存不会失效。
	//if _, err = cache.RedisClient.Get("CourseCap:" + v.CourseID).Result(); err != nil {
	//	//缓存失效了
	//	var tc model.TCourse
	//	model.DB.First(&tc, cid)
	//	cache.RedisClient.Set("CourseCap:"+v.CourseID, tc.CourseStock, 0)
	//}
	var num int64
	num, err = cache.RedisClient.Decr("CourseCap:" + v.CourseID).Result()
	if num < 0 {
		courseCapCache[v.CourseID] = false
		res.Code = vo.CourseNotAvailable
		return
	}
	//抢课成功，快速返回
	mp := make(map[string]interface{})
	json, _ := json.Marshal(*v)
	mp["StudentCourseObj"] = string(json)
	if _, err = cache.RedisClient.XAdd(&redis.XAddArgs{
		Stream:       "BookCourseStream",
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values:       mp,
	}).Result(); err != nil {
		util.Log().Error("BookCourse XAdd Error : %v \n", err)
		//redis通讯出错，记录到差错表中，等待批处理
		record := model.TRecord{
			StudentID: sid,
			CourseID:  cid,
		}
		if err = model.DB.Create(&record).Error; err != nil {
			util.Log().Error("Create Record Error : %v\n", err)
		}
		res.Code = vo.UnknownError
		return
	}
	CacheStudentCourse(scStr)
	res.Code = vo.OK
	return
}

func (m *StudentService) GetStudentCourse(v *vo.GetStudentCourseRequest) (res vo.GetStudentCourseResponse) {
	//先查redis redis中的key格式 GetStudentCourse - StudentID
	str, err := cache.RedisClient.HGet("GetStudentCourse", v.StudentID).Result()
	if str != "" {
		res.Code = vo.OK
		var slice []vo.TCourse
		json.Unmarshal([]byte(str), &slice)
		res.Data.CourseList = slice
		return
	}

	//再查数据库
	var member model.TMember
	sid, _ := strconv.ParseInt(v.StudentID, 10, 64)
	//判断学生是否存在
	if err = model.DB.First(&member, sid).Error; err != nil {
		res.Code = vo.StudentNotExisted
		return
	}
	//判断学生类型
	if vo.UserType(member.UserType) != vo.Student {
		res.Code = vo.StudentNotExisted
	}

	courses := make([]model.StudentCourse, 0)
	if err = model.DB.Where("STUDENT_ID = ?", sid).Find(&courses).Error; err != nil {
		util.Log().Error("GetStudentCourse Query StudentCourse Error : %v \n", err)
		res.Code = vo.UnknownError
		return
	}
	//判断学生是否有课程
	if len(courses) == 0 {
		res.Code = vo.StudentHasNoCourse
		return
	}
	var queryList = make([]model.TCourse, len(courses))
	var courseList = make([]vo.TCourse, len(courses))
	for i := 0; i < len(courses); i++ {
		if err = model.DB.Where("course_id = ?", courses[i].CourseID).Find(&queryList[i]).Error; err != nil {
			res.Code = vo.CourseNotExisted
			return
		}
	}
	for i := 0; i < len(courseList); i++ {
		courseList[i].CourseID = strconv.FormatInt(queryList[i].CourseID, 10)
		courseList[i].TeacherID = strconv.FormatInt(queryList[i].TeacherID, 10)
		courseList[i].Name = queryList[i].Name
	}
	//缓存到redis中
	data, _ := json.Marshal(courseList)
	cache.RedisClient.HSet("GetStudentCourse", v.StudentID, string(data))
	res.Code = vo.OK
	res.Data.CourseList = courseList
	return
}

// CacheStudentCourse 抢课成功信息写入缓存
func CacheStudentCourse(scStr string) {
	studentCourseCache[scStr] = false
	cache.RedisClient.Set("StuCourse:"+scStr, 1, 15*time.Minute)
}
