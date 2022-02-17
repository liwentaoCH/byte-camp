package model

type TRecord struct {
	ID        int   `gorm:"primary_key"`
	StudentID int64 `gorm:"type:bigint(20);not null;comment:学生ID"`
	CourseID  int64 `gorm:"type:bigint(20);not null;comment:课程ID"`
}

func (TRecord) TableName() string {
	return "t_record"
}
