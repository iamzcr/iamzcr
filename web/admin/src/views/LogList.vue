<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, useMessage } from 'naive-ui'
import { logApi } from '../api'

const logs = ref<any[]>([])
const loading = ref(false)
const pagination = ref({ page: 1, pageSize: 10, itemCount: 0 })
const message = useMessage()

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
  { title: '用户名', key: 'username', width: 120 },
  { title: 'IP', key: 'ip', width: 140 },
  { title: '内容', key: 'content', ellipsis: { tooltip: true } },
  { title: '类型', key: 'type', width: 80 },
  { title: '时间', key: 'create_time', width: 180, render: (row: any) => formatDate(row.create_time) },
  { title: '操作', key: 'actions', width: 120, render: (row: any) => h(NButton, { size: 'small', type: 'error', onClick: () => deleteLog(row.id) }, () => '删除') }
]

async function loadLogs() {
  loading.value = true
  try {
    const res = await logApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    logs.value = res.data.data.list || res.data.data
    pagination.value.itemCount = res.data.data.total || 0
  } finally {
    loading.value = false
  }
}

async function deleteLog(id: number) {
  await logApi.delete(id)
  message.success('删除成功')
  loadLogs()
}

onMounted(loadLogs)
</script>

<template>
  <div>
    <n-data-table :columns="columns" :data="logs" :loading="loading" remote :pagination="pagination" @update:page="pagination.page = $event; loadLogs()" />
  </div>
</template>
