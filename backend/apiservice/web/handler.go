package web

import (
	json2 "encoding/json"
	"fmt"
	"github/Outsider565/gostudy/apiservice/db"
	"github/Outsider565/gostudy/apiservice/utils"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func reportError(c *gin.Context, errCode int, errStr string) {
	log.Println(errStr)
	c.JSON(http.StatusBadRequest, gin.H{"code": errCode, "msg": errStr})
}

// EmptyClassroomsHandle from database using GET methods
// API: /api/empty_classroom?building=H2&date=2021-06-01&index=6,8,9
func EmptyClassroomsHandle(c *gin.Context, con *db.Connection) {
	queryBuilding := c.Query("building")
	ok := con.IsBuildingNo(queryBuilding)
	if !ok {
		errStr := fmt.Sprintf("%s is not in Building No list", queryBuilding)
		reportError(c, ParameterError, errStr)
		return
	}
	date, _ := c.GetQuery("date")
	if date == "" {
		date = time.Now().Format(utils.StdDateFormat)
	} else {
		_, err := time.Parse(utils.StdDateFormat, date) //检查字符串是否符合规定格式
		if err != nil {
			reportError(c, ParameterError, err.Error())
			return
		}
	}
	indexStr, _ := c.GetQuery("index")
	var indexList []int
	if indexStr == "" {
		indexList = append(indexList, utils.GetClassIndex(time.Now()))
	} else {
		indexStrList := strings.SplitN(indexStr, ",", -1)
		for _, val := range indexStrList {
			index, err := strconv.Atoi(val)
			if err != nil {
				reportError(c, ParameterError, err.Error()) //如果教室号不对直接退出
				return
			}
			indexList = append(indexList, index)
		}
	}
	hash := make(map[string]bool)
	for i, index := range indexList {
		emptyClassrooms, err := con.EmptyClassroomsIdx(queryBuilding, date, index)
		if err != nil {
			reportError(c, DatabaseError, err.Error())
			return
		}
		if i == 0 {
			for _, classroom := range emptyClassrooms {
				hash[classroom] = false
			}
		} else {
			for _, classroom := range emptyClassrooms {
				_, ok := hash[classroom]
				if !ok { // classroom in list don't even exist in classroom list 0, just ignore
					continue
				} else {
					hash[classroom] = true
				}
			}
			for key, val := range hash {
				if !val {
					delete(hash, key)
				}
			}
			for key := range hash {
				hash[key] = false
			}
		}

	}
	emptyClassrooms := make([]string, len(hash))
	i := 0
	for k := range hash {
		emptyClassrooms[i] = k
		i++
	}
	c.JSON(http.StatusOK, gin.H{
		"empty_classroom": emptyClassrooms,
		"date":            date,
		"index":           indexStr,
		"num":             len(emptyClassrooms),
	})
}

//RequestCrawlHandle using start_date(included) and end_date(included) using get methods, return state
//API: /admin/request_crawl?start_date=2021-06-01&end_date=2021-06-21
func RequestCrawlHandle(c *gin.Context, con *db.Connection) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		errStr := "start_date and end_date can't be empty string"
		reportError(c, ParameterError, errStr)
		return
	}
	start, err := time.Parse(utils.StdDateFormat, startDate)
	if err != nil {
		errStr := fmt.Sprintf("Can not parse start_date: %s", startDate)
		reportError(c, ParameterError, errStr)
		return
	}
	end, err := time.Parse(utils.StdDateFormat, endDate)
	if err != nil {
		errStr := fmt.Sprintf("Can not parse end_date: %s", endDate)
		reportError(c, ParameterError, errStr)
		return
	}
	if end.Before(start) {
		errStr := fmt.Sprintf("end_date: %s is earlier than start_date", startDate)
		reportError(c, ParameterError, errStr)
		return
	}
	log.Println("Crawl between:", start, end)
	var tArr []time.Time
	for idx := 0; start.AddDate(0, 0, idx).Before(end); idx++ {
		tArr = append(tArr, start.AddDate(0, 0, idx))
	}
	err = con.CrawlClasses(tArr)
	log.Println("Crawl complete!")
	if err != nil {
		reportError(c, DatabaseError, fmt.Sprintf("Crawl error %s", err))
		return
	}
	c.String(http.StatusOK, "Crawl successfully")
	return
}

