<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NCard, NForm, NFormItem, NInput, NButton, NSpace, NSpin } from 'naive-ui'
import { websiteApi } from '../api'

const loading = ref(false)
const settings = ref<Record<string, string>>({})

async function loadSettings() {
  loading.value = true
  try {
    const res = await websiteApi.get()
    settings.value = res.data.data
  } finally {
    loading.value = false
  }
}

async function saveSettings() {
  loading.value = true
  try {
    await websiteApi.update(settings.value)
    alert('保存成功')
  } finally {
    loading.value = false
  }
}

onMounted(loadSettings)
</script>

<template>
  <n-spin :show="loading">
    <n-card title="网站设置">
      <n-form label-placement="left">
        <n-form-item v-for="(_, key) in settings" :key="key" :label="key">
          <n-input v-model:value="settings[key]" />
        </n-form-item>
        <n-form-item>
          <n-space>
            <n-button type="primary" @click="saveSettings" :loading="loading">保存</n-button>
          </n-space>
        </n-form-item>
      </n-form>
    </n-card>
  </n-spin>
</template>