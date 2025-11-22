<script setup>
import { ref, onMounted, computed } from 'vue'
import PluginRenderer from './components/PluginRenderer.vue'
import AuthContainer from './components/AuthContainer.vue'
import AIAssistant from './components/AIAssistant.vue'
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
  const authStatus = authService.isAuthenticated()
  console.log('App组件 - 认证状态检查结果:', authStatus)
  isAuthenticated.value = authStatus
  if (authStatus) {
    currentUser.value = authService.getCurrentUser()
    console.log('App组件 - 当前用户信息:', currentUser.value)
  }
}

// 调试功能：强制清除认证状态（可通过控制台调用）
window.clearWeaveAuth = () => {
  console.log('执行强制清除认证状态操作')
  authService.clearAuthData()
  isAuthenticated.value = false
  currentUser.value = null
  selectedPlugin.value = null
  console.log('认证状态已重置，请刷新页面查看登录界面')
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
  if (key === 'services') {
    selectedSection.value = 'services'
    selectedPlugin.value = null
  } else if (key === 'plugins') {
    selectedSection.value = 'plugins'
    // 如果没有选中插件，默认选择第一个
    if (!selectedPlugin.value && availablePlugins.value.length > 0) {
      selectedPlugin.value = availablePlugins.value[0].name
    }
  } else if (key === 'teams') {
    selectedSection.value = 'teams'
    selectedPlugin.value = null
  } else if (key === 'personal') {
    selectedSection.value = 'personal'
    selectedPlugin.value = null
  } else if (key === 'security') {
    selectedSection.value = 'security'
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
          <div class="brand-section">
            <div class="brand-container">
              <img src="/logo.png" alt="Weave Logo" class="brand-icon" />
              <div class="brand-text">
                <h1>Weave</h1>
                <p>高性能、高效率、安全可靠的插件开发/服务聚合平台</p>
              </div>
            </div>
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
                  <svg viewBox="0 0 24 24" width="24" height="24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M12 12C14.2091 12 16 10.2091 16 8C16 5.79086 14.2091 4 12 4C9.79086 4 8 5.79086 8 8C8 10.2091 9.79086 12 12 12Z" stroke="currentColor" stroke-width="1.5"/>
                    <path d="M19 20C19 16.134 12 14 12 14C12 14 5 16.134 5 20V21H19V20Z" stroke="currentColor" stroke-width="1.5"/>
                  </svg>
                </span>
                <span class="user-name">{{ currentUser?.username }}</span>
                <span class="dropdown-arrow" :class="{ 'rotate': showMenu }">
                  <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M6 9L12 15L18 9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </span>
              </button>
              <div 
                class="dropdown"
                :class="{ 'show': showMenu }"
                id="user-dropdown"
                role="menu"
                aria-labelledby="user-menu"
              >
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
        <aside class="sidebar">
          <!-- 主菜单 -->
          <nav class="sidebar-nav">
            <ul class="menu-list">
              <li>
                <button 
                  @click="handleMenuSelect('services')"
                  :class="{ 'active': selectedSection === 'services' }"
                  class="menu-item"
                >
                  <span class="menu-icon">🏢</span>
                  <span class="menu-label">Services</span>
                </button>
              </li>
              <li>
                <button 
                  @click="handleMenuSelect('plugins')"
                  :class="{ 'active': selectedSection === 'plugins' }"
                  class="menu-item"
                >
                  <span class="menu-icon">🔌</span>
                  <span class="menu-label">Plugins</span>
                </button>
              </li>
              <li>
                <button 
                  @click="handleMenuSelect('teams')"
                  :class="{ 'active': selectedSection === 'teams' }"
                  class="menu-item"
                >
                  <span class="menu-icon">👥</span>
                  <span class="menu-label">团队</span>
                </button>
              </li>
              <li>
                <button 
                  @click="handleMenuSelect('personal')"
                  :class="{ 'active': selectedSection === 'personal' }"
                  class="menu-item"
                >
                  <span class="menu-icon">👤</span>
                  <span class="menu-label">个人</span>
                </button>
              </li>
              <li>
                <button 
                  @click="handleMenuSelect('security')"
                  :class="{ 'active': selectedSection === 'security' }"
                  class="menu-item"
                >
                  <span class="menu-icon">🔒</span>
                  <span class="menu-label">安全中心</span>
                </button>
              </li>
            </ul>
          </nav>
        </aside>
        
        <!-- 右侧内容区域 -->
        <div class="content-area">
          <header class="content-header">
            <h2>
              {{ 
                selectedSection === 'services' ? 'Services' :
                selectedSection === 'plugins' ? 'Plugins' :
                selectedSection === 'teams' ? 'Treams' :
                selectedSection === 'personal' ? 'Personal' :
                selectedSection === 'security' ? 'Security' :
                '内容'
              }}
            </h2>
          </header>
          
          <main class="content-body">
            <transition name="content-switch" mode="out-in" appear>
              <!-- Services 内容 -->
              <div v-if="selectedSection === 'services'" class="services-content">
                <div class="service-card">
                  <h3>🏢 Service</h3>
                  <p>集研发、聚合、管理为一体</p>
                </div>
              </div>
              
              <!-- Plugins 内容 -->
              <div v-else-if="selectedSection === 'plugins'" class="plugins-content">
                <div class="plugins-layout">
                  <!-- 插件列表区域 -->
                  <div class="plugins-sidebar">
                    <div class="plugins-header">
                      <h3>插件列表</h3>
                      <div class="plugins-tools">
                        <input v-model="pluginKeyword" type="text" class="plugins-search" placeholder="搜索插件..." />
                        <span class="plugins-count">共 {{ availablePlugins.length }}，匹配 {{ filteredPlugins.length }}</span>
                      </div>
                    </div>
                    <div class="plugins-list">
                      <button 
                        v-for="pluginInfo in filteredPlugins" 
                        :key="pluginInfo.name"
                        @click="selectPlugin(pluginInfo.name)"
                        :class="{ 'active': selectedPlugin === pluginInfo.name }"
                        :title="pluginInfo.info?.description || pluginInfo.name"
                        class="plugins-item"
                      >
                        <span class="item-title">{{ pluginInfo.name }}</span>
                        <span v-if="pluginInfo.info?.description" class="item-desc">{{ pluginInfo.info.description }}</span>
                        <span v-if="pluginInfo.info?.noteCount !== undefined" class="item-badge">{{ pluginInfo.info.noteCount }}</span>
                      </button>
                    </div>
                  </div>
                  
                  <!-- 插件内容区域 -->
                  <div class="plugins-main">
                    <div v-if="selectedPlugin" class="plugin-renderer-container">
                      <PluginRenderer 
                        :plugin-name="selectedPlugin"
                        :plugin-manager="pluginManager"
                        ref="pluginRendererRef"
                      />
                    </div>
                    <div v-else class="no-plugin-selected">
                      <div class="empty-state">
                        <span class="empty-icon">🔌</span>
                        <h3>请选择一个插件</h3>
                        <p>从左侧插件列表中选择要使用的插件</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Teams 内容 -->
              <div v-else-if="selectedSection === 'teams'" class="teams-content">
                <TeamsCenter />
              </div>
              
              <!-- Personal 内容 -->
              <div v-else-if="selectedSection === 'personal'" class="personal-content">
                <UserCenter 
                  :current-user="currentUser"
                  @updated-user="currentUser = $event"
                />
              </div>
              
              <!-- Security 内容 -->
              <div v-else-if="selectedSection === 'security'" class="security-content">
                <div class="security-card">
                  <h3>🔒 安全中心</h3>
                  <p>权限管控、认证加密、沙盒环境等</p>
                </div>
              </div>
            </transition>
          </main>
        </div>
      </main>
      
      <footer class="app-footer">
        <div class="footer-content">
          <!-- <div class="footer-left">
            <span class="brand">Weave</span>
            <span class="sep">·</span>
          </div> -->
          <div class="footer-right">
            <span class="version">Weave v{{ appVersion }}</span>
            <a class="github-link" href="https://github.com/liaotxcn/Weave" target="_blank" rel="noopener noreferrer" aria-label="GitHub">
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
      
      <!-- AI智能助手（仅在登录后显示） -->
      <AIAssistant v-if="isAuthenticated" />
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

/* 头部样式 - 优化版本 */
.app-header {
  background: linear-gradient(135deg, 
    rgba(99, 102, 241, 0.95) 0%, 
    rgba(79, 70, 229, 0.95) 25%,
    rgba(67, 56, 202, 0.95) 75%,
    rgba(55, 48, 163, 0.95) 100%
  );
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  color: white;
  padding: var(--space-3) 0;
  box-shadow: 
    0 4px 20px rgba(0, 0, 0, 0.08),
    0 1px 3px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.app-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, 
    rgba(255, 255, 255, 0.1) 0%, 
    transparent 50%, 
    rgba(255, 255, 255, 0.05) 100%
  );
  pointer-events: none;
}

.app-header:hover {
  box-shadow: 
    0 8px 32px rgba(0, 0, 0, 0.12),
    0 2px 8px rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
  transform: translateY(-1px);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--space-5);
  min-height: 72px;
  position: relative;
  z-index: 1;
}

