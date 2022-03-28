package db

import (
	json2 "encoding/json"
	"fmt"
	"github.com/go-pg/pg/extra/pgdebug/v10"
	"github.com/go-pg/pg/v10"
	"github.com/go-redis/redis/v8"
	"github/Outsider565/gostudy/apiservice/crawl"
	"github/Outsider565/gostudy/apiservice/utils"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
	"time"
)

// Connection is the port of all database operation, which means underline db is private, and also rdb
type Connection struct {
	db  *pg.DB
	rdb *redis.Client
	ctx context.Context
}

func (con *Connection) FlushRdb() error {
	_, err := con.rdb.FlushDB(con.ctx).Result()
	return err
}

// CrawlClasses crawl from web and insert course at the time of given dateArr
// dateArr can be set at any time, only the date will be considered (other than hour or minute)
// only return the first error it encounters
func (con *Connection) CrawlClasses(dateArr []time.Time) error {
	con.rdb.FlushDB(con.ctx)
	log.Println("Flushing redis database")
	errChan := make(chan error)
	for _, date := range dateArr {
		date := date
		go func() {
			err := con.CrawlClassesDate(date)
			if err != nil {
				log.Println("Error in CrawlClasses", err)
			}
			errChan <- err
		}()
	}
	length := len(dateArr)
	for i := 0; i < length; i++ {
		err := <-errChan
		if err != nil {
			return err
		}
	}
	return nil
}
func (con *Connection) CrawlClassesDate(date time.Time) error {
	db := con.db
	dateString := date.Format(utils.StdDateFormat)
	for _, building := range buildings {
		buildingNo := building.BuildingNo
		classrooms := crawl.Crawl(buildingNo, dateString)
		var locateIns []LocateIn
		err := db.Model((*LocateIn)(nil)).Where("building_no=?", buildingNo).Select(&locateIns)
		if err != nil {
			return err
		}
		if len(locateIns) == 0 {
			for _, classroom := range classrooms {
				classroomIns := &Classroom{
					ClassroomNo: classroom.Name,
					Capacity:    classroom.SeatNum,
				}
				_, err = db.Model(classroomIns).OnConflict("DO NOTHING").Insert()
				if err != nil {
					return err
				}
				locateIns := &LocateIn{
					ClassroomNo: classroom.Name,
					BuildingNo:  buildingNo,
				}
				_, err := db.Model(locateIns).
					Insert()
				if err != nil {
					return err
				}
			}
		}
		for _, classroom := range classrooms {
			var lastClass = crawl.Class{}
			for idx, class := range classroom.Classes {
				if class.CourseId == "" {
					if !class.Empty() {
						fmt.Println(class)
					}
					continue
				}
				if class != lastClass {
					classIns := &Course{
						CourseId:    class.CourseId,
						Year:        date.Year(),
						Semester:    utils.GetSemester(date),
						CourseName:  class.CourseName,
						StudentNum:  class.StudentNum,
						TeacherName: class.Teacher,
						Type:        class.Type,
					}
					if strings.Contains(classIns.Type, "本") {
						classIns.Type = "本"
					} else if strings.Contains(classIns.Type, "硕") {
						classIns.Type = "硕"
					} else if strings.Contains(classIns.Type, "博") {
						classIns.Type = "博"
					} else if strings.Contains(classIns.Type, "二专") {
						classIns.Type = "二专"
					} else {
						classIns.Type = "其他"
					}
					_, err = db.Model(classIns).
						OnConflict("(course_id, year, semester) DO UPDATE").
						Set("course_name=EXCLUDED.course_name,student_num=EXCLUDED.student_num," +
								"teacher_name=EXCLUDED.teacher_name,type=EXCLUDED.type").
						Insert() //Our school might change course info so update every time on conflict
					if err != nil {
						fmt.Println(classIns)
						return err
					}
				}
				lastClass = class
				teachInIns := &TeachIn{
					CourseId:    class.CourseId,
					Year:        date.Year(),
					Semester:    utils.GetSemester(date),
					ClassroomNo: classroom.Name,
					Date:        dateString,
					ClassIndex:  idx + 1,
				}
				_, err = db.Model(teachInIns).
					OnConflict("DO NOTHING").
					Insert() //As every item is primary key, we can simply do nothing when conflict
				if err != nil {
					fmt.Println(teachInIns)
					return err
				}
			}
		}
	}
	return nil
}

