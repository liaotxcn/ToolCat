<template>
  <div class="register-container">
    <h2 class="form-title">注册</h2>
    <form class="auth-form" @submit.prevent="handleRegister">
      <!-- 用户名输入 -->
      <div class="form-group">
        <label for="username">用户名</label>
        <div class="input-wrapper">
          <input 
            v-model="username" 
            type="text" 
            id="username" 
            required 
            minlength="3" 
            maxlength="50" 
            placeholder="3-50个字符" 
            @input="clearError"
            :class="{ 'input-error': usernameInvalid && username }"
          />
        </div>
        <div v-if="usernameInvalid && username" class="input-hint input-hint-error">用户名至少需要3个字符</div>
        <div v-else-if="username && !usernameInvalid" class="input-hint input-hint-success">用户名可用</div>
      </div>
      
      <!-- 邮箱输入 -->
      <div class="form-group">
        <label for="email">邮箱</label>
        <div class="input-wrapper">
          <input 
            v-model="email" 
            type="email" 
            id="email" 
            required 
            placeholder="name@example.com" 
            @input="clearError"
            :class="{ 'input-error': emailInvalid && email }"
          />
        </div>
        <div v-if="emailInvalid && email" class="input-hint input-hint-error">请输入有效的邮箱地址</div>
        <div v-else-if="email && !emailInvalid" class="input-hint input-hint-success">邮箱格式正确</div>
      </div>
      
      <!-- 密码输入 -->
      <div class="form-group">
        <label for="password">密码</label>
        <div class="password-wrap">
          <input 
            :type="showPassword ? 'text' : 'password'" 
            v-model="password" 
            id="password" 
            required 
            minlength="6" 
            placeholder="至少6个字符" 
            @input="clearError"
            :class="{ 'input-error': password.length > 0 && password.length < 6 }"
          />
          <button type="button" class="toggle-psw" @click="showPassword = !showPassword" :aria-pressed="showPassword" :title="showPassword ? '隐藏密码' : '显示密码'" aria-label="切换密码可见性">
            <svg class="eye-icon" viewBox="0 0 20 20" width="18" height="18" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M2 10c2.5-4.5 6-6.5 8-6.5s5.5 2 8 6.5c-2.5 4.5-6 6.5-8 6.5S4.5 14.5 2 10z" fill="none" stroke="currentColor" stroke-width="1.5" />
              <circle cx="10" cy="10" r="3" fill="currentColor" />
              <path v-if="!showPassword" d="M4 4L16 16" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
            </svg>
          </button>
        </div>
        
        <!-- 密码强度指示器 -->
        <div v-if="password" class="pw-strength">
          <div class="strength-label">密码强度：</div>
          <div class="strength-bar-container">
            <div class="strength-bar">
              <div 
                class="strength-progress" 
                :class="passwordLevel"
                :style="{ width: getStrengthWidth() }"
              ></div>
            </div>
            <span class="strength-text" :class="passwordLevel">{{ passwordLabel }}</span>
          </div>
          
          <!-- 密码强度提示 -->
          <div class="strength-hints">
            <div class="hint-item" :class="{ 'hint-passed': password.length >= 6 }">
              <span class="hint-icon">{{ password.length >= 6 ? '✓' : '•' }}</span>
              <span>至少6个字符</span>
            </div>
            <div class="hint-item" :class="{ 'hint-passed': /[A-Z]/.test(password) }">
              <span class="hint-icon">{{ /[A-Z]/.test(password) ? '✓' : '•' }}</span>
              <span>包含大写字母</span>
            </div>
            <div class="hint-item" :class="{ 'hint-passed': /[a-z]/.test(password) }">
              <span class="hint-icon">{{ /[a-z]/.test(password) ? '✓' : '•' }}</span>
              <span>包含小写字母</span>
            </div>
            <div class="hint-item" :class="{ 'hint-passed': /\d/.test(password) }">
              <span class="hint-icon">{{ /\d/.test(password) ? '✓' : '•' }}</span>
              <span>包含数字</span>
            </div>
            <div class="hint-item" :class="{ 'hint-passed': /[^\w]/.test(password) }">
              <span class="hint-icon">{{ /[^\w]/.test(password) ? '✓' : '•' }}</span>
              <span>包含特殊字符</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 确认密码 -->
      <div class="form-group">
        <label for="confirmPassword">确认密码</label>
        <div class="password-wrap">
          <input 
            :type="showConfirmPassword ? 'text' : 'password'" 
            v-model="confirmPassword" 
            id="confirmPassword" 
            required 
            minlength="6" 
            placeholder="再次输入密码" 
            @input="clearError"
            :class="{ 'input-error': passwordMismatch && confirmPassword }"
          />
          <button type="button" class="toggle-psw" @click="showConfirmPassword = !showConfirmPassword" :aria-pressed="showConfirmPassword" :title="showConfirmPassword ? '隐藏密码' : '显示密码'" aria-label="切换密码可见性">
            <svg class="eye-icon" viewBox="0 0 20 20" width="18" height="18" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M2 10c2.5-4.5 6-6.5 8-6.5s5.5 2 8 6.5c-2.5 4.5-6 6.5-8 6.5S4.5 14.5 2 10z" fill="none" stroke="currentColor" stroke-width="1.5" />
              <circle cx="10" cy="10" r="3" fill="currentColor" />
              <path v-if="!showConfirmPassword" d="M4 4L16 16" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
            </svg>
          </button>
        </div>
        <div v-if="passwordMismatch && confirmPassword" class="input-hint input-hint-error">两次输入的密码不一致</div>
        <div v-else-if="confirmPassword && !passwordMismatch && password" class="input-hint input-hint-success">密码一致</div>
      </div>
      
      <!-- 错误消息 -->
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      
      <!-- 注册按钮 -->
      <button class="primary-btn" type="submit" :disabled="loading || !canRegister">
        <span v-if="loading" class="loading-spinner"></span>
        {{ loading ? '注册中...' : '注册' }}
      </button>
    </form>
    
    <p class="switch-tip">
      已有账号？
      <button class="link-btn" type="button" @click="switchToLogin">返回登录</button>
    </p>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { authService } from '../services/auth'