.brand-section {
  flex: 1;
  display: flex;
  align-items: center;
}

.brand-container {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  position: relative;
}

.brand-icon {
  width: 56px;
  height: 56px;
  object-fit: contain;
  border-radius: 12px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 
    0 4px 12px rgba(0, 0, 0, 0.15),
    0 2px 4px rgba(0, 0, 0, 0.1);
  position: relative;
  overflow: hidden;
}

.brand-icon::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, 
    rgba(255, 255, 255, 0.2) 0%, 
    transparent 50%
  );
  border-radius: 12px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.brand-icon:hover::before {
  opacity: 1;
}

.app-header:hover .brand-icon {
  transform: rotate(3deg) scale(1.08);
  box-shadow: 
    0 8px 24px rgba(0, 0, 0, 0.2),
    0 4px 8px rgba(0, 0, 0, 0.15);
}

.brand-text {
  display: flex;
  flex-direction: column;
  position: relative;
}

.brand-text h1 {
  font-size: 2rem;
  font-weight: var(--font-weight-bold);
  margin: 0;
  background: linear-gradient(135deg, 
    #ffffff 0%, 
    rgba(255, 255, 255, 0.95) 50%,
    rgba(255, 255, 255, 0.85) 100%
  );
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  letter-spacing: -0.02em;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  position: relative;
}

