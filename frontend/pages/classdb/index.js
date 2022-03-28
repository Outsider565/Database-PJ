import Toast from "../../src/dist/toast/toast";
const app = getApp()

// pages/classdb/index.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    searchValue: '',
    showRelate: false,
    showClsDetail: false,
    isSearch: false,
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.setData({ userToken: app.globalData.userToken })
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

  /**
   * 生命周期函数--监听页面显示
   */
  onShow: function () {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({
        active: 1
      })
    }
  },

  /**
   * 处理搜索聚集
   */
  onSearch: function() {
    this.onSearchChange({detail: this.data.searchValue});
  },
  /**
   * 处理搜索清除
   */
  onSearchClear: function() {
    this.setData({
      showRelate: false
    });
  },
  /**
   * 处理搜索改变，进行模糊搜索
   */
  onSearchChange: function(e) {
    var searchKey = e.detail
    // 清空搜索框
    if(searchKey == ""){
      this.setData({ showRelate: false });
      return;
    }
    this.setData({
      showRelate: true,
      searchValue: searchKey
    });

    if(this.checkSearchType(searchKey) == 'CourseName')
      this.getClsByCourseName(searchKey)
    else if(this.checkSearchType(searchKey) == 'CourseId')
      this.getClsByCourseId(searchKey)
  },
  getClsByCourseName: function(course_name){
    let reqUrl = app.globalData.baseUrl + 'course_search';
    var _this = this;
    wx.request({
      url: reqUrl,
      method: 'GET',
      dataType: 'json',
      timeout: 2000,
      data: {
        max_num: 12,
        course_name: course_name
      },
      header: { 'content-type': 'application/json' }, // 默认值 
      success (res) {
          let Res = res.data.courses;
          if(!Res)
            return // 无结果
          for(let i = 0; i < Res.length; ++i){
            if(!Res[i].TeacherName || Res[i].TeacherName.length == 0 || Res[i].TeacherName.length > 12)
              Res[i].TeacherName = "暂无"
          }
          _this.setData({ relatedClsList: Res })
          return Res
      },
      fail: err => { 
        console.log('API GET fail', err) 
      }
    })
  },
  getClsByCourseId: function(courseId){
    let reqUrl = app.globalData.baseUrl + 'course_search';
    var _this = this;
    wx.request({
      url: reqUrl,
      method: 'GET',
      dataType: 'json',
      timeout: 2000,
      data: {
        max_num: 12, 
        course_id: courseId
      },
      header: { 'content-type': 'application/json' }, // 默认值 
      success (res) {
          _this.setData({ relatedClsList: res.data.courses })
          return res.data.courses
      },
      fail: err => { 
        console.log('API GET fail', err) 
      }
    })
  },
  onCellTap: function(e){
    // 获取上课时间和地点
    let reqUrl = app.globalData.baseUrl + 'course_time_loc';
    let idx = parseInt(e.currentTarget.id);
    let curCourse = this.data.relatedClsList[idx];
    var _this = this;
    var weekDict = app.globalData.weekDict
    var dayDict = app.globalData.dayDict
    wx.request({
      url: reqUrl,
      method: 'GET',
      timeout: 2000,
      data: {
        course_id: curCourse.CourseId,
        type: 0
      },
      header: { 'content-type': 'application/json' }, // 默认值 
      success (res) {
        var Res = res.data.res;
        var course_day_idx_list = [], course_time_idx_list = [];
        var course_time_list = [], location_list = [];
        // 解析返回的时间
        for(let day in Res){
          if(Res[day].Indexes){
            course_day_idx_list.push(dayDict[day]);
            course_time_idx_list.push(
              Res[day].Indexes.map(function (index) { return index - 1})
            )
            let dayStr = weekDict[day] + ' ' + Res[day].Indexes.join('节 ') + '节';
            if(!course_time_list || dayStr != course_time_list[0])
              course_time_list.push(dayStr);
            if(!location_list || Res[day].Location != location_list[0])
              location_list.push(Res[day].Location);
          }
        }
        // 课程上课日 Index 列表
        curCourse.CourseDayIdxList = course_day_idx_list;
        // 课程上课时间 Index 列表
        curCourse.CourseTimeIdxList = course_time_idx_list;
        // 课程上课时间 String 显示
        curCourse.CourseTimeStr = course_time_list.join(', ');
        // 课程上课地点 String 显示
        curCourse.Location = location_list.join(' ');
        // 设置课程类型
        curCourse.__type__ = 'Course'

        // 设置数据，展示detail界面
        _this.setData({
          curCls: curCourse,
          showClsDetail: true
        })
      },
      fail: err => { 
        console.log('API GET fail', err) 
      }
    })
  },
  hideClsDetail: function(){
    this.setData({ 
      showClsDetail: false,
    })
  },
  addCls2DB: function(course){
    let reqUrl = app.globalData.pUrl + 'take_course';
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
        console.log(res)
      },
      fail: err => { 
        console.log('API GET fail', err) 
      }
    })
  },
  onAddClsBtnTap: function(){
    if(app.globalData.isLogin === false){
      Toast({
        type: 'fail',
        forbidClick: true,
        message: '请先登录',
        duration: 800
      });
      setTimeout(() => {
        wx.switchTab({
          url: '/pages/user/index',
        })
      }, 800);
      return;
    }
    let data = this.data.curCls;
    
    app.globalData.userCourseList.push(data);
    app.globalData.isAddingCourse = true;

    this.addCls2DB(this.data.curCls);
    this.setData({ showClsDetail: false })
    wx.switchTab({ url: '/pages/curriculum/index' })
  },
  checkSearchType: function(searchKey){
    if((searchKey >= 'a' && searchKey <= 'z') || searchKey >= 'A' && searchKey <= 'Z')
      return 'CourseId';
    else 
      return 'CourseName';
  }
})