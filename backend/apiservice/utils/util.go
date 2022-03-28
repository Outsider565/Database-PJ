package utils

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"time"
)

const StdDateFormat = "2006-01-02"

func GetSemester(t time.Time) (res string) {
	m := t.Month()
	switch {
	case m <= 6:
		res = "Spring"
	case m <= 8:
		res = "Summer"
	case m <= 12:
		res = "Autumn"
	default:
		panic("cannot get semester")
	}
	return
}

type dayTime struct {
	hour   int
	minute int
}

func getDayTime(h, m int) dayTime {
	return dayTime{
		hour:   h,
		minute: m,
	}
}
func (d *dayTime) getMinute() int {
	return d.hour*60 + d.minute
}
func (d *dayTime) inRange(dLow, dHigh dayTime) bool {
	return d.getMinute() >= dLow.getMinute() && d.getMinute() <= dHigh.getMinute()
}
func GetClassIndex(t time.Time) (res int) {
	//这里其实可以算做一个小bug或者feature，假设目前这门课在下课，我们确实可以进去，但是显然两节课之间的课件不应该标为空教室，为了让服务先跑起来这样写
	d := getDayTime(t.Hour(), t.Minute())
	switch {
	case d.inRange(getDayTime(8, 0), getDayTime(8, 55)):
		res = 1
	case d.inRange(getDayTime(8, 55), getDayTime(9, 55)):
		res = 2
	case d.inRange(getDayTime(9, 55), getDayTime(10, 50)):
		res = 3
	case d.inRange(getDayTime(10, 50), getDayTime(11, 45)):
		res = 4
	case d.inRange(getDayTime(11, 45), getDayTime(12, 30)):
		res = 5
	case d.inRange(getDayTime(13, 30), getDayTime(14, 25)):
		res = 6
	case d.inRange(getDayTime(14, 25), getDayTime(15, 25)):
		res = 7
	case d.inRange(getDayTime(15, 25), getDayTime(16, 20)):
		res = 8
	case d.inRange(getDayTime(16, 20), getDayTime(17, 15)):
		res = 9
	case d.inRange(getDayTime(17, 15), getDayTime(18, 00)):
		res = 10
	case d.inRange(getDayTime(18, 30), getDayTime(19, 25)):
		res = 11
	case d.inRange(getDayTime(19, 25), getDayTime(20, 20)):
		res = 12
	case d.inRange(getDayTime(20, 20), getDayTime(21, 15)):
		res = 13
	case d.inRange(getDayTime(21, 15), getDayTime(22, 00)):
		res = 14
	default:
		res = 0
		if os.Getenv("DEBUG") == "TRUE" {
			log.Println("Get index return 0 time:", t)
		}
	}
	return
}
func CheckDate(date string) (bool, error) {
	_, err := time.Parse(StdDateFormat, date)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
func Hashstr(Txt string) string {
	h := sha256.New()
	h.Write([]byte(Txt))
	bs := h.Sum(nil)
	sh := string(fmt.Sprintf("%x\n", bs))
	return sh
}
