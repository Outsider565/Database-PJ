// pages/user/index.js
import Toast from '../../src/dist/toast/toast';
var app = getApp();
import { parseFormat, parseTimeData } from './utils';
Page({

  /**
   * 页面的初始数据
   */
  data: {
    isLogin: app.globalData.isLogin,
    isStudy: app.globalData.isStudy,
    userName: "",
    userCls: "",
    userImgsrc: "../../src/icon/userImg/user.png" ,
    studyTime: null,
    timeData: {},
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    // 加载字体
    wx.loadFontFace({
      family: 'ZCOOL_W05',
      source: 'url(../../../../src/font/ZCOOL_Normal/Zcool-W05.ttf)',
      success: res => {
       // console.log('font load success', res)
      },
      fail: err => {
        console.log('font load fail', err)
      }
    });
    wx.loadFontFace({
      family: 'ZCOOL_h',
      source: 'url(../../../../src/font/ZCOOL_Kuaile/Zcool_h.ttf)',
      success: res => {
       // console.log('font load success', res)
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
        active: 3
      })
    }
    // 总学习时长
    if(app.globalData.isLogin)
      this.setTotalStudyTime(app.globalData.userToken);
      
    this.setData({
      isStudy: app.globalData.isStudy,
      userCls: app.globalData.curCls
    });
  },

  onChange(e) {
    this.setData({
      timeData: e.detail,
    });
  },

  /**
   * 获取用户信息的回调函数
   */
  clickLogin(e){
    wx.clearStorage()
    //console.log(JSON.parse(e.detail.rawData))
    if(!e.detail.rawData){
      Toast.fail("登录失败");
      return;
    }
    app.globalData.userInfo=JSON.parse(e.detail.rawData);
    
    var _this = this;
    wx.login({
      success(res) {
        _this.getBackAuth(res.code, app.globalData.userInfo);
        if (res.code) {
          let userInfo = app.globalData.userInfo;
          // 设置信息
          app.globalData.isLogin = true;
          _this.setData({ isLogin: app.globalData.isLogin, })
          _this.setData({
            userName: userInfo.nickName,
            userImgsrc: userInfo.avatarUrl,
          })
          Toast.success("登录成功");
        } else {
          console.log('登录失败！' + res.errMsg)
          Toast.fail("登录失败");
        }
      }
    })
  },
  // 获取后端验证
  getBackAuth: function(code, userInfo){
    var post_data = {
      code: code,
      avatar_url: userInfo.avatarUrl,
      nick_name: userInfo.nickName,
      gender: userInfo.gender,
      city: userInfo.city,
      province: userInfo.province,
      country: userInfo.country
    }
    console.log(post_data)
    let reqUrl = app.globalData.authUrl + 'wx';
    var _this = this;
    wx.request({
      url: reqUrl,
      method: 'POST',
      timeout: 2000,
      data: post_data,
      header: { 
        'content-type': 'application/json',
        'data': JSON.stringify(post_data)
      },
      success (res) {
        app.globalData.userToken = res.data.token
        _this.fetchCourse(res.data.token);
        _this.setTotalStudyTime(res.data.token);
      },
      fail: err => { 
        console.log('API GET fail', err) 
      }
    })
  },
  fetchCourse: function(token){
    let reqUrl = app.globalData.pUrl + 'all_course';
    var _this = this;
    wx.request({
      url: reqUrl,
      method: 'GET',
      timeout: 2000,
      data: {
        year: new Date().getFullYear()
      },
      header: { 
        'content-type': 'application/json',
        'token': token
      }, // 默认值 
      success (res) {
        app.globalData.userCourseList = res.data.courses;
        _this.processCourse(app.globalData.userCourseList);
      },
      fail: err => { 
        console.log('用户课程获取错误', err) 
      }
    })
  },
  processCourse: function(courseList){
    let reqUrl = app.globalData.baseUrl + 'course_time_loc';
    var weekDict = app.globalData.weekDict
    var dayDict = app.globalData.dayDict
    for(let i = 0; i < courseList.length; ++i){
      var curCourse = courseList[i];
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

          app.globalData.userCourseList[i] = curCourse;
        },
        fail: err => { 
          console.log('API GET fail', err) 
        }
      })
    }
  },
  // 停止学习
  stopStudy: function() {
    app.globalData.isStudy = false;
    this.setData({isStudy: false});
    this.sendStudySig2Db();
    this.setTotalStudyTime(app.globalData.userToken);
    Toast.success("时长已计入");
  },
  sendStudySig2Db: function(){
    let reqUrl = app.globalData.pUrl + 'end_study';
    wx.request({
      url: reqUrl,
      method: 'GET',
      timeout: 2000,
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
  setTotalStudyTime: function(token){
    let reqUrl = app.globalData.pUrl + 'study_time';
    var _this = this;
    wx.request({
      url: reqUrl,
      method: 'GET',
      timeout: 2000,
      header: { 
        'content-type': 'application/json',
        'token': token
      }, // 默认值 
      success (res) {
        let time = res.data.total_time;
        if(time < 60)
          _this.setData({ studyTime: '不足一分钟' })
        else
          _this.setData({ studyTime: parseFormat("HH 时 mm 分", parseTimeData(time * 1000)) })
        _this.setRate(time);
      },
      fail: err => { 
        console.log('API GET fail', err) 
      }
    })
  },
  setRate: function(time){
    var hour = time / 3600;
    if(hour < 1)
        this.setData({ Rate: '寝室自习王' });
      else if(hour < 3)
        this.setData({ Rate: '拿A突击手' });
      else if(hour < 10)
        this.setData({ Rate: '黄金自习手' });
      else if(hour < 50)
        this.setData({ Rate: '三教久留' });
      else if(hour < 100)
        this.setData({ Rate: '三教王者' });
      else if(hour < 250)
        this.setData({ Rate: '轻松保研' });
      else
        this.setData({ Rate: 'GPA 4.0' });
  },
  handleUserInfoClick: function() {
    Toast('暂未实现，敬请期待');
  },
  handleSettingClick: function() {
    Toast('暂未实现，敬请期待');
  },
  handleAboutClick: function() {
    Toast('王少文、谢子飏开发');
  }
})