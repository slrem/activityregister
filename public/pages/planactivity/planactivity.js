// planactivity.js
var app = getApp();

Page({
  data: {
    activedate:'',
    activetime:'',
    emptyTitle:''
  },
  formSubmit:function (e) {
    var that= this;

    var acms = that.data.activedate + " " + that.data.activetime;
    var date = new Date();
    var activedate = `${date.getFullYear()}-${(date.getMonth() + 1) < 10 ? '0' + (date.getMonth() + 1) : (date.getMonth() + 1)}-${date.getDate()}`,
      activetime = `${date.getHours()}:${date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()}`;
    var now = activedate + " " + activetime;

    if (e.detail.value.activityname.trim()==''){
      that.setData({
        emptyTitle: '标题不得为空'
      })    
      return;
    }else{
      that.setData({
        emptyTitle: ''
      })
    }

    wx.request({
      url: `${app.globalData.domainName}/api/createactivtty`, //仅为示例，并非真实的接口地址
      data: {
        title: e.detail.value.activityname,
        descriptiom: e.detail.value.activitycontainer,

        act_start_time: acms,
        act_end_time: acms,
        join_start_time: now,
        join_end_time: acms,

        uid: that.data.userInfo.uid,        
        session: that.data.userInfo.session
      },
      method:'POST',
      header: {
        'content-type': 'application/x-www-form-urlencoded'
      },
      success: function (res) {
        if (res.data){
          if (res.data.Ret==0){
            console.log('成功')
            wx.switchTab({
              url: '/pages/issueactivity/issueactivity'
            })
          }else{
            app.getErrorBox(res.fata.Msg)
          }
        } else {
          app.getErrorBox()
        }
      }
    })
   
  },
  bindActivedateChange: function (e) {
    this.setData({
      activedate: e.detail.value
    })
  },
  bindActiveTimeChange: function (e) {
    this.setData({
      activetime: e.detail.value
    })
  },
  onLoad: function (options) {
    var that = this;
    //调用应用实例的方法获取全局数据
    app.getUserInfo(function (userInfo) {
      //更新数据
      that.setData({
        userInfo: userInfo,
        activityId: options.activityId
      })
    })
    var date = new Date();
    var activedate = `${date.getFullYear()}-${(date.getMonth() + 1) < 10 ? '0' + (date.getMonth() + 1) : (date.getMonth() + 1)}-${date.getDate()}`,
      activetime = `${date.getHours()}:${date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()}`;
    that.setData({
      activedate: activedate,
      activetime: activetime
    })
  },
})