package crawl

import (
	"fmt"
	_ "io"
	_ "io/ioutil"
	"log"
	_ "net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gocolly/colly"
)

const baseUrl string = "http://10.64.130.6/"

type Class struct {
	CourseId   string `json:"course_id"`
	CourseName string `json:"course_name"`
	Type       string `json:"type"`
	Teacher    string `json:"teacher"`
	StudentNum int    `json:"student_num"`
}

type ClassRoom struct {
	Name    string    `json:"name"`
	SeatNum int       `json:"seat_num"`
	Classes [14]Class `json:"classes"`
}

func (c Class) Empty() bool {
	if c.CourseName == "" && c.CourseId == "" && c.Type == "" {
		return true
	}
	return false
}

func (c Class) String() (s string) {
	if c.Empty() {
		return "NULL"
	}
	s = fmt.Sprintf("%s*%s*%s*%s*%d", c.CourseId, c.CourseName, c.Type, c.Teacher, c.StudentNum)
	return
}
func (c ClassRoom) String() (s string) {
	s = fmt.Sprintf("Name: %s\nNum: %d\nCourse: %s\n", c.Name, c.SeatNum, c.Classes)
	return
}
func parseClass(text string) (c Class) {
	//need improvement but not important:如果在同一时间、同一教室有两门课怎么办，如周一三教唐荣堂的毛概
	//思路：每门课的test分成五个部分，Type：全中文的描述；课程号：确定的13位ascii;课程名字，在方括号中的;老师名字，在圆括号中的；学生人数，在花括号中的
	//课程有可能为空（此时无课）
	//
	if len(text) == 0 {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic occured text is %s, err is %s\n", text, r)
			panic(r)
		}
	}()
	textRune := []rune(text)
	textRuneLen := len(textRune)
	runeIdx := 0
	for textRune[runeIdx] > unicode.MaxASCII {
		runeIdx++
	}
	if runeIdx >= textRuneLen { //防止越界
		return
	}
	c.Type = string(textRune[0:runeIdx]) //就是课程名字前你的"本""排考"等标志
	idIdx := runeIdx
	for idIdx < textRuneLen && string(textRune[idIdx]) != "[" {
		idIdx++
	}
	if idIdx == textRuneLen { //防止越界
		//如果是临时借的教室，那么前五位是课程临借编号，后面是课程名
		c.CourseId = string(textRune[runeIdx : runeIdx+5])
		c.CourseName = string(textRune[runeIdx+5 : idIdx])
		return
	}
	c.CourseId = string(textRune[runeIdx:idIdx])
	nameIdx := idIdx
	for nameIdx < textRuneLen && string(textRune[nameIdx]) != "]" {
		nameIdx++
	}
	c.CourseName = string(textRune[idIdx+1 : nameIdx])
	text = strings.ReplaceAll(strings.ReplaceAll(text, "）", ")"), "（", "(")
	if strings.Contains(string(textRune[nameIdx:textRuneLen]), "[") {
		return
	}
	teacherIdx := nameIdx + 1
	if teacherIdx >= textRuneLen { //防止越界
		return
	}
	for teacherIdx < textRuneLen && string(textRune[teacherIdx]) != "{" {
		teacherIdx++
	}
	if teacherIdx-1 < nameIdx+2 {
		return
	}
	c.Teacher = string(textRune[nameIdx+2 : teacherIdx-1])
	numIdx := teacherIdx
	if strings.Contains(c.Teacher, "[") {
		log.Println(string(textRune))
		c.Teacher = ""
	}
	if numIdx >= textRuneLen { //防止越界
		return
	}
	for numIdx < textRuneLen && string(textRune[numIdx]) != "人" {
		numIdx++
	}
	num, err := strconv.Atoi(string(textRune[teacherIdx+1 : numIdx]))
	if err != nil {
		log.Println(err, text)
	}
	c.StudentNum = num
	return
}
func crawlFromUrl(url string) []ClassRoom {
	classrooms := make([]ClassRoom, 0, 200)
	c := colly.NewCollector()
	c.SetRequestTimeout(time.Minute)
	c.OnHTML("table[id='statusTable_0']", func(table *colly.HTMLElement) {
		table.ForEach("tr", func(row int, line *colly.HTMLElement) {
			if row < 3 {
				return
			}
			tempClassroom := ClassRoom{}
			idxShift := 4 //因为教室那栏有时候可能有4个column，三个是麦克风还有标志
			line.ForEach("td", func(idx int, class *colly.HTMLElement) {
				if class.Attr("style") == "" {
					return
				}
				text := strings.TrimSpace(
					strings.ReplaceAll(
						strings.ReplaceAll(strings.ReplaceAll(class.Text, "\n", ""), " ", ""),
						"&nbsp",
						"",
					),
				)
				text = strings.Map(func(r rune) rune {
					if unicode.IsPrint(r) {
						return r
					}
					return -1
				}, text)
				switch {
				case idx == 0:
					tempClassroom.Name = text
					if len(class.ChildAttr("table", "border")) == 0 {
						idxShift -= 3
					}
				case idx == idxShift:
					seatNum, err := strconv.Atoi(text)
					if err != nil {
						log.Printf("Cannot convert '%s'", text)
					}
					tempClassroom.SeatNum = seatNum
				default:
					idx -= idxShift
					tempClass := parseClass(text)
					tempClassroom.Classes[idx-1] = tempClass
				}
			})
			if tempClassroom.Name != "" {
				classrooms = append(classrooms, tempClassroom)
			}
		})
	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting: ", r.URL.String())
	})
	err := c.Visit(url)
	if err != nil {
		log.Println("Error when visiting ", url, err)
	}
	c.Wait()
	return classrooms
}

// Crawl from building string and day string
// building string should be one of H2, H3NoH30, H4, H5, H6, HGD, HGX, HQ, J, Z, F
// day string should be formatted as YYYY-MM-DD
func Crawl(building string, day string) []ClassRoom {
	url := fmt.Sprintf("%s?b=%s&day=%s", baseUrl, building, day)
	return crawlFromUrl(url)
}
