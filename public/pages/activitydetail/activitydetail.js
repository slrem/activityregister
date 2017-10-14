// activitydetail.js
var app = getApp();
Page({
  data: {
    userInfo: {},
    activityDetails:{},
    activityId:null
  },

  joinActive:function(event){
    var that = this;
    wx.request({
      url: `${app.globalData.domainName}/api/joinactivity`,
      data: {
        activityid: event.currentTarget.id,
        uid: that.data.userInfo.uid,
        session: that.data.userInfo.session,
        formid: event.detail.formId
      },
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      method: 'POST',
      success: function (res) {
        if (res.data) {
          if (res.data.Ret == 0) {
            that.updataData();
          }else{
            app.getErrorBox(res.data.Msg)
          }
        }else{
          app.getErrorBox()
        }
      }
    })
  },
  onPullDownRefresh: function () {
    var that = this;
    that.updataData();
  },
  cancelJoin: function (event){
    var that = this;

    wx.request({
      url: `${app.globalData.domainName}/api/canceljoinactivity`,
      data: {
        activityid: event.currentTarget.id,
        uid: that.data.userInfo.uid,
        session: that.data.userInfo.session,
      },
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      method: 'POST',
      success: function (res) {
        if (res.data) {
          if (res.data.Ret == 0) {
            that.updataData();
          } else {
            app.getErrorBox(res.data.Msg)
          }
        } else {
          app.getErrorBox()
        }
      }
    })
  },
  updataData:function (){
    var that = this;

    wx.request({
      url: `${app.globalData.domainName}/api/getactivitybyid`,
      data: {
        activityid: that.data.activityId,
        uid: that.data.userInfo.uid,
        session: that.data.userInfo.session
      },
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      method: 'POST',
      success: function (res) {
        wx.stopPullDownRefresh()
        if (res.data) {
          if (res.data.Ret == 0) {
            var nowGap = new Date() - new Date(res.data.Data.ActivityStartTime.replace(/\-/g, "/"))
            res.data.Data.nowGap = nowGap;
            that.setData({
              activityDetails: res.data.Data
            })
          } else {
            app.getErrorBox(res.data.Msg)
          }
        } else {
          app.getErrorBox()
        }
      }

    })
  },

  /* 生命周期函数--监听页面加载*/
  onLoad: function (options) {
    var that = this;

    //调用应用实例的方法获取全局数据
    app.getUserInfo(function (userInfo) {
      //更新数据
      that.setData({
        userInfo: userInfo,
        activityId: options.activityId
      })

      that.updataData();
    })
  },


  /* 用户点击右上角分享*/
  onShareAppMessage: function (res) {
    var that = this;
    if (res.from === 'button') {
      // 来自页面内转发按钮
    }
    return {
      title: '活动报名',
      path: `/pages/activitydetail/activitydetail?activityId=${that.data.activityId}`,
      success: function (res) {
        // 转发成功
      },
      fail: function (res) {
        // 转发失败
      }
    }
  }
})