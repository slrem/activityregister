// myactivity.js
var app = getApp();
Page({
  data: {
    userInfo: {},
    activities: null,
    loadding:true
  },

  hrefDetail: function (event){
    wx.navigateTo({
      url: `/pages/activitydetail/activitydetail?activityId=${event.currentTarget.id}`
    })
  },
  updateData(){
    var that = this;
    wx.request({
      url: `${app.globalData.domainName}/api/getjoinactivity`,
      data: {
        uid: that.data.userInfo.uid,
        session: that.data.userInfo.session
      },
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      method: 'POST',
      success: function (res) {
        wx.stopPullDownRefresh(); 
        if (res.data) {
          if (res.data.Ret == 0) {
            if (res.data.Data == null) {
              that.setData({
                activities: null,
                loadding: false
              })
            } else {
              that.setData({
                activities: res.data.Data,
                loadding: false
              })
            }

          } else {
            app.getErrorBox(res.data.Msg);
          }
        } else {
          app.getErrorBox();
        }
      }
    })
  },
  /*生命周期函数--监听页面加载*/
  onLoad: function (options) {
    var that = this;
  },
  onShow: function () {
    var that = this;
    //调用应用实例的方法获取全局数据
    app.getUserInfo(function (userInfo) {
      //更新数据
      that.setData({
        userInfo: userInfo
      })
      that.updateData();
    })
  },
  onPullDownRefresh: function () {
    var that = this;
    that.updateData();
  }
})