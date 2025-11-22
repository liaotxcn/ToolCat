<template>
  <div class="user-center">
    <div class="welcome-header">
      <div class="welcome-content">
        <div class="user-avatar-large">
          <svg viewBox="0 0 24 24" width="48" height="48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 12C14.2091 12 16 10.2091 16 8C16 5.79086 14.2091 4 12 4C9.79086 4 8 5.79086 8 8C8 10.2091 9.79086 12 12 12Z" stroke="currentColor" stroke-width="1.5"/>
            <path d="M19 20C19 16.134 12 14 12 14C12 14 5 16.134 5 20V21H19V20Z" stroke="currentColor" stroke-width="1.5"/>
          </svg>
        </div>
        <div class="welcome-info">
          <h2>个人中心</h2>
          <p>管理您的账户信息、安全设置和活动记录</p>
        </div>
      </div>
    </div>

    <div class="cards">
      <section class="card profile-card">
        <div class="card-header">
          <div class="card-icon">
            <svg viewBox="0 0 24 24" width="20" height="20" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M20 21V19C20 17.9391 19.5786 16.9217 18.8284 16.1716C18.0783 15.4214 17.0609 15 16 15H8C6.93913 15 5.92172 15.4214 5.17157 16.1716C4.42143 16.9217 4 17.9391 4 19V21" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M12 11C14.2091 11 16 9.20914 16 7C16 4.79086 14.2091 3 12 3C9.79086 3 8 4.79086 8 7C8 9.20914 9.79086 11 12 11Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <h3>基本资料</h3>
        </div>
        
        <div class="form-content">
          <div class="form-group">
            <label class="form-label">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 12C14.2091 12 16 10.2091 16 8C16 5.79086 14.2091 4 12 4C9.79086 4 8 5.79086 8 8C8 10.2091 9.79086 12 12 12Z" stroke="currentColor" stroke-width="1.5"/>
                <path d="M19 20C19 16.134 12 14 12 14C12 14 5 16.134 5 20V21H19V20Z" stroke="currentColor" stroke-width="1.5"/>
              </svg>
              用户名
            </label>
            <div class="input-wrapper">
              <input 
                type="text" 
                :value="user?.username" 
                readonly 
                class="form-input readonly"
              />
              <div class="input-suffix">
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 15L15 12L12 9L12 15Z" fill="currentColor"/>
                  <path d="M9 12L12 15L12 9L9 12Z" fill="currentColor"/>
                </svg>
              </div>
            </div>
          </div>

          <div class="form-group" :class="{ 'has-error': emailInvalid }">
            <label class="form-label">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M4 4H20C21.1 4 22 4.9 22 6V18C22 19.1 21.1 20 20 20H4C2.9 20 2 19.1 2 18V6C2 4.9 2.9 4 4 4Z" stroke="currentColor" stroke-width="1.5"/>
                <path d="M22 6L12 13L2 6" stroke="currentColor" stroke-width="1.5"/>
              </svg>
              邮箱地址
            </label>
            <div class="input-wrapper">
              <input 
                v-model="email" 
                type="email" 
                placeholder="name@example.com" 
                class="form-input"
                :class="{ 'error': emailInvalid }"
                @input="clearStatus" 
              />
              <div class="input-suffix" v-if="email && !emailInvalid">
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M9 12L11 14L15 10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </div>
            </div>
            <div v-if="emailInvalid" class="field-message error">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 8V12M12 16H12.01M22 12C22 17.5228 17.5228 22 12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2C17.5228 2 22 6.47715 22 12Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              请输入有效的邮箱地址
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 8V12L15 15M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              创建时间
            </label>
            <div class="input-wrapper">
              <input 
                type="text" 
                :value="formatDate(user?.created_at)" 
                readonly 
                class="form-input readonly"
              />
            </div>
          </div>
        </div>

        <div class="card-actions">
          <button 
            class="btn btn-primary" 
            @click="updateProfile" 
            :disabled="updating || !canSave"
          >
            <span v-if="updating" class="btn-loading">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 2V6M12 18V22M4.22 4.22L6.34 6.34M17.66 17.66L19.78 19.78M2 12H6M18 12H22M4.22 19.78L6.34 17.66M17.66 6.34L19.78 4.22" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              保存中...
            </span>
            <span v-else>
              <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M5 13L9 17L19 7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              保存资料
            </span>
          </button>
        </div>
      </section>

      <section class="card security-card">
        <div class="card-header">
          <div class="card-icon">
            <svg viewBox="0 0 24 24" width="20" height="20" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 22C16.4183 22 20 18.4183 20 14C20 9.58172 16.4183 6 12 6C7.58172 6 4 9.58172 4 14C4 18.4183 7.58172 22 12 22Z" stroke="currentColor" stroke-width="1.5"/>
              <path d="M12 14V16M12 10H12.01" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M12 2V6M22 12H18M6 12H2M19.07 4.93L16.24 7.76M7.76 16.24L4.93 19.07M19.07 19.07L16.24 16.24M7.76 7.76L4.93 4.93" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <h3>安全设置</h3>
        </div>

        <div class="form-content">
          <div class="form-group" :class="{ 'has-error': passwordMismatch }">
            <label class="form-label">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M7 10V7C7 4.23858 9.23858 2 12 2C14.7614 2 17 4.23858 17 7V10M5 21H19C20.1046 21 21 20.1046 21 19V12C21 10.8954 20.1046 10 19 10H5C3.89543 10 3 10.8954 3 12V19C3 20.1046 3.89543 21 5 21Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              新密码
            </label>
            <div class="input-wrapper">
              <input 
                v-model="newPassword" 
                type="password" 
                placeholder="至少6个字符" 
                class="form-input"
                @input="clearStatus" 
              />
            </div>
            
            <div v-if="newPassword" class="password-strength">
              <div class="strength-bar">
                <div class="strength-fill" :class="passwordLevel"></div>
              </div>
              <div class="strength-info">
                <span class="strength-label" :class="passwordLevel">{{ passwordLabel }}</span>
                <span class="strength-tips">
                  <span v-if="newPassword.length < 6">• 至少6个字符</span>
                  <span v-if="!/[A-Z]/.test(newPassword)">• 包含大写字母</span>
                  <span v-if="!/[a-z]/.test(newPassword)">• 包含小写字母</span>
                  <span v-if="!/\d/.test(newPassword)">• 包含数字</span>
                  <span v-if="!/[^\w]/.test(newPassword)">• 包含特殊字符</span>
                </span>
              </div>
            </div>
          </div>

          <div class="form-group" :class="{ 'has-error': passwordMismatch }">
            <label class="form-label">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M9 12L11 14L15 10M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              确认密码
            </label>
            <div class="input-wrapper">
              <input 
                v-model="confirmNewPassword" 
                type="password" 
                placeholder="再次输入新密码" 
                class="form-input"
                :class="{ 'error': passwordMismatch }"
                @input="clearStatus" 
              />
              <div class="input-suffix" v-if="confirmNewPassword && !passwordMismatch">
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M9 12L11 14L15 10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </div>
            </div>
            <div v-if="passwordMismatch" class="field-message error">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 8V12M12 16H12.01M22 12C22 17.5228 17.5228 22 12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2C17.5228 2 22 6.47715 22 12Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              两次输入的密码不一致
            </div>
          </div>
        </div>

        <div class="card-actions">
          <button 
            class="btn btn-secondary" 
            @click="updatePassword" 
            :disabled="updating || !canUpdatePassword"
          >
            <span v-if="updating" class="btn-loading">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 2V6M12 18V22M4.22 4.22L6.34 6.34M17.66 17.66L19.78 19.78M2 12H6M18 12H22M4.22 19.78L6.34 17.66M17.66 6.34L19.78 4.22" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              提交中...
            </span>
            <span v-else>
              <svg viewBox="0 0 24 24" width="16" height="16" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 22C16.4183 22 20 18.4183 20 14C20 9.58172 16.4183 6 12 6C7.58172 6 4 9.58172 4 14C4 18.4183 7.58172 22 12 22Z" stroke="currentColor" stroke-width="1.5"/>
                <path d="M12 14V16M12 10H12.01" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              更新密码
            </span>
          </button>
        </div>
      </section>
    </div>

    <section class="card activity-card">
      <div class="card-header">
        <div class="card-icon">
          <svg viewBox="0 0 24 24" width="20" height="20" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <h3>近期活动</h3>
        <span class="card-subtitle">审计日志记录</span>
      </div>
      
      <div class="activity-content">
        <div v-if="auditLogs.length === 0" class="empty-state">
          <div class="empty-icon">
            <svg viewBox="0 0 24 24" width="48" height="48" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <p>暂无活动记录</p>
        </div>
        
        <div v-else class="activity-list">
          <div v-for="log in auditLogs" :key="log.id" class="activity-item">
            <div class="activity-header">
              <div class="activity-action">
                <span class="action-badge" :class="getActionClass(log.action)">
                  {{ log.action }}
                </span>
                <span class="resource-info">{{ log.resource_type }} #{{ log.resource_id }}</span>
              </div>
              <div class="activity-time">{{ formatDate(log.created_at) }}</div>
            </div>
            <div v-if="log.new_value || log.old_value" class="activity-details">
              <div v-if="log.old_value" class="change-item old">
                <svg viewBox="0 0 24 24" width="14" height="14" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M3 12H21M3 12L9 6M3 12L9 18" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                <span>变更前: {{ short(log.old_value) }}</span>
              </div>
              <div v-if="log.new_value" class="change-item new">
                <svg viewBox="0 0 24 24" width="14" height="14" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M21 12H3M21 12L15 6M21 12L15 18" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                <span>变更后: {{ short(log.new_value) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 全局消息提示 -->
    <div v-if="errorMessage" class="global-message error">
      <svg viewBox="0 0 24 24" width="20" height="20" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M12 8V12M12 16H12.01M22 12C22 17.5228 17.5228 22 12 22C6.47715 22 2 17.5228 2 12C2 6.47715 6.47715 2 12 2C17.5228 2 22 6.47715 22 12Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      {{ errorMessage }}
    </div>
    
    <div v-if="successMessage" class="global-message success">
      <svg viewBox="0 0 24 24" width="20" height="20" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M9 12L11 14L15 10M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      {{ successMessage }}
    </div>
  </div>
