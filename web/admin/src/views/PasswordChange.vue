<script setup lang="ts">
import { ref } from 'vue'
import { NCard, NForm, NFormItem, NInput, NButton, NSpace, useMessage } from 'naive-ui'
import { adminApi } from '../api'

const message = useMessage()
const loading = ref(false)
const form = ref({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

async function changePassword() {
  if (!form.value.old_password) {
    message.error('请输入原密码')
    return
  }
  if (!form.value.new_password) {
    message.error('请输入新密码')
    return
  }
  if (form.value.new_password !== form.value.confirm_password) {
    message.error('两次输入的密码不一致')
    return
  }
  if (form.value.new_password.length < 6) {
    message.error('密码长度不能少于6位')
    return
  }

  loading.value = true
  try {
    await adminApi.changePassword({
      old_password: form.value.old_password,
      new_password: form.value.new_password
    })
    message.success('密码修改成功')
    form.value = { old_password: '', new_password: '', confirm_password: '' }
  } catch (e: any) {
    message.error(e.response?.data?.msg || '修改失败')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <n-card title="修改密码">
    <n-form :model="form" label-placement="left" label-width="100">
      <n-form-item label="原密码">
        <n-input v-model:value="form.old_password" type="password" placeholder="请输入原密码" />
      </n-form-item>
      <n-form-item label="新密码">
        <n-input v-model:value="form.new_password" type="password" placeholder="请输入新密码" />
      </n-form-item>
      <n-form-item label="确认密码">
        <n-input v-model:value="form.confirm_password" type="password" placeholder="请再次输入新密码" />
      </n-form-item>
      <n-form-item>
        <n-space>
          <n-button type="primary" @click="changePassword" :loading="loading">确认修改</n-button>
        </n-space>
      </n-form-item>
    </n-form>
  </n-card>
</template>
