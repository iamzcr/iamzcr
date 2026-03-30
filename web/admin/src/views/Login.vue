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
  <div class="login-container">
    <n-card title="管理员登录" class="login-card">
      <n-form
        :model="formValue"
        :rules="rules"
        label-placement="left"
      >
        <n-form-item path="username" label="用户名">
          <n-input 
            v-model:value="formValue.username" 
            placeholder="请输入用户名" 
            @keyup.enter="handleLogin"
          />
        </n-form-item>
        <n-form-item path="password" label="密码">
          <n-input 
            v-model:value="formValue.password" 
            type="password"
            placeholder="请输入密码" 
            @keyup.enter="handleLogin"
          />
        </n-form-item>
        <n-form-item>
          <n-button 
            type="primary" 
            :loading="loading" 
            block
            @click="handleLogin"
          >
            登录
          </n-button>
        </n-form-item>
      </n-form>
    </n-card>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f0f2f5;
}

.login-card {
  width: 400px;
}
</style>
