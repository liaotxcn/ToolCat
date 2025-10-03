<script setup>
import { ref, onMounted, computed } from 'vue'
import PluginRenderer from './components/PluginRenderer.vue'
import AuthContainer from './components/AuthContainer.vue'
import pluginManager from './pluginManager.js'
import HelloPlugin from './plugins/HelloPlugin.js'
import NotePlugin from './plugins/NotePlugin.js'
import { authService } from './services/auth.js'

// 认证状态
const isAuthenticated = ref(false)
const currentUser = ref(null)

// 可用插件列表
const availablePlugins = ref([])
const selectedPlugin = ref(null)

// 初始化应用
onMounted(() => {
  // 检查用户是否已登录
  checkAuthentication()
  
  // 注册插件
  registerPlugins()
  
  // 获取可用插件信息
  updateAvailablePlugins()
})

// 注册插件
const registerPlugins = () => {
  // 注册Hello插件
  const helloPlugin = new HelloPlugin()
  pluginManager.registerPlugin('hello', helloPlugin)
  
  // 注册Note插件
  const notePlugin = new NotePlugin()
  pluginManager.registerPlugin('note', notePlugin)
  
  // 这里可以注册更多插件
}

// 更新可用插件列表
const updateAvailablePlugins = () => {
  const plugins = pluginManager.getAllPlugins()
  availablePlugins.value = Object.keys(plugins).map(key => {
    const plugin = plugins[key]
    return {
      name: key,
      info: plugin.getInfo()
    }
  })
  
  // 默认选择第一个插件
  if (availablePlugins.value.length > 0 && !selectedPlugin.value) {
    selectedPlugin.value = availablePlugins.value[0].name
  }
}

// 选择插件
const selectPlugin = (pluginName) => {
  selectedPlugin.value = pluginName
}

// 检查用户认证状态
const checkAuthentication = () => {
  if (authService.isAuthenticated()) {
    isAuthenticated.value = true
    currentUser.value = authService.getCurrentUser()
  }
}

// 处理认证成功
const handleAuthSuccess = () => {
  checkAuthentication()
}

// 处理用户登出
const handleLogout = () => {
  authService.logout()
  isAuthenticated.value = false
  currentUser.value = null
  selectedPlugin.value = null
}
</script>

<template>
  <div class="app">
    <!-- 用户未登录时显示登录/注册界面 -->
    <AuthContainer v-if="!isAuthenticated" @auth-success="handleAuthSuccess" />
    
    <!-- 用户已登录时显示主应用界面 -->
    <template v-else>
      <header class="app-header">
        <div class="header-content">
          <div>
            <h1>ToolCat</h1>
            <p>集成多种工具的高效工具箱</p>
          </div>
          <div class="user-info">
            <span class="welcome-message">欢迎，{{ currentUser?.username }}</span>
            <button class="logout-btn" @click="handleLogout">退出登录</button>
          </div>
        </div>
      </header>
      
      <main class="app-main">
        <div class="plugin-selection">
          <h2>插件列表</h2>
          <div class="plugin-list">
            <button 
              v-for="pluginInfo in availablePlugins" 
              :key="pluginInfo.name"
              @click="selectPlugin(pluginInfo.name)"
              :class="{ 'active': selectedPlugin === pluginInfo.name }"
            >
              {{ pluginInfo.name }}
            </button>
          </div>
        </div>
        
        <div class="plugin-content">
          <h2>插件内容</h2>
          <div class="plugin-renderer-container">
            <PluginRenderer 
              v-if="selectedPlugin"
              :plugin-name="selectedPlugin"
              :plugin-manager="pluginManager"
            />
            <div v-else class="no-plugin-selected">
              请选择一个插件
            </div>
          </div>
        </div>
      </main>
      
      <footer class="app-footer">
        <p>ToolCat - 集成多种工具的高效工具箱</p>
      </footer>
    </template>
  </div>
</template>

<style>
/* 应用主样式 */
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background-color: #f5f5f5;
  color: #333;
  line-height: 1.6;
}

.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1rem 0;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.welcome-message {
  font-size: 1rem;
  opacity: 0.9;
}

.logout-btn {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
}

.app-header h1 {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.app-header p {
  font-size: 1.1rem;
  opacity: 0.9;
}

.app-main {
  flex: 1;
  display: flex;
  max-width: 1200px;
  margin: 2rem auto;
  width: 100%;
  gap: 2rem;
  padding: 0 1rem;
}

.plugin-selection {
  width: 250px;
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.plugin-selection h2 {
  font-size: 1.3rem;
  margin-bottom: 1rem;
  color: #333;
  border-bottom: 2px solid #667eea;
  padding-bottom: 0.5rem;
}

.plugin-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.plugin-list button {
  padding: 0.8rem 1rem;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  text-align: left;
  transition: all 0.3s ease;
  font-size: 1rem;
}

.plugin-list button:hover {
  background: #f8f9fa;
  border-color: #667eea;
}

.plugin-list button.active {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.plugin-content {
  flex: 1;
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.plugin-content h2 {
  font-size: 1.3rem;
  margin-bottom: 1rem;
  color: #333;
  border-bottom: 2px solid #667eea;
  padding-bottom: 0.5rem;
}

.plugin-renderer-container {
  min-height: 400px;
  padding: 1rem;
  border: 1px solid #eee;
  border-radius: 4px;
  background: #fafafa;
}

.no-plugin-selected {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #999;
  font-size: 1.1rem;
}

.app-footer {
  background: #333;
  color: white;
  text-align: center;
  padding: 1rem 0;
  margin-top: auto;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-main {
    flex-direction: column;
  }
  
  .plugin-selection {
    width: 100%;
  }
  
  .plugin-list {
    flex-direction: row;
    flex-wrap: wrap;
  }
  
  .plugin-list button {
    flex: 1;
    min-width: 120px;
    text-align: center;
  }
}
</style>
