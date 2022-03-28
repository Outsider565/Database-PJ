// app.js

App({
  onLaunch() {
    // Nothing
  },
  globalData: {
    isLogin: false, // 是否登录
    isStudy: false, // 是否学习
    isAddingCourse: false, // 是否在添加课程
    curCls: '',
    userInfo: {},
    userCourseList: [], // 用户课表
    userToken: "",  // 后端服务器鉴权 token
    addCourseId: "",
    baseUrl: "",
    authUrl: "",
    pUrl: "",
    buildingList: ["H2", "H3NoH30", "H4", "H5", "H6", "HGX"],
    weekDict: {
      'Sunday': '周日',
      'Monday' : '周一',
      'Tuesday': '周二',
      'Wednesday': '周三',
      'Thursday': '周四',
      'Friday': '周五',
      'Saturday': '周六',
    },
    dayDict: {
      'Monday' : 1,
      'Tuesday': 2,
      'Wednesday': 3,
      'Thursday': 4,
      'Friday': 5,
      'Saturday': 6,
      'Sunday': 0
    }
  }
})
