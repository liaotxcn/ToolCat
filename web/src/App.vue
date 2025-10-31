<script setup>
import { ref, onMounted, computed } from 'vue'
import PluginRenderer from './components/PluginRenderer.vue'
import AuthContainer from './components/AuthContainer.vue'
import pluginManager from './pluginManager.js'
import HelloPlugin from './plugins/HelloPlugin.js'
import NotePlugin from './plugins/NotePlugin.js'
import { authService } from './services/auth.js'
import UserCenter from './components/UserCenter.vue'
import TeamsCenter from './components/TeamsCenter.vue'
const appVersion = '1.0.0'

// 认证状态
const isAuthenticated = ref(false)
const currentUser = ref(null)
const showMenu = ref(false)

// 可用插件列表
const availablePlugins = ref([])
const selectedPlugin = ref(null)
const selectedSection = ref('plugins')
const pluginRendererRef = ref(null)
// 新增：侧栏搜索关键词与过滤列表
const pluginKeyword = ref('')
const filteredPlugins = computed(() => {
  const kw = pluginKeyword.value.trim().toLowerCase()
  if (!kw) return availablePlugins.value
  return availablePlugins.value.filter(p => {
    const name = (p.name || '').toLowerCase()
    const desc = (p.info?.description || '').toLowerCase()
    return name.includes(kw) || desc.includes(kw)
  })
})

// 引用插件渲染器（重复声明已移除）

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
  selectedSection.value = 'plugins'
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
  // 登录成功后刷新当前插件，使其加载用户数据
  if (pluginRendererRef.value && typeof pluginRendererRef.value.refreshPlugin === 'function') {
    pluginRendererRef.value.refreshPlugin()
  }
}

