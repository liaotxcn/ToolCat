import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // 先配置更具体的路径规则，确保它们优先匹配
      // 为LLM聊天API添加单独的代理配置
      '/plugins/LLMChat/api/chat': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        timeout: 180000, // 为聊天API设置3分钟超时，确保足够处理复杂请求
        // 添加keep-alive和重试配置来防止socket hang up
        headers: {
          Connection: 'keep-alive'
        },
        followRedirects: true,
        rewrite: (path) => path
      },
      // LLM相关其他API也使用较长超时
      '/plugins/LLMChat': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        timeout: 120000, // 2分钟超时
        headers: {
          Connection: 'keep-alive'
        }
      },
      // 其他代理配置
      '/auth': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        timeout: 120000
      },
      '/plugins': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        timeout: 120000
      },
      '/api/v1': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        timeout: 120000
      },
      '/health': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        timeout: 120000
      }
    },
    // 增加服务器响应头，改善连接稳定性
    headers: {
      'Connection': 'keep-alive',
      'Keep-Alive': 'timeout=65'
    },
    // 增加服务器超时设置
    timeout: 180000
  }
})