.brand-text h1::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background: linear-gradient(90deg, 
    rgba(255, 255, 255, 0.8) 0%, 
    transparent 100%
  );
  transition: width 0.3s ease;
}

.app-header:hover .brand-text h1::after {
  width: 100%;
}

.brand-highlight {
  position: relative;
  background: linear-gradient(135deg, 
    #ffd700 0%, 
    #ffb800 50%, 
    #ffa500 100%
  );
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  font-weight: 800;
  margin-left: 4px;
  text-shadow: 0 1px 3px rgba(255, 183, 0, 0.3);
}

.brand-text p {
  font-size: var(--font-size-sm);
  opacity: 0;
  margin: 4px 0 0;
  color: rgba(255, 255, 255, 0.85);
  letter-spacing: 0.01em;
  font-weight: var(--font-weight-medium);
  transform: translateY(-4px);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.app-header:hover .brand-text p {
  opacity: 1;
  transform: translateY(0);
}

/* 用户菜单样式 - 优化版本 */
.user-menu {
  position: relative;
}

.menu-trigger {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-2) var(--space-4);
  background: linear-gradient(135deg, 
    rgba(255, 255, 255, 0.15) 0%, 
    rgba(255, 255, 255, 0.08) 100%
  );
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: var(--radius-xl);
  color: white;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  position: relative;
  overflow: hidden;
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.menu-trigger::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, 
    transparent 0%, 
    rgba(255, 255, 255, 0.2) 50%, 
    transparent 100%
  );
  transition: left 0.6s ease;
}

.menu-trigger:hover::before {
  left: 100%;
}

.menu-trigger:hover {
  background: linear-gradient(135deg, 
    rgba(255, 255, 255, 0.25) 0%, 
    rgba(255, 255, 255, 0.15) 100%
  );
  border-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
  box-shadow: 
    0 8px 25px rgba(0, 0, 0, 0.15),
    0 4px 12px rgba(0, 0, 0, 0.1);
}

.menu-trigger.active {
  background: linear-gradient(135deg, 
    rgba(255, 255, 255, 0.3) 0%, 
    rgba(255, 255, 255, 0.2) 100%
  );
  border-color: rgba(255, 255, 255, 0.4);
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, 
    rgba(255, 255, 255, 0.2) 0%, 
    rgba(255, 255, 255, 0.1) 100%
  );
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.menu-trigger:hover .user-avatar {
  border-color: rgba(255, 255, 255, 0.5);
  transform: scale(1.05);
}

