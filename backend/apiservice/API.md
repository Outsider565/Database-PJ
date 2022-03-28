# API README

### **空闲教室**

- **API示例：** /api/empty_classroom?building=``` para```&date=```para```&index=```para```
- **URL:**   /api/empty_classroom
- **请求方式：** GET
- **参数列表：**

|   Para   |                           Type                           |
| :------: | :------------------------------------------------------: |
| building | H2 \ H3NoH30 \ H4 \ H5 \ H6 \ HGD \ HGX \ HQ \ J \ Z \ F |
|   date   |                        YYYY-MM-DD                        |
|  index   |                       Range(1, 14)                       |



### 用户信息

- **API示例：** /api/add_user?openId=```para```&nickName=``` para```&gender=```para```&index=```para```
- **URL:**  /api/add_user
- **请求方式：** POST/GET (get时仅有openId参数，方法为GET)
- **参数列表：**

|   Para   |   Type   |             Example              |
| :------: | :------: | :------------------------------: |
|  openid  | char(32) | 001G9LFa1rPA9B0fOmIa1dITJr1G9LFK |
| nickName | Varchar  |               Xzy                |
|  Gender  | int 2/1/0  |               - -                |
|   city   | Varchar  |       Pudong New District        |
| province | Varchar  |             Shanghai             |
| country  | Varchar  |              China               |

- **返回值说明：**
  - 形式：JSON
  - 内容：返回相应用户信息

Note: 0 for man, 1 for women, 2 for others

### 用户学习时长

- **API示例：** /api/study_time?openId=``` para```&timeLen=```para```
- **URL:**   /api/study_time
- **请求方式：** POST/GET (get时仅有openId参数，方法为GET)
- **参数列表：**

|  Para   |     Type     |             Example              |
| :-----: | :----------: | :------------------------------: |
| openid  |   char(32)   | 001G9LFa1rPA9B0fOmIa1dITJr1G9LFK |
| timeLen | unsigned int |   12123123123 【单位 second】    |

- **返回值说明：**
  - 形式：JSON
  - 内容：返回用户学习时长



### 用户开始学习

- **API示例：** /api/start_study?openId=``` para```&timeStamp=```para```&location=```para```
- **URL:**   /api/start_study

- **请求方式：** POST/GET (get时仅有openId参数，方法为GET)
- **参数列表：**

|   Para    |   Type    |                           Example                            |
| :-------: | :-------: | :----------------------------------------------------------: |
|  openid   | char(32)  |               001G9LFa1rPA9B0fOmIa1dITJr1G9LFK               |
| timeStamp | timeStamp |                          1622712992                          |
| location  |   Char    | H2 \ H3NoH30 \ H4 \ H5 \ H6 \ HGD \ HGX \ HQ \ J \ Z \ F \ R (R代表其他) |

- **返回值说明：**
  - 形式：JSON
  - 内容：返回用户学习开始Unix Time



### 课程查询

- **API示例：** /api/search_course?searchKey=```para```    or /api/search_course?courseId=```para```
- **URL:**   /api/search_course
- **请求方式：** GET
- **参数列表：**

|   Para    |  Type   |            Example            |
| :-------: | :-----: | :---------------------------: |
| searchKey | Varchar | 计算\生物\马克思 （模糊搜索） |
| courseId  | Varchar |   PTSS110067 (课程代号查询)   |

- **返回值说明：**

  返回一个匹配课表 JSON，其中每个

  | Return Value |     Key     |  Type  |              Example               |
  | :----------: | :---------: | :----: | :--------------------------------: |
  |    课程名    |   clsName   |  Char  |             计算机系统             |
  |   课程代码   |    clsNo    |  Char  |             COMP130144             |
  |    教师名    | teacherName |  Char  |                张凯                |
  |   开课时间   |   clsTime   | Object |              (见示例)              |
  |   课程类型   |    Type     |  Int   | 0 文/1 理 /2 工/3 医/4 体育/5 思政 |

  错误或失败返回 JSON 其带一个 errorMsg

- **返回格式样例**

`````json
{
  {
  	"clsName": "计算机系统",
  	"clsNo" : "COMP130144",
  	"teacherName" : "张凯",
  	"clsTime": {
  							"周一" : [8,9],
								"周三" : [6,7,8]
								// 看起来比较阴间，如果你有什么好设计，可以改一下
								},
		"type" : 2
	},
	{
    ...
  },
}
`````



### 用户课表添加

- **API示例：** /api/add_course?courseId=```para```&type=```para```
- **URL:**   /api/add_course
- **请求方式：** POST
- **参数列表：**

|   Para   |  Type   |          Example          |
| :------: | :-----: | :-----------------------: |
| courseId | Varchar | PTSS110067 (课程代号查询) |
|   Type   |   Int   |      0 选课 / 1 旁听      |

