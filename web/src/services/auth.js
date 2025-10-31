import axios from 'axios'

// 创建axios实例
const api = axios.create({
  // 支持环境变量配置后端API地址，开发环境为空则使用Vite代理
  baseURL: (typeof import.meta !== 'undefined' && import.meta.env && import.meta.env.VITE_API_BASE) ? import.meta.env.VITE_API_BASE : '',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true
})

console.log('Axios实例已创建，baseURL:', api.defaults.baseURL)

// 小工具：确保已获取CSRF Cookie
async function ensureCsrf() {
  const hasCookie = /(?:^|; )XSRF-TOKEN=([^;]+)/.test(document.cookie)
  if (!hasCookie) {
    try {
      await api.get('/health', { withCredentials: true })
    } catch (_) {}
  }
}

// 请求拦截器 - 添加token 与 CSRF
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    // 为需要修改服务器状态的请求自动附加CSRF头
    const method = (config.method || 'get').toLowerCase()
    if (['post', 'put', 'delete', 'patch'].includes(method)) {
      const match = document.cookie.match(/(?:^|; )XSRF-TOKEN=([^;]+)/)
      if (match) {
        config.headers['X-CSRF-Token'] = decodeURIComponent(match[1])
      } else {
        // 当尚未从Cookie中拿到时，尝试使用已缓存的令牌
        const cached = sessionStorage.getItem('csrf_token')
        if (cached) {
          config.headers['X-CSRF-Token'] = cached
        }
      }
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理token过期等情况，并捕获CSRF令牌
api.interceptors.response.use(
  response => {
    // 捕获响应头中的CSRF令牌，作为备份缓存
    try {
      const csrfHeader = response.headers && (response.headers['x-csrf-token'] || response.headers['X-CSRF-Token'])
      if (csrfHeader) {
        sessionStorage.setItem('csrf_token', csrfHeader)
      }
    } catch (_) {}
    console.log('API请求成功:', response.config.url, response.data)
    return response.data
  },
  error => {
    if (error.response) {
      // 处理HTTP错误
      console.error('API请求错误:', error.config?.url, error.response.status, error.response.data)
      switch (error.response.status) {
        case 401:
          // token过期或无效，清除localStorage中的token
          localStorage.removeItem('token')
          localStorage.removeItem('userInfo')
          break
        default:
          // 其他错误情况
          console.error('API请求错误:', error.response.data)
      }
    } else if (error.request) {
      // 请求发出但没有收到响应
      console.error('API请求无响应:', error.config?.url)
      console.error('网络错误详情:', error.message)
    } else {
      // 设置请求时发生错误
      console.error('API请求配置错误:', error.message)
    }
    return Promise.reject(error)
  }
)

// 认证相关API方法
export const authService = {
  // 用户注册
  register: async (userData) => {
    try {
      await ensureCsrf()
      const response = await api.post('/auth/register', userData)
      return response
    } catch (error) {
      // 针对CSRF或第一次缺Cookie的情况重试一次
      if (error?.response?.status === 403) {
        try {
          await ensureCsrf()
          const response = await api.post('/auth/register', userData)
          return response
        } catch (e2) {
          console.error('注册失败(重试后):', e2)
          throw e2
        }
      }
      console.error('注册失败:', error)
      throw error
    }
  },

  // 用户登录
  login: async (credentials) => {
    try {
      await ensureCsrf()
      const response = await api.post('/auth/login', credentials)
      // 保存token和用户信息到localStorage（后端返回access_token与refresh_token）
      if (response.access_token && response.user) {
        localStorage.setItem('token', response.access_token)
        if (response.refresh_token) {
          localStorage.setItem('refresh_token', response.refresh_token)
        }
        localStorage.setItem('userInfo', JSON.stringify(response.user))
      }
      return response
    } catch (error) {
      // 针对CSRF或第一次缺Cookie的情况重试一次
      if (error?.response?.status === 403) {
        try {
          await ensureCsrf()
          const response = await api.post('/auth/login', credentials)
          if (response.access_token && response.user) {
            localStorage.setItem('token', response.access_token)
            if (response.refresh_token) {
              localStorage.setItem('refresh_token', response.refresh_token)
            }
            localStorage.setItem('userInfo', JSON.stringify(response.user))
          }
          return response
        } catch (e2) {
          console.error('登录失败(重试后):', e2)
          throw e2
        }
      }
      console.error('登录失败:', error)
      throw error
    }
  },

  // 用户登出
  logout: () => {
    // 清除localStorage中的token和用户信息
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  },

  // 获取当前登录用户信息
  getCurrentUser: () => {
    const userInfo = localStorage.getItem('userInfo')
    return userInfo ? JSON.parse(userInfo) : null
  },

  // 检查用户是否已登录
  isAuthenticated: () => {
    const token = localStorage.getItem('token')
    return !!token
  }
}

export default api