.user-name {
  font-weight: var(--font-weight-medium);
  color: rgba(255, 255, 255, 0.95);
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.dropdown-arrow {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  color: rgba(255, 255, 255, 0.8);
}

.dropdown-arrow.rotate {
  transform: rotate(180deg);
}

/* 下拉菜单样式 */
.dropdown {
  position: absolute;
  top: calc(100% + var(--space-2));
  right: 0;
  background: white;
  border-radius: var(--radius-lg);
  box-shadow: 
    0 10px 40px rgba(0, 0, 0, 0.15),
    0 2px 10px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(0, 0, 0, 0.05);
  min-width: 200px;
  z-index: 1000;
  overflow: hidden;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  opacity: 0;
  transform: translateY(-10px) scale(0.95);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: top right;
}

.dropdown.show {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  width: 100%;
  padding: var(--space-3) var(--space-4);
  border: none;
  background: transparent;
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  text-align: left;
  position: relative;
  overflow: hidden;
}

.dropdown-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: var(--primary);
  transform: scaleY(0);
  transition: transform 0.2s ease;
}

.dropdown-item:hover {
  background: linear-gradient(90deg, 
    rgba(99, 102, 241, 0.08) 0%, 
    rgba(99, 102, 241, 0.04) 100%
  );
  color: var(--primary);
}

.dropdown-item:hover::before {
  transform: scaleY(1);
}

.dropdown-item.logout-item {
  color: var(--error);
}

.dropdown-item.logout-item:hover {
  background: linear-gradient(90deg, 
    rgba(239, 68, 68, 0.08) 0%, 
    rgba(239, 68, 68, 0.04) 100%
  );
  color: var(--error);
}

.dropdown-item.logout-item:hover::before {
  background: var(--error);
}

.dropdown-divider {
  height: 1px;
  background: linear-gradient(90deg, 
    transparent 0%, 
    rgba(0, 0, 0, 0.08) 50%, 
    transparent 100%
  );
  margin: var(--space-1) 0;
}

/* 主内容区域样式 */
.app-main {
  flex: 1;
  display: flex;
  min-height: 0;
  background: var(--color-background);
}

/* 侧边栏样式 */
.sidebar {
  width: 280px;
  background: white;
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.04);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  z-index: 10;
  overflow: hidden;
}

/* 添加侧边栏装饰元素 */
.sidebar::after {
  content: '';
  position: absolute;
  top: 0;
  right: 0;
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1), transparent);
  border-radius: 0 0 0 100%;
  z-index: 0;
}

.sidebar-nav {
  padding: var(--space-4) 0;
  border-bottom: 1px solid var(--border);
  position: relative;
  z-index: 1;
}

.menu-list {
  list-style: none;
  margin: 0;
  padding: var(--space-2);
  border-radius: var(--radius-md);
  margin: 0 var(--space-3);
  overflow: hidden;
  background: var(--background);
  backdrop-filter: blur(8px);
}

.menu-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  width: 100%;
  padding: var(--space-3) var(--space-4);
  border: none;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  text-align: left;
  border-radius: var(--radius-md);
  position: relative;
  overflow: hidden;
}

/* 优化的波纹点击效果 */
.menu-item {
  position: relative;
  overflow: hidden;
}

.menu-item::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-radius: 50%;
  background: rgba(99, 102, 241, 0.4);
  transform: translate(-50%, -50%);
  transition: width 0.4s cubic-bezier(0.4, 0, 0.2, 1), 
              height 0.4s cubic-bezier(0.4, 0, 0.2, 1),
              opacity 0.6s ease;
  opacity: 0;
  pointer-events: none;
}

.menu-item:active::after {
  width: 400px;
  height: 400px;
  opacity: 1;
}

/* 添加额外的点击状态效果 */
.menu-item:active {
  transform: translateX(1px) translateY(1px);
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: transform 0.1s ease, box-shadow 0.1s ease;
}

.menu-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: var(--primary);
  transform: scaleY(0);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 0 var(--radius-full) var(--radius-full) 0;
}

.menu-item:hover {
  background: linear-gradient(90deg, 
    rgba(99, 102, 241, 0.1) 0%, 
    rgba(99, 102, 241, 0.05) 100%
  );
  color: var(--primary);
  transform: translateX(2px);
}

.menu-item:hover::before {
  transform: scaleY(0.9);
}

