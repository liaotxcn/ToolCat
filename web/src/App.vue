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
              <button class="menu-trigger" @click="showMenu = !showMenu">
                欢迎，{{ currentUser?.username }} ▾
              </button>
              <div v-if="showMenu" class="dropdown">
                <button @click="handleMenuSelect('teams')">协作团队</button>
                <button @click="handleMenuSelect('personal')">个人中心</button>
                <button @click="handleMenuSelect('logout')">退出登录</button>
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
          <div class="footer-left">
            <span class="brand">ToolCat</span>
            <span class="sep">·</span>
          </div>
          <div class="footer-right">
            <span class="version">v{{ appVersion }}</span>
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
  padding: var(--space-4) 0;
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
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--space-1);
  background: linear-gradient(to right, #ffffff, rgba(255, 255, 255, 0.8));
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  letter-spacing: -0.025em;
}

.header-content p {
  font-size: var(--font-size-sm);
  opacity: 0.9;
  margin: 0;
}

/* 用户菜单样式 */
.user-menu {
  position: relative;
}

.menu-trigger {
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(4px);
  border: 1px solid rgba(255, 255, 255, 0.25);
  color: #fff;
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-lg);
  cursor: pointer;
  font-size: var(--font-size-sm);
  transition: var(--transition-all);
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-weight: var(--font-weight-medium);
}

.menu-trigger:hover {
  background: rgba(255, 255, 255, 0.25);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.menu-trigger:active {
  transform: translateY(0);
}

.dropdown {
  position: absolute;
  right: 0;
  top: calc(100% + var(--space-2));
  background: var(--bg-primary);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  min-width: 200px;
  padding: var(--space-1);
  z-index: 200;
  opacity: 0;
  transform: translateY(-4px);
  animation: fadeInUp 0.2s ease forwards;
}

@keyframes fadeInUp {
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

.dropdown button:active {
  transform: scale(0.98);
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
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-medium);
  color: var(--text-tertiary);
  margin-top: auto;
  padding: var(--space-4) 0;
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
  color: var(--primary-700);
  font-size: var(--font-size-base);
}

.footer-left .sep {
  color: var(--border-medium);
  display: inline;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.version {
  font-size: var(--font-size-sm);
  color: var(--primary-700);
  background: var(--primary-50);
  border: 1px solid var(--primary-200);
  border-radius: var(--radius-full);
  padding: var(--space-1) var(--space-2);
  font-weight: var(--font-weight-medium);
}

.github-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  background: var(--bg-primary);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-full);
  width: 36px;
  height: 36px;
  transition: var(--transition-all);
}

.github-link:hover {
  color: var(--primary-700);
  border-color: var(--primary-300);
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
    gap: var(--space-2);
  }
  
  .footer-left {
    flex-direction: column;
    gap: var(--space-1);
  }
  
  .footer-left .sep {
    display: none;
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
}
</style>

