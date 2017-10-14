// issueactivity.js
var app = getApp();

Page({
  data: {
    userInfo:'',
    activities: [],
    initHeight:'180rpx',
    templActivities:'loadding',

    cancelbtnWidth:'200rpx'
  },

  hrefDetail: function (event) {
    var that = this;
    
    wx.navigateTo({
      url: `/pages/activitydetail/activitydetail?activityId=${event.currentTarget.id}`
    })
  },

  cancelActivity:function(event){
    var that = this;
    
    wx.request({
      url: `${app.globalData.domainName}/api/cancelactivity`,
      method: 'POST',
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      data: {
        activityid: event.currentTarget.id,
        uid: that.data.userInfo.uid,
        session: that.data.userInfo.session
      },
      success: function (res) {
        if (res.data) {
          if (res.data.Ret == 0) {

            var list = that.data.activities;
            for (var i in list) {
              if (list[i].ID == event.currentTarget.id) {
                  list[i].translateX = '0px';
                  list[i].State = -1;
              }
            }
            that.setData({
              activities: list
            });

          }else{
            app.getErrorBox(res.data.Msg);
          }

        }else{
          app.getErrorBox();
        }
      }
    })
  },

  activityRegisterclick: function (event) {
    wx.navigateTo({
      url: '/pages/planactivity/planactivity'
    })
  },
  onPullDownRefresh: function () {
    var that =this;
    that.updateData();
  },
  updateData:function(){
    var that =this;
    app.getUserInfo(function (userInfo) {
      //更新数据
      that.setData({
        userInfo: userInfo
      })

      wx.request({
        url: `${app.globalData.domainName}/api/getcreateactivity`, //仅为示例，并非真实的接口地址
        data: {
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
              if (res.data.Data == null) {
                that.setData({
                  activities: [],
                  templActivities: '你还没有活动呢，快快添加吧'
                })
              } else {
                for (var i in res.data.Data) {
                  res.data.Data[i].SatrtTime = `${res.data.Data[i].SatrtTime.slice(5, 7)}月${
                    res.data.Data[i].SatrtTime.slice(8, 10)}日${
                    res.data.Data[i].SatrtTime.slice(10, 16)}`;

                  res.data.Data[i].translateX = '0px';
                }


                that.setData({
                  activities: res.data.Data
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

    })
  },
  /*生命周期函数--监听页面加载*/
  onLoad: function (options) {
    var that =this;
  },
  onShow: function(){
    var that = this;
    //调用应用实例的方法获取全局数据
    that.updateData();
  },
  touchS: function (e) {
    if (e.touches.length == 1) {
      this.setData({
        startX: e.touches[0].clientX,
        startY: e.touches[0].clientY
      });
    }
  }, 
  touchM:function(e){
    var that = this;

    if (e.touches.length == 1) {
      var moveX = e.touches[0].clientX;
      var disX = that.data.startX - moveX;
      
      //获取删除按钮宽度
      var res = wx.getSystemInfoSync();
      var cancelbtnWidth = parseInt(that.data.cancelbtnWidth) * res.screenWidth/750;

      var dixY = e.touches[0].clientY - that.data.startY;
      var cancelbtnHeight = parseInt(that.data.initHeight) * res.screenWidth / 750;
      
    
      var translateX;
      if (disX <= cancelbtnWidth/2 ) {
        translateX = '0px';
      } else if (disX > 0 && dixY < cancelbtnHeight) { 
        translateX = `-${cancelbtnWidth}px`;
      }

      var list = that.data.activities;

      for(var i in list){
        if (list[i].State==0){
          if (i == e.currentTarget.dataset.index) {
            list[i].translateX = translateX;
          } else {
            list[i].translateX = '0px';
          }
        }
      }

      that.setData({
        activities: list
      });
    }

  }





})