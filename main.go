package main

import (
	"camp-course-selection/cache"
	"camp-course-selection/common/util"
	"camp-course-selection/conf"
	"camp-course-selection/model"
	"camp-course-selection/server"
	"camp-course-selection/vo"
	"encoding/json"
	"github.com/go-redis/redis"
	"strconv"
)

func main() {
	// 从配置文件读取配置
	conf.Init()
	go listen()
	// 装载路由
	r := server.NewRouter()
	r.Run(":3000")
}

//消费数据
func listen() {
	for {
		val, err := cache.RedisClient.XRead(&redis.XReadArgs{
			Streams: []string{"BookCourseStream", "0"},
			Count:   1,
			Block:   0,
		}).Result()
		if err != nil {
			util.Log().Error("Listen XRead Error : %v\n", err)
			continue
		}
		if len(val) == 0 {
			continue
		}
		if len(val[0].Messages) == 0 {
			continue
		}
		bookCourseJson := val[0].Messages[0].Values["StudentCourseObj"]
		var bookCourseVo vo.BookCourseRequest
		json.Unmarshal([]byte(bookCourseJson.(string)), &bookCourseVo)
		sid, _ := strconv.ParseInt(bookCourseVo.StudentID, 10, 64)
		cid, _ := strconv.ParseInt(bookCourseVo.CourseID, 10, 64)
		sc := model.StudentCourse{
			StudentID: sid,
			CourseID:  cid,
		}
		//先查如果有，就不插入
		count := int64(0)
		model.DB.Model(&model.StudentCourse{}).Where("STUDENT_ID = ? AND COURSE_ID = ?", sid, cid).Count(&count)
		if count > 0 {
			//删除消息
			cache.RedisClient.XDel("BookCourseStream", val[0].Messages[0].ID)
			continue
		}
		//双主键约束保证幂等性
		if err = model.DB.Create(&sc).Error; err != nil {
			util.Log().Error("Create StudentCourse Error : %v\n ", err)
		} else {
			//课程容量减1
			model.DB.Exec("UPDATE t_course SET COURSE_STOCK = COURSE_STOCK - 1 where COURSE_ID = ?", cid)
			//删除课表缓存
			cache.RedisClient.HDel("GetStudentCourse", bookCourseVo.StudentID)
			//删除消息
			cache.RedisClient.XDel("BookCourseStream", val[0].Messages[0].ID)
		}

	}
}