</template>

<script>
import api, { authService } from '../services/auth'

export default {
  name: 'UserCenter',
  props: {
    currentUser: { type: Object, default: null }
  },
  data() {
    return {
      user: null,
      email: '',
      newPassword: '',
      confirmNewPassword: '',
      auditLogs: [],
      updating: false,
      errorMessage: '',
      successMessage: ''
    }
  },
  mounted() {
    this.loadUser()
  },
  computed: {
    // 新增：是否允许保存资料
    canSave() {
      return !!this.email && !this.emailInvalid
    },
    // 新增：邮箱是否非法
    emailInvalid() {
      if (!this.email) return false
      const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      return !re.test(this.email)
    },
    // 新增：密码强度等级样式
    passwordLevel() {
      const n = this.newPassword || ''
      let score = 0
      if (n.length >= 6) score++
      if (/[A-Z]/.test(n)) score++
      if (/[a-z]/.test(n)) score++
      if (/\d/.test(n)) score++
      if (/[^\w]/.test(n)) score++
      if (score <= 2) return 'weak'
      if (score === 3 || score === 4) return 'medium'
      return 'strong'
    },
    // 新增：密码强度文案
    passwordLabel() {
      return this.passwordLevel === 'strong' ? '强' : (this.passwordLevel === 'medium' ? '中' : '弱')
    },
    // 新增：两次密码是否不一致
    passwordMismatch() {
      return !!this.confirmNewPassword && this.newPassword !== this.confirmNewPassword
    },
    // 新增：是否允许更新密码
    canUpdatePassword() {
      return !!this.newPassword && this.newPassword.length >= 6 && !this.passwordMismatch
    }
  },
  methods: {
    // 新增：输入时清除顶部成功/错误提示，避免阻碍表单操作
    clearStatus() {
      this.errorMessage = ''
      this.successMessage = ''
    },
    // 新增：获取操作类型对应的样式类
    getActionClass(action) {
      const actionMap = {
        'UPDATE': 'UPDATE',
        'CREATE': 'CREATE', 
        'DELETE': 'DELETE'
      }
      return actionMap[action] || 'UPDATE'
    },
    async loadUser() {
      try {
        const cur = authService.getCurrentUser()
        if (!cur || !cur.id) {
          this.errorMessage = '未获取到当前用户信息，请重新登录'
          return
        }
        const res = await api.get(`/api/v1/users/${cur.id}`)
        this.user = res
        this.email = res.email || ''
        this.auditLogs = Array.isArray(res.audit_logs) ? res.audit_logs.slice(0, 20) : []
        this.errorMessage = ''
      } catch (e) {
        this.errorMessage = e?.response?.data?.message || '加载用户信息失败'
      }
    },
    async updateProfile() {
      try {
        this.updating = true
        this.successMessage = ''
        this.errorMessage = ''
        const cur = authService.getCurrentUser()
        if (!cur || !cur.id) throw new Error('未登录')
        const payload = { email: this.email }
        const updated = await api.put(`/api/v1/users/${cur.id}`, payload)
        this.user = updated
        this.$emit('updated-user', updated)
        this.successMessage = '资料已更新'
      } catch (e) {
        this.errorMessage = e?.response?.data?.message || '更新失败'
      } finally {
        this.updating = false
      }
    },
    async updatePassword() {
      try {
        this.updating = true
        this.successMessage = ''
        this.errorMessage = ''
        if (!this.newPassword || this.newPassword.length < 6) {
          this.errorMessage = '新密码至少6个字符'
          return
        }
        if (this.newPassword !== this.confirmNewPassword) {
          this.errorMessage = '两次输入的密码不一致'
          return
        }
        const cur = authService.getCurrentUser()
        if (!cur || !cur.id) throw new Error('未登录')
        const payload = { password: this.newPassword }
        const updated = await api.put(`/api/v1/users/${cur.id}`, payload)
        this.user = updated
        this.$emit('updated-user', updated)
        this.newPassword = ''
        this.confirmNewPassword = ''
        this.successMessage = '密码已更新'
      } catch (e) {
        this.errorMessage = e?.response?.data?.message || '密码更新失败'
      } finally {
        this.updating = false
      }
    },
    formatDate(dt) {
      if (!dt) return '-'
      try { return new Date(dt).toLocaleString() } catch (_) { return dt }
    },
    short(text) {
      if (!text) return ''
      const s = String(text)
      return s.length > 180 ? s.slice(0, 180) + '…' : s
    }
  }
}
</script>

