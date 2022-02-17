package model

type StudentCourse struct {
	StudentID int64 `gorm:"primarykey;type:bigint(20);not null;comment:学生ID"`
	CourseID  int64 `gorm:"primarykey;type:bigint(20);not null;comment:课程ID"`
}

func (StudentCourse) TableName() string {
	return "t_student_course"
}
