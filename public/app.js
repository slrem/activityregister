//app.js
var loginStatus = true;
App({
  onLaunch: function () {
    //调用API从本地缓存中获取数据
    var logs = wx.getStorageSync('logs') || []
    logs.unshift(Date.now())
    wx.setStorageSync('logs', logs)
  },
  getErrorBox:function(msg){
    msg=msg?msg:'';
    wx.showModal({
      content: '请求错误 ' + msg,
      showCancel: false,
      success: function (res) {
        if (res.confirm) {
          console.log('用户点击确定')
        } else if (res.cancel) {
          console.log('用户点击取消')
        }
      }
    })
  },
  getUserInfo: function (cb){
    var that = this
    if (this.globalData.userInfo){
      typeof cb == "function" && cb(this.globalData.userInfo)
    }else{

      if (!loginStatus){
        wx.openSetting({
          success: function (data) {
            if (data) {
              if (data.authSetting["scope.userInfo"] == true) {
                loginStatus = true;
                wx.getUserInfo({
                  withCredentials: false,
                  success: function (data) {
                    console.log("授权成功返回数据");
                  },
                  fail: function () {
                    console.info("授权失败返回数据");
                  }});
              }
            }
          },
          fail: function () {
            console.info("设置失败返回数据");
          }});
      }else{
        wx.login({
          success: function (res) {
            var code = res.code;

            wx.getUserInfo({
              withCredentials: true,
              success: function (res) {
                that.globalData.userInfo = res.userInfo

                wx.request({
                  url: `${that.globalData.domainName}/login`, //仅为示例，并非真实的接口地址
                  method: 'POST',
                  data: {
                    code: code,
                    encryptedData: res.encryptedData,
                    iv: res.iv,
                  },
                  header: {
                    'content-type': 'application/x-www-form-urlencoded'
                  },
                  success: function (res) {
                    if (res.data){
                      if (res.data.Ret == 0) {

                        that.globalData.userInfo.session = res.data.Data.Session;
                        that.globalData.userInfo.uid = res.data.Data.Uid;

                        typeof cb == "function" && cb(that.globalData.userInfo)

                      } else {
                        that.getErrorBox(res.data.Msg)
                      }
                    }else{
                      that.getErrorBox();
                    }
                  },
                  fail:function(){
                    that.getErrorBox('网络环境较差');
                  }
                })
              },
              fail: function () {
                loginStatus=false;
                wx.showModal({
                  content: '关闭授权后可能会影响使用小程序的部分功能，请确认',
                  confirmText:'关闭授权',
                  cancelText:'取消',
                  success: function (res) {
                    if (res.confirm) {
                      console.log('用户点击确定')
                    } else if (res.cancel) {
                      wx.openSetting({
                        success: function (data) {
                          if (data) {
                            if (data.authSetting["scope.userInfo"] == true) {
                              loginStatus = true;
                              wx.getUserInfo({
                                  withCredentials: false,
                                  success: function (data) {
                                    console.info("成功获取用户返回数据");
                                    console.info(data.userInfo);
                                  },
                                  fail: function () {
                                    console.info("授权失败返回数据");
                                  }
                                });
                            }
                          } 
                        },
                        fail: function () {
                          console.info("设置失败返回数据");
                        }});
                    }
                  }
                });

              }
            })
          },
          fail: function () {
            that.getErrorBox('网络环境较差');
          }
        })
      }

    }
  },

  globalData:{
    userInfo: null,
    // domainName: 'http://192.168.9.151:7777'
    domainName: 'https://hd.gomydodo.com'
  }
})