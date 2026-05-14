<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NInput, NButton, NPagination, useMessage } from 'naive-ui'
import { messageApi } from '../api'

const msg = useMessage()
const messages = ref<any[]>([])
const loading = ref(false)
const submitting = ref(false)
const pagination = ref({ page: 1, pageSize: 3, itemCount: 0 })
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
    const res = await messageApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    const data = res.data.data
    messages.value = data.list || data || []
    pagination.value.itemCount = data.total || 0
  } finally {
    loading.value = false
  }
}

function onPageChange(page: number) {
  pagination.value.page = page
  loadMessages()
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
    pagination.value.page = 1
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
  <div class="messages-layout">
    <aside class="form-panel">
      <div class="panel-header">
        <span class="panel-icon">></span>
        <h2>写留言</h2>
      </div>
      <div class="form-body">
        <div class="field">
          <label class="field-label">昵称</label>
          <n-input v-model:value="form.name" placeholder="你的昵称" size="large" />
        </div>
        <div class="field">
          <label class="field-label required">邮箱</label>
          <n-input v-model:value="form.email" placeholder="your@email.com" size="large" />
        </div>
        <div class="field">
          <label class="field-label">网址</label>
          <n-input v-model:value="form.url" placeholder="https://" size="large" />
        </div>
        <div class="field">
          <label class="field-label required">内容</label>
          <n-input v-model:value="form.content" type="textarea" :rows="5" placeholder="说点什么..." size="large" />
        </div>
        <n-button
          type="primary"
          :loading="submitting"
          @click="submitMessage"
          size="large"
          class="submit-btn"
          block
        >
          <template v-if="!submitting">
            <span class="btn-arrow">→</span>
            提交留言
          </template>
          <template v-else>提交中...</template>
        </n-button>
      </div>
    </aside>

    <main class="list-panel">
      <div class="panel-header">
        <span class="count-badge" v-if="pagination.itemCount">{{ pagination.itemCount }} 条留言</span>
        <span class="count-badge" v-else>暂无留言</span>
      </div>

      <div class="message-list" v-if="messages.length > 0">
        <article v-for="m in messages" :key="m.id" class="msg-card">
          <div class="msg-top">
            <div class="msg-avatar">{{ (m.name || '匿').charAt(0) }}</div>
            <div class="msg-meta">
              <span class="msg-name">{{ m.name || '匿名' }}</span>
              <time class="msg-time">{{ formatDate(m.create_time) }}</time>
            </div>
          </div>
          <p class="msg-body">{{ m.content }}</p>
          <a v-if="m.url" :href="m.url" target="_blank" class="msg-url">↗ {{ m.url }}</a>
        </article>

        <div class="pagination-wrap">
          <n-pagination
            v-model:page="pagination.page"
            :page-size="pagination.pageSize"
            :item-count="pagination.itemCount"
            @update:page="onPageChange"
          />
        </div>
      </div>

      <div class="empty-state" v-else-if="!loading">
        <div class="empty-graphic">
          <span class="empty-bracket">{</span>
          <span class="empty-dots">···</span>
          <span class="empty-bracket">}</span>
        </div>
        <p>还没有留言，来做第一个留言的人吧</p>
      </div>
    </main>
  </div>
</template>

<style scoped>
.messages-layout {
  display: flex;
  gap: 0;
  max-width: 1040px;
  margin: 0 auto;
  min-height: calc(100vh - 180px);
}

/* ── left panel ── */
.form-panel {
  width: 360px;
  flex-shrink: 0;
  padding: 48px 36px;
  border-right: 1px solid var(--border);
  position: sticky;
  top: 80px;
  align-self: flex-start;
  max-height: calc(100vh - 120px);
  overflow-y: auto;
}

.form-panel .panel-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 36px;
}

.form-panel .panel-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
  color: var(--text-h);
  letter-spacing: 0.5px;
}

.panel-icon {
  font-family: var(--mono);
  font-size: 14px;
  color: var(--accent);
  font-weight: 700;
}

.form-body {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.field-label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 6px;
}

.field-label.required::after {
  content: ' *';
  color: #d03050;
}

.submit-btn {
  margin-top: 8px;
  font-weight: 600;
  letter-spacing: 0.5px;
  border-radius: 8px;
  height: 44px;
  font-size: 15px;
}

.btn-arrow {
  font-family: var(--mono);
  margin-right: 4px;
  font-size: 16px;
}

/* ── right panel ── */
.list-panel {
  flex: 1;
  padding: 48px 40px;
  min-width: 0;
}

.list-panel .panel-header {
  margin-bottom: 32px;
}

.count-badge {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  letter-spacing: 0.5px;
}

.message-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

/* ── message card ── */
.msg-card {
  padding: 24px 0;
  border-bottom: 1px solid var(--border);
  transition: background 0.2s;
}

.msg-card:first-child {
  padding-top: 0;
}

.msg-card:hover {
  background: var(--card-bg);
  margin: 0 -16px;
  padding-left: 16px;
  padding-right: 16px;
  border-radius: 8px;
}

.msg-top {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.msg-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--accent-bg);
  color: var(--accent);
  font-weight: 700;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-family: var(--mono);
}

.msg-meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.msg-name {
  font-weight: 600;
  font-size: 14px;
  color: var(--text-h);
}

.msg-time {
  font-size: 12px;
  color: var(--text-muted);
  font-family: var(--mono);
}

.msg-body {
  margin: 0;
  color: var(--text);
  line-height: 1.8;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 15px;
}

.msg-url {
  display: inline-block;
  margin-top: 8px;
  font-size: 12px;
  color: var(--accent);
  text-decoration: none;
  font-family: var(--mono);
}

.msg-url:hover {
  text-decoration: underline;
}

.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 32px;
}

/* ── empty state ── */
.empty-state {
  text-align: center;
  padding: 64px 0;
  color: var(--text-muted);
}

.empty-state p {
  font-size: 14px;
  margin-top: 16px;
}

.empty-graphic {
  font-family: var(--mono);
  font-size: 32px;
  color: var(--border);
  letter-spacing: 4px;
}

.empty-bracket {
  color: var(--text-muted);
}

.empty-dots {
  color: var(--border);
}

/* ── responsive ── */
@media (max-width: 768px) {
  .messages-layout {
    flex-direction: column;
  }

  .form-panel {
    width: 100%;
    position: static;
    border-right: none;
    border-bottom: 1px solid var(--border);
    padding: 32px 20px;
    max-height: none;
  }

  .list-panel {
    padding: 32px 20px;
  }
}
</style>
