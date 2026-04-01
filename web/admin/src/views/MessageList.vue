<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSwitch } from 'naive-ui'
import { messageApi } from '../api'

const messages = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const pagination = ref({ page: 1, pageSize: 10, itemCount: 0 })
const editingMessage = ref({ id: 0, ip: '', name: '', email: '', url: '', is_reply: 0, content: '' })

function formatDate(time: number | string) {
  if (!time) return '-'
  const ts = typeof time === 'string' ? parseInt(time) : time
  const date = new Date(ts * 1000)
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const ho = String(date.getHours()).padStart(2, '0')
  const min = String(date.getMinutes()).padStart(2, '0')
  const s = String(date.getSeconds()).padStart(2, '0')
  return `${y}-${m}-${d} ${ho}:${min}:${s}`
}

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: 'IP', key: 'ip', width: 140 },
  { title: '姓名', key: 'name', width: 100 },
  { title: '邮箱', key: 'email', width: 180 },
  { title: 'URL', key: 'url', ellipsis: { tooltip: true } },
  { title: '内容', key: 'content', ellipsis: { tooltip: true } },
  { title: '已回复', key: 'is_reply', width: 80, render: (row: any) => h('span', row.is_reply === 1 ? '是' : '否') },
  { title: '时间', key: 'create_time', width: 180, render: (row: any) => formatDate(row.create_time) },
  { title: '操作', key: 'actions', width: 120, render: (row: any) => h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑') }
]

async function loadMessages() {
  loading.value = true
  try {
    const res = await messageApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    messages.value = res.data.data.list || res.data.data
    pagination.value.itemCount = res.data.data.total || 0
  } finally {
    loading.value = false
  }
}

function openEdit(row?: any) {
  editingMessage.value = row ? { ...row } : { id: 0, ip: '', name: '', email: '', url: '', is_reply: 0, content: '' }
  showModal.value = true
}

async function saveMessage() {
  if (editingMessage.value.id) {
    await messageApi.update(editingMessage.value.id, editingMessage.value)
  }
  showModal.value = false
  loadMessages()
}

onMounted(loadMessages)
</script>

<template>
  <div>
    <n-data-table :columns="columns" :data="messages" :loading="loading" remote :pagination="pagination" @update:page="pagination.page = $event; loadMessages()" />
    <n-modal v-model:show="showModal" preset="card" title="留言管理" style="width: 500px">
      <n-form :model="editingMessage">
        <n-form-item label="姓名">
          <n-input v-model:value="editingMessage.name" />
        </n-form-item>
        <n-form-item label="邮箱">
          <n-input v-model:value="editingMessage.email" />
        </n-form-item>
        <n-form-item label="URL">
          <n-input v-model:value="editingMessage.url" />
        </n-form-item>
        <n-form-item label="内容">
          <n-input v-model:value="editingMessage.content" type="textarea" :rows="3" />
        </n-form-item>
        <n-form-item label="已回复">
          <n-switch v-model:value="editingMessage.is_reply" :checked-value="1" :unchecked-value="0">
            <template #checked>是</template>
            <template #unchecked>否</template>
          </n-switch>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="saveMessage">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>
