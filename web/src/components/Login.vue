<template>
  <div class="login-container">
    <h2 class="form-title">登录</h2>
    <form class="auth-form" @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="username">用户名</label>
        <div class="input-wrapper">
          <input v-model="username" type="text" id="username" required placeholder="请输入用户名" autofocus @input="clearError" />
        </div>
        <div v-if="usernameInvalid && username" class="input-hint">用户名不能为空</div>
      </div>
      
      <div class="form-group">
        <label for="password">密码</label>
        <div class="password-wrap">
          <input :type="showPassword ? 'text' : 'password'" v-model="password" id="password" required placeholder="请输入密码" @input="clearError" />
          <button type="button" class="toggle-psw" @click="showPassword = !showPassword" :aria-pressed="showPassword" :title="showPassword ? '隐藏密码' : '显示密码'" aria-label="切换密码可见性">
            <svg class="eye-icon" viewBox="0 0 20 20" width="18" height="18" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M2 10c2.5-4.5 6-6.5 8-6.5s5.5 2 8 6.5c-2.5 4.5-6 6.5-8 6.5S4.5 14.5 2 10z" fill="none" stroke="currentColor" stroke-width="1.5" />
              <circle cx="10" cy="10" r="3" fill="currentColor" />
              <path v-if="!showPassword" d="M4 4L16 16" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
            </svg>
          </button>
        </div>
        <div v-if="passwordInvalid && password" class="input-hint">密码不能为空</div>
      </div>
      
      <div class="assist">
        <label class="remember">
          <input type="checkbox" v-model="rememberMe" class="checkbox-custom" />
          <span class="checkbox-label">记住我</span>
        </label>
        <button class="link-btn" type="button" title="暂未实现接口" @click="showForgotPasswordTooltip = true">
          忘记密码？
        </button>
      </div>
      
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      
      <button class="primary-btn" type="submit" :disabled="loading || !canLogin">
        <span v-if="loading" class="loading-spinner"></span>
        {{ loading ? '登录中...' : '登录' }}
      </button>
    </form>
    
    <p class="switch-tip">
      还没有账号？
      <button class="link-btn" type="button" @click="switchToRegister">立即注册</button>
    </p>
    
    <!-- 忘记密码提示 -->
    <div v-if="showForgotPasswordTooltip" class="tooltip" @click.self="showForgotPasswordTooltip = false">
      <div class="tooltip-content">
        <p>该功能暂未实现</p>
        <button class="tooltip-close" @click="showForgotPasswordTooltip = false">
          <svg viewBox="0 0 20 20" width="16" height="16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"/>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { authService } from '../services/auth'

export default {
  name: 'Login',
  emits: ['login-success', 'switch-to-register'],
  setup(props, { emit }) {
    // 响应式数据
    const username = ref('')
    const password = ref('')
    const rememberMe = ref(true)
    const showPassword = ref(false)
    const loading = ref(false)
    const errorMessage = ref('')
    const showForgotPasswordTooltip = ref(false)

    // 计算属性
    const usernameInvalid = computed(() => {
      return !(username.value && username.value.trim().length > 0)
    })
    
    const passwordInvalid = computed(() => {
      return !(password.value && password.value.length > 0)
    })
    
    const canLogin = computed(() => {
      return !usernameInvalid.value && !passwordInvalid.value
    })

    // 方法
    const clearError = () => {
      errorMessage.value = ''
    }
    
    const handleLogin = async () => {
      if (!canLogin.value) return
      errorMessage.value = ''
      
      try {
        loading.value = true
        const payload = {
          username: username.value.trim(),
          password: password.value,
          remember_me: rememberMe.value
        }

        const response = await authService.login(payload)

        if (response && response.user) {
          emit('login-success', response.user)
        } else {
          errorMessage.value = response?.message || '登录失败，请稍后重试'
        }
      } catch (error) {
        const data = error?.response?.data || {}
        errorMessage.value = data?.message || '登录失败，请检查账号或网络'
      } finally {
        loading.value = false
      }
    }
    
    const switchToRegister = () => {
      emit('switch-to-register')
    }

    return {
      username,
      password,
      rememberMe,
      showPassword,
      loading,
      errorMessage,
      showForgotPasswordTooltip,
      usernameInvalid,
      passwordInvalid,
      canLogin,
      clearError,
      handleLogin,
      switchToRegister
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  flex-direction: column;
  gap: 0;
  position: relative;
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

.input-hint {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
  margin-top: 4px;
  padding-left: 4px;
}

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

/* 辅助选项 */
.assist {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 8px;
}

.remember {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  position: relative;
}

.checkbox-custom {
  position: absolute;
  opacity: 0;
  cursor: pointer;
  height: 0;
  width: 0;
}

.checkbox-label {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  padding-left: 24px;
  position: relative;
}

.checkbox-label::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 18px;
  height: 18px;
  border: 2px solid var(--border-medium);
  border-radius: var(--radius);
  background: var(--bg-primary);
  transition: all var(--transition-fast);
}

.checkbox-custom:checked + .checkbox-label::before {
  background: var(--primary-600);
  border-color: var(--primary-600);
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 20 20' fill='%23ffffff'%3E%3Cpath fill-rule='evenodd' d='M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd'/%3E%3C/svg%3E");
  background-size: 12px;
  background-repeat: no-repeat;
  background-position: center;
}

.checkbox-custom:focus + .checkbox-label::before {
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
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

/* 提示框 */
.tooltip {
  position: absolute;
  top: 120px; /* 调整位置使其显示在"忘记密码"下方 */
  right: 0;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

.tooltip-content {
  background: var(--text-primary);
  color: var(--bg-primary);
  padding: 12px 16px;
  border-radius: var(--radius-lg);
  font-size: var(--font-size-sm);
  box-shadow: var(--shadow-lg);
  min-width: 180px;
  position: relative;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tooltip-content::before {
  content: '';
  position: absolute;
  top: -6px;
  right: 12px;
  width: 12px;
  height: 12px;
  background: var(--text-primary);
  transform: rotate(45deg);
}

.tooltip-close {
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 4px;
  border-radius: var(--radius);
  transition: color var(--transition-fast);
  display: flex;
  align-items: center;
  justify-content: center;
}

.tooltip-close:hover {
  color: var(--bg-primary);
}

.tooltip-content p {
  margin: 0;
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
  
  .assist {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .switch-tip {
    margin-top: 20px;
  }
}
</style>