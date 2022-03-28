import Toast from "../../src/dist/toast/toast";
const app = getApp()

Page({
  /* 初始数据 */
  data: {
    // 用户 Token or Maybe called cookies
    userToken: app.globalData.userToken,
    // 周课程总表 
    WeekCourseList: Array(7).fill([]),

    dayPickerShow: false,
    timePicker1Show: false,
    timePicker2Show: false,
    showAddCls: false,
    dayColumns: ["周日", "周一", "周二", "周三", "周四", "周五", "周六"],
    timeColumns: ["早八1", "早八2", "第三节", "第四节", "第五节", "下午1", "下午2", "下午3", "下午4", "晚课1", "晚课2", "晚课3", "晚课4"],
    
    
    customCourseName: "",
    // Time Picker
    customCourseDay: {value: "周一", index: 1},
    customCourseTimeIndex: [0, 1],
    customCourseTime: "早八1 ~ 早八2"
  },
  onLoad: function (options) {
    wx.loadFontFace({
      family: 'ZCOOL_h',
      source: 'url(../../../../src/font/ZCOOL_Kuaile/Zcool_h.ttf)',
      success: res => {
        console.log('font load success', res)
      },
      fail: err => {
        console.log('font load fail', err)
      }
    })
  },
  onShow: function () {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({
        active: 2
      })
    }
    // 每次展示刷新，以防深夜过点
    this.setData({ 
      courseEntityList: this.data.WeekCourseList[new Date().getDay()],
      currentDate: '星期'+'日一二三四五六'.charAt(new Date().getDay())
    });
    
    if(app.globalData.userCourseList.length > 0){
      wx.showLoading({ title: '加载中', })
      this.addCourse();
    }
  },
  onInput(event) {
    this.setData({
      currentDate: event.detail,
    });
  },
  addCls: function() {
    if(!app.globalData.isLogin){
      Toast({
        type: 'fail',
        message: '请先登录',
        duration: 1000
      });
      setTimeout(() => {
        wx.switchTab({ url: '/pages/user/index' })
      }, 1000);
      return;
    }
    this.setData({ showAddCls: true })
  },

  hideAddCls: function(){
    this.setData({ showAddCls: false })
  },
  onAddClsBtnTap: function(){
    if(app.globalData.isLogin === false){
      Toast({
        type: 'fail',
        message: '请先登录',
        duration: 800
      });
      setTimeout(wx.switchTab({
        url: '/pages/user/index',
      }), 800)
    }
    let courseName = this.data.customCourseName;
    if(!courseName){
      Toast({
        type: 'fail',
        message: '请输入课程名',
        duration: 1500
      });
      return;
    }
    // 调整为加课模式
    app.globalData.isAddingCourse = true; 
    
    var timeCol = this.data.timeColumns;
    let course_day_idx_list = [this.data.customCourseDay.index];
    let course_time_idx_list = [this.data.customCourseTimeIndex];
    let courseTime = this.data.customCourseTime;
    var CourseEntity = {
      CourseName: courseName,
      CourseId: "",
      Teacher: "暂无",
      CourseDayIdxList: course_day_idx_list,
      CourseTimeIdxList: course_time_idx_list,
      CourseTimeStr: courseTime,
      CourseTimeTag: [timeCol[this.data.customCourseTimeIndex[0]], timeCol[this.data.customCourseTimeIndex.slice(-1)]],
      __type__: 'Custom'
    };

    if(this.checkAddCourseValid(CourseEntity.CourseDayIdxList, CourseEntity.CourseTimeIdxList)){
      this.data.WeekCourseList[course_day_idx_list[0]].push(CourseEntity)
  
      this.setData({ courseEntityList: this.data.WeekCourseList[new Date().getDay()] })
      if(app.globalData.isAddingCourse === true){
        Toast({
          type: 'success',
          message: '添加成功',
          duration: 1500
        });
        app.globalData.isAddingCourse = false;
      }
    }
    this.setData({ showAddCls: false });
  },
  onDisplayDayPicker: function(e) {
    this.setData({ 
      dayPickerShow: true,
      showAddCls: false
    });
  },
  onDisplayTimePicker: function(e) {
    this.setData({ 
      timePicker1Show: true,
      showAddCls: false
    });
  },
  onDayCancel: function(){
    this.setData({ 
      dayPickerShow: false,
      showAddCls: true 
    });
  },
  onDayConfirm: function(e){
    this.setData({ 
      dayPickerShow: false,
      showAddCls: true,
      customCourseDay: e.detail
    });
  },
  onTime1Cancel: function(){
    this.setData({ 
      timePicker1Show: false,
      showAddCls: true 
    });
  },
  onTime1Confirm: function(e){
    // 开始时段选择结束，展示结束时段选择
    this.setData({ 
      timePicker1Show: false,
      timePicker2Show: true,
      customCourseStTime: e.detail
    });
  },
  onTime2Cancel: function(){
    this.setData({ 
      timePicker2Show: false,
      showAddCls: true
    });
  },
  onTime2Confirm: function(e){
    this.setData({ 
      timePicker2Show: false,
      showAddCls: true,
      customCourseEdTime: e.detail
    });
    // 检查合法性，开始时间小于结束时间
    if(this.checkCustomTimeValid()){
      let startTime = this.data.customCourseStTime;
      let endTime = this.data.customCourseEdTime;
      var custom_course_time_index = [];
      for(let i = startTime.index; i <= endTime.index; ++i)
        custom_course_time_index.push(i);

      this.setData({ 
        customCourseTime: startTime.value + " ~ " + endTime.value,
        customCourseTimeIndex: custom_course_time_index
      });
    }
    else{
      Toast({
        type: 'fail',
        message: '开始时间不能晚于结束时间',
        duration: 1500
      });
      this.setData({ 
        customCourseStTime: "",
        customCourseEdTime: ""
      });
    }
  },
  checkCustomTimeValid: function(){
    let startTime = this.data.customCourseStTime;
    let endTime = this.data.customCourseEdTime;
    if(!startTime || !endTime)
      return false;
    
    return startTime.index <= endTime.index;
  },
  onCourseNameChange: function(e){
    this.setData({ customCourseName: e.detail })
  },
  onDeleteCourse: function(e){
    if(e.detail != 'right')
      return;
    let delIndex = e.target.id;
    let courseDayIdx = new Date().getDay();
    var dropped_course = this.data.WeekCourseList[courseDayIdx].splice(delIndex, 1);
    if(dropped_course[0].__type__ == 'Course')
      this.rmInDB(dropped_course[0]);
    this.setData({ courseEntityList: this.data.WeekCourseList[courseDayIdx] })
  },
  rmInDB: function(course){
    let reqUrl = app.globalData.pUrl + 'drop_course';
    var _this = this;
    wx.request({
      url: reqUrl,
      method: 'GET',
      timeout: 2000,
      data: {
        course_id: course.CourseId,
        year: course.Year,
        semester: course.Semester
      },
      header: { 
        'content-type': 'application/json',
        'token': app.globalData.userToken
      }, // 默认值 
      success (res) {
        // console.log(res)
      },
      fail: err => { 
        console.log('课程删除错误', err) 
      }
    })
  },
  setAddCourse: function(course){
    let entityList = this.data.courseEntityList;
    // 检查唯一性
    for(let entity of entityList){
      if(entity.CourseId === course.CourseId){
        if(app.globalData.isAddingCourse === true){
          Toast({
            type: 'fail',
            message: '不能重复添加',
            duration: 1000
          });
          app.globalData.isAddingCourse = false;
        }
        return;
      }
    }

    if(this.checkAddCourseValid(course.CourseDayIdxList, course.CourseTimeIdxList)){
      // 插入课程至指定日期的表
      var todayIdx = new Date().getDay();
      for(let [idx, dayIdx] of course.CourseDayIdxList.entries()){
        var timeList = course.CourseTimeIdxList[idx];
        var timeCol = this.data.timeColumns;
        if(dayIdx == todayIdx)
          course.CourseTimeTag = [timeCol[timeList[0]], timeCol[timeList.slice(-1)]];
        this.data.WeekCourseList[dayIdx].push(course)
      }
      // 更新前端数据
      this.setData({ courseEntityList: this.data.WeekCourseList[todayIdx] });
      if(app.globalData.isAddingCourse === true){
        Toast({
          type: 'success',
          message: '添加课程成功',
          duration: 1100
        });
        app.globalData.isAddingCourse = false;
      }
    }
    
  },
  addCourse: function(){
    while(app.globalData.userCourseList.length > 0){
      var course = app.globalData.userCourseList.splice(0, 1);
      this.setAddCourse(course[0]);
    }
    wx.hideLoading();
  },
  checkAddCourseValid: function(dayIdxList, timeIdxList){
    if(dayIdxList.length != timeIdxList.length)
      return false;
    var len = dayIdxList.length;
    var weekList = this.data.WeekCourseList;
    for(let i = 0; i < len; ++i){
      let eList = weekList[dayIdxList[i]];
      for(let entity of eList){
        for(let courseTimeList of entity.CourseTimeIdxList){
          if(!this.checkTime(courseTimeList, timeIdxList[i])){
            if(app.globalData.isAddingCourse){
              Toast({
                type: 'fail',
                message: '课程时间冲突',
                duration: 1000
              });
            }
            return false;
          }
        }
      }
    }
    return true;
  },
  checkTime: function(timeList1, timeList2){
    // 判断两个有序数组之间是否有重复元素，使用 Merge 思想，O(n)
    var i = 0, j = 0;
    var limit = timeList1.length > timeList2.length ? timeList2.length : timeList1.length;
    for(let k = 0; k < limit; ++k){
      if(timeList1[i] > timeList2[j])
        j++;
      else if(timeList1[i] < timeList2[j])
        i++;
      else
        return false; // 相等
    }
    return true;
  }
})