func (con *Connection) EmptyClassroomsIdx(
	building string,
	dateString string,
	classIndex int,
) ([]string, error) {
	// 缓存策略如下：
	// 查询一次我们就缓存，直到下一次插入或者爬取前，这个数据都是有效的
	// 如果下一次插入或者爬取，我们就把清空数据库
	// 如果本次查询cache miss，则执行下面的语句，将缓存结果写入redis中
	db := con.db
	rdb := con.rdb
	ctx := con.ctx
	type key struct {
		Building   string `json:"building"`
		DateString string `json:"date_string"`
		ClassIndex int    `json:"class_index"`
	}
	type classroomVal struct {
		ClassroomNo []string `json:"classroom_no"`
	}
	json, err := json2.Marshal(key{Building: building, DateString: dateString, ClassIndex: classIndex})
	if err != nil {
		return nil, err
	}
	val, err := rdb.Get(ctx, string(json)).Result()
	if err == nil {
		var temp classroomVal
		err := json2.Unmarshal(([]byte)(val), &temp)
		if err != nil {
			log.Println(err)
		} else if len(temp.ClassroomNo) != 0 {
			if os.Getenv("DEBUG") == "TRUE" {
				log.Println("redis hit: ", temp.ClassroomNo)
			}
			return temp.ClassroomNo, nil
		}
	} else if err != redis.Nil {
		return nil, err
	}
	var classroomNo []string
	allClassRooms := db.Model().
		Column("classroom_no").
		Distinct().
		Table("classrooms").
		Join("NATURAL JOIN \"locate_ins\"").
		Where("building_no=?", building)
	usedClassRooms := db.Model((*TeachIn)(nil)).
		Column("classroom_no").
		Distinct().
		Join("NATURAL JOIN \"classrooms\" NATURAL JOIN \"locate_ins\"").
		Where("date=?", dateString).Where("class_index=?", classIndex).
		Where("building_no=?", building)
	err = db.Model().
		TableExpr("(?) AS res", allClassRooms.Except(usedClassRooms)).
		Order("classroom_no asc").
		Select(&classroomNo)
	jsonVal, err := json2.Marshal(classroomVal{ClassroomNo: classroomNo})
	if err != nil {
		log.Println(err)
	}
	err = rdb.Set(ctx, string(json), string(jsonVal), 0).Err()
	if err != nil {
		log.Println(err)
	}
	return classroomNo, err
}

func (con *Connection) ReInit() error {
	err := con.DropAllTables()
	if err != nil {
		return err
	}
	err = con.CreateSchema()
	if err != nil {
		return err
	}
	err = con.InitBuildingInfo()
	if err != nil {
		return err
	}
	err = con.InitAdmin()
	if err != nil {
		return err
	}
	return nil
}
func (con *Connection) SearchCourse(maxNum int, courseId, year, semester, courseName, teacherName string, ifStrict bool) ([]Course, error) {
	var courses []Course
	whereClause := "type IN ('本','硕','博','二专')"
	if !ifStrict {
		if courseId != "" {
			whereClause += fmt.Sprintf(" AND \"course\".\"course_id\" LIKE '%s%%'", courseId)
		}
		if year != "" {
			whereClause += fmt.Sprintf(" AND \"course\".\"year\"=%s", year)
		}
		if semester != "" {
			whereClause += fmt.Sprintf(" AND \"course\".\"semester\"='%s'", semester)
		}
		if teacherName != "" {
			whereClause += fmt.Sprintf(" AND \"course\".\"teacher_name\" LIKE '%s%%'", teacherName)
		}
	} else {
		if courseId != "" {
			whereClause += fmt.Sprintf(" AND \"course\".\"course_id\"='%s'", courseId)
		}
		if year != "" {
			whereClause += fmt.Sprintf(" AND \"course\".\"year\"=%s", year)
		}
		if courseName != "" {
			whereClause += fmt.Sprintf(" AND \"course\".\"course_name\"='%s'", courseName)
		}
		if semester != "" {
			whereClause += fmt.Sprintf(" AND \"course\".\"semester\"='%s'", semester)
		}
		if teacherName != "" {
			whereClause += fmt.Sprintf(" AND \"course\".\"teacher_name\"='%s'", teacherName)
		}
	}
	var err error
	if !ifStrict && courseName != "" {
		err = con.db.Model((*Course)(nil)).
			Where(whereClause).
			Where("similarity(course_name,?)>0.2", courseName).
			OrderExpr("year, similarity(course_name,?) DESC", courseName).
			Limit(maxNum).
			Select(&courses)
	} else {
		err = con.db.Model((*Course)(nil)).
			Where(whereClause).
			Limit(maxNum).
			Order("year").
			Select(&courses)
	}
	log.Println(courses)
	if err != nil {
		return nil, err
	} else {
		return courses, nil
	}
}

