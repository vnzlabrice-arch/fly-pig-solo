// 飞猪出行 - 网约车主页逻辑

(function() {
  'use strict';

  // 状态管理
  const state = {
    isLoggedIn: false,
    userId: null,
    token: null,
    userInfo: null,
    fromAddress: '',
    fromLng: 116.397428,
    fromLat: 39.90923,
    toAddress: '',
    toLng: 0,
    toLat: 0,
    selectedCarType: 1,
    couponId: null,
    // 地图相关
    map: null,
    geolocation: null,
    geocoder: null,
    markers: [],
    walking: null,
  };

  // 默认城市（北京市）
  const DEFAULT_CITY = '北京';
  const DEFAULT_CENTER = [116.397428, 39.90923]; // 北京天安门

  // DOM 元素
  const elements = {
    toast: document.getElementById('toast'),
    userPanel: document.getElementById('userPanel'),
    panelOverlay: document.getElementById('panelOverlay'),
    closePanelBtn: document.getElementById('closePanelBtn'),
    userBtn: document.getElementById('userBtn'),
    loginBtn: document.getElementById('loginBtn'),
    logoutBtn: document.getElementById('logoutBtn'),
    userAvatar: document.getElementById('userAvatar'),
    userName: document.getElementById('userName'),
    userPhone: document.getElementById('userPhone'),
    fromAddress: document.getElementById('fromAddress'),
    toAddress: document.getElementById('toAddress'),
    toInput: document.querySelector('.to-input'),
    clearToBtn: document.getElementById('clearToBtn'),
    destinationPanel: document.getElementById('destinationPanel'),
    searchInput: document.getElementById('searchInput'),
    panelContent: document.getElementById('panelContent'),
    carTypeSection: document.getElementById('carTypeSection'),
    carTypes: document.querySelectorAll('.car-type'),
    couponBar: document.getElementById('couponBar'),
    couponText: document.getElementById('couponText'),
    couponModal: document.getElementById('couponModal'),
    couponOverlay: document.getElementById('couponOverlay'),
    closeCouponBtn: document.getElementById('closeCouponBtn'),
    couponList: document.getElementById('couponList'),
    noCouponBtn: document.getElementById('noCouponBtn'),
    remarkSection: document.getElementById('remarkSection'),
    remarkInput: document.getElementById('remarkInput'),
    confirmSection: document.getElementById('confirmSection'),
    confirmBtn: document.getElementById('confirmBtn'),
    btnPrice: document.getElementById('btnPrice'),
    locationBtn: document.getElementById('locationBtn'),
    fromLocation: document.getElementById('fromLocation'),
    navItems: document.querySelectorAll('.nav-item'),
    realNameMenuItem: document.getElementById('realNameMenuItem'),
    addressMenuItem: document.getElementById('addressMenuItem'),
    couponMenuItem: document.getElementById('couponMenuItem'),
    profileMenuItem: document.getElementById('profileMenuItem'),
  };

  // 模拟优惠券数据（实际项目中应从API获取）
  const mockCoupons = [
    {
      id: 1,
      name: '新人专享券',
      value: 5,
      condition: '满10元可用',
      expireTime: '2026-05-31',
      tags: ['新人专享'],
      minAmount: 10,
      status: 'available'
    },
    {
      id: 2,
      name: '限时8折券',
      value: 20,
      condition: '最高抵扣20元',
      expireTime: '2026-05-15',
      tags: ['限时优惠'],
      minAmount: 0,
      discount: 0.8,
      status: 'available'
    },
    {
      id: 3,
      name: '雨天出行券',
      value: 3,
      condition: '满8元可用',
      expireTime: '2026-05-20',
      tags: ['天气专项'],
      minAmount: 8,
      status: 'available'
    }
  ];

  // Toast 提示
  function showToast(message, type = 'success') {
    elements.toast.textContent = message;
    elements.toast.className = `toast show ${type}`;
    setTimeout(() => {
      elements.toast.className = 'toast';
    }, 2500);
  }

  // 初始化高德地图
  async function initMap() {
    const mapContainer = document.getElementById('amap-container');
    const mapLoading = document.getElementById('mapLoading');
    
    // 检查 AMap 是否加载成功
    if (typeof AMap === 'undefined') {
      console.error('高德地图 SDK 未加载');
      if (mapLoading) {
        mapLoading.innerHTML = '<span style="color: #999;">地图加载失败，请刷新重试</span>';
      }
      return;
    }

    try {
      // 先加载需要的插件
      await new Promise((resolve, reject) => {
        AMap.plugin(['AMap.Geolocation', 'AMap.Walking', 'AMap.Geocoder'], () => {
          resolve();
        }, (err) => {
          console.error('插件加载失败:', err);
          reject(err);
        });
      });

      // 创建地图实例
      state.map = new AMap.Map('amap-container', {
        zoom: 14,
        center: DEFAULT_CENTER,
        viewMode: '2D',
        mapStyle: 'amap://styles/normal',
      });

      // 添加地图加载完成事件
      state.map.on('complete', function() {
        if (mapLoading) {
          mapLoading.classList.add('hidden');
          setTimeout(() => mapLoading.style.display = 'none', 300);
        }
      });

      // 初始化定位插件
      state.geolocation = new AMap.Geolocation({
        enableHighAccuracy: true,
        timeout: 10000,
        buttonOffset: new AMap.Pixel(10, 20),
        zoomToAccuracy: true,
        buttonPosition: 'RB'
      });
      state.map.addControl(state.geolocation);

      // 初始化步行路径规划
      state.walking = new AMap.Walking({
        map: state.map,
        panel: null,
      });

      // 初始化逆地理编码
      state.geocoder = new AMap.Geocoder({
        city: DEFAULT_CITY,
      });

      // 地图点击事件 - 选择目的地
      state.map.on('click', function(e) {
        const lng = e.lnglat.getLng();
        const lat = e.lnglat.getLat();

        // 清除之前的终点标记
        clearMarkers();

        // 添加新的终点标记
        addMarker(lng, lat, '目的地', 'end');

        // 逆地理编码 - 将坐标转为地址
        state.geocoder.getAddress([lng, lat], function(status, result) {
          if (status === 'complete' && result.regeocode) {
            const address = result.regeocode.formattedAddress;
            // 简化地址（去掉省市区）
            const shortAddress = address
              .replace(result.regeocode.addressComponent.province, '')
              .replace(result.regeocode.addressComponent.city, '')
              .replace(result.regeocode.addressComponent.district, '')
              .trim() || '未知位置';

            state.toAddress = shortAddress;
            state.toLng = lng;
            state.toLat = lat;
            elements.toAddress.value = state.toAddress;
            elements.clearToBtn.style.display = 'block';

            // 规划步行路线
            planWalkingRoute(lng, lat, state.toAddress);

            // 显示车型选择和确认区域
            elements.carTypeSection.style.display = 'block';
            elements.couponBar.style.display = 'flex';
            elements.remarkSection.style.display = 'flex';
            elements.confirmSection.style.display = 'flex';
            
            // 更新优惠券可用性
            updateCouponAvailability();
          }
        });
      });

      // 默认定位
      autoLocate();

    } catch (e) {
      console.error('地图初始化失败:', e);
      if (mapLoading) {
        mapLoading.innerHTML = '<span style="color: #999;">地图初始化失败</span>';
      }
    }
  }

  // 自动定位
  function autoLocate() {
    // 优先使用浏览器原生定位（不需要 Key）
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          // 转换坐标（浏览器是 GPS 坐标，高德是火星坐标，这里简化处理）
          state.fromLng = position.coords.longitude;
          state.fromLat = position.coords.latitude;
          
          // 如果有高德地图，进行坐标转换
          if (typeof AMap !== 'undefined' && state.map) {
            convertGPSToAMap(state.fromLng, state.fromLat);
          } else {
            // 没有高德地图，直接使用浏览器坐标
            state.fromAddress = '当前位置';
            elements.fromAddress.value = state.fromAddress;
            setDefaultLocation();
          }
        },
        (error) => {
          console.log('浏览器定位失败:', error.message);
          // 降级到高德定位或默认位置
          if (state.geolocation) {
            state.geolocation.getCurrentPosition(function(status, result) {
              if (status === 'complete') {
                onLocationSuccess(result);
              } else {
                console.log('高德定位也失败，使用默认位置');
                setDefaultLocation();
              }
            });
          } else {
            setDefaultLocation();
          }
        },
        { enableHighAccuracy: true, timeout: 10000, maximumAge: 60000 }
      );
    } else if (state.geolocation) {
      // 没有浏览器定位，使用高德定位
      state.geolocation.getCurrentPosition(function(status, result) {
        if (status === 'complete') {
          onLocationSuccess(result);
        } else {
          console.log('高德定位失败，使用默认位置');
          setDefaultLocation();
        }
      });
    } else {
      // 都没有，使用默认位置
      setDefaultLocation();
    }
  }

  // GPS 坐标转高德坐标（火星坐标系）
  function convertGPSToAMap(gpsLng, gpsLat) {
    AMap.plugin('AMap.Geocoder', function() {
      const geocoder = new AMap.Geocoder();
      
      // 使用高德逆地理编码获取准确地址
      geocoder.getAddress([gpsLng, gpsLat], function(status, result) {
        if (status === 'complete' && result.regeocode) {
          const addr = result.regeocode;
          state.fromAddress = addr.formattedAddress || (
            addr.addressComponent.province +
            addr.addressComponent.city +
            addr.addressComponent.district +
            (addr.aois && addr.aois[0] ? addr.aois[0].name : '')
          );
          elements.fromAddress.value = state.fromAddress;
        }
        
        // 移动地图中心
        state.map.setCenter([gpsLng, gpsLat]);
        addMarker(gpsLng, gpsLat, state.fromAddress, 'start');
      });
    });
  }

  // 定位成功回调
  function onLocationSuccess(result) {
    const { position, addressComponent, formatted } = result;
    
    state.fromLng = position.lng;
    state.fromLat = position.lat;
    state.fromAddress = formatted || (
      addressComponent.province + 
      addressComponent.city + 
      addressComponent.district + 
      addressComponent.township
    );
    
    elements.fromAddress.value = state.fromAddress;
    
    // 移动地图中心
    state.map.setCenter([state.fromLng, state.fromLat]);
    
    // 添加起点标记
    addMarker(state.fromLng, state.fromLat, state.fromAddress, 'start');
  }

  // 设置默认位置
  function setDefaultLocation() {
    state.fromLng = DEFAULT_CENTER[0];
    state.fromLat = DEFAULT_CENTER[1];
    state.fromAddress = DEFAULT_CITY + '市天安门广场';
    
    elements.fromAddress.value = state.fromAddress;
    
    state.map.setCenter(DEFAULT_CENTER);
    addMarker(state.fromLng, state.fromLat, state.fromAddress, 'start');
  }

  // 添加标记
  function addMarker(lng, lat, title, type) {
    if (!state.map) return;

    const markerContent = document.createElement('div');
    const markerIcon = document.createElement('div');
    
    if (type === 'start') {
      markerIcon.className = 'marker-start';
      markerIcon.innerHTML = '<div class="marker-inner green"></div><div class="marker-pulse"></div>';
    } else {
      markerIcon.className = 'marker-end';
      markerIcon.innerHTML = '<div class="marker-inner red"></div><div class="marker-pulse"></div>';
    }
    markerContent.appendChild(markerIcon);

    const marker = new AMap.Marker({
      content: markerContent,
      position: [lng, lat],
      title: title,
      extData: { type, title },
    });

    state.map.add(marker);
    state.markers.push(marker);

    return marker;
  }

  // 清除所有标记
  function clearMarkers() {
    if (!state.map) return;
    
    state.markers.forEach(marker => {
      state.map.remove(marker);
    });
    state.markers = [];
  }

  // 清除路线
  function clearRoute() {
    if (!state.map) return;
    // 清除步行路线（通过移除 polyline）
    const overlays = state.map.getAllOverlays();
    overlays.forEach(overlay => {
      if (overlay instanceof AMap.Polyline || overlay instanceof AMap.Marker) {
        // 保留标记
      }
    });
  }

  // 搜索地址
  async function searchAddress(keyword) {
    if (!keyword || !state.map) return [];

    return new Promise((resolve) => {
      AMap.plugin('AMap.AutoComplete', function() {
        const autoOptions = {
          city: DEFAULT_CITY,
        };
        const autoComplete = new AMap.AutoComplete(autoOptions);
        
        autoComplete.search(keyword, function(status, result) {
          if (status === 'complete' && result.tips) {
            resolve(result.tips.filter(tip => tip.location));
          } else {
            resolve([]);
          }
        });
      });
    });
  }

  // 规划步行路线
  function planWalkingRoute(toLng, toLat, toAddress) {
    if (!state.walking || !state.fromLng) return;

    state.walking.clear();
    
    state.walking.search([state.fromLng, state.fromLat], [toLng, toLat], function(status, result) {
      if (status === 'complete' && result.routes && result.routes.length > 0) {
        const route = result.routes[0];
        const distance = (route.distance / 1000).toFixed(1);
        const duration = Math.ceil(route.time / 60);
        
        // 更新 ETA
        updateETA(duration);
        
        // 更新优惠券可用性显示
        updateCouponAvailability();
      }
    });
  }

  // 更新优惠券可用性显示
  function updateCouponAvailability() {
    const currentPrice = getCurrentPrice();
    const availableCoupons = mockCoupons.filter(c => c.status === 'available' && currentPrice >= c.minAmount);
    
    if (availableCoupons.length > 0) {
      elements.couponText.textContent = `有 ${availableCoupons.length} 张优惠券可用`;
      elements.couponBar.style.background = 'linear-gradient(135deg, #FFF8E1 0%, #FFF3E0 100%)';
    } else {
      elements.couponText.textContent = '暂无可用优惠券';
      elements.couponBar.style.background = 'linear-gradient(135deg, #F5F5F5 0%, #EEEEEE 100%)';
    }
  }

  // 更新 ETA
  function updateETA(minutes) {
    document.querySelectorAll('.eta').forEach(el => {
      el.textContent = minutes;
    });
  }

  // 初始化状态
  function initState() {
    state.token = localStorage.getItem('user_token') || '';
    
    // 严格解析 user_id
    let savedUserId = localStorage.getItem('user_id');
    if (savedUserId) {
      savedUserId = String(savedUserId).trim();
      // 处理异常格式如 "2:1" 或 ":1"
      if (savedUserId.includes(':')) {
        savedUserId = savedUserId.split(':')[0];
      }
      // 只保留数字
      savedUserId = savedUserId.replace(/[^0-9]/g, '');
      state.userId = parseInt(savedUserId, 10) || 0;
    }
    
    // 只有同时有有效的 token 和 userId 才算登录
    state.isLoggedIn = Boolean(state.token && state.userId > 0);
  }

  // 生成用户头像 URL
  function generateAvatar(nickname, phone) {
    // 如果有真实头像，直接使用
    if (state.userInfo && state.userInfo.Avatar) {
      return state.userInfo.Avatar;
    }
    
    // 生成唯一标识（使用昵称或手机号）
    const seed = encodeURIComponent(nickname || phone || 'user');
    
    // 使用 DiceBear 随机头像服务
    const style = ['avataaars', 'bottts', 'personas', 'shapes', 'identicon'][Math.floor(Math.random() * 5)];
    return `https://api.dicebear.com/7.x/${style}/svg?seed=${seed}&backgroundColor=ff6b00,ff8c00,ff4500,ffa500,e55c00`;
  }

  // 设置头像
  function setAvatar(avatarUrl, nickname) {
    if (avatarUrl.startsWith('http')) {
      elements.userAvatar.innerHTML = `<img src="${avatarUrl}" alt="头像" onerror="this.parentElement.textContent='${nickname.charAt(0).toUpperCase()}'">`;
    } else {
      elements.userAvatar.textContent = nickname.charAt(0).toUpperCase();
    }
  }

  // 更新用户界面
  function updateUserUI() {
    if (state.isLoggedIn && state.userInfo) {
      const nickname = state.userInfo.Nickname || state.userInfo.nickname || '用户';
      const phone = state.userInfo.Phone || state.userInfo.phone || '';
      const avatar = state.userInfo.Avatar || '';
      
      elements.userName.textContent = nickname;
      elements.userPhone.textContent = phone ? phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2') : '-';
      
      // 设置头像
      const avatarUrl = avatar || generateAvatar(nickname, phone);
      setAvatar(avatarUrl, nickname);
      
      elements.loginBtn.style.display = 'none';
      elements.logoutBtn.style.display = 'flex';
    } else {
      elements.userName.textContent = '未登录';
      elements.userPhone.textContent = '-';
      elements.userAvatar.textContent = '👤';
      elements.loginBtn.style.display = 'flex';
      elements.logoutBtn.style.display = 'none';
    }
  }

  // 加载用户信息
  async function loadUserInfo() {
    if (!state.isLoggedIn) return;

    try {
      const resp = await window.UserApi.getUserDetail({ user_id: state.userId });
      if (window.UserApi.isSuccessCode(resp.code) && resp.data) {
        state.userInfo = resp.data;
        updateUserUI();
      }
    } catch (e) {
      console.error('加载用户信息失败:', e);
    }
  }

  // 打开用户面板
  function openUserPanel() {
    elements.userPanel.classList.add('active');
  }

  // 关闭用户面板
  function closeUserPanel() {
    elements.userPanel.classList.remove('active');
  }

  // 退出登录
  function logout() {
    localStorage.removeItem('user_token');
    localStorage.removeItem('user_id');
    state.token = '';
    state.userId = null;
    state.userInfo = null;
    state.isLoggedIn = false;
    updateUserUI();
    closeUserPanel();
    showToast('已退出登录');
  }

  // 跳转到登录页
  function goToLogin() {
    closeUserPanel();
    window.location.href = './login.html?returnTo=' + encodeURIComponent(window.location.pathname.split('/').pop());
  }

  // 显示/隐藏目的地面板
  function toggleDestinationPanel(show) {
    if (show) {
      elements.destinationPanel.classList.add('active');
      elements.searchInput.focus();
    } else {
      elements.destinationPanel.classList.remove('active');
    }
  }

  // 选择目的地
  function selectDestination(title, desc) {
    state.toAddress = title;
    elements.toAddress.value = state.toAddress;
    elements.clearToBtn.style.display = 'block';
    toggleDestinationPanel(false);
    
    // 搜索目的地坐标
    geocodeAddress(state.toAddress).then(coords => {
      if (coords) {
        state.toLng = coords.lng;
        state.toLat = coords.lat;
        
        // 添加终点标记
        clearMarkers();
        addMarker(state.fromLng, state.fromLat, state.fromAddress, 'start');
        addMarker(state.toLng, state.toLat, state.toAddress, 'end');
        
        // 规划步行路线
        planWalkingRoute(state.toLng, state.toLat, state.toAddress);
      }
    });
    
    // 显示车型选择和确认区域
    elements.carTypeSection.style.display = 'block';
    elements.couponBar.style.display = 'flex';
    elements.remarkSection.style.display = 'flex';
    elements.confirmSection.style.display = 'flex';
    
    // 更新优惠券可用性
    updateCouponAvailability();
  }

  // 地址转坐标
  async function geocodeAddress(address) {
    if (!address || typeof AMap === 'undefined') return null;
    
    return new Promise((resolve) => {
      AMap.plugin('AMap.Geocoder', function() {
        const geocoder = new AMap.Geocoder({
          city: DEFAULT_CITY,
        });
        
        geocoder.getLocation(address, function(status, result) {
          if (status === 'complete' && result.geocodes.length > 0) {
            const location = result.geocodes[0].location;
            resolve({ lng: location.lng, lat: location.lat });
          } else {
            // 如果精确搜索失败，使用关键词搜索
            searchAddress(address).then(results => {
              if (results && results.length > 0) {
                const loc = results[0].location;
                resolve({ lng: loc.lng, lat: loc.lat });
              } else {
                resolve(null);
              }
            });
          }
        });
      });
    });
  }

  // 选择车型
  function selectCarType(carTypeEl) {
    elements.carTypes.forEach(el => el.classList.remove('active'));
    carTypeEl.classList.add('active');
    
    state.selectedCarType = parseInt(carTypeEl.dataset.type, 10);
    let price = parseFloat(carTypeEl.dataset.price);
    
    // 如果有选择优惠券，应用优惠
    if (state.couponId) {
      const coupon = mockCoupons.find(c => c.id === state.couponId);
      if (coupon) {
        if (coupon.discount) {
          // 折扣券
          const discount = price * (1 - coupon.discount);
          price = price - Math.min(discount, coupon.value);
        } else {
          // 满减券
          price = price - coupon.value;
        }
        price = Math.max(0, price);
      }
    }
    
    elements.btnPrice.textContent = '¥' + price.toFixed(0);
  }

  // 打开优惠券弹窗
  function openCouponModal() {
    if (!state.isLoggedIn) {
      showToast('请先登录', 'error');
      goToLogin();
      return;
    }
    
    elements.couponModal.classList.add('active');
    loadCoupons();
  }

  // 关闭优惠券弹窗
  function closeCouponModal() {
    elements.couponModal.classList.remove('active');
  }

  // 加载优惠券列表
  function loadCoupons() {
    const currentPrice = getCurrentPrice();
    
    // 过滤可用优惠券（满足使用条件）
    const availableCoupons = mockCoupons.filter(coupon => {
      return coupon.status === 'available' && currentPrice >= coupon.minAmount;
    });
    
    const unavailableCoupons = mockCoupons.filter(coupon => {
      return coupon.status === 'available' && currentPrice < coupon.minAmount;
    });
    
    // 更新优惠券栏显示
    if (availableCoupons.length > 0) {
      elements.couponText.textContent = `有 ${availableCoupons.length} 张优惠券可用`;
      elements.couponBar.style.background = 'linear-gradient(135deg, #FFF8E1 0%, #FFF3E0 100%)';
    } else {
      elements.couponText.textContent = '暂无可用优惠券';
      elements.couponBar.style.background = 'linear-gradient(135deg, #F5F5F5 0%, #EEEEEE 100%)';
    }
    
    // 渲染优惠券列表
    if (availableCoupons.length === 0 && unavailableCoupons.length === 0) {
      elements.couponList.innerHTML = `
        <div class="no-coupons">
          <div class="no-coupons-icon">🎫</div>
          <div class="no-coupons-text">暂无可用优惠券</div>
          <div class="no-coupons-tip">下单后可获得更多优惠</div>
        </div>
      `;
      return;
    }
    
    let html = '';
    
    // 可用优惠券
    availableCoupons.forEach(coupon => {
      const isSelected = state.couponId === coupon.id;
      html += `
        <div class="coupon-item ${isSelected ? 'selected' : ''}" data-id="${coupon.id}">
          <div class="coupon-item-header">
            <div class="coupon-value">
              ${coupon.discount ? coupon.discount * 10 + '折' : '¥' + coupon.value}
              ${!coupon.discount ? '<span>立减</span>' : ''}
            </div>
            <div class="coupon-info-block">
              <div class="coupon-name">${coupon.name}</div>
              <div class="coupon-condition">${coupon.condition}</div>
            </div>
          </div>
          <div class="coupon-item-footer">
            <span>有效期至 ${coupon.expireTime}</span>
            ${coupon.tags.map(tag => `<span class="coupon-tag">${tag}</span>`).join('')}
          </div>
        </div>
      `;
    });
    
    // 不可用优惠券
    unavailableCoupons.forEach(coupon => {
      html += `
        <div class="coupon-item disabled" data-id="${coupon.id}">
          <div class="coupon-item-header">
            <div class="coupon-value">
              ${coupon.discount ? coupon.discount * 10 + '折' : '¥' + coupon.value}
            </div>
            <div class="coupon-info-block">
              <div class="coupon-name">${coupon.name}</div>
              <div class="coupon-condition">还差 ¥${(coupon.minAmount - currentPrice).toFixed(0)}可用</div>
            </div>
          </div>
          <div class="coupon-item-footer">
            <span>有效期至 ${coupon.expireTime}</span>
            ${coupon.tags.map(tag => `<span class="coupon-tag">${tag}</span>`).join('')}
          </div>
        </div>
      `;
    });
    
    elements.couponList.innerHTML = html;
    
    // 绑定点击事件
    elements.couponList.querySelectorAll('.coupon-item:not(.disabled)').forEach(item => {
      item.addEventListener('click', () => {
        selectCoupon(parseInt(item.dataset.id, 10));
      });
    });
  }

  // 选择优惠券
  function selectCoupon(couponId) {
    // 如果点击已选中的优惠券，则取消选择
    if (state.couponId === couponId) {
      state.couponId = null;
    } else {
      state.couponId = couponId;
    }
    
    // 更新优惠券列表显示
    elements.couponList.querySelectorAll('.coupon-item').forEach(item => {
      item.classList.remove('selected');
      if (parseInt(item.dataset.id, 10) === state.couponId) {
        item.classList.add('selected');
      }
    });
    
    // 更新按钮价格
    updatePriceWithCoupon();
    
    // 更新优惠券栏显示
    if (state.couponId) {
      const coupon = mockCoupons.find(c => c.id === state.couponId);
      if (coupon) {
        elements.couponText.innerHTML = `<span class="selected-coupon-tag">已选: ${coupon.name}</span>`;
      }
    } else {
      const availableCoupons = mockCoupons.filter(c => c.status === 'available' && getCurrentPrice() >= c.minAmount);
      elements.couponText.textContent = availableCoupons.length > 0 ? `有 ${availableCoupons.length} 张优惠券可用` : '暂无可用优惠券';
    }
    
    // 延迟关闭弹窗
    setTimeout(() => {
      closeCouponModal();
    }, 300);
  }

  // 不使用优惠券
  function useNoCoupon() {
    state.couponId = null;
    updatePriceWithCoupon();
    const availableCoupons = mockCoupons.filter(c => c.status === 'available' && getCurrentPrice() >= c.minAmount);
    elements.couponText.textContent = availableCoupons.length > 0 ? `有 ${availableCoupons.length} 张优惠券可用` : '暂无可用优惠券';
    closeCouponModal();
  }

  // 获取当前价格
  function getCurrentPrice() {
    const activeCar = document.querySelector('.car-type.active');
    return activeCar ? parseFloat(activeCar.dataset.price) : 12;
  }

  // 更新价格（应用优惠券）
  function updatePriceWithCoupon() {
    const currentPrice = getCurrentPrice();
    let finalPrice = currentPrice;
    
    if (state.couponId) {
      const coupon = mockCoupons.find(c => c.id === state.couponId);
      if (coupon) {
        if (coupon.discount) {
          // 折扣券
          const discount = currentPrice * (1 - coupon.discount);
          finalPrice = currentPrice - Math.min(discount, coupon.value);
        } else {
          // 满减券
          finalPrice = currentPrice - coupon.value;
        }
        finalPrice = Math.max(0, finalPrice);
      }
    }
    
    elements.btnPrice.textContent = '¥' + finalPrice.toFixed(0);
  }

  // 确认呼叫
  async function confirmRide() {
    if (!state.isLoggedIn) {
      showToast('请先登录', 'error');
      goToLogin();
      return;
    }

    if (!state.fromAddress) {
      showToast('请选择出发地', 'error');
      return;
    }

    if (!state.toAddress) {
      showToast('请选择目的地', 'error');
      return;
    }

    elements.confirmBtn.disabled = true;
    elements.confirmBtn.querySelector('.btn-text').textContent = '叫车中...';

    try {
      // 确保 userId 是有效数字
      let validUserId = state.userId;
      if (!validUserId || validUserId <= 0) {
        showToast('用户未登录或登录已过期', 'error');
        goToLogin();
        return;
      }
      
      const resp = await window.UserApi.addPassengerOrder({
        userId: validUserId,
        passengerName: state.userInfo?.Nickname || state.userInfo?.nickname || '乘客',
        passengerPhone: state.userInfo?.Phone || state.userInfo?.phone || '',
        startAddress: state.fromAddress,
        startLbg: 116.4, // 注意：这里字段名拼写为 startLbg（后端有个拼写错误）
        startLat: 39.9,
        endAddress: state.toAddress,
        endLng: state.toLng,
        endLat: state.toLat,
        carType: state.selectedCarType,
        couponId: state.couponId || 0,
        remark: elements.remarkInput.value || '',
      });

      if (window.UserApi.isSuccessCode(resp.code)) {
        showToast('叫车成功，正在为您派单...');
        // 重置状态
        resetOrderState();
      } else {
        showToast(resp.message || '叫车失败，请重试', 'error');
      }
    } catch (e) {
      showToast('网络错误，请重试', 'error');
    } finally {
      elements.confirmBtn.disabled = false;
      elements.confirmBtn.querySelector('.btn-text').textContent = '确认呼叫';
    }
  }

  // 重置订单状态
  function resetOrderState() {
    state.toAddress = '';
    state.toLng = 0;
    state.toLat = 0;
    state.couponId = null;
    elements.toAddress.value = '';
    elements.clearToBtn.style.display = 'none';
    elements.remarkInput.value = '';
    elements.carTypeSection.style.display = 'none';
    elements.couponBar.style.display = 'none';
    elements.remarkSection.style.display = 'none';
    elements.confirmSection.style.display = 'none';
    
    // 重置优惠券栏显示
    elements.couponText.textContent = '有 2 张优惠券可用';
    elements.couponBar.style.background = 'linear-gradient(135deg, #FFF8E1 0%, #FFF3E0 100%)';
    
    // 清除地图标记和路线
    clearMarkers();
    if (state.walking) {
      state.walking.clear();
    }
    if (state.map && state.fromLng) {
      addMarker(state.fromLng, state.fromLat, state.fromAddress, 'start');
    }
    
    // 重置 ETA
    updateETA(5);
  }

  // 获取当前位置（浏览器原生）
  function getCurrentLocation() {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const lat = position.coords.latitude;
          const lng = position.coords.longitude;
          // 模拟获取地址
          state.fromAddress = '我的位置 (' + lat.toFixed(4) + ', ' + lng.toFixed(4) + ')';
          elements.fromAddress.value = state.fromAddress;
          showToast('已定位到您的位置');
        },
        (error) => {
          console.error('获取位置失败:', error);
          // 默认地址
          state.fromAddress = '北京市朝阳区';
          elements.fromAddress.value = state.fromAddress;
          showToast('已使用默认位置');
        }
      );
    } else {
      state.fromAddress = '北京市朝阳区';
      elements.fromAddress.value = state.fromAddress;
      showToast('浏览器不支持定位');
    }
  }

  // 绑定事件
  function bindEvents() {
    // 用户面板
    elements.userBtn.addEventListener('click', openUserPanel);
    elements.closePanelBtn.addEventListener('click', closeUserPanel);
    elements.panelOverlay.addEventListener('click', closeUserPanel);
    elements.loginBtn.addEventListener('click', goToLogin);
    elements.logoutBtn.addEventListener('click', logout);

    // 面板菜单导航
    elements.realNameMenuItem.addEventListener('click', () => {
      closeUserPanel();
      window.location.href = './profile.html#realName';
    });
    elements.addressMenuItem.addEventListener('click', () => {
      closeUserPanel();
      showToast('常用地址功能开发中');
    });
    elements.couponMenuItem.addEventListener('click', () => {
      closeUserPanel();
      window.location.href = './coupons.html';
    });
    elements.profileMenuItem.addEventListener('click', () => {
      closeUserPanel();
      window.location.href = './profile.html';
    });

    // 优惠券弹窗
    elements.couponBar.addEventListener('click', openCouponModal);
    elements.couponOverlay.addEventListener('click', closeCouponModal);
    elements.closeCouponBtn.addEventListener('click', closeCouponModal);
    elements.noCouponBtn.addEventListener('click', useNoCoupon);

    // 目的地输入
    elements.toInput.addEventListener('focus', () => toggleDestinationPanel(true));
    elements.clearToBtn.addEventListener('click', () => {
      resetOrderState();
    });

    // 搜索输入
    elements.searchInput.addEventListener('input', (e) => {
      const query = e.target.value.toLowerCase();
      const items = elements.panelContent.querySelectorAll('.suggestion-item');
      items.forEach(item => {
        const title = item.querySelector('.suggestion-title').textContent.toLowerCase();
        const desc = item.querySelector('.suggestion-desc').textContent.toLowerCase();
        if (title.includes(query) || desc.includes(query)) {
          item.style.display = 'flex';
        } else {
          item.style.display = 'none';
        }
      });
    });

    // 建议项点击
    elements.panelContent.addEventListener('click', (e) => {
      const item = e.target.closest('.suggestion-item');
      if (item) {
        const title = item.querySelector('.suggestion-title').textContent;
        const desc = item.querySelector('.suggestion-desc').textContent;
        selectDestination(title, desc);
      }
    });

    // 车型选择
    elements.carTypes.forEach(el => {
      el.addEventListener('click', () => selectCarType(el));
    });

    // 确认呼叫
    elements.confirmBtn.addEventListener('click', confirmRide);

    // 定位
    elements.locationBtn.addEventListener('click', autoLocate);
    elements.fromLocation.addEventListener('click', autoLocate);

    // 底部导航
    elements.navItems.forEach((item, index) => {
      item.addEventListener('click', () => {
        elements.navItems.forEach(i => i.classList.remove('active'));
        item.classList.add('active');
        
        if (index === 1) {
          window.location.href = './order.html';
        } else if (index === 2) {
          window.location.href = './coupons.html';
        } else if (index === 3) {
          window.location.href = './profile.html';
        }
      });
    });

    // 点击其他地方关闭面板
    document.addEventListener('click', (e) => {
      if (!e.target.closest('.order-card') && !e.target.closest('.usermodel-panel')) {
        toggleDestinationPanel(false);
      }
    });
  }

  // 初始化
  async function init() {
    initState();
    updateUserUI();
    bindEvents();
    
    // 初始化地图
    await initMap();
    
    // 加载用户信息
    if (state.isLoggedIn) {
      loadUserInfo();
      // 初始化WebSocket连接
      if (typeof window.initWebSocket === 'function') {
        window.initWebSocket(state.userId);
      }
    }
  }

  // 启动
  init();
})();
