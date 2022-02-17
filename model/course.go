package model

// Course 课程模型

type TCourse struct {
	CourseID    int64  `gorm:"primarykey;type:bigint(20);not null;comment:课程ID"`
	Name        string `gorm:"type:varchar(32);not null;comment:课程名"`
	TeacherID   int64  `gorm:"type:bigint(20);index;comment:授课教师ID"`
	CourseStock int    `gorm:"type:int(11);not null;comment:课程容量"`
}

func (TCourse) TableName() string {
	return "t_course"
}