type Class struct {
	tableName   struct{} `pg:",discard_unknown_columns"`
	CourseId    string
	Year        int
	Semester    string
	CourseName  string
	StudentNum  int
	TeacherName string
	Type        string
	ClassIndex  int
}

func (con *Connection) GetClassesOfDate(date, classroomNo string) ([]Class, error) {
	var classes []Class
	err := con.db.Model().Table("courses").Join("NATURAL JOIN \"teach_ins\"").Where("date=?", date).Where("\"teach_ins\".\"classroom_no\"=?", classroomNo).Order("class_index").Select(&classes)
	return classes, err
}
func (con *Connection) AllBuildings() ([]Building, error) {
	type val struct {
		Buildings []Building `json:"buildings"`
	}
	var buildings []Building
	result, err := con.rdb.Get(con.ctx, "all_buildings").Result()
	if err != nil && err != redis.Nil {
		log.Println(err)
		return nil, err
	} else if err != redis.Nil {
		var temp val
		err = json2.Unmarshal([]byte(result), &temp)
		if err != nil {
			log.Println(err)
		} else {
			//Cache hit了，直接返回
			if os.Getenv("DEBUG") == "TRUE" {
				log.Println("redis hit: ", temp.Buildings)
			}
			return temp.Buildings, nil
		}
	}
	//如果Cache Miss了，那么走pg，并且把数据插进来
	err = con.db.Model((*Building)(nil)).Select(&buildings)
	json, err := json2.Marshal(val{Buildings: buildings})
	if err != nil {
		log.Println(err)
	}
	con.rdb.Set(con.ctx, "all_buildings", string(json), 0)
	return buildings, err
}

type ClassState struct {
	Indexes  []int
	Location string
}

func (con *Connection) GetCourseTimeAndLoc(courseId string) (map[string]*ClassState, error) {
	var teaches []TeachIn
	var temp = map[string][]bool{
		"Monday":    make([]bool, 15),
		"Tuesday":   make([]bool, 15),
		"Wednesday": make([]bool, 15),
		"Thursday":  make([]bool, 15),
		"Friday":    make([]bool, 15),
		"Saturday":  make([]bool, 15),
		"Sunday":    make([]bool, 15),
	}
	err := con.db.Model((*TeachIn)(nil)).Where("course_id=?", courseId).Select(&teaches)
	var res = map[string]*ClassState{
		"Monday":    {},
		"Tuesday":   {},
		"Wednesday": {},
		"Thursday":  {},
		"Friday":    {},
		"Saturday":  {},
		"Sunday":    {},
	}
	if err != nil {
		return nil, err
	}
	for _, class := range teaches {
		t, err := time.Parse(utils.StdDateFormat, class.Date)
		if err != nil {
			return nil, err
		}
		day := t.Weekday().String()
		temp[day][class.ClassIndex] = true
		res[day].Location = class.ClassroomNo
	}

	for day, bitmap := range temp {
		for idx, val := range bitmap {
			if val {
				res[day].Indexes = append(res[day].Indexes, idx)
			}
		}
	}
	return res, nil
}

//ConnectDB connect postgresql database and redis database, return a connection object to manipulate data
func ConnectDB() Connection {
	db := pg.Connect(&pg.Options{
		Addr:     "db:5432",
		User:     "postgres",
		Password: "passwd",
		Database: "gostudy",
	})
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Connect to database successfully: ", db)
	if os.Getenv("DEBUG") == "TRUE" {
		db.AddQueryHook(&pgdebug.DebugHook{
			Verbose: true,
		})
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "passwd",
		DB:       0,
	})
	res, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	} else {
		log.Println("Ping redis database", res)
	}
	return Connection{db: db, rdb: rdb, ctx: ctx}
}
func (con *Connection) CloseDB() {
	err := con.rdb.Close()
	if err != nil {
		log.Println(err)
	}
	err = con.db.Close()
	if err != nil {
		log.Println(err)
	}
}