//ReInitHandle everything, drop all tables, create all tables and constraint, insert building info
//API: /admin/re_init
func ReInitHandle(c *gin.Context, con *db.Connection) {
	err := con.ReInit()
	if err != nil {
		reportError(c, DatabaseError, err.Error())
		return
	}
	c.String(http.StatusOK, "ReInitHandle successfully")
	return
}

//FlushRdbHandle flush the redis database
//Will be called by Django
//API: /admin/flush_rdb
func FlushRdbHandle(c *gin.Context, con *db.Connection) {
	err := con.FlushRdb()
	if err != nil {
		reportError(c, DatabaseError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": Success,
	})
}

// CourseSearchHandle SearchFussNameHandle search for course fuss name, return top n most similar ones
//API: /api/course_search?max_num=5&course_name=example&course_id=example&year=example&semester=example&teacher_name=example
//API: /api/course_search_strict?max_num=5&course_name=example&course_id=example&year=example&semester=example&teacher_name=example
//num must be a valid integer and can't be omitted
func CourseSearchHandle(c *gin.Context, con *db.Connection, ifStrict bool) {
	num, err := strconv.Atoi(c.Query("max_num"))
	{
		if err != nil {
			reportError(c, ParameterError, err.Error())
			return
		}
	}
	courses, err := con.SearchCourse(num, c.Query("course_id"), c.Query("year"), c.Query("semester"), c.Query("course_name"), c.Query("teacher_name"), ifStrict)
	if err != nil {
		reportError(c, ParameterError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
		"num":     len(courses),
	})
	return
}

// DistanceHandle returns a slice of distances between given coordinates and buildings
//API: /api/distance?lat=31.297664&lon=121.504083
func DistanceHandle(c *gin.Context, con *db.Connection) {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		reportError(c, ParameterError, err.Error())
		return
	}
	lon, err := strconv.ParseFloat(c.Query("lon"), 64)
	if err != nil {
		reportError(c, ParameterError, err.Error())
		return
	}
	buildings, err := con.AllBuildings()
	if err != nil {
		reportError(c, DatabaseError, err.Error())
		return
	}
	var distanceSquare float64
	type buildingInfo map[string]interface{}
	var res []buildingInfo
	for _, building := range buildings {
		distanceSquare = (building.Latitude-lat)*(building.Latitude-lat) + (building.Longitude-lon)*(building.Longitude-lon)
		var temp = buildingInfo{
			"building_no":   building.BuildingNo,
			"building_name": building.BuildingName,
			"distance":      math.Sqrt(distanceSquare),
		}
		res = append(res, temp)
	}
	c.JSON(http.StatusOK, gin.H{"res": res})
	return
}

