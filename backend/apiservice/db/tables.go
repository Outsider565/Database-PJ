package db

import (
	"time"
)

const (
	man   = 1
	woman = 2
	other = 0
)

var _ = man
var _ = woman
var _ = other
var allModels = []interface{}{
	(*Building)(nil),
	(*LocateIn)(nil),
	(*Classroom)(nil),
	(*TeachIn)(nil),
	(*Course)(nil),
	(*User)(nil),
	(*Study)(nil),
	(*Takes)(nil),
	(*Admin)(nil),
}

type Building struct {
	BuildingNo   string `pg:"building_no,pk"`
	BuildingName string `pg:"building_name"`
	Latitude     float64
	Longitude    float64
}
type LocateIn struct {
	ClassroomNo string `pg:"classroom_no,pk"`
	BuildingNo  string `pg:"building_no,fk"`
}
type Classroom struct {
	ClassroomNo string `pg:"classroom_no,pk"`
	Capacity    int
}
type TeachIn struct {
	CourseId    string `pg:",pk"`
	Year        int    `pg:",pk"`
	Semester    string `pg:",pk"`
	ClassroomNo string `pg:",pk"` //一堂课如果是实验或者考试的话就可能有多个教室
	Date        string `pg:"type:date,pk"`
	ClassIndex  int    `pg:",pk"`
}
type Course struct {
	CourseId    string `pg:",pk"`
	Year        int    `pg:",pk"`
	Semester    string `pg:",pk"`
	CourseName  string
	StudentNum  int
	TeacherName string
	Type        string
}

type User struct {
	UserId    string `pg:"type:uuid,default:gen_random_uuid(),pk,on_delete:CASCADE"`
	OpenId    string `pg:",unique"`
	AvatarUrl string
	NickName  string
	Gender    int8
	City      string
	Province  string
	Country   string
}

type Study struct {
	StudyId     int    `pg:"type:serial,pk"`
	UserId      string `pg:"type:uuid"`
	StartTime   time.Time `pg:",notnull"`
	TimeLen     float64 `pg:"default:0"` // The unit is second
	ClassroomNo string
}

type Takes struct {
	UserId   string `pg:"type:uuid,pk"`
	CourseId string `pg:",pk"`
	Year     int    `pg:",pk"`
	Semester string `pg:",pk"`
}

type Admin struct {
	AdminId      string `pg:",pk"`
	PasswordHash string `pg:",notnull"`
}