func (con *Connection) CheckUserIdPassword(userId, password string) bool {
	passwordHash := utils.Hashstr(password)
	adminer := new(Admin)
	err := con.db.Model(adminer).Where("admin_id=?", userId).Select()
	if err != nil {
		return false
	} else {
		if adminer.PasswordHash == passwordHash {
			return true
		} else {
			return false
		}
	}
}
func (con *Connection) UseridByOpenid(openId string, avatarUrl, nickName string,
	gender int, city, province, country string) (string, error) {
	user := new(User)
	newUser := &User{
		OpenId:    openId,
		AvatarUrl: avatarUrl,
		NickName:  nickName,
		Gender:    int8(gender),
		City:      city,
		Province:  province,
		Country:   country,
	}
	err := con.db.Model(user).Where("open_id=?", openId).Select()
	if user.UserId != "" {
		//如果已有用户，那么更新其他项
		newUser.UserId = user.UserId
		_, err2 := con.db.Model(newUser).WherePK().Update()
		if err2 != nil {
			return "", err2
		}
		return user.UserId, nil
	}
	//否则直接插入新用户，生成新的Userid
	_, err = con.db.Model(newUser).Insert()
	if err != nil {
		return "", err
	} else {
		log.Println(newUser)
		return newUser.UserId, nil
	}
}
func (con *Connection) TakeCourse(userId, courseId string, year int, semester string) error {
	_, err := con.db.Model(&Takes{
		UserId:   userId,
		CourseId: courseId,
		Year:     year,
		Semester: semester,
	}).OnConflict("DO NOTHING").Insert()
	return err
}
func (con *Connection) DropCourse(userId, courseId string, year int, semester string) error {
	_, err := con.db.Model(&Takes{
		UserId:   userId,
		CourseId: courseId,
		Year:     year,
		Semester: semester,
	}).WherePK().Delete()
	return err
}
func (con *Connection) getLatestStudy(userId string) (*Study, error) {
	s := new(Study)
	err := con.db.Model(s).Where("?=?", pg.Ident("user_id"), userId).OrderExpr("start_time DESC").Limit(1).Select()
	fmt.Println(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

type HasStartError struct{}

func (HasStartError) Error() string {
	return "There is a start study record for this account, please end it first"
}
func (con *Connection) StartStudy(userId, classroomNo string) error {
	s, err := con.getLatestStudy(userId)
	if s != nil {
		if err != nil {
			return err
		}
		if s.TimeLen == -1 {
			// if last insert change it to -1, just change it back.
			// And mark it done
			return HasStartError{}
		}
	}
	s = &Study{
		UserId:      userId,
		StartTime:   time.Now(),
		TimeLen:     -1, //Set it to -1 for EndStudy to detect
		ClassroomNo: classroomNo,
	}
	_, err = con.db.Model(s).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}
	return nil
}
func (con *Connection) IsBuildingNo(buildingNo string) bool {
	buildings, err := con.AllBuildings()
	if err != nil {
		return false
	}
	ok := false
	for _, building := range buildings {
		if building.BuildingNo == buildingNo {
			ok = true
			break
		}
	}
	return ok
}

type NotStartError struct{}

func (NotStartError) Error() string {
	return "There is no start study records for this account"
}
func (con *Connection) EndStudy(userId string) (float64, error) {
	s, err := con.getLatestStudy(userId)
	if err != nil {
		return -1, err
	}
	if s.TimeLen != -1 { //If it has not been set by Start Study
		return -1, NotStartError{}
	}
	t := time.Now().Sub(s.StartTime).Seconds()
	s.TimeLen = t
	_, err = con.db.Model(s).WherePK().Update()
	return t, err
}
func (con *Connection) TotalStudyTime(UserId string) (float64, error) {
	var res float64
	err := con.db.Model((*Study)(nil)).ColumnExpr("sum(time_len)").Where("user_id=?", UserId).Select(&res)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (con *Connection)AllCourse(userId string,year int,semester string)([]Course,error){
	var res []Course
	err:=con.db.Model((*Course)(nil)).Join("NATURAL JOIN takes").Where("user_id=?", userId).Select(&res)
	if err!=nil{
		return nil,err
	} else{
		return res,nil
	}
}