// 处理用户登出
const handleLogout = () => {
  authService.logout()
  isAuthenticated.value = false
  currentUser.value = null
  selectedPlugin.value = null
}
const handleMenuSelect = (key) => {
  if (key === 'teams') {
    selectedSection.value = 'teams'
    selectedPlugin.value = null
  } else if (key === 'personal') {
    selectedSection.value = 'personal'
    selectedPlugin.value = null
  } else if (key === 'logout') {
    handleLogout()
  }
  showMenu.value = false
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
            <p>高性能、高效率、插件化易扩展的插件开发/服务聚合平台</p>
          </div>
          <div class="user-info">
            <div class="user-menu">
              <button 
                class="menu-trigger"
                :class="{ 'active': showMenu }"
                @click="showMenu = !showMenu"
                aria-expanded="showMenu"
                aria-haspopup="true"
                aria-controls="user-dropdown"
              >
                <span class="user-avatar">
                  <svg viewBox="0 0 24 24" width="20" height="20" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M12 12C14.2091 12 16 10.2091 16 8C16 5.79086 14.2091 4 12 4C9.79086 4 8 5.79086 8 8C8 10.2091 9.79086 12 12 12Z" stroke="currentColor" stroke-width="1.5"/>
                    <path d="M19 20C19 16.134 12 14 12 14C12 14 5 16.134 5 20V21H19V20Z" stroke="currentColor" stroke-width="1.5"/>
                  </svg>
                </span>
                <span class="user-name">欢迎，{{ currentUser?.username }}</span>
                <span class="dropdown-arrow" :class="{ 'rotate': showMenu }">
                  <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M6 9L12 15L18 9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </span>
              </button>
              <div 
                v-show="showMenu" 
                class="dropdown"
                id="user-dropdown"
                role="menu"
                aria-labelledby="user-menu"
              >
                <button 
                  @click="handleMenuSelect('teams')"
                  class="dropdown-item"
                  role="menuitem"
                >
                  <svg viewBox="0 0 24 24" width="18" height="18" class="item-icon" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M17 21V19C17 17.9391 16.5786 16.9217 15.8284 16.1716C15.0783 15.4214 14.0609 15 13 15H5C3.93913 15 2.92172 15.4214 2.17157 16.1716C1.42143 16.9217 1 17.9391 1 19V21" stroke="currentColor" stroke-width="1.5"/>
                    <path d="M9 13C11.2091 13 13 11.2091 13 9C13 6.79086 11.2091 5 9 5C6.79086 5 5 6.79086 5 9C5 11.2091 6.79086 13 9 13Z" stroke="currentColor" stroke-width="1.5"/>
                    <path d="M23 21V19C23 17.9391 22.5786 16.9217 21.8284 16.1716C21.0783 15.4214 20.0609 15 19 15C18.6118 15 18.2379 15.0583 17.8858 15.1716C17.9395 15.3562 17.9757 15.548 17.994 15.7456C18.7175 15.4198 19.5183 15.25 20.34 15.25C21.9503 15.25 23.375 15.8893 24 17V17C23.375 18.1107 21.9503 18.75 20.34 18.75C19.5183 18.75 18.7175 18.5802 17.994 18.2544C17.9757 18.452 17.9395 18.6438 17.8858 18.8284C18.2379 18.9417 18.6118 19 19 19C20.0609 19 21.0783 18.5786 21.8284 17.8284C22.5786 17.0783 23 16.0609 23 15" stroke="currentColor" stroke-width="1.5"/>
                  </svg>
                  <span>协作团队</span>
                </button>
                <button 
                  @click="handleMenuSelect('personal')"
                  class="dropdown-item"
                  role="menuitem"
                >
                  <svg viewBox="0 0 24 24" width="18" height="18" class="item-icon" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M20 21V19C20 17.9391 19.5786 16.9217 18.8284 16.1716C18.0783 15.4214 17.0609 15 16 15H8C6.93913 15 5.92172 15.4214 5.17157 16.1716C4.42143 16.9217 4 17.9391 4 19V21" stroke="currentColor" stroke-width="1.5"/>
                    <path d="M12 11C14.2091 11 16 9.20914 16 7C16 4.79086 14.2091 3 12 3C9.79086 3 8 4.79086 8 7C8 9.20914 9.79086 11 12 11Z" stroke="currentColor" stroke-width="1.5"/>
                  </svg>
                  <span>个人中心</span>
                </button>
                <div class="dropdown-divider"></div>
                <button 
                  @click="handleMenuSelect('logout')"
                  class="dropdown-item logout-item"
                  role="menuitem"
                >
                  <svg viewBox="0 0 24 24" width="18" height="18" class="item-icon" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M9 21H5C3.89543 21 3 20.1046 3 19V5C3 3.89543 3.89543 3 5 3H9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M16 17L21 12L16 7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M21 12H9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  <span>退出登录</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </header>
      
      <main class="app-main">
        <div class="plugin-selection">
          <h2>Plugins/Service</h2>
          <!-- 新增：侧栏工具（搜索 + 统计） -->
          <div class="sidebar-tools">
            <input v-model="pluginKeyword" type="text" class="sidebar-search" placeholder="搜索插件..." />
            <span class="sidebar-count">共 {{ availablePlugins.length }}，匹配 {{ filteredPlugins.length }}</span>
          </div>
          <div class="plugin-list">
            <!-- 插件列表：使用过滤后的列表，并展示描述与计数 -->
            <button 
              v-for="pluginInfo in filteredPlugins" 
              :key="pluginInfo.name"
              @click="selectPlugin(pluginInfo.name)"
              :class="{ 'active': selectedSection === 'plugins' && selectedPlugin === pluginInfo.name }"
              :title="pluginInfo.info?.description || pluginInfo.name"
              aria-label="选择插件"
            >
              <span class="item-title">{{ pluginInfo.name }}</span>
              <span v-if="pluginInfo.info?.description" class="item-desc">{{ pluginInfo.info.description }}</span>
              <span v-if="pluginInfo.info?.noteCount !== undefined" class="item-badge">{{ pluginInfo.info.noteCount }}</span>
            </button>
          </div>
        </div>
        
        <div class="plugin-content">
          <h2>{{ selectedSection === 'personal' ? '个人中心' : (selectedSection === 'teams' ? '协作团队' : 'Content') }}</h2>
          <div class="plugin-renderer-container">
            <transition name="plugin-switch" mode="out-in" appear>
              <TeamsCenter v-if="selectedSection === 'teams'" />
              <UserCenter 
                v-else-if="selectedSection === 'personal'"
                :current-user="currentUser"
                @updated-user="currentUser = $event"
              />
              <PluginRenderer 
                v-else-if="selectedPlugin"
                :plugin-name="selectedPlugin"
                :plugin-manager="pluginManager"
                ref="pluginRendererRef"
              />
              <div v-else class="no-plugin-selected">
                请选择一个插件/服务
              </div>
            </transition>
          </div>
        </div>
      </main>
      
      <footer class="app-footer">
        <div class="footer-content">
          <!-- <div class="footer-left">
            <span class="brand">ToolCat</span>
            <span class="sep">·</span>
          </div> -->
          <div class="footer-right">
            <span class="version">ToolCat v{{ appVersion }}</span>
            <a class="github-link" href="https://github.com/liaotxcn/ToolCat" target="_blank" rel="noopener noreferrer" aria-label="GitHub">
              <svg class="github-icon" viewBox="0 0 16 16" width="18" height="18" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38
                0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53
                .63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95
                0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27
                .68 0 1.36.09 2 .27 1.53-1.03 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15
                0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.19 0 .21.15.46.55.38
                A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
              </svg>
            </a>
          </div>
        </div>
      </footer>
    </template>
  </div>
</template>

<style scoped>
/* 应用主样式 */
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* 头部样式 */
.app-header {
  background: linear-gradient(135deg, var(--primary-600) 0%, var(--primary-800) 100%);
  color: white;
  padding: var(--space-6) 0;
  box-shadow: var(--shadow-md);
  position: sticky;
  top: 0;
  z-index: 100;
  backdrop-filter: blur(8px);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--space-4);
}

