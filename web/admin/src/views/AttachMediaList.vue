<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NTag, NSpace, useMessage } from 'naive-ui'
import { attachMediaApi } from '../api'

const records = ref<any[]>([])
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
  const mi = String(date.getMinutes()).padStart(2, '0')
  const s = String(date.getSeconds()).padStart(2, '0')
  return `${y}-${m}-${d} ${ho}:${mi}:${s}`
}

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '附件ID', key: 'attach_id', width: 80 },
  { title: '平台ID', key: 'platform_id', width: 80 },
  { title: '媒体ID', key: 'media_id', ellipsis: { tooltip: true } },
  { title: '媒体URL', key: 'media_url', ellipsis: { tooltip: true } },
  { title: '状态', key: 'status', width: 80, render: (row: any) => {
    if (row.status === 1) return h(NTag, { type: 'success', size: 'small' }, () => '已同步')
    if (row.status === 2) return h(NTag, { type: 'error', size: 'small' }, () => '失败')
    return h(NTag, { type: 'default', size: 'small' }, () => '待处理')
  }},
  { title: '错误信息', key: 'error_msg', ellipsis: { tooltip: true }, width: 150 },
  { title: '更新时间', key: 'update_time', width: 180, render: (row: any) => formatDate(row.update_time) },
  { title: '操作', key: 'actions', width: 100, render: (row: any) => h(NSpace, () => [
    h(NButton, { size: 'small', type: 'error', onClick: () => deleteRecord(row.id) }, () => '删除')
  ])}
]

async function loadRecords() {
  loading.value = true
  try {
    const res = await attachMediaApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    records.value = res.data.data.list || res.data.data
    pagination.value.itemCount = res.data.data.total || 0
  } finally {
    loading.value = false
  }
}

async function deleteRecord(id: number) {
  await attachMediaApi.delete(id)
  message.success('删除成功')
  loadRecords()
}

onMounted(loadRecords)
</script>

<template>
  <div>
    <n-data-table 
      :columns="columns" 
      :data="records" 
      :loading="loading"
      remote
      :pagination="pagination"
      @update:page="pagination.page = $event; loadRecords()"
    />
  </div>
</template>