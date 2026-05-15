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
  <div class="guestbook-wrap">
    <div class="guestbook-card">
      <div class="card-accent"></div>

      <div class="guestbook-layout">
        <section class="form-section">
          <div class="section-head">
            <span class="head-icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
                <path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/>
              </svg>
            </span>
            <h2>写留言</h2>
          </div>

          <div class="form-body">
            <div class="field-row">
              <div class="field">
                <label class="field-label">昵称</label>
                <n-input v-model:value="form.name" placeholder="你的昵称" size="large" :input-props="{ autocomplete: 'name' }" />
              </div>
              <div class="field">
                <label class="field-label required">邮箱</label>
                <n-input v-model:value="form.email" placeholder="your@email.com" size="large" :input-props="{ autocomplete: 'email' }" />
              </div>
            </div>
            <div class="field">
              <label class="field-label">网址</label>
              <n-input v-model:value="form.url" placeholder="https://" size="large" :input-props="{ autocomplete: 'url' }" />
            </div>
            <div class="field field-content">
              <label class="field-label required">内容</label>
              <n-input
                v-model:value="form.content"
                type="textarea"
                :rows="8"
                placeholder="分享你的想法、建议或者随便聊聊..."
                size="large"
                class="content-input"
              />
              <span class="char-hint" v-if="form.content">{{ form.content.length }} 字</span>
            </div>
            <n-button
              type="primary"
              :loading="submitting"
              @click="submitMessage"
              size="large"
              class="submit-btn"
            >
              <template v-if="!submitting">提交留言</template>
              <template v-else>提交中...</template>
            </n-button>
          </div>
        </section>

        <section class="list-section">
          <div class="section-head">
            <span class="head-icon">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
              </svg>
            </span>
            <h2>留言</h2>
            <span class="count-badge" v-if="pagination.itemCount">{{ pagination.itemCount }}</span>
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
              <a v-if="m.url" :href="m.url" target="_blank" rel="noopener" class="msg-url">↗ {{ m.url }}</a>
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
            <div class="empty-graphic">—— · · · ——</div>
            <p>还没有留言，来做第一个留言的人吧</p>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<style scoped>
.guestbook-wrap {
  padding: 48px 24px 80px;
  display: flex;
  justify-content: center;
}

.guestbook-card {
  width: 100%;
  background: var(--card-bg);
  border: 1px solid var(--border);
  border-radius: 16px;
  position: relative;
  overflow: hidden;
  box-shadow: var(--shadow);
}

.card-accent {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--accent), #4ecb71 40%, var(--accent) 80%);
}

.guestbook-layout {
  display: grid;
  grid-template-columns: 2fr 3fr;
}

/* ── sections ── */
.form-section {
  padding: 44px 40px 40px;
  border-right: 1px solid var(--border);
  background: var(--bg);
  min-width: 280px;
}

.list-section {
  flex: 1;
  padding: 44px 40px 40px;
  min-width: 0;
}

.section-head {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 32px;
}

.section-head h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: var(--text-h);
  letter-spacing: 0.5px;
}

.head-icon {
  color: var(--accent);
  display: flex;
  align-items: center;
}

.count-badge {
  margin-left: auto;
  font-size: 12px;
  font-weight: 700;
  color: var(--text-muted);
  background: var(--card-bg);
  padding: 2px 10px;
  border-radius: 10px;
  border: 1px solid var(--border);
}

/* ── form ── */
.form-body {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.field-row {
  display: flex;
  gap: 14px;
}

.field-row .field {
  flex: 1;
  min-width: 0;
}

.field-label {
  display: block;
  font-size: 11px;
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

.field-content {
  position: relative;
}

.char-hint {
  position: absolute;
  right: 4px;
  bottom: -18px;
  font-size: 11px;
  color: var(--text-muted);
  font-family: var(--mono);
}

.submit-btn {
  margin-top: 6px;
  font-weight: 600;
  letter-spacing: 0.5px;
  border-radius: 10px;
  height: 46px;
  font-size: 15px;
  width: 100%;
}

/* ── message list ── */
.message-list {
  display: flex;
  flex-direction: column;
}

.msg-card {
  padding: 22px 0;
  border-bottom: 1px solid var(--border);
  transition: background 0.2s;
}

.msg-card:first-child {
  padding-top: 0;
}

.msg-card:hover {
  background: var(--bg);
  margin: 0 -16px;
  padding-left: 16px;
  padding-right: 16px;
  border-radius: 10px;
}

.msg-top {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}

.msg-avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: var(--accent-bg);
  color: var(--accent);
  font-weight: 700;
  font-size: 13px;
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
}

.msg-name {
  font-weight: 600;
  font-size: 14px;
  color: var(--text-h);
}

.msg-time {
  font-size: 11px;
  color: var(--text-muted);
  font-family: var(--mono);
}

.msg-body {
  margin: 0;
  color: var(--text);
  line-height: 1.85;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 14.5px;
}

.msg-url {
  display: inline-block;
  margin-top: 6px;
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
  margin-top: 28px;
}

/* ── empty ── */
.empty-state {
  text-align: center;
  padding: 72px 0;
  color: var(--text-muted);
}

.empty-state p {
  font-size: 14px;
  margin-top: 12px;
}

.empty-graphic {
  font-size: 13px;
  color: var(--border);
  letter-spacing: 4px;
}

/* ── responsive ── */
@media (max-width: 768px) {
  .guestbook-wrap {
    padding: 24px 12px 48px;
  }

  .guestbook-layout {
    grid-template-columns: 1fr;
    display: block;
  }

  .form-section {
    width: 100%;
    min-width: unset;
    border-right: none;
    border-bottom: 1px solid var(--border);
    padding: 36px 24px 32px;
  }

  .list-section {
    padding: 36px 24px 40px;
  }

  .field-row {
    flex-direction: column;
  }

  .msg-card:hover {
    margin: 0 -12px;
    padding-left: 12px;
    padding-right: 12px;
  }
}
</style>
