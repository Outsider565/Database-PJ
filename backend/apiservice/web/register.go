package web

import (
	"github.com/gin-gonic/gin"
	"github/Outsider565/gostudy/apiservice/db"
)

func RegisterAllAPI(r *gin.RouterGroup, con *db.Connection) {
	r.GET("/empty_classroom", func(c *gin.Context) {
		EmptyClassroomsHandle(c, con)
	})
	r.GET("/course_search", func(c *gin.Context) {
		CourseSearchHandle(c, con, false)
	})
	r.GET("/course_search_strict", func(c *gin.Context) {
		CourseSearchHandle(c, con, true)
	})
	r.GET("/distance", func(c *gin.Context) {
		DistanceHandle(c, con)
	})
	r.GET("/course_time_loc", func(c *gin.Context) {
		CourseTimeHandle(c, con)
	})
	r.GET("/get_classes", func(c *gin.Context) {
		GetClassesHandle(c, con)
	})
}
func RegisterAdminAPI(r *gin.RouterGroup, con *db.Connection) {
	r.Use(AdminAuthMiddleware())
	r.GET("/request_crawl", func(c *gin.Context) {
		RequestCrawlHandle(c, con)
	})
	r.GET("/re_init", func(c *gin.Context) {
		ReInitHandle(c, con)
	})
	r.GET("/flush_rdb", func(c *gin.Context) {
		FlushRdbHandle(c, con)
	})
}
func RegisterAuth(r *gin.RouterGroup, con *db.Connection) {
	r.POST("/admin", func(c *gin.Context) {
		AdminAuthHandle(c, con)
	})
	r.POST("/wx", func(c *gin.Context) {
		WxAuthHandle(c, con)
	})
}
func RegisterPrivateAPI(r *gin.RouterGroup, con *db.Connection) {
	r.Use(WxAuthMiddleware())
	r.GET("/take_course", func(c *gin.Context) {
		TakeCourseHandle(c, con)
	})
	r.GET("/drop_course", func(c *gin.Context) {
		DropCourseHandle(c, con)
	})
	r.GET("/start_study", func(c *gin.Context) {
		StartStudyHandle(c, con)
	})
	r.GET("/end_study", func(c *gin.Context) {
		EndStudyHandle(c, con)
	})
	r.GET("/study_time", func(c *gin.Context) {
		StudyTimeHandle(c, con)
	})
	r.GET("/all_course",func(c *gin.Context){
		AllCourseHandle(c, con)
	})
}
