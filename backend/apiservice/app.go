package main

import (
	"fmt"
	"github/Outsider565/gostudy/apiservice/db"
	"github/Outsider565/gostudy/apiservice/utils"
	"github/Outsider565/gostudy/apiservice/web"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

//这是一个测试函数，插入今天未来的两个星期的课作为测试
func doTestInsert(c db.Connection) {
	t, err := time.Parse(utils.StdDateFormat, "2021-03-01")
	if err != nil {
		panic(err)
	}
	var tArr []time.Time
	for idx := 0; idx < 130; idx++ {
		tArr = append(tArr, t.AddDate(0, 0, idx))
	}
	if err := c.CrawlClasses(tArr); err != nil {
		log.Println(err)
	}
}

func main() {
	con := db.ConnectDB()
	defer con.CloseDB()
	if os.Getenv("TEST_REINIT") == "TRUE" {
		err := con.DropAllTables()
		if err != nil {
			panic(err)
		}
	} else if os.Getenv("TEST_REINIT") != "FALSE" {
		panic(fmt.Sprintf("TEST_REINIT can't be %s\n", os.Getenv("TEST_REINIT")))
	}
	err := con.CreateSchema()
	if err != nil {
		panic(err)
	}
	err = con.InitBuildingInfo()
	if err != nil {
		panic(err)
	}
	err = con.InitAdmin()
	if err != nil {
		panic(err)
	}
	if os.Getenv("TEST_INSERT") == "TRUE" {
		go doTestInsert(con)
	} else if os.Getenv("TEST_INSERT") != "FALSE" {
		panic(fmt.Sprintf("TEST_INSERT can't be %s\n", os.Getenv("TEST_INSERT")))
	}
	r := gin.Default()
	apiRouter := r.Group("/api")
	web.RegisterAllAPI(apiRouter, &con)
	authRouter := r.Group("/auth")
	web.RegisterAuth(authRouter, &con)
	adminRouter := r.Group("/admin_api")
	web.RegisterAdminAPI(adminRouter, &con)
	pApiRouter := r.Group("/p_api")
	web.RegisterPrivateAPI(pApiRouter, &con)
	err = r.Run(":1234")
	if err != nil {
		panic(err)
	}
}