export default {
  name: 'Register',
  emits: ['register-success', 'switch-to-login'],
  setup(props, { emit }) {
    // 响应式数据
    const username = ref('')
    const email = ref('')
    const password = ref('')
    const confirmPassword = ref('')
    const showPassword = ref(false)
    const showConfirmPassword = ref(false)
    const loading = ref(false)
    const errorMessage = ref('')

    // 计算属性
    const usernameInvalid = computed(() => {
      return !(username.value && username.value.trim().length >= 3)
    })
    
    const emailInvalid = computed(() => {
      if (!email.value) return false
      const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      return !re.test(email.value)
    })
    
    const passwordLevel = computed(() => {
      const n = password.value || ''
      let score = 0
      if (n.length >= 6) score++
      if (/[A-Z]/.test(n)) score++
      if (/[a-z]/.test(n)) score++
      if (/\d/.test(n)) score++
      if (/[^\w]/.test(n)) score++
      if (score <= 2) return 'weak'
      if (score === 3 || score === 4) return 'medium'
      return 'strong'
    })
    
    const passwordLabel = computed(() => {
      return passwordLevel.value === 'strong' ? '强' : (passwordLevel.value === 'medium' ? '中' : '弱')
    })
    
    const passwordMismatch = computed(() => {
      return !!confirmPassword.value && password.value !== confirmPassword.value
    })
    
    const canRegister = computed(() => {
      return !usernameInvalid.value && !emailInvalid.value && !!password.value && 
             password.value.length >= 6 && !passwordMismatch.value
    })

    // 方法
    const clearError = () => {
      errorMessage.value = ''
    }
    
    const getStrengthWidth = () => {
      const n = password.value || ''
      let score = 0
      if (n.length >= 6) score++
      if (/[A-Z]/.test(n)) score++
      if (/[a-z]/.test(n)) score++
      if (/\d/.test(n)) score++
      if (/[^\w]/.test(n)) score++
      
      // 计算宽度百分比
      return `${(score / 5) * 100}%`
    }
    
    const handleRegister = async () => {
      if (!canRegister.value) return
      errorMessage.value = ''
      
      try {
        loading.value = true
        const payload = {
          username: username.value.trim(),
          email: email.value.trim(),
          password: password.value,
          confirm_password: confirmPassword.value
        }

        const response = await authService.register(payload)

        if (response && response.user) {
          emit('register-success', response.user)
        } else {
          errorMessage.value = response?.message || '注册失败，请稍后重试'
        }
      } catch (error) {
        const data = error?.response?.data || {}
        errorMessage.value = data?.message || '注册失败，请检查输入或网络'
      } finally {
        loading.value = false
      }
    }
    
    const switchToLogin = () => {
      emit('switch-to-login')
    }

    return {
      username,
      email,
      password,
      confirmPassword,
      showPassword,
      showConfirmPassword,
      loading,
      errorMessage,
      usernameInvalid,
      emailInvalid,
      passwordLevel,
      passwordLabel,
      passwordMismatch,
      canRegister,
      clearError,
      getStrengthWidth,
      handleRegister,
      switchToLogin
    }
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  flex-direction: column;
  gap: 0;
  max-width: 400px;
  width: 100%;
  margin: 0 auto;
  padding: 24px;
  background: white;
  border-radius: var(--radius-xl);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.04), 0 2px 8px rgba(0, 0, 0, 0.02);
}

.form-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 28px 0;
  text-align: center;
  letter-spacing: -0.02em;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  position: relative;
}

.form-group label {
  font-weight: 500;
  color: var(--text-primary);
  font-size: var(--font-size-sm);
  letter-spacing: 0.02em;
  margin-bottom: 2px;
}