// CourseTimeHandle gives the course_time by given course_id
// API: /api/course_time_loc?course_id=?
func CourseTimeHandle(c *gin.Context, con *db.Connection) {
	courseId := c.Query("course_id")
	if courseId == "" {
		reportError(c, ParameterError, "courseId can't be empty string")
		return
	}
	t, err := con.GetCourseTimeAndLoc(courseId)
	if err != nil {
		reportError(c, DatabaseError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": t})
	return
}

// AdminAuthHandle return authentication token by posting user_id and password
// API: /auth/admin
// Post data:
// 			user_id: string, required
//			password: string, required
// As we are only using the https connection for auth, we can hash it at the server side
func AdminAuthHandle(c *gin.Context, con *db.Connection) {
	type UserInfo struct {
		UserId   string `form:"user_id" json:"user_id" xml:"user_id" yaml:"user_id" binding:"required"`
		Password string `form:"password" json:"password" xml:"password" yaml:"password" binding:"required"`
	}
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		reportError(c, ParameterError, "Invalid parameters")
		return
	}
	if con.CheckUserIdPassword(user.UserId, user.Password) {
		// 生成Token
		tokenString, _ := GenToken(user.UserId)
		c.JSON(http.StatusOK, gin.H{
			"err":   nil,
			"token": tokenString,
		})
		return
	}
	reportError(c, AuthenticationError, "Invalid Username and password")
	return
}

// WxAuthHandle return authentication token by posting wechat open id and other data
// API: /auth/wx
// Post data:
// 			user_id: string, required
//			avatar_url: string, required
//			nick_name: string, required
//			gender: int, required
// 			city: string, required
// 			province: string, required
// 			country: string, required
func WxAuthHandle(c *gin.Context, con *db.Connection) {
	type UserInfo struct {
		Code      string `form:"code" json:"code" xml:"code" yaml:"code" binding:"required"`
		AvatarUrl string `form:"avatar_url" json:"avatar_url" xml:"avatar_url" yaml:"avatar_url" binding:"required"`
		NickName  string `form:"nick_name" json:"nick_name" xml:"nick_name" yaml:"nick_name" binding:"required"`
		Gender    int    `form:"gender" json:"gender" xml:"gender" yaml:"gender" binding:"required"`
		City      string `form:"city" json:"city" xml:"city" yaml:"city" binding:"required"`
		Province  string `form:"province" json:"province" xml:"province" yaml:"province" binding:"required"`
		Country   string `form:"country" json:"country" xml:"country" yaml:"country" binding:"required"`
	}
	var user UserInfo
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		reportError(c, ParameterError, "Invalid data: "+err.Error())
		panic("hello")
	}
	err = json2.Unmarshal(jsonData, &user)
	if err != nil {
		reportError(c, ParameterError, "Invalid parameters: "+err.Error())
		panic("hello")
	}
	log.Println(user)
	var userId string
	if os.Getenv("TEST_WX") == "FALSE" {
		resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", os.Getenv("APPID"), os.Getenv("APPSECRET"), user.Code))
		if err != nil {
			reportError(c, AuthenticationError, err.Error())
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(resp.Body)
		type wxResp struct {
			SessionKey string `json:"session_key"`
			Openid     string `json:"openid"`
			Errcode    int    `json:"errcode"`
			Errmsg     string `json:"errmsg"`
		}
		res := new(wxResp)
		body, err := ioutil.ReadAll(resp.Body)
		err = json2.Unmarshal(body, res)
		if err != nil {
			reportError(c, AuthenticationError, err.Error())
			return
		}
		if res.Errcode != 0 {
			reportError(c, AuthenticationError, res.Errmsg)
			return
		}
		userId, err = con.UseridByOpenid(res.Openid, user.AvatarUrl, user.NickName, user.Gender, user.City, user.Province, user.Country)
		if err != nil {
			reportError(c, AuthenticationError, err.Error())
			return
		}
	} else if os.Getenv("TEST_WX") == "TRUE" {
		if user.Code != "test" {
			reportError(c, AuthenticationError, "In TEST_WX MODE, the only user code must be test")
			return
		}
		userId = db.TestUser.UserId
	} else {
		panic("Invalid env TEST_WX: " + os.Getenv("TEST_WX"))
	}
	// 生成Token
	tokenString, _ := GenToken(userId)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
	return
}

// GetClassesHandle get the classes of given classroom at a given date
// API: /api/get_classes?date=2021-06-08&classroom_no=H3109
func GetClassesHandle(c *gin.Context, con *db.Connection) {
	date := c.Query("date")
	ok, err := utils.CheckDate(date)
	if !ok && err != nil {
		reportError(c, ParameterError, err.Error())
		return
	}
	classroomNo := c.Query("classroom_no")
	if classroomNo == "" {
		reportError(c, ParameterError, "classroom_no can't be null")
		return
	}
	classes, err := con.GetClassesOfDate(date, classroomNo)
	if err != nil {
		reportError(c, DatabaseError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": classes})
}

//ATTENTION!!!
//Every p_api need token

func getUser(c *gin.Context) string {
	userIdVal, ok := c.Get("user_id")
	if !ok {
		reportError(c, AuthenticationError, "Can't get user properly")
		return ""
	}
	userId, ok := userIdVal.(string)
	if !ok {
		reportError(c, AuthenticationError, "Can't get user properly")
		return ""
	}
	return userId
}

// TakeCourseHandle let user take courses by course_id
// API: /p_api/take_course?course_id=example&year=example&semester=example
func TakeCourseHandle(c *gin.Context, con *db.Connection) {
	courseId := c.Query("course_id")
	if courseId == "" {
		reportError(c, ParameterError, "course id can't be null")
		return
	}
	yearStr := c.Query("year")
	if yearStr == "" {
		reportError(c, ParameterError, "year can't be null")
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		reportError(c, ParameterError, err.Error())
		return
	}
	semester := c.Query("semester")
	if semester == "" {
		reportError(c, ParameterError, "semester can't be null")
		return
	}
	userId := getUser(c)
	if userId == "" {
		return
	}
	courses, err := con.SearchCourse(1, courseId, yearStr, semester, "", "", true)
	if err != nil {
		reportError(c, DatabaseError, err.Error())
	}
	if len(courses) == 0 {
		reportError(c, DatabaseError, "No corresponding courseId in database")
		return
	}
	err = con.TakeCourse(userId, courseId, year, semester)
	if err != nil {
		reportError(c, DatabaseError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": Success,
	})
}

// DropCourseHandle let user drop courses by course_id
// API: /p_api/drop_course?course_id=example&year=example&semester=example
// WARNING: drop_course will not send an error if it tries to delete an not taken course
// It will simply do nothing
func DropCourseHandle(c *gin.Context, con *db.Connection) {
	courseId := c.Query("course_id")
	if courseId == "" {
		reportError(c, ParameterError, "course id can't be null")
		return
	}
	yearStr := c.Query("year")
	if yearStr == "" {
		reportError(c, ParameterError, "year can't be null")
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		reportError(c, ParameterError, err.Error())
		return
	}
	semester := c.Query("semester")
	if semester == "" {
		reportError(c, ParameterError, "semester can't be null")
		return
	}
	userId := getUser(c)
	if userId == "" {
		return
	}
	courses, err := con.SearchCourse(1, courseId, yearStr, semester, "", "", true)
	if err != nil {
		reportError(c, DatabaseError, err.Error())
	}
	if len(courses) == 0 {
		reportError(c, DatabaseError, "No corresponding courseId in database")
		return
	}
	err = con.DropCourse(userId, courseId, year, semester)
	if err != nil {
		reportError(c, DatabaseError, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": Success,
	})
}

// AllCourseHandle returns all course of a given user(by posting userid)
// API: /p_api/all_course?year=2021&semester=Spring
// If year OR semester is unset, it will use current time to determine the year and semester
func AllCourseHandle(c *gin.Context, con *db.Connection) {
	userId := getUser(c)
	if userId == "" {
		return
	}
	yearStr := c.Query("year")
	semester := c.Query("semester")
	var year int
	var err error
	if yearStr == "" || semester == "" {
		year = time.Now().Year()
		semester = utils.GetSemester(time.Now())
	} else {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			reportError(c, ParameterError, "Cannot parse year")
			return
		}
	}
	res, err := con.AllCourse(userId, year, semester)
	if err != nil {
		reportError(c, DatabaseError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    Success,
		"courses": res,
	})
}

// StartStudyHandle set start_study status
// API: /p_api/start_study?classroom_no=example
func StartStudyHandle(c *gin.Context, con *db.Connection) {
	userId := getUser(c)
	if userId == "" {
		return
	}
	classroomNo := c.Query("classroom_no")
	err := con.StartStudy(userId, classroomNo)
	if err != nil {
		reportError(c, DatabaseError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": Success,
	})
}

// EndStudyHandle end the study and returns the study time (in seconds)
// API: /p_api/end_study
func EndStudyHandle(c *gin.Context, con *db.Connection) {
	userId := getUser(c)
	if userId == "" {
		return
	}
	t, err := con.EndStudy(userId)
	if err != nil {
		reportError(c, DatabaseError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     Success,
		"duration": t,
	})
}

// StudyTimeHandle calculate the total time of study (in seconds)
// API: /p_api/study_time
func StudyTimeHandle(c *gin.Context, con *db.Connection) {
	userId := getUser(c)
	if userId == "" {
		return
	}
	t, err := con.TotalStudyTime(userId)
	if err != nil {
		reportError(c, DatabaseError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":       Success,
		"total_time": t,
	})
}