.header-content h1 {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--space-1);
  background: linear-gradient(to right, #ffffff, rgba(255, 255, 255, 0.8));
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  letter-spacing: -0.025em;
}

.header-content p {
  font-size: var(--font-size-base);
  opacity: 0.9;
  margin: 0;
}

/* 用户菜单样式 */
.user-menu {
  position: relative;
}

.menu-trigger {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-2) var(--space-3);
  background-color: rgba(255, 255, 255, 0.1);
  border: 1px solid transparent;
  border-radius: var(--radius-lg);
  color: white;
  cursor: pointer;
  transition: var(--transition-all);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  backdrop-filter: blur(4px);
}

.menu-trigger:hover {
  background-color: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.menu-trigger.active {
  background-color: rgba(255, 255, 255, 0.25);
  border-color: rgba(255, 255, 255, 0.4);
}

.menu-trigger:active {
  transform: translateY(0);
}

.user-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  flex-shrink: 0;
}

.user-name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 150px;
}

.dropdown-arrow {
  transition: transform 0.2s ease;
  flex-shrink: 0;
}

.dropdown-arrow.rotate {
  transform: rotate(180deg);
}

.dropdown {
  position: absolute;
  right: 0;
  top: calc(100% + var(--space-2));
  margin-top: var(--space-1);
  background: var(--bg-primary);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-xl);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  min-width: 200px;
  z-index: 200;
  animation: fadeInUp 0.2s ease forwards;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dropdown::before {
  content: '';
  position: absolute;
  right: var(--space-3);
  top: -6px;
  width: 12px;
  height: 12px;
  background: var(--bg-primary);
  border-top: 1px solid var(--border-medium);
  border-left: 1px solid var(--border-medium);
  transform: rotate(45deg);
}

.dropdown button {
  width: 100%;
  text-align: left;
  background: var(--bg-primary);
  border: none;
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-lg);
  color: var(--text-primary);
  cursor: pointer;
  font-size: var(--font-size-sm);
  transition: var(--transition-all);
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.dropdown button:hover {
  background: var(--bg-secondary);
  color: var(--primary-600);
}

.dropdown button:focus {
  outline: 2px solid var(--primary-600);
  outline-offset: -2px;
}

.dropdown button:active {
  transform: scale(0.98);
}

.item-icon {
  flex-shrink: 0;
  opacity: 0.7;
}

.dropdown button:hover .item-icon {
  opacity: 1;
}

