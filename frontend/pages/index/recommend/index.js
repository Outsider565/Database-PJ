// pages/recommend/index.js
import Dialog from '../../../src/dist/dialog/dialog';
import Toast from '../../../src/dist/toast/toast';

const app = getApp()

Page({

  /**
   * 页面的初始数据
   */
  data: {
    buildingList: app.globalData.buildingList,
    // 教室基础分值，越高代表越乐色
    bias: {
      H2: 1.5, H3NoH30: 0.6, H4: 1.3,
      H5: 0.5, H6: 0.6, HGX: 1
    }
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.loadFontFunc();        // 加载字体
    this.getCurPos();         // 获取距离并计算推荐教室
  },
  onShow: function () {
    wx.showLoading({
      title: '加载中',
    })
  },
  loadFontFunc: function(){
    wx.loadFontFace({
      family: 'ZCOOL_h',
      source: 'url(../../../../../src/font/ZCOOL_Kuaile/Zcool_h.ttf)',
      success: res => {
        console.log('font load success', res)
      },
      fail: err => {
        console.log('font load fail', err)
      }
    })
  },
  // 获取当前位置信息
  getCurPos: function(){
    var _this = this;
    wx.getLocation({
      type: 'wgs84',
      isHighAccuracy: true,
      success (res) {
        let latitude = res.latitude
        let longitude = res.longitude
        _this.getDistance({latitude, longitude});
      },
      fail: err => { console.log('API GET fail', err) },
      complete (res) {
        //return null;
      }
      })
  },
  // 获取与教学楼距离
  getDistance(pos){
    var _this = this;
    //调用距离计算接口
    let reqUrl = app.globalData.baseUrl + 'distance';
    wx.request({
      url: reqUrl,
      method: 'GET',
      dataType: 'json',
      timeout: 2000,
      data: { 
        lat: pos.latitude,
        lon: pos.longitude
      },
      header: { 'content-type': 'application/json' }, // 默认值 
      success (res) {
          var resList = res.data.res;
          var dis = {};
          for (var i = 0; i < resList.length; i++) //将返回数据存入dis数组，
            dis[_this.data.buildingList[i]] = resList[i].distance;
          _this.setData({ distance: dis }) //设置并更新distance数据
          _this.generateRec();  // 教学楼打分
      },
      fail: err => { console.log('API GET fail', err) }
    })
  },
  generateRec: function(){
    var sortArray = []
    var distanceData = this.data.distance
    var bias = this.data.bias
    for (var key in distanceData) {
      sortArray.push([ key, distanceData[key], bias[key] ])
    }
    sortArray.sort(function compare(kv1, kv2) { return kv1[1] * kv1[2] - kv2[1] * kv2[2] })
    this.getEmptyClassroom(sortArray)
  },
  getEmptyClassroom: function(buildingArray){
    let reqUrl = app.globalData.baseUrl + 'empty_classroom';
    var _this = this;
    var emptyCls;
    var tmpList = [];

    for(var i = 0; i < 2; i++){
      wx.request({
        url: reqUrl,
        method: 'GET',
        dataType: 'json',
        timeout: 2000,
        data: { building: buildingArray[i][0] },
        header: { 'content-type': 'application/json' }, // 默认值 
        success (res) {
          emptyCls = res.data.empty_classroom;
        },
        fail: err => { console.log('API GET fail', err) },
        complete (res) {
          if(!emptyCls) // failed
            return;
          if(emptyCls.length >= 6)
            emptyCls = emptyCls.slice(0,6);
          tmpList.push.apply(tmpList,emptyCls)
          if(tmpList.length >= 8){
            _this.setData({ clsList: tmpList })
            wx.hideLoading()
          }
        }
      })
    }
  },
  startStudy: function(e){
    var class_room = e.target.id;
    Dialog.confirm({
      title: '开始自习',
      message: '是否要进入 ' + class_room + ' 自习',
      className: 'dialog',
      confirmButtonColor: '#4a69bd'
    })
      .then(() => {
        // on confirm 开始计时学习
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
        }
        else if(app.globalData.isStudy === true){
          Toast({
            type: 'fail',
            forbidClick: true,
            message: '已经在学习了\n    别卷了',
            duration: 1500
          });
        }
        else{
          this.sendStudySig2DB(class_room);
          Toast('开始学习~');
          app.globalData.isStudy = true;
          app.globalData.curCls = class_room
          setTimeout(() => {
            wx.switchTab({
              url: '/pages/user/index',
            })
          }, 1500);
        }
      })
      .catch(() => {
        // on cancel
      });
  },
  sendStudySig2DB: function(roomId){
    let reqUrl = app.globalData.pUrl + 'start_study';
    var _this = this;
    wx.request({
      url: reqUrl,
      method: 'GET',
      timeout: 2000,
      data: {
        classroom_no: roomId
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
  }
})