.input-wrapper,
.password-wrap {
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
}

/* 图标已移除 */

.auth-form input {
  width: 100%;
  padding: 16px;
  border: 1px solid var(--border-light);
  border-radius: var(--radius-lg);
  font-size: 16px;
  font-family: inherit;
  line-height: 1.5;
  transition: all var(--transition-normal);
  background: var(--bg-primary);
  color: var(--text-primary);
  box-sizing: border-box;
}

input::placeholder {
  color: var(--text-tertiary);
  font-size: 14px;
  opacity: 0.8;
  transition: all var(--transition-fast);
}

input:focus {
  outline: none;
  border-color: var(--primary-500);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
  background: var(--bg-primary);
}

input:focus::placeholder {
  opacity: 0.6;
  transform: translateX(2px);
}

/* 图标相关的聚焦样式已移除 */

input.input-error {
  border-color: var(--error);
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
}

/* 密码输入框特殊样式 */
.password-wrap input {
  padding-right: 52px;
}

.toggle-psw {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-tertiary);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all var(--transition-fast);
  z-index: 2;
}

.toggle-psw:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.toggle-psw:focus {
  outline: none;
  background: var(--bg-secondary);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

.eye-icon {
  display: block;
}

/* 输入提示 */
.input-hint {
  font-size: var(--font-size-xs);
  margin-top: 2px;
  transition: all var(--transition-fast);
}

.input-hint-error {
  color: var(--error);
  font-weight: var(--font-weight-medium);
}

.input-hint-success {
  color: var(--success);
  font-weight: var(--font-weight-medium);
}

/* 密码强度指示器 */
.pw-strength {
  margin-top: 8px;
}

.strength-label {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  margin-bottom: 4px;
  font-weight: var(--font-weight-medium);
}

.strength-bar-container {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.strength-bar {
  flex: 1;
  height: 8px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-full);
  overflow: hidden;
  position: relative;
}

.strength-progress {
  height: 100%;
  transition: all var(--transition-normal);
  border-radius: var(--radius-full);
  position: relative;
}

.strength-progress::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, transparent 0%, rgba(255, 255, 255, 0.2) 50%, transparent 100%);
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

.strength-progress.weak {
  background: var(--error);
}

.strength-progress.medium {
  background: var(--warning);
}

.strength-progress.strong {
  background: var(--success);
}

.strength-text {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  min-width: 24px;
}

.strength-text.weak {
  color: var(--error);
}

.strength-text.medium {
  color: var(--warning);
}

.strength-text.strong {
  color: var(--success);
}

/* 密码强度提示 */
.strength-hints {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 6px;
  font-size: var(--font-size-xs);
}

@media (max-width: 480px) {
  .strength-hints {
    grid-template-columns: 1fr;
  }
}

.hint-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text-tertiary);
}

.hint-passed {
  color: var(--success);
  font-weight: var(--font-weight-medium);
}

.hint-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--bg-tertiary);
  font-size: 10px;
}

.hint-passed .hint-icon {
  background: var(--success-100);
  color: var(--success);
}

/* 提交按钮 */
.primary-btn {
  padding: 14px 24px;
  border: none;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, var(--primary-600) 0%, var(--primary-700) 100%);
  color: white;
  cursor: pointer;
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  transition: all var(--transition-normal);
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 20px;
  min-height: 48px;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.2);
}

.primary-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-600) 100%);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
  transform: translateY(-1px);
}

.primary-btn:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 2px 6px rgba(99, 102, 241, 0.2);
}

.primary-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* 加载动画 */
.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid var(--bg-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 错误消息 */
.error-message {
  background: var(--error-100);
  color: var(--error-700);
  border: 1px solid var(--error);
  padding: 10px 12px;
  border-radius: var(--radius-lg);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
}

.error-message::before {
  content: '';
  width: 16px;
  height: 16px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 20 20' fill='%23dc2626'%3E%3Cpath fill-rule='evenodd' d='M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z' clip-rule='evenodd'/%3E%3C/svg%3E") no-repeat center center;
  background-size: contain;
}

/* 切换提示 */
.switch-tip {
  text-align: center;
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  margin-top: 24px;
}

.link-btn {
  font-size: var(--font-size-sm);
  background: none;
  border: none;
  color: var(--primary-600);
  cursor: pointer;
  font-weight: 500;
  padding: 4px 8px;
  border-radius: var(--radius);
  transition: all var(--transition-fast);
  text-decoration: none;
  display: inline-flex;
  align-items: center;
}

.link-btn:hover {
  background: var(--primary-50);
  color: var(--primary-700);
  transform: translateY(-1px);
}

/* 响应式调整 */
@media (max-width: 480px) {
  .form-title {
    font-size: var(--font-size-lg);
    margin-bottom: 20px;
  }
  
  .auth-form {
    gap: 14px;
  }
  
  .switch-tip {
    margin-top: 20px;
  }
}
</style>