.logout-item {
  color: var(--danger-600);
}

.logout-item:hover {
  background-color: var(--danger-50);
  color: var(--danger-700);
}

/* 主内容区域 */
.app-main {
  flex: 1;
  display: flex;
  max-width: 1200px;
  margin: var(--space-6) auto;
  width: 100%;
  gap: var(--space-6);
  padding: 0 var(--space-4);
}

/* 插件选择区域 */
.plugin-selection {
  width: 280px;
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  padding: var(--space-4);
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-light);
  transition: var(--transition-all);
  position: sticky;
  top: calc(var(--space-4) * 5);
  height: fit-content;
  max-height: calc(100vh - 120px);
  overflow-y: auto;
}

.plugin-selection:hover {
  box-shadow: var(--shadow-lg);
  transform: translateY(-2px);
}

.plugin-selection h2 {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-3);
  color: var(--text-primary);
}

/* 侧栏工具区样式 */
.sidebar-tools {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  margin-bottom: var(--space-4);
}

.sidebar-search {
  width: 100%;
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-lg);
  font-size: var(--font-size-sm);
  background: var(--bg-secondary);
  transition: var(--transition-all);
}

.sidebar-search:focus {
  outline: none;
  border-color: var(--primary-500);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
  background: var(--bg-primary);
}

.sidebar-count {
  color: var(--text-tertiary);
  font-size: var(--font-size-xs);
  text-align: center;
  padding: var(--space-1) 0;
}

/* 插件列表样式 */
.plugin-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.plugin-list button {
  position: relative;
  padding: var(--space-3);
  border: 1px solid var(--border-light);
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  text-align: left;
  transition: var(--transition-all);
  font-size: var(--font-size-sm);
  color: var(--text-primary);
  display: grid;
  grid-template-columns: 1fr auto;
  grid-template-rows: auto auto;
  row-gap: 2px;
  overflow: hidden;
}

.plugin-list button::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 3px;
  background: transparent;
  transition: var(--transition-colors);
}

.plugin-list button:hover {
  background: var(--bg-secondary);
  border-color: var(--primary-300);
  transform: translateX(2px);
}

.plugin-list button.active {
  background: var(--primary-50);
  color: var(--primary-700);
  border-color: var(--primary-200);
  box-shadow: var(--shadow-sm);
}

.plugin-list button.active::before {
  background: var(--primary-500);
}

.plugin-list button:active {
  transform: translateY(1px);
}

.item-title {
  grid-column: 1 / 2;
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-base);
}

.item-desc {
  grid-column: 1 / 2;
  color: var(--text-tertiary);
  font-size: var(--font-size-xs);
  line-height: 1.4;
}

.item-badge {
  grid-column: 2 / 3;
  align-self: start;
  justify-self: end;
  background: var(--primary-100);
  color: var(--primary-700);
  border: 1px solid var(--primary-200);
  border-radius: var(--radius-full);
  padding: 0 var(--space-2);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  min-width: 20px;
  text-align: center;
}

/* 插件内容区域 */
.plugin-content {
  flex: 1;
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  padding: var(--space-5);
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-light);
  min-height: 500px;
}

.plugin-content h2 {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-4);
  color: var(--text-primary);
}

/* 插件渲染器容器 */
.plugin-renderer-container {
  min-height: 420px;
  padding: var(--space-4);
  border: 2px dashed var(--border-medium);
  border-radius: var(--radius-lg);
  background: var(--bg-secondary);
  transition: var(--transition-all);
  position: relative;
  overflow: hidden;
}

.plugin-renderer-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--primary-400), var(--primary-600), var(--primary-400));
  transform: scaleX(0);
  transform-origin: left;
  transition: transform 0.3s ease;
}

.plugin-renderer-container:hover {
  background: var(--bg-overlay);
  border-color: var(--border-dark);
  box-shadow: var(--shadow-sm) inset;
}

.plugin-renderer-container:hover::before {
  transform: scaleX(1);
}

/* 内容切换动画 */
.plugin-switch-enter-active, 
.plugin-switch-leave-active {
  transition: opacity var(--transition-normal), transform var(--transition-normal);
}

