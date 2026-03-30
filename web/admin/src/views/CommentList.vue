<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton } from 'naive-ui'
import { commentApi } from '../api'

const comments = ref<any[]>([])
const loading = ref(false)
const pagination = ref({ page: 1, pageSize: 10, total: 0 })

function formatDate(time: number | string) {
  if (!time) return '-'
  const date = new Date(time)
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
  { title: '文章ID', key: 'aid', width: 80 },
  { title: '作者', key: 'name', width: 100 },
  { title: '邮箱', key: 'email', width: 180 },
  { title: '内容', key: 'content', ellipsis: { tooltip: true } },
  { title: '时间', key: 'create_time', width: 180, render: (row: any) => formatDate(row.create_time) },
  { title: '操作', key: 'actions', width: 100, render: (row: any) => h(NButton, { size: 'small', type: 'error', onClick: () => deleteComment(row.id) }, () => '删除') }
]

async function loadComments() {
  loading.value = true
  try {
    const res = await commentApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    comments.value = res.data.data.list || res.data.data
    pagination.value.total = res.data.data.total || 0
  } finally {
    loading.value = false
  }
}

async function deleteComment(id: number) {
  if (confirm('确定要删除吗？')) {
    await commentApi.delete(id)
    loadComments()
  }
}

onMounted(loadComments)
</script>

<template>
  <div>
    <n-data-table :columns="columns" :data="comments" :loading="loading" :pagination="pagination" @update:page="pagination.page = $event; loadComments()" />
  </div>
</template>