.menu-item.active {
  background: linear-gradient(90deg, 
    rgba(99, 102, 241, 0.15) 0%, 
    rgba(99, 102, 241, 0.08) 100%
  );
  color: var(--primary);
  font-weight: var(--font-weight-semibold);
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.15);
}

.menu-item.active::before {
  transform: scaleY(1);
}

/* 菜单项进入动画 */
.menu-item {
  animation: slideIn 0.3s ease forwards;
  opacity: 0;
  transform: translateY(10px);
}

.menu-item:nth-child(1) { animation-delay: 0.1s; }
.menu-item:nth-child(2) { animation-delay: 0.15s; }
.menu-item:nth-child(3) { animation-delay: 0.2s; }
.menu-item:nth-child(4) { animation-delay: 0.25s; }
.menu-item:nth-child(5) { animation-delay: 0.3s; }

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.menu-icon {
  font-size: 1.3em;
  width: 28px;
  text-align: center;
  transition: transform 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 28px;
  background: rgba(99, 102, 241, 0.05);
  border-radius: var(--radius-md);
}

.menu-item:hover .menu-icon {
  transform: scale(1.1);
  background: rgba(99, 102, 241, 0.15);
}

.menu-item.active .menu-icon {
  background: rgba(99, 102, 241, 0.2);
}

.menu-label {
  flex: 1;
  position: relative;
  transition: color 0.3s ease;
}

/* 优化侧边栏整体阴影和边界 */
.sidebar:hover {
  box-shadow: 0 0 25px rgba(0, 0, 0, 0.06);
}

/* 优化菜单分组 */
.menu-list li:not(:last-child) {
  margin-bottom: var(--space-1);
}

/* 插件内容区域样式 */
.plugins-content {
  height: 100%;
}

.plugins-layout {
  display: flex;
  height: 100%;
  gap: var(--space-4);
}

.plugins-sidebar {
  width: 320px;
  min-width: 280px;
  background: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 200px);
}

.plugins-header {
  margin-bottom: var(--space-4);
  padding-bottom: var(--space-3);
  border-bottom: 1px solid var(--color-border);
}

.plugins-header h3 {
  margin: 0 0 var(--space-3) 0;
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text-primary);
}

