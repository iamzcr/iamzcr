<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NInput, NButton, NSpace, NForm, NFormItem, NCard, NDivider, NEmpty, NMessageProvider, useMessage } from 'naive-ui'
import { messageApi } from '../api'

const msg = useMessage()
const messages = ref<any[]>([])
const loading = ref(false)
const submitting = ref(false)
const form = ref({ name: '', email: '', url: '', content: '' })

function formatDate(time: number) {
  if (!time) return '-'
  const date = new Date(time * 1000)
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const h = String(date.getHours()).padStart(2, '0')
  const min = String(date.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${d} ${h}:${min}`
}

async function loadMessages() {
  loading.value = true
  try {
    const res = await messageApi.list()
    messages.value = res.data.data || []
  } finally {
    loading.value = false
  }
}

async function submitMessage() {
  if (!form.value.email || !form.value.content) {
    msg.warning('请输入邮箱和留言内容')
    return
  }
  submitting.value = true
  try {
    await messageApi.create(form.value)
    msg.success('留言成功')
    form.value.content = ''
    loadMessages()
  } catch (e: any) {
    msg.error(e.response?.data?.message || '留言失败')
  } finally {
    submitting.value = false
  }
}

onMounted(loadMessages)
</script>

<template>
  <n-message-provider>
    <div class="page-wrap messages-page">
    <n-card title="留言板" class="form-card">
      <n-form :model="form" label-placement="top">
        <n-space vertical>
          <n-grid :cols="2" :x-gap="16" style="width: 100%">
            <n-grid-item>
              <n-form-item label="昵称">
                <n-input v-model:value="form.name" placeholder="你的昵称（选填）" />
              </n-form-item>
            </n-grid-item>
            <n-grid-item>
              <n-form-item label="邮箱 *">
                <n-input v-model:value="form.email" placeholder="your@email.com" />
              </n-form-item>
            </n-grid-item>
          </n-grid>
          <n-form-item label="网址">
            <n-input v-model:value="form.url" placeholder="https://（选填）" />
          </n-form-item>
          <n-form-item label="内容 *">
            <n-input v-model:value="form.content" type="textarea" :rows="4" placeholder="说点什么吧..." />
          </n-form-item>
          <n-button type="primary" :loading="submitting" @click="submitMessage" block>
            {{ submitting ? '提交中...' : '提交留言' }}
          </n-button>
        </n-space>
      </n-form>
    </n-card>

    <n-divider />

    <div v-if="messages.length > 0">
      <n-card v-for="m in messages" :key="m.id" class="message-card" size="small">
        <div class="msg-header">
          <span class="msg-name">{{ m.name || '匿名' }}</span>
          <span class="msg-time">{{ formatDate(m.create_time) }}</span>
        </div>
        <div class="msg-content">{{ m.content }}</div>
        <div class="msg-meta" v-if="m.url">
          <a :href="m.url" target="_blank">{{ m.url }}</a>
        </div>
      </n-card>
    </div>
    <n-empty v-else-if="!loading" description="暂无留言，来做第一个留言的人吧" />
  </div>
  </n-message-provider>
</template>

<style scoped>
.messages-page {
  max-width: 680px;
  margin: 0 auto;
}

.form-card {
  margin-bottom: 16px;
}

.message-card {
  margin-bottom: 12px;
}

.msg-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.msg-name {
  font-weight: 600;
  color: var(--text);
}

.msg-time {
  font-size: 12px;
  color: var(--text-secondary);
}

.msg-content {
  color: var(--text);
  line-height: 1.7;
  white-space: pre-wrap;
}

.msg-meta {
  margin-top: 8px;
  font-size: 12px;
}

.msg-meta a {
  color: var(--primary);
}
</style>
