import { createStore } from 'vuex'

export default createStore({
  state: {
    // web服务器地址（本地开发）
    backendUrl: 'http://127.0.0.1:8888',
    wsUrl: 'ws://127.0.0.1:8888',
    // 生产环境地址
    // backendUrl: 'https://123.56.164.220:8000',
    // wsUrl: 'wss://123.56.164.220:8000',
    // 信令服务器地址
    // signalUrl: 'wss://127.0.0.1:8001',
    userInfo: (sessionStorage.getItem('userInfo') && JSON.parse(sessionStorage.getItem('userInfo'))) || {},
    token: sessionStorage.getItem('token') || '',
    masterKey: null, // 主密钥（仅存储在内存中，不持久化）
    socket: null,
  },
  getters: {
  },
  mutations: {
    setUserInfo(state, userInfo) {
      state.userInfo = userInfo;
      sessionStorage.setItem('userInfo', JSON.stringify(userInfo));
      // 同时保存 token
      if (userInfo.token) {
        state.token = userInfo.token;
        sessionStorage.setItem('token', userInfo.token);
      }
    },
    setMasterKey(state, masterKey) {
      state.masterKey = masterKey;
      console.log('主密钥已设置到 Vuex store（仅内存）');
    },
    cleanUserInfo(state) {
      state.userInfo = {};
      state.token = '';
      state.masterKey = null; // 清除主密钥
      sessionStorage.removeItem('userInfo');
      sessionStorage.removeItem('token');
      console.log('用户信息和主密钥已清除');
    }
  },
  actions: {
  },
  modules: {
  }
})
