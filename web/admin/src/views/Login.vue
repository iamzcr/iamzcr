<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
import { authApi } from '../api'

const router = useRouter()
const message = useMessage()

const loading = ref(false)

const formValue = ref({
  username: '',
  password: ''
})

const rules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  password: { required: true, message: '请输入密码', trigger: 'blur' }
}

async function handleLogin() {
  loading.value = true
  
  try {
    const res = await authApi.login(formValue.value)
    if (res.data.code === 0) {
      localStorage.setItem('admin_token', res.data.data.token)
      localStorage.setItem('admin_info', JSON.stringify(res.data.data))
      message.success('登录成功')
      router.push('/')
    } else {
      message.warning(res.data.message || '登录失败')
    }
  } catch (e: any) {
    console.error('Login error:', e)
    const errorMessage = e.response?.data?.message || e.message || '登录失败，请检查网络'
    message.error(errorMessage)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-bg"></div>
    <div class="login-card-wrap">
      <div class="login-brand">
        <span class="login-brand-icon">☰</span>
        <span class="login-brand-text">Blog Admin</span>
      </div>
      <n-card class="login-card" :bordered="false">
        <n-form
          :model="formValue"
          :rules="rules"
          label-placement="left"
          label-width="56px"
          size="large"
        >
          <n-form-item path="username" label="账户">
            <n-input 
              v-model:value="formValue.username" 
              placeholder="用户名" 
              @keyup.enter="handleLogin"
            />
          </n-form-item>
          <n-form-item path="password" label="密码">
            <n-input 
              v-model:value="formValue.password" 
              type="password"
              placeholder="密码" 
              @keyup.enter="handleLogin"
            />
          </n-form-item>
          <n-form-item>
            <n-button 
              type="primary" 
              :loading="loading" 
              block
              size="large"
              @click="handleLogin"
            >
              登 录
            </n-button>
          </n-form-item>
        </n-form>
      </n-card>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  width: 100%;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  overflow: hidden;
  background: var(--page-bg);
}

.login-bg {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(ellipse 80% 50% at 50% -20%, rgba(24, 160, 88, 0.06), transparent),
    radial-gradient(ellipse 60% 40% at 80% 80%, rgba(24, 160, 88, 0.04), transparent);
  pointer-events: none;
}

.login-card-wrap {
  position: relative;
  z-index: 1;
  width: 400px;
  animation: loginIn 0.5s ease both;
}

@keyframes loginIn {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.login-brand {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-bottom: 28px;
}

.login-brand-icon {
  font-size: 28px;
  color: var(--accent);
}

.login-brand-text {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-h);
  letter-spacing: 0.02em;
}

.login-card {
  border-radius: 14px;
  box-shadow:
    0 1px 3px rgba(0, 0, 0, 0.04),
    0 4px 16px rgba(0, 0, 0, 0.06),
    0 20px 60px rgba(0, 0, 0, 0.04);
  border: 1px solid var(--card-border);
}

.login-card :deep(.n-card__content) {
  padding: 28px 32px 24px;
}

.login-card :deep(.n-form-item) {
  margin-bottom: 18px;
}

.login-card :deep(.n-form-item:last-child) {
  margin-bottom: 0;
  margin-top: 8px;
}

.login-card :deep(.n-form-item-label) {
  color: var(--text);
  font-weight: 500;
}

.login-card :deep(.n-button--primary-type) {
  --n-height: 44px;
  --n-font-size: 15px;
  --n-font-weight: 600;
  letter-spacing: 0.3em;
}
</style>
