// pages/map/map.js
// 引入SDK核心类
var QQMapWX = require('../../libs/qqmap-wx-jssdk.js');
var util = require('../../common/utils/util')
// 实例化API核心类
var qqmapsdk = new QQMapWX({
  key: 'RS2BZ-XGFEB-UWBU3-JRDH5-IMW3H-AJFQN' // 必填
});

Page({  
  data: { 
    latitude: 23.412840,  
    longitude: 116.633490,  
    currentLatitude: 0,
    currentLongitude: 0,
    markers: [],  
    
    dragPosition: 0, // 拖动窗口的初始位置
    startY: 0 ,// 触摸开始时的Y坐标
    initialDragPosition: 0, // 记录拖动开始时的初始位置+

    tabList: [
      {
        id: '1',
          title: '公共交通',
          src:'../../icon/大巴.png',
          url:''
      },
      {
        id: '2',
          title: '动车',
          src:'../../icon/动车.png',
          url:''
      },
      {
        id: '3',
          title: '景点',
          src:'../../icon/景点.png',
          url:''
      },
      {
        id: '4',
          title: '日程安排',
          src:'../../icon/日历.png',
          url:''
      },
      {
        id: '5',
          title: '酒店',
          src:'../../icon/酒店.png',
          url:''
      },
      {
        id: '6',
          title: '美食',
          src:'../../icon/夏日饮品.png',
          url:''
      },
      {
        id: '7',
          title: '购物',
          src:'../../icon/购物.png',
          url:''
      },
      {
        id: '8',
          title: '更多',
          src:'../../icon/泳圈.png',
          url:''
      }
    ],
    
    moodOptions: ["探索感", "放松感", "治愈感", "冒险感", "社交感", "浪漫感", "空白感", "回忆感"],
    moodIndex: 2,
    selectedMood: "治愈感",
    searchKeyword: "",
    searchSuggestList: [],
    recommendPlaces: [],
    isRecommendLoading: false,
    recommendError: "",
    polyline: [],
    includePoints: [],
    currentDestination: null,
    currentRouteResult: null,
    currentMode: "driving",
    routeDistance: 0,
    routeDuration: 0,
    routeInfoText: "",
    showSaveButton: false,
    
    scrollHeight: 0,
    isFull: false,
    scrollTop: 0,
    windowHeight: 0,
    isDragging: false
    
  },  
    
  onLoad() {  
//获取位置
    this.getLocation();
//--------------------------------
    const systemInfo = wx.getSystemInfoSync();
    const windowHeight = systemInfo.windowHeight; // 获取当前设备的窗口高度
    const initialPos = windowHeight - 100;
    this.setData({
    dragPosition: initialPos, // 让窗口初始位置在距离底部100px的位置
    initialDragPosition: initialPos, // 初始化初始位置+
    scrollHeight: windowHeight - 40,
    windowHeight: windowHeight,
    isFull: false,
    isDragging: false
  });
    
  },  

  onShow() {
    var footID = wx.getStorageSync('openFootID')
    if (footID) {
      wx.removeStorageSync('openFootID')
      this.loadFootByID(String(footID))
      return
    }

    var scenicMode = wx.getStorageSync('scenicMode')
    if (scenicMode) {
      wx.removeStorageSync('scenicMode')
      if (this.data.latitude && this.data.longitude) {
        this.loadScenicRecommend()
      } else {
        this._needScenic = true
      }
    }
  },

  onMoodChange(e) {
    var idx = Number(e.detail.value || 0);
    var mood = this.data.moodOptions[idx] || this.data.selectedMood;
    this.setData({
      moodIndex: idx,
      selectedMood: mood
    });
    this.refreshRecommend();
  },

  formatDuration(mins) {
    var m = Number(mins || 0)
    if (!m || m < 0) return ''
    if (m < 60) return `${Math.round(m)}分钟`
    var h = Math.floor(m / 60)
    var mm = Math.round(m % 60)
    if (mm === 0) return `${h}小时`
    return `${h}小时${mm}分钟`
  },

  formatDistance(meters) {
    var d = Number(meters || 0)
    if (!d || d < 0) return ''
    if (d < 1000) return `${Math.round(d)}m`
    var km = d / 1000
    return `${(Math.round(km * 10) / 10).toFixed(1)}km`
  },

  updateSaveButtonVisibility() {
    const visible = !!(this.data.polyline && this.data.polyline.length && !this.data.isFull)
    this.setData({ showSaveButton: visible })
  },

  onSaveFootTap() {
    this.saveCurrentFoot()
  },

  onSearchInput(e) {
    var keyword = (e && e.detail && e.detail.value) || ''
    this.setData({ searchKeyword: keyword })
  },

  onSearchConfirm() {
    var keyword = (this.data.searchKeyword || '').trim()
    if (!keyword) return
    if (!this.data.currentLatitude || !this.data.currentLongitude) {
      wx.showToast({ title: '请先开启定位权限', icon: 'none' })
      this.getLocation()
      return
    }
    this._confirmSearch = true
    this.fetchSuggest(keyword)
  },

  fetchSuggest(keyword) {
    var that = this
    keyword = (keyword || '').trim()
    if (!keyword) {
      that.setData({ searchSuggestList: [] })
      return
    }
    qqmapsdk.getSuggestion({
      keyword: keyword,
      location: {
        latitude: that.data.currentLatitude || that.data.latitude,
        longitude: that.data.currentLongitude || that.data.longitude
      },
      region: '全国',
      region_fix: 0,
      page_size: 10,
      success: function (_raw, data) {
        var list = (data && (data.suggestSimplify || data.suggestResult)) || []
        var normalized = list.map(function (x, idx) {
          return {
            id: x.id || String(idx),
            title: x.title || x.name || '',
            address: x.address || '',
            latitude: Number(x.latitude || (x.location && x.location.lat) || 0),
            longitude: Number(x.longitude || (x.location && x.location.lng) || 0),
          }
        }).filter(function (x) { return x.title && x.latitude && x.longitude })
        that.setData({ searchSuggestList: normalized })
        if (that._confirmSearch) {
          that._confirmSearch = false
          if (normalized.length) {
            that.applySuggestItem(normalized[0], true)
            that.showRouteAction(normalized[0])
          } else {
            wx.showToast({ title: '未找到目的地', icon: 'none' })
          }
        }
      },
      fail: function (err) {
        that._confirmSearch = false
        that.setData({ searchSuggestList: [] })
        var msg = (err && err.message) || '搜索失败'
        wx.showToast({ title: msg, icon: 'none' })
      }
    })
  },

  applySuggestItem(item, keepList) {
    if (!item) return
    this.setData({
      searchSuggestList: keepList ? (this.data.searchSuggestList || []) : [],
      searchKeyword: item.title,
    })

    var fromLat = this.data.currentLatitude || this.data.latitude
    var fromLng = this.data.currentLongitude || this.data.longitude
    var markers = [{
      id: 1,
      latitude: fromLat,
      longitude: fromLng,
      width: 50,
      height: 50,
      iconPath: '/icon/086_定位.png'
    }, {
      id: 2,
      latitude: item.latitude,
      longitude: item.longitude,
      width: 32,
      height: 32,
      iconPath: '/icon/旗帜.png'
    }]
    this.setData({
      markers: markers,
      currentDestination: { latitude: item.latitude, longitude: item.longitude, name: item.title, address: item.address },
      includePoints: [
        { latitude: fromLat, longitude: fromLng },
        { latitude: item.latitude, longitude: item.longitude },
      ],
      latitude: item.latitude,
      longitude: item.longitude,
    })
  },

  onTapSuggest(e) {
    var idx = Number(e.currentTarget.dataset.index || 0)
    var item = (this.data.searchSuggestList || [])[idx]
    if (!item) return
    this.applySuggestItem(item, false)
    this.showRouteAction(item)
  },

  showRouteAction(item) {
    var that = this
    wx.showActionSheet({
      itemList: ['查看地点', '驾车路线', '公交路线', '步行路线', '一键保存足迹'],
      success(res) {
        if (res.tapIndex === 0) {
          wx.openLocation({
            latitude: item.latitude,
            longitude: item.longitude,
            name: item.title || item.name || '',
            address: item.address || '',
            scale: 16
          })
          return
        }
        if (res.tapIndex === 1) {
          that.planRouteTo({ latitude: item.latitude, longitude: item.longitude, name: item.title || item.name || '' }, 'driving')
            .catch((err) => wx.showToast({ title: (err && err.message) || '规划失败', icon: 'none' }))
          return
        }
        if (res.tapIndex === 2) {
          that.planRouteTo({ latitude: item.latitude, longitude: item.longitude, name: item.title || item.name || '' }, 'transit')
            .catch((err) => wx.showToast({ title: (err && err.message) || '规划失败', icon: 'none' }))
          return
        }
        if (res.tapIndex === 3) {
          that.planRouteTo({ latitude: item.latitude, longitude: item.longitude, name: item.title || item.name || '' }, 'walking')
            .catch((err) => wx.showToast({ title: (err && err.message) || '规划失败', icon: 'none' }))
          return
        }
        if (res.tapIndex === 4) {
          that.planRouteTo({ latitude: item.latitude, longitude: item.longitude, name: item.title || item.name || '' }, 'driving')
            .then(() => that.saveCurrentFoot())
            .catch((err) => wx.showToast({ title: (err && err.message) || '保存失败', icon: 'none' }))
        }
      }
    })
  },

//--------------------------------获取位置 
  //初版
  // getLocation: function () {  
  //   const that = this;  
  //   wx.getLocation({  
  //     type: 'gcj02', // 返回可以用于wx.openLocation的经纬度  
  //     success: function (res) {  
  //       that.setData({  
  //         latitude: res.latitude,  
  //         longitude: res.longitude,  
  //         markers: [{  
  //           id: 1,  
  //           latitude: res.latitude,  
  //           longitude: res.longitude,  
  //           width: 50,  
  //           height: 50,  
  //           iconPath: '/resources/location.png' // 可选：自定义标记图标  
  //         }]  
  //       });  
  //     },  
  //     fail: function (error) {  
  //       console.error("获取位置失败：" + error.errMsg);  
  //       wx.showToast({  
  //         title: '获取位置失败',  
  //         icon: 'none'  
  //       });  
  //     }  
  //   });  
  // },

  //-----修改后加入授权
  getLocation: function () {
    const that = this;
  
    wx.getSetting({
      success(res) {
        const locationSetting = res.authSetting['scope.userLocation'];
        if (locationSetting) {
          that.fetchLocation();
          return;
        }

        wx.showModal({
          title: '位置权限',
          content: '需要获取您的位置信息，用于在地图上显示您的当前位置，是否允许？',
          confirmText: '允许',
          cancelText: '暂不',
          success(modalRes) {
            if (!modalRes.confirm) return;

            wx.authorize({
              scope: 'scope.userLocation',
              success() {
                that.fetchLocation();
              },
              fail() {
                wx.showModal({
                  title: '提示',
                  content: '未开启位置权限，请在设置中开启后再使用定位功能',
                  confirmText: '去设置',
                  cancelText: '取消',
                  success(res2) {
                    if (!res2.confirm) return;
                    wx.openSetting({
                      success(settingRes) {
                        if (settingRes.authSetting['scope.userLocation']) {
                          that.fetchLocation();
                        } else {
                          wx.showToast({
                            title: '未授权，无法获取位置',
                            icon: 'none'
                          });
                        }
                      }
                    });
                  }
                });
              }
            });
          }
        });
      },
      fail(err) {
        console.error('获取设置失败:', err);
        wx.showToast({
          title: '获取设置失败，请稍后重试',
          icon: 'none'
        });
      }
    });
  },
  
  // 封装获取地理位置的逻辑
  fetchLocation: function () {
    const that = this;
    wx.getLocation({
      type: 'gcj02', // 返回可以用于 wx.openLocation 的经纬度
      success(res) {
        that.setData({
          latitude: res.latitude,
          longitude: res.longitude,
          currentLatitude: res.latitude,
          currentLongitude: res.longitude,
          markers: [{
            id: 1,
            latitude: res.latitude,
            longitude: res.longitude,
            width: 50,
            height: 50,
            iconPath: '/icon/086_定位.png'
          }]
        });
        that.setData({ includePoints: [{ latitude: res.latitude, longitude: res.longitude }] })
        that.updateSaveButtonVisibility()
        if (that._needScenic) {
          that._needScenic = false
          that.loadScenicRecommend()
          return
        }
        that.refreshRecommend();
      },
      fail(error) {
        console.error("获取位置失败：" + error.errMsg);
        wx.showToast({
          title: '获取位置失败',
          icon: 'none'
        });
      }
    });
  },
  
  fetchRecommendKeywords(mood) {
    return util.request({
      url: '/travel/recommend',
      method: 'GET',
      data: { mood: mood },
    }).then(res => {
      if (!res || !res.data || res.data.code !== 200) {
        throw new Error((res && res.data && res.data.msg) || '获取推荐失败')
      }
      return res.data.data || {}
    })
  },

  qqSearch(keyword, region) {
    var that = this;
    return new Promise(function (resolve, reject) {
      qqmapsdk.search({
        keyword: keyword,
        location: {
          latitude: that.data.latitude,
          longitude: that.data.longitude
        },
        region: region || '全国',
        distance: 5000,
        page_size: 10,
        page_index: 1,
        success: function (_raw, data) {
          var list = (data && (data.searchSimplify || data.searchResult)) || [];
          resolve(list);
        },
        fail: function (err) {
          reject(err);
        }
      });
    });
  },

  loadScenicRecommend(keyword, retried) {
    var that = this
    var kw = keyword || '景点'
    if (!retried) {
      that.setData({ isRecommendLoading: true, recommendError: "", searchSuggestList: [] })
    }
    qqmapsdk.search({
      keyword: kw,
      location: {
        latitude: that.data.latitude,
        longitude: that.data.longitude
      },
      region: '全国',
      distance: 8000,
      page_size: 10,
      page_index: 1,
      success: function (_raw, data) {
        var list = (data && (data.searchSimplify || data.searchResult)) || []
        var userLat = that.data.latitude
        var userLng = that.data.longitude
        var mapped = list.map(function (x) {
          var lat = Number(x.latitude || (x.location && x.location.lat) || 0)
          var lng = Number(x.longitude || (x.location && x.location.lng) || 0)
          var dist = lat && lng ? that.calculateDistanceKM(userLat, userLng, lat, lng) : 0
          return {
            name: x.title || x.name || '',
            address: x.address || '',
            latitude: lat,
            longitude: lng,
            distanceKM: dist,
          }
        }).filter(function (x) { return x.name && x.latitude && x.longitude })
        mapped.sort(function (a, b) { return a.distanceKM - b.distanceKM })
        mapped = mapped.slice(0, 8)

        var places = mapped.map(function (x) {
          var kmText = (Math.round(x.distanceKM * 10) / 10).toFixed(1)
          return {
            name: x.name,
            reason: `附近景点 · ${kmText}km · ${x.address || ''}`,
            image: '../../icon/旅游景点.png',
            latitude: x.latitude,
            longitude: x.longitude,
            address: x.address
          }
        })

        if (!places.length) {
          if (!retried) {
            that.loadScenicRecommend('景区', true)
            return
          }
          that.setData({ isRecommendLoading: false, recommendPlaces: [], recommendError: "附近暂无景点推荐" })
          wx.showToast({ title: '附近暂无景点推荐', icon: 'none' })
          return
        }

        var markers = [{
          id: 1,
          latitude: userLat,
          longitude: userLng,
          width: 50,
          height: 50,
          iconPath: '/icon/086_定位.png'
        }]
        for (var i = 0; i < mapped.length; i++) {
          markers.push({
            id: i + 2,
            latitude: mapped[i].latitude,
            longitude: mapped[i].longitude,
            width: 32,
            height: 32,
            iconPath: '/icon/旗帜.png'
          })
        }

        that.setData({
          recommendPlaces: places,
          markers: markers,
          polyline: [],
          currentDestination: null,
          currentRouteResult: null,
          isRecommendLoading: false,
        })
      },
      fail: function () {
        that.setData({ isRecommendLoading: false, recommendError: "获取景点推荐失败" })
        wx.showToast({ title: '获取景点推荐失败', icon: 'none' })
      }
    })
  },

  calculateDistanceKM(lat1, lng1, lat2, lng2) {
    var R = 6371;
    var dLat = (lat2 - lat1) * Math.PI / 180;
    var dLng = (lng2 - lng1) * Math.PI / 180;
    var a = Math.sin(dLat / 2) * Math.sin(dLat / 2) +
      Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
      Math.sin(dLng / 2) * Math.sin(dLng / 2);
    var c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    return R * c;
  },

  decodePolyline(coors) {
    if (!Array.isArray(coors) || coors.length < 4) return []
    var kr = 1000000
    for (var i = 2; i < coors.length; i++) {
      coors[i] = Number(coors[i - 2]) + Number(coors[i]) / kr
    }
    var pts = []
    for (var j = 0; j < coors.length; j += 2) {
      var lat = Number(coors[j])
      var lng = Number(coors[j + 1])
      if (!lat || !lng) continue
      pts.push({ latitude: lat, longitude: lng })
    }
    return pts
  },

  loadFootByID(footID) {
    var that = this
    if (!footID) return
    util.request({ url: `/travel/foot/show/${footID}`, method: 'GET' })
      .then(res => {
        var foot = res && res.data && res.data.data
        if (!foot) throw new Error('足迹不存在')

        var origin = String(foot.Origin || '')
        var originParts = origin.split(',')
        var originLat = Number(originParts[0] || 0)
        var originLng = Number(originParts[1] || 0)

        var destLat = 0
        var destLng = 0
        try {
          var destArr = JSON.parse(foot.Destinations || '[]')
          if (Array.isArray(destArr) && destArr.length) {
            var last = String(destArr[destArr.length - 1] || '')
            var lastParts = last.split(',')
            destLat = Number(lastParts[0] || 0)
            destLng = Number(lastParts[1] || 0)
          }
        } catch (e) {}

        var routeObj = null
        try {
          routeObj = JSON.parse(foot.RouteResult || '{}')
        } catch (e) {}

        var route = routeObj && routeObj.result && routeObj.result.routes && routeObj.result.routes[0]
        var coors = route && route.polyline
        var points = that.decodePolyline((coors || []).slice ? coors : [])

        var markers = []
        if (originLat && originLng) {
          markers.push({
            id: 1,
            latitude: originLat,
            longitude: originLng,
            width: 50,
            height: 50,
            iconPath: '/icon/086_定位.png'
          })
        }
        if (destLat && destLng) {
          markers.push({
            id: 2,
            latitude: destLat,
            longitude: destLng,
            width: 32,
            height: 32,
            iconPath: '/icon/旗帜.png'
          })
        }

        var distance = route && route.distance
        var duration = route && route.duration
        var infoText = ''
        var distText = that.formatDistance(distance)
        var durText = that.formatDuration(duration)
        if (durText || distText) {
          infoText = `${foot.Mode || 'driving'} · ${durText} · ${distText}`
        }

        that.setData({
          latitude: originLat || that.data.latitude,
          longitude: originLng || that.data.longitude,
          markers: markers.length ? markers : that.data.markers,
          polyline: points.length ? [{ points: points, color: "#3b82f6", width: 6 }] : [],
          currentDestination: destLat && destLng ? { latitude: destLat, longitude: destLng, name: foot.Title || '足迹' } : null,
          currentRouteResult: routeObj,
          currentMode: foot.Mode || 'driving',
          routeDistance: Number(distance || 0),
          routeDuration: Number(duration || 0),
          routeInfoText: infoText,
        })
        that.updateSaveButtonVisibility()
      })
      .catch(() => {
        wx.showToast({ title: '打开足迹失败', icon: 'none' })
      })
  },

  planRouteTo(item, mode) {
    var that = this
    if (!item || !item.latitude || !item.longitude) {
      return Promise.reject(new Error('目的地无效'))
    }
    var fromLat = that.data.currentLatitude || that.data.latitude
    var fromLng = that.data.currentLongitude || that.data.longitude
    if (!fromLat || !fromLng) {
      return Promise.reject(new Error('请先开启定位权限'))
    }
    return new Promise(function (resolve, reject) {
      qqmapsdk.direction({
        mode: mode || 'driving',
        from: {
          latitude: fromLat,
          longitude: fromLng,
        },
        to: {
          latitude: item.latitude,
          longitude: item.longitude,
        },
        success: function (raw, data) {
          resolve({ raw: raw, data: data })
        },
        fail: function (err) {
          reject(err)
        },
      })
    }).then(function (res) {
      var full = res && res.raw
      var route = full && full.result && full.result.routes && full.result.routes[0]
      var coors = route && route.polyline
      var points = that.decodePolyline(Array.isArray(coors) ? coors.slice() : [])
      if (points.length) {
        var distance = route && route.distance
        var duration = route && route.duration
        var infoText = ''
        var distText = that.formatDistance(distance)
        var durText = that.formatDuration(duration)
        if (durText || distText) {
          infoText = `${mode || 'driving'} · ${durText} · ${distText}`
        }
        var startPoint = points[0]
        var endPoint = points[points.length - 1]
        that.setData({
          polyline: [{
            points: points,
            color: "#3b82f6",
            width: 6,
          }],
          currentDestination: item,
          currentRouteResult: full,
          currentMode: mode || 'driving',
          routeDistance: Number(distance || 0),
          routeDuration: Number(duration || 0),
          routeInfoText: infoText,
          includePoints: [startPoint, endPoint],
          markers: [{
            id: 1,
            latitude: fromLat,
            longitude: fromLng,
            width: 50,
            height: 50,
            iconPath: '/icon/086_定位.png'
          }, {
            id: 2,
            latitude: item.latitude,
            longitude: item.longitude,
            width: 32,
            height: 32,
            iconPath: '/icon/旗帜.png'
          }],
        })
        that.updateSaveButtonVisibility()
        setTimeout(function () {
          try {
            var ctx = wx.createMapContext('map', that)
            ctx.includePoints({
              points: points,
              padding: [60, 60, 240, 60]
            })
          } catch (e) {}
        }, 0)
      } else {
        throw new Error('未获取到路线')
      }
      return res
    })
  },

  saveCurrentFoot() {
    var dest = this.data.currentDestination
    var routeResult = this.data.currentRouteResult
    var mode = this.data.currentMode || 'driving'
    if (!dest || !routeResult) {
      wx.showToast({ title: '请先规划路线', icon: 'none' })
      return
    }
    var fromLat = this.data.currentLatitude || this.data.latitude
    var fromLng = this.data.currentLongitude || this.data.longitude
    var origin = fromLat + ',' + fromLng
    var destinations = [dest.latitude + ',' + dest.longitude]
    var destinationName = dest.name || dest.title || '目的地'
    var that = this
    that.resolveOriginName(fromLat, fromLng)
      .then(function (originName) {
        return util.request({
          url: '/travel/foot/create',
          method: 'POST',
          header: { 'Content-Type': 'application/json' },
          data: {
            title: destinationName,
            origin: origin,
            origin_name: originName,
            destinations: destinations,
            destination_names: [destinationName],
            mode: mode,
            routeResult: routeResult,
          },
        })
      })
      .then(() => {
        wx.showToast({ title: '足迹已保存', icon: 'success' })
      })
      .catch(() => {
        wx.showToast({ title: '保存失败', icon: 'none' })
      })
  },

  resolveOriginName(lat, lng) {
    var that = this
    if (that._originNameCache && that._originNameCache.name) {
      return Promise.resolve(that._originNameCache.name)
    }
    return new Promise(function (resolve) {
      qqmapsdk.reverseGeocoder({
        location: { latitude: lat, longitude: lng },
        success: function (_raw, data) {
          var full = data
          var addr = full && full.result && full.result.address
          var comp = full && full.result && full.result.address_component
          var city = (comp && comp.city) || ''
          var district = (comp && comp.district) || ''
          var name = addr || (city + district) || '我的位置'
          that._originNameCache = { name: name }
          resolve(name)
        },
        fail: function () {
          resolve('我的位置')
        }
      })
    })
  },

  refreshRecommend() {
    var that = this;
    if (that.data.isRecommendLoading) {
      return;
    }
    that.setData({ isRecommendLoading: true, recommendError: "" });

    util.ensureToken()
      .then(function () {
        return that.fetchRecommendKeywords(that.data.selectedMood);
      })
      .then(function (payload) {
        var keywords = payload.keywords || [];
        var region = payload.region || '汕头';
        var sliced = keywords.slice(0, 4);
        var tasks = sliced.map(function (kw) {
          return that.qqSearch(kw, region).then(function (items) {
            return items.map(function (x) {
              return {
                keyword: kw,
                title: x.title || x.name || '',
                address: x.address || '',
                latitude: Number(x.latitude || (x.location && x.location.lat) || 0),
                longitude: Number(x.longitude || (x.location && x.location.lng) || 0),
              };
            });
          });
        });
        return Promise.all(tasks).then(function (groups) {
          return { groups: groups, region: region, keywords: sliced };
        });
      })
      .then(function (result) {
        var userLat = that.data.latitude;
        var userLng = that.data.longitude;
        var seen = {};
        var merged = [];
        for (var i = 0; i < result.groups.length; i++) {
          var group = result.groups[i] || [];
          for (var j = 0; j < group.length; j++) {
            var item = group[j];
            if (!item.title || !item.latitude || !item.longitude) {
              continue;
            }
            var key = (item.title + '|' + item.address + '|' + item.latitude + '|' + item.longitude).toLowerCase();
            if (seen[key]) {
              continue;
            }
            seen[key] = true;
            item.distanceKM = that.calculateDistanceKM(userLat, userLng, item.latitude, item.longitude);
            merged.push(item);
          }
        }
        merged.sort(function (a, b) {
          return a.distanceKM - b.distanceKM;
        });
        merged = merged.slice(0, 10);

        var places = merged.map(function (x) {
          var kmText = (Math.round(x.distanceKM * 10) / 10).toFixed(1);
          return {
            name: x.title,
            reason: x.keyword + ' · ' + kmText + 'km · ' + (x.address || ''),
            image: '../../icon/地球仪.png',
            latitude: x.latitude,
            longitude: x.longitude,
            address: x.address
          };
        });

        var markers = [{
          id: 1,
          latitude: userLat,
          longitude: userLng,
          width: 50,
          height: 50,
          iconPath: '/icon/086_定位.png'
        }];
        for (var k = 0; k < merged.length; k++) {
          markers.push({
            id: k + 2,
            latitude: merged[k].latitude,
            longitude: merged[k].longitude,
            width: 32,
            height: 32,
            iconPath: '/icon/旗帜.png'
          });
        }

        that.setData({
          recommendPlaces: places,
          markers: markers
        });
        that.setData({ isRecommendLoading: false });
      })
      .catch(function (err) {
        var msg = (err && err.message) || '获取推荐失败';
        that.setData({ recommendError: msg, isRecommendLoading: false });
        wx.showToast({ title: msg, icon: 'none' });
      });
  },

  onTapRecommend(e) {
    var idx = Number(e.currentTarget.dataset.index || 0);
    var item = (this.data.recommendPlaces || [])[idx];
    if (!item || !item.latitude || !item.longitude) {
      return;
    }
    var that = this
    wx.showActionSheet({
      itemList: ['查看地点', '规划驾车路线', '规划公交路线', '规划步行路线', '保存足迹'],
      success(res) {
        if (res.tapIndex === 0) {
          wx.openLocation({
            latitude: item.latitude,
            longitude: item.longitude,
            name: item.name || '',
            address: item.address || '',
            scale: 16
          })
          return
        }
        if (res.tapIndex === 1) {
          that.planRouteTo(item, 'driving').catch((err) => wx.showToast({ title: (err && err.message) || '规划失败', icon: 'none' }))
          return
        }
        if (res.tapIndex === 2) {
          that.planRouteTo(item, 'transit').catch((err) => wx.showToast({ title: (err && err.message) || '规划失败', icon: 'none' }))
          return
        }
        if (res.tapIndex === 3) {
          that.planRouteTo(item, 'walking').catch((err) => wx.showToast({ title: (err && err.message) || '规划失败', icon: 'none' }))
          return
        }
        if (res.tapIndex === 4) {
          that.saveCurrentFoot()
        }
      }
    })
  },
  
//----地图选点------------------

  // 打开腾讯地图选点插件
  // openLocationPicker() {
  //   const chooseLocation = requirePlugin("chooseLocation");

  //   // 跳转到插件页面
  //   wx.navigateTo({
  //     url: "plugin://chooseLocation/index",
  //   });

  //   // 监听插件返回的数据
  //   wx.onAppShow((res) => {
  //     if (res && res.referrerInfo && res.referrerInfo.extraData) {
  //       const locationInfo = res.referrerInfo.extraData;

  //       // 更新页面数据中的纬度和经度
  //       this.setData({
  //         latitude: locationInfo.latitude,
  //         longitude: locationInfo.longitude,
  //       });

  //       console.log("更新后的经纬度：", this.data.latitude, this.data.longitude);
  //     }
  //   });
  // },

//--------------------------------- 

touchStart(e) {
  const touch = e.touches[0] || e.changedTouches[0];
  this.setData({
    startY: touch.clientY,
    initialDragPosition: this.data.dragPosition,
    isDragging: true
  });
},

touchMove(e) {
  const touch = e.touches[0] || e.changedTouches[0];
  const moveY = touch.clientY;
  const deltaY = this.data.startY - moveY;

  if (this.data.isFull) {
    if (deltaY > 0 || this.data.scrollTop > 0) {
      return;
    }
  }

  let newDragPosition = this.data.initialDragPosition - deltaY;
  const maxPos = (this.data.windowHeight || wx.getSystemInfoSync().windowHeight) - 100;
  
  newDragPosition = Math.max(0, Math.min(newDragPosition, maxPos));

  this.setData({
      dragPosition: newDragPosition
  });
  this.updateSaveButtonVisibility()
},

touchEnd() {
  const maxPos = (this.data.windowHeight || wx.getSystemInfoSync().windowHeight) - 100;

  let closestPosition = 0;
  if (this.data.dragPosition > maxPos / 2) {
      closestPosition = maxPos;
  } else {
      closestPosition = 0;
  }

  this.setData({
      dragPosition: closestPosition,
      isDragging: false,
      isFull: closestPosition === 0
  });
  this.updateSaveButtonVisibility()
},

onScroll(e) {
  this.setData({ scrollTop: e.detail.scrollTop });
}

});
