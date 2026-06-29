import { createStore } from 'vuex'

function trimTrailingSlash(value) {
  return (value || '').replace(/\/+$/, '')
}

function resolveBrowserOrigin() {
  if (typeof window === 'undefined' || !window.location) {
    return 'http://localhost:8000'
  }
  return window.location.origin
}

function toWebSocketOrigin(value) {
  try {
    const url = new URL(value)
    url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:'
    return trimTrailingSlash(url.origin)
  } catch (error) {
    return ''
  }
}

function resolveBackendUrls() {
  const backendUrl = trimTrailingSlash(
    process.env.VUE_APP_API_BASE_URL || resolveBrowserOrigin()
  )
  const wsUrl = trimTrailingSlash(
    process.env.VUE_APP_WS_BASE_URL || toWebSocketOrigin(backendUrl)
  )

  return {
    backendUrl,
    wsUrl,
  }
}

const { backendUrl, wsUrl } = resolveBackendUrls()

export default createStore({
  state: {
    backendUrl,
    wsUrl,
    userInfo: (sessionStorage.getItem('userInfo') && JSON.parse(sessionStorage.getItem('userInfo'))) || {},
    socket: null,
  },
  getters: {
  },
  mutations: {
    setUserInfo(state, userInfo) {
      state.userInfo = userInfo;
      sessionStorage.setItem('userInfo', JSON.stringify(userInfo));
    },
    cleanUserInfo(state) {
      state.userInfo = {};
      sessionStorage.removeItem('userInfo');
    }
  },
  actions: {
  },
  modules: {
  }
})
