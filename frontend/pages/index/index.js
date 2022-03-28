// index.js
// 获取应用实例
const app = getApp()

Page({
  data: {
    circleSz: 90, // 圆环大小
    volumeData: {},
    buildingList: app.globalData.buildingList,
    upboundData: {
      H2: 32, H3NoH30: 30, H4: 30,
      H5: 63, H6: 51, HGX: 46
    },
    gradientColor: {
      '0%': '#6a89cc',
      '100%': '#3B3B98',
    },
    userInfo: {},
    hasUserInfo: false,
    canIUse: wx.canIUse('button.open-type.getUserInfo'),
    canIUseGetUserProfile: false,
    canIUseOpenData: wx.canIUse('open-data.type.userAvatarUrl') && wx.canIUse('open-data.type.userNickName') // 如需尝试获取用户信息可改为false
  },
  // 事件处理函数
  bindViewTap() {
    wx.navigateTo({
      url: '../logs/logs'
    })
  },
  onLoad() {
    if (wx.getUserProfile) {
      this.setData({
        canIUseGetUserProfile: true
      })
    }
    wx.showLoading({ title: '加载中', })

    /* Load Font */
    this.loadFontFunc();
    /* init Classroom Volume */
    this.getCurEmptyCls();
    /* Refresh Classroom Volume */
    setInterval(this.getCurEmptyCls, 5000);
  },
  onShow() {
    if (typeof this.getTabBar === 'function' && this.getTabBar())
      this.getTabBar().setData({ active: 0 })
  },
  ReCommend(){
    wx.navigateTo({
      url: './recommend/index',
    })
  },
  CheckCls(){
    wx.navigateTo({
      url: './check_classroom/index',
    })
  },
  loadFontFunc: function(){
    wx.loadFontFace({
      family: 'ZCOOL_h',
      source: 'url(../../../../src/font/ZCOOL_Kuaile/Zcool_h.ttf)',
      success: res => {
        //console.log('font load success', res)
      },
      fail: err => {
        console.log('font load fail', err)
      }
    })
  },
  getCurEmptyCls: function(){
    let reqUrl = app.globalData.baseUrl + 'empty_classroom';
    var that = this;
    for(let building of this.data.buildingList){
      wx.request({
        url: reqUrl,
        method: 'GET',
        dataType: 'json',
        timeout: 2000,
        data: { building: building },
        header: { 'content-type': 'application/json' }, // 默认值 
        success (res) {
            that.setVolumeData(building, res.data.num)
            if(building == 'HGX')
              wx.hideLoading()
        },
        fail: err => { 
          console.log('API GET fail', err) 
        }
      })
    }
  },
  setVolumeData: function(building, num) {
    this.data.volumeData[building] = num;
    this.setData({
      volumeData: this.data.volumeData
    })
  },
  checkEmpty: function(e){
    let reqUrl = app.globalData.baseUrl + 'empty_classroom';
    console.log(reqUrl)
    wx.request({
      url: reqUrl,
      method: 'GET',
      dataType: 'json',
      timeout: 2000,
      data: { building: e.target.id, },
      header: { 'content-type': 'application/json' }, // 默认值 
      success (res) {
        console.log(res.data)
        wx.navigateTo({
          url: './empty_classroom/index?eclsList=' + res.data.empty_classroom + '&building=' + e.target.id,
        })
      },
      fail: err => { console.log('API GET fail', err) }
    })
  }
})