.plugins-tools {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.plugins-search {
  width: 100%;
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  background: var(--color-background);
  color: var(--color-text-primary);
  transition: all 0.2s ease;
}

.plugins-search:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.plugins-count {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
}

.plugins-list {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.plugins-item {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-background);
  color: var(--color-text-primary);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  text-align: left;
  width: 100%;
  will-change: transform, background-color, border-color;
}

.plugins-item:hover {
  background: var(--color-background-hover);
  border-color: var(--color-primary-light);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.plugins-item.active {
  background: var(--color-primary-light);
  color: var(--color-text-primary);
  border-color: var(--color-primary);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
  font-weight: 600;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-color: var(--color-primary);
}

.plugins-item .item-title {
  font-weight: 500;
  font-size: var(--font-size-sm);
  margin-bottom: var(--space-1);
}

.plugins-item .item-desc {
  font-size: var(--font-size-xs);
  opacity: 0.8;
  line-height: 1.4;
}

.plugins-item .item-badge {
  align-self: flex-end;
  background: rgba(99, 102, 241, 0.2);
  color: var(--color-primary);
  font-size: var(--font-size-xs);
  padding: 2px 6px;
  border-radius: var(--radius-sm);
  margin-top: var(--space-1);
}

.plugins-item.active .item-badge {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.plugins-main {
  flex: 1;
  background: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  overflow: auto;
  max-height: calc(100vh - 200px);
}

.plugin-renderer-container {
  height: 100%;
  width: 100%;
}

/* 空状态样式 */
.no-plugin-selected {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-state {
  text-align: center;
  color: var(--color-text-secondary);
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: var(--space-4);
  opacity: 0.5;
}

.empty-state h3 {
  margin: 0 0 var(--space-2) 0;
  font-size: var(--font-size-lg);
  color: var(--color-text-primary);
}

.empty-state p {
  margin: 0;
  font-size: var(--font-size-sm);
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .plugins-layout {
    flex-direction: column;
  }
  
  .plugins-sidebar {
    width: 100%;
    max-height: 300px;
  }
  
  .plugins-main {
    max-height: none;
  }
}

@media (max-width: 768px) {
  .plugins-sidebar {
    padding: var(--space-3);
  }
  
  .plugins-main {
    padding: var(--space-3);
  }
}

/* 原有的插件子菜单样式保留用于向后兼容 */
.plugin-submenu {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.submenu-header {
  padding: var(--space-4);
  border-bottom: 1px solid var(--color-border);
}

.submenu-header h3 {
  margin: 0 0 var(--space-3) 0;
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
}

.sidebar-tools {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.sidebar-search {
  width: 100%;
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  transition: all 0.2s ease;
}

.sidebar-search:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.sidebar-count {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

.plugin-list {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-2);
}

.plugin-item {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: var(--space-1);
  width: 100%;
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: white;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: var(--space-2);
  text-align: left;
}

.plugin-item:hover {
  border-color: var(--color-primary);
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.1);
}

.plugin-item.active {
  border-color: var(--color-primary);
  background: linear-gradient(135deg, 
    rgba(99, 102, 241, 0.05) 0%, 
    rgba(99, 102, 241, 0.02) 100%
  );
  box-shadow: 0 2px 12px rgba(99, 102, 241, 0.15);
}

.item-title {
  font-weight: var(--font-weight-medium);
  color: var(--color-text-primary);
}

.item-desc {
  font-size: var(--font-size-xs);
  color: var(--color-text-secondary);
  line-height: 1.4;
}

.item-badge {
  background: var(--color-primary);
  color: white;
  font-size: var(--font-size-xs);
  padding: 2px 6px;
  border-radius: var(--radius-sm);
  margin-left: auto;
}

/* 内容区域样式 */
.content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.content-header {
  padding: var(--space-5) var(--space-6) var(--space-3);
  border-bottom: 1px solid var(--color-border);
  background: white;
}

.content-header h2 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
}

.content-body {
  flex: 1;
  padding: var(--space-6);
  overflow-y: auto;
}

/* 内容卡片样式 */
.services-content,
.plugins-content,
.security-content {
  height: 100%;
}

.service-card,
.security-card {
  background: white;
  border-radius: var(--radius-lg);
  padding: var(--space-6);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border: 1px solid var(--color-border);
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.service-card h3,
.security-card h3 {
  margin: 0 0 var(--space-3) 0;
  font-size: 1.25rem;
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
}

.service-card p,
.security-card p {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-base);
}

.no-plugin-selected {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--color-text-tertiary);
}

.empty-state {
  text-align: center;
  padding: var(--space-8);
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: var(--space-4);
  display: block;
  opacity: 0.6;
}

.empty-state h3 {
  margin: 0 0 var(--space-2) 0;
  font-size: 1.25rem;
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
}

.empty-state p {
  margin: 0;
  color: var(--color-text-tertiary);
  font-size: var(--font-size-sm);
}

/* 过渡动画 */
.content-switch-enter-active,
.content-switch-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.content-switch-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.content-switch-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

/* 底部样式 */
.app-footer {
  background: white;
  border-top: 1px solid var(--color-border);
  padding: var(--space-3) 0;
  margin-top: auto;
}

.footer-content {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--space-5);
}

.footer-right {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.version {
  font-size: var(--font-size-sm);
  color: var(--color-text-tertiary);
  font-weight: var(--font-weight-medium);
}

.github-link {
  display: flex;
  align-items: center;
  color: var(--color-text-tertiary);
  transition: all 0.2s ease;
  text-decoration: none;
}

.github-link:hover {
  color: var(--color-primary);
  transform: translateY(-2px);
}

.github-icon {
  transition: all 0.2s ease;
}

.github-link:hover .github-icon {
  transform: scale(1.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 240px;
  }
  
  .content-body {
    padding: var(--space-4);
  }
  
  .brand-text h1 {
    font-size: 1.5rem;
  }
  
  .brand-text p {
    display: none;
  }
}

@media (max-width: 640px) {
  .app-main {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid var(--color-border);
  }
  
  .sidebar-nav {
    padding: var(--space-2) 0;
  }
  
  .menu-list {
    display: flex;
    overflow-x: auto;
    padding: 0 var(--space-2);
  }
  
  .menu-item {
    flex-shrink: 0;
    min-width: 80px;
    flex-direction: column;
    gap: var(--space-1);
    padding: var(--space-2);
  }
  
  .menu-label {
    font-size: var(--font-size-xs);
  }
  
  .plugin-submenu {
    display: none;
  }
  
  .content-header {
    padding: var(--space-3) var(--space-4) var(--space-2);
  }
  
  .content-body {
    padding: var(--space-3);
  }
}
</style>

