// pages/index/empty_classroom/index.js
const app = getApp()
import Dialog from '../../../src/dist/dialog/dialog';
import Toast from '../../../src/dist/toast/toast';

Page({
  /**
   * 页面的初始数据
   */
  data: {
    timePickerShow: false,
    currentTime: new Date().getTime(),
    currentDate: new Date().getDate(),
    maxDate: new Date().setDate(new Date().getDate() + 14),
    // floorList
    floorList: ["全部", "一楼", "二楼", "三楼", "四楼", "五楼"],
    floorDict: {
      "H2": 3, "H3NoH30": 4, "H4": 5,
      "H5": 4, "H6": 5, "HGX": 5
    },
    floorRadio: 0,
    // timeList
    timeList: ["全部", "上午", "下午", "晚间"],
    timeRadio: 0,
    // datePicker
    datePickerShow: false,
    date: '',
    // checkBox
    clsChoocedIdx: 0,
    clsCheckBox: Array(13).fill(false),
    tagClass: '',
    // Empty Cls by Index
    indexList: [0],
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    var date = new Date(this.data.currentTime);
    var floorList = this.data.floorList;
    var len = this.data.floorDict[options.building] + 1;
    var emptyList = options.eclsList.split(',');
    var tagClass = 'tags'
    if(options.building === 'HGX')
      tagClass = 'htags'

    this.setData({
      building: options.building,
      empty_classroom: emptyList,
      floorList: floorList.slice(0, len),
      date: this.formatDate(date),
      date_query: this.formatter(date),
      tagClass: tagClass
    });
    this.loadFontFunc();

    // prefetch all day empty
    this.prefetchEmpty(null);
  },
  loadFontFunc: function(){
    wx.loadFontFace({
      family: 'ZCOOL_h',
      source: 'url(../../../../../src/font/ZCOOL_Normal/Zcool-W02.ttf)',
      success: res => {
        console.log('font load success', res)
      },
      fail: err => {
        console.log('font load fail', err)
      }
    })
  },
  onDisplayDatePicker: function(e) {
    this.setData({ datePickerShow: true });
  },
  onDatePickerClose: function(e) {
    this.setData({ datePickerShow: false });
  },
  onDatePickerConfirm: function(e){
    this.setData({
      datePickerShow: false,
      date: this.formatDate(e.detail),
      date_query: this.formatter(e.detail)
    });
    this.getEmptyClassroom(this.data.indexList, e.detail, this.data.floorRadio);
  },
  onFloorRadioChoose: function(e) {
    let targetFloor = e.detail;
    this.setData({ floorRadio: targetFloor });
    this.getEmptyClassroom(this.data.indexList, this.data.date_query, targetFloor);
  },
  ontimeRadioChoose: function(e) {
    // 节次选择和时段选择互斥
    this.setData({ 
      timeRadio: 0,
      clsCheckBox: Array(13).fill(false)
    })

    let timeRangeIdx = e.detail;
    this.setData({ timeRadio: timeRangeIdx });
    if(timeRangeIdx === 0){  // 全天
      this.getEmptyClassroom([0], this.data.date_query, this.data.floorRadio);
      this.setData({ indexList: [0]});
    } else if(timeRangeIdx === 1) { // 上午
      this.getEmptyClassroom([1,2,3,4,5], this.data.date_query, this.data.floorRadio)
      this.setData({ indexList: [1,2,3,4,5]});
    } else if(timeRangeIdx === 2) { // 下午
      this.getEmptyClassroom([6,7,8,9,10], this.data.date_query, this.data.floorRadio)
      this.setData({ indexList: [6,7,8,9,10]});
    } else if(timeRangeIdx === 3){ // 晚上
      this.getEmptyClassroom([11,12,13], this.data.date_query, this.data.floorRadio)
      this.setData({ indexList: [11,12,13]});
    }
  },
  onClsCheckBoxChoose: function(e) {
    // 节次选择和时段选择互斥
    this.setData({ timeRadio: 0 })

    var tmpList = this.data.clsCheckBox;
    var idx = parseInt(e.detail);
    var indexList = [];
    tmpList[idx] = !tmpList[idx]; // toggle
    for(let i = 0; i < 12; ++i){
      if(tmpList[i])
        indexList.push(i);
    }
    this.setData({ 
      clsCheckBox: tmpList,
      indexList: indexList
    });
    this.getEmptyClassroom(indexList, this.data.date_query, this.data.floorRadio);
  },
  checkIndex: function(idx){
    return true
  },
  chooseTimeTap: function() {
    this.setData({timePickerShow: true})
  },
  getEmptyClassroom: function(indexList, date_query, floor){
    if(date_query == null)
      date_query = this.data.currentTime
    let index_query = indexList.toString();

    // 先向缓存请求，极联 API，这个方法需要设计缓存，并且需要进行 evict
    this.getEmptyListFromCache(index_query, date_query, floor)
  },
  getEmptyListFromCache: function(indexList, date_query, floor){
    var date_query = this.formatter(date_query);
    var _this = this;
    wx.getStorageInfo({
      success (res) {
        _this.setData({storageSize: res.currentSize})
      }
    })
    let idx_query = indexList.toString();

    try {
      // 从缓存中读取
      let key = this.building + '&' + idx_query + '&' + date_query;
      let res = wx.getStorageSync(key);
      if(!res || this.data.storageSize > 2048)  // 若本地未缓存或者缓存超限制，尝试从服务器api获取
        this.requestEmptyClsByIdx(indexList, date_query, floor)
      else{
        // 筛选楼层
        var Res = res;
        if(floor && floor != 0){  
          var tmpList = [];
          var floorStrIdx;
          if(this.data.building == 'H3NoH30')
            floorStrIdx = 2;
          else  
            floorStrIdx = this.data.building.length
          for(let i = 0; i < Res.length; ++i){
            if(parseInt(Res[i][floorStrIdx]) === floor)
              tmpList.push(Res[i]);
          }

          this.setData({ empty_classroom: tmpList })
        }
        else 
          this.setData({ empty_classroom: Res })
        // 设置最终
        
      }
    } catch (e) {
      console.log("Get Empty Class By Index From Cache Failed\n", e)
      return NULL
    }
  },
  requestEmptyClsByIdx: function(indexList, date_query, floor){
    let reqUrl = app.globalData.baseUrl + 'empty_classroom';
    var _this = this;
    var idx_query = indexList.toString();
    wx.request({
      url: reqUrl,
      method: 'GET',
      dataType: 'json',
      timeout: 2000,
      data: {
        building: this.data.building,
        date: date_query,
        index: idx_query
      },
      header: { 'content-type': 'application/json' }, // 默认值 
      success (res) {
          // 插入缓存
          let key = _this.building + '&' + idx_query  + '&' + date_query;
          wx.setStorageSync(key, res.data.empty_classroom);
          
          // 筛选楼层
          var Res = res.data.empty_classroom;
          if(floor && floor != 0){  
            var tmpList = [];
            var floorStrIdx;
            if(_this.data.building == 'H3NoH30')
              floorStrIdx = 2;
            else  
              floorStrIdx = _this.data.building.length
            for(let i = 0; i < Res.length; ++i){
              if(parseInt(Res[i][floorStrIdx]) === floor)
                tmpList.push(Res[i]);
            }
            _this.setData({ empty_classroom: tmpList })
          }
          else 
            _this.setData({ empty_classroom: Res })
      },
      fail: err => { 
        console.log('API GET fail', err) 
      }
    })
  },
  prefetchEmpty: function(date){
    this.getEmptyClassroom([0], date, this.data.floorRadio);
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
  },
  formatDate(date) {
    date = new Date(date);
    return `${date.getYear() + 1900} 年 ${date.getMonth() + 1} 月 ${date.getDate()} 日 `;
  },
  formatter(date) {
    var d = new Date(date),
    month = '' + (d.getMonth() + 1),
    day = '' + d.getDate(),
    year = d.getFullYear();
   
    if (month.length < 2) month = '0' + month;
    if (day.length < 2) day = '0' + day;
   
    return [year, month, day].join('-');
  }
})