.plugin-switch-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

.plugin-switch-enter-to {
  opacity: 1;
  transform: translateY(0);
}

.plugin-switch-leave-from {
  opacity: 1;
  transform: translateY(0);
}

.plugin-switch-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* 无插件选择状态 */
.no-plugin-selected {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  color: var(--text-muted);
  font-size: var(--font-size-lg);
  text-align: center;
  gap: var(--space-3);
}

.no-plugin-selected::before {
  content: '🔍';
  font-size: 3rem;
  opacity: 0.5;
}

/* 页脚样式 */
.app-footer {
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-700) 100%);
  border-top: 3px solid var(--primary-400);
  color: white;
  margin-top: auto;
  padding: var(--space-2) 0;
  position: relative;
  overflow: hidden;
}

.app-footer::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at 30% 30%, rgba(255, 255, 255, 0.15) 0%, transparent 50%),
              radial-gradient(circle at 70% 70%, rgba(255, 255, 255, 0.1) 0%, transparent 60%);
  z-index: 0;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--space-4);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  gap: var(--space-3);
  position: relative;
  z-index: 1;
}

.footer-left {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  flex-wrap: wrap;
  justify-content: center;
}

.footer-left .brand {
  font-weight: var(--font-weight-bold);
  color: white;
  font-size: var(--font-size-lg);
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
  transition: transform 0.2s ease;
}

.footer-left .brand:hover {
  transform: scale(1.05);
}

.footer-left .sep {
  color: rgba(255, 255, 255, 0.7);
  display: inline;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.version {
  font-size: var(--font-size-sm);
  color: var(--primary-900);
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.7);
  border-radius: var(--radius-full);
  padding: var(--space-1) var(--space-2);
  font-weight: var(--font-weight-medium);
  backdrop-filter: blur(4px);
  transition: all 0.2s ease;
}

.version:hover {
  background: white;
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.github-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: white;
  background: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: var(--radius-full);
  width: 36px;
  height: 36px;
  transition: var(--transition-all);
  backdrop-filter: blur(4px);
}

.github-link:hover {
  color: var(--primary-900);
  background: rgba(255, 255, 255, 0.9);
  border-color: rgba(255, 255, 255, 0.7);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.github-icon {
  display: block;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .app-main {
    gap: var(--space-4);
  }
  
  .plugin-selection {
    width: 250px;
  }
}

@media (max-width: 768px) {
  .app-main {
    flex-direction: column;
    margin: var(--space-4) auto;
    gap: var(--space-4);
  }
  
  .plugin-selection {
    width: 100%;
    position: static;
    max-height: none;
  }
  
  .header-content {
    padding: 0 var(--space-3);
  }
  
  .header-content h1 {
    font-size: var(--font-size-xl);
  }
  
  .plugin-content {
    padding: var(--space-3);
  }
  
  .plugin-list button {
    flex: 1;
    min-width: 100%;
  }
  
  .footer-content {
    gap: var(--space-3);
    padding: 0 var(--space-3);
  }
  
  .footer-left {
    flex-direction: column;
    gap: var(--space-2);
  }
  
  .footer-left .sep {
    display: none;
  }
  
  .app-footer {
     padding: var(--space-1) 0;
   }
}

@media (max-width: 480px) {
  .app-header {
    padding: var(--space-3) 0;
  }
  
  .app-main {
    padding: 0 var(--space-2);
  }
  
  .plugin-selection,
  .plugin-content {
    padding: var(--space-3);
    border-radius: var(--radius-lg);
  }
  
  .menu-trigger {
    font-size: var(--font-size-xs);
    padding: var(--space-1) var(--space-2);
  }
  
  .dropdown {
    min-width: 160px;
  }
  
  .dropdown button {
    font-size: var(--font-size-xs);
    padding: var(--space-2) var(--space-2);
  }
  
  .footer-content {
    padding: 0 var(--space-2);
  }
  
  .app-footer {
     padding: var(--space-1) 0;
   }
  
  .footer-left .brand {
    font-size: var(--font-size-base);
  }
}
</style>

