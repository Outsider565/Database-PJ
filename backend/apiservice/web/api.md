package web // import "github/Outsider565/gostudy/apiservice/web"


CONSTANTS

const (
	Success             = 0
	ParameterError      = 1
	DatabaseError       = 2
	AuthenticationError = 3
)

FUNCTIONS

func AdminAuthHandle(c *gin.Context, con *db.Connection)
    AdminAuthHandle return authentication token by posting user_id and password
    API: /auth/admin Post data:

        user_id: string, required
        password: string, required

    As we are only using the https connection for auth, we can hash it at the
    server side

func AdminAuthMiddleware() func(c *gin.Context)
func AllCourseHandle(c *gin.Context, con *db.Connection)
    AllCourseHandle returns all course of a given user(by posting userid) API:
    /p_api/all_course?year=2021&semester=Spring If year OR semester is unset, it
    will use current time to determine the year and semester

func CourseSearchHandle(c *gin.Context, con *db.Connection, ifStrict bool)
    CourseSearchHandle SearchFussNameHandle search for course fuss name, return
    top n most similar ones API:
    /api/course_search?max_num=5&course_name=example&course_id=example&year=example&semester=example&teacher_name=example
    API:
    /api/course_search_strict?max_num=5&course_name=example&course_id=example&year=example&semester=example&teacher_name=example
    num must be a valid integer and can't be omitted

func CourseTimeHandle(c *gin.Context, con *db.Connection)
    CourseTimeHandle gives the course_time by given course_id API:
    /api/course_time_loc?course_id=?

func DistanceHandle(c *gin.Context, con *db.Connection)
    DistanceHandle returns a slice of distances between given coordinates and
    buildings API: /api/distance?lat=31.297664&lon=121.504083

func DropCourseHandle(c *gin.Context, con *db.Connection)
    DropCourseHandle let user drop courses by course_id API:
    /p_api/drop_course?course_id=example&year=example&semester=example WARNING:
    drop_course will not send an error if it tries to delete an not taken course
    It will simply do nothing

func EmptyClassroomsHandle(c *gin.Context, con *db.Connection)
    EmptyClassroomsHandle from database using GET methods API:
    /api/empty_classroom?building=H2&date=2021-06-01&index=6,8,9

func EndStudyHandle(c *gin.Context, con *db.Connection)
    EndStudyHandle end the study and returns the study time (in seconds) API:
    /p_api/end_study

func FlushRdbHandle(c *gin.Context, con *db.Connection)
    FlushRdbHandle flush the redis database Will be called by Django API:
    /admin/flush_rdb

func GenToken(userId string) (string, error)
func GetClassesHandle(c *gin.Context, con *db.Connection)
    GetClassesHandle get the classes of given classroom at a given date API:
    /api/get_classes?date=2021-06-08&classroom_no=H3109

func ReInitHandle(c *gin.Context, con *db.Connection)
    ReInitHandle everything, drop all tables, create all tables and constraint,
    insert building info API: /admin/re_init

func RegisterAdminAPI(r *gin.RouterGroup, con *db.Connection)
func RegisterAllAPI(r *gin.RouterGroup, con *db.Connection)
func RegisterAuth(r *gin.RouterGroup, con *db.Connection)
func RegisterPrivateAPI(r *gin.RouterGroup, con *db.Connection)
func RequestCrawlHandle(c *gin.Context, con *db.Connection)
    RequestCrawlHandle using start_date(included) and end_date(included) using
    get methods, return state API:
    /admin/request_crawl?start_date=2021-06-01&end_date=2021-06-21

func StartStudyHandle(c *gin.Context, con *db.Connection)
    StartStudyHandle set start_study status API:
    /p_api/start_study?classroom_no=example

func StudyTimeHandle(c *gin.Context, con *db.Connection)
    StudyTimeHandle calculate the total time of study (in seconds) API:
    /p_api/study_time

func TakeCourseHandle(c *gin.Context, con *db.Connection)
    TakeCourseHandle let user take courses by course_id API:
    /p_api/take_course?course_id=example&year=example&semester=example

func WxAuthHandle(c *gin.Context, con *db.Connection)
    WxAuthHandle return authentication token by posting wechat open id and other
    data API: /auth/wx Post data:

        user_id: string, required
        avatar_url: string, required
        nick_name: string, required
        gender: int, required
        city: string, required
        province: string, required
        country: string, required

func WxAuthMiddleware() func(c *gin.Context)

TYPES

type Claims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

func ParseToken(tokenString string) (*Claims, error)