<style scoped>
/* 全局样式 */
.user-center {
  display: flex;
  flex-direction: column;
  gap: 24px;
  padding: 0;
  max-width: 1200px;
  margin: 0 auto;
}

/* 欢迎头部 */
.welcome-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  padding: 32px;
  color: white;
  position: relative;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.3);
}

.welcome-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at 20% 80%, rgba(255, 255, 255, 0.1) 0%, transparent 50%);
  pointer-events: none;
}

.welcome-content {
  display: flex;
  align-items: center;
  gap: 20px;
  position: relative;
  z-index: 1;
}

.user-avatar-large {
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.3);
  transition: transform 0.3s ease;
}

.user-avatar-large:hover {
  transform: scale(1.05);
}

.welcome-info h2 {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 600;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.welcome-info p {
  margin: 0;
  opacity: 0.9;
  font-size: 16px;
}

/* 卡片布局 */
.cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(380px, 1fr));
  gap: 24px;
}

/* 卡片样式 */
.card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #667eea, #764ba2);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.card:hover::before {
  opacity: 1;
}

/* 卡片头部 */
.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f3f4f6;
}

.card-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.card-subtitle {
  color: #6b7280;
  font-size: 14px;
  margin-left: auto;
}

/* 表单样式 */
.form-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group.has-error {
  position: relative;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-label svg {
  color: #6b7280;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.form-input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  font-size: 14px;
  transition: all 0.3s ease;
  background: #fff;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-input.error {
  border-color: #ef4444;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
}

.form-input.readonly {
  background: #f9fafb;
  color: #6b7280;
  cursor: not-allowed;
}

.input-suffix {
  position: absolute;
  right: 12px;
  color: #10b981;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 字段消息 */
.field-message {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  margin-top: 4px;
}

.field-message.error {
  color: #ef4444;
}

.field-message svg {
  flex-shrink: 0;
}

/* 密码强度 */
.password-strength {
  margin-top: 8px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.strength-bar {
  height: 6px;
  background: #e5e7eb;
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 8px;
}

.strength-fill {
  height: 100%;
  border-radius: 3px;
  transition: all 0.3s ease;
  width: 0%;
}

.strength-fill.weak {
  width: 33%;
  background: #ef4444;
}

.strength-fill.medium {
  width: 66%;
  background: #f59e0b;
}

.strength-fill.strong {
  width: 100%;
  background: #10b981;
}

.strength-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.strength-label {
  font-size: 12px;
  font-weight: 500;
}

.strength-label.weak {
  color: #ef4444;
}

.strength-label.medium {
  color: #f59e0b;
}

.strength-label.strong {
  color: #10b981;
}

.strength-tips {
  font-size: 11px;
  color: #6b7280;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

/* 按钮样式 */
.card-actions {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #f3f4f6;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
  position: relative;
  overflow: hidden;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
}

.btn-primary:active:not(:disabled) {
  transform: translateY(0);
}

.btn-secondary {
  background: #fff;
  color: #374151;
  border: 2px solid #e5e7eb;
}

.btn-secondary:hover:not(:disabled) {
  background: #f9fafb;
  border-color: #d1d5db;
  transform: translateY(-1px);
}

.btn-loading {
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

/* 活动卡片 */
.activity-card {
  grid-column: 1 / -1;
}

.activity-content {
  min-height: 200px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px;
  color: #6b7280;
}

.empty-icon {
  margin-bottom: 16px;
  opacity: 0.3;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.activity-item {
  padding: 16px;
  background: #f9fafb;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  transition: all 0.3s ease;
}

.activity-item:hover {
  background: #fff;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transform: translateY(-1px);
}

.activity-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.activity-action {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.action-badge {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.action-badge.UPDATE {
  background: #dbeafe;
  color: #1e40af;
}

.action-badge.CREATE {
  background: #d1fae5;
  color: #065f46;
}

.action-badge.DELETE {
  background: #fee2e2;
  color: #991b1b;
}

.resource-info {
  color: #6b7280;
  font-size: 13px;
}

.activity-time {
  color: #9ca3af;
  font-size: 12px;
  white-space: nowrap;
}

.activity-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #e5e7eb;
}

.change-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 6px;
}

.change-item.old {
  background: #fef2f2;
  color: #991b1b;
}

.change-item.new {
  background: #f0fdf4;
  color: #166534;
}

/* 全局消息 */
.global-message {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 16px 20px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
  font-weight: 500;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  animation: slideIn 0.3s ease;
  max-width: 400px;
}

.global-message.error {
  background: #fef2f2;
  color: #991b1b;
  border: 1px solid #fecaca;
}

.global-message.success {
  background: #f0fdf4;
  color: #166534;
  border: 1px solid #bbf7d0;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-center {
    padding: 0 16px;
    gap: 20px;
  }
  
  .welcome-header {
    padding: 24px;
  }
  
  .welcome-content {
    flex-direction: column;
    text-align: center;
  }
  
  .welcome-info h2 {
    font-size: 24px;
  }
  
  .cards {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .card {
    padding: 20px;
  }
  
  .global-message {
    right: 16px;
    left: 16px;
    max-width: none;
  }
  
  .activity-header {
    flex-direction: column;
    gap: 8px;
  }
  
  .activity-action {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .welcome-header {
    padding: 20px;
  }
  
  .user-avatar-large {
    width: 60px;
    height: 60px;
  }
  
  .welcome-info h2 {
    font-size: 20px;
  }
  
  .welcome-info p {
    font-size: 14px;
  }
  
  .card {
    padding: 16px;
  }
  
  .btn {
    width: 100%;
    justify-content: center;
  }
}
</style>