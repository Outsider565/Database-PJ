// pages/index/check_classroom/index.js
import Toast from '../../../src/dist/toast/toast';
const app = getApp()

Page({

  /**
   * 页面的初始数据
   */
  data: {
    ClsList: Array(13).fill({CourseName: "空闲教室"}),
    timeColumns: ["早八1", "早八2", "第三节", "第四节", "第五节", "下午1", "下午2", "下午3", "下午4", "晚课1", "晚课2", "晚课3", "晚课4"],
    isSearch: false,
    searchValue: ""
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {

  },
  onSearch: function(e){
    let RoomId = e.detail
    if(!RoomId.length){
      Toast({
        type: 'fail',
        message: '搜索不能为空',
        duration: 1000
      });
      return;
    }
    this.setData({ searchValue: RoomId })

    // 开始搜索
    this.searchCourse(RoomId);
    console.log(this.data.ClsList)
  },
  onSearchClear: function(){
    this.setData({ isSearch: false });
  },
  searchCourse: function(RoomId){
    let reqUrl = app.globalData.baseUrl + 'get_classes';
    var _this = this;
    wx.request({
      url: reqUrl,
      method: 'GET',
      dataType: 'json',
      timeout: 2000,
      data: { 
        date: _this.formatter(new Date().getTime()),
        classroom_no: RoomId
      },
      header: { 'content-type': 'application/json' }, // 默认值 
      success (res) {
        console.log(res)
        let Res = res.data.res;
        if(Res === null){
          Toast({
            type: 'fail',
            message: '不存在的教室',
            duration: 1000
          });
          return 
        }

        let clsList = _this.data.ClsList;
        for(let i = 0; i < Res.length; ++i)
          clsList[Res[i].ClassIndex - 1] = Res[i];
        _this.setData({ 
          ClsList: clsList,
          isSearch: true
        });
      },
      fail: err => { console.log('API GET fail', err) }
    })
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