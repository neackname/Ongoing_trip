// app.js
const util = require('./common/utils/util')

App({
  onLaunch() {
    // 展示本地存储能力
    const logs = wx.getStorageSync('logs') || []
    logs.unshift(Date.now())
    wx.setStorageSync('logs', logs)
  },

  globalData: {
    userInfo: {
      avatar:'https://img0.baidu.com/it/u=3204281136,1911957924&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=500',
      nickname:'微信用户',
      motto:'个性签名',
      gender:'未填',
      telephone:'未填',
    },
//--------
    isLogin: false, // 全局登录状态
    token: null,    // 用户 token
  },

  checkLogin() {
    const token = util.getToken()
    if (token) {
      this.globalData.isLogin = true;
      this.globalData.token = token;
    } else {
      this.globalData.isLogin = false;
      this.globalData.token = null;
    }
    return this.globalData.isLogin;
  },
})
