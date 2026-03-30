<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NDataTable, NInput } from 'naive-ui'
import { readApi, articleApi } from '../api'

const reads = ref<any[]>([])
const loading = ref(false)
const articleMap = ref<Record<number, string>>({})
const aid = ref('')

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
  { title: '文章ID', key: 'aid', width: 80 },
  { title: '文章标题', key: 'title', render: (row: any) => articleMap.value[row.aid] || '未知' },
  { title: '来源', key: 'referer', ellipsis: { tooltip: true } },
  { title: 'IP', key: 'ip', width: 140 },
  { title: '时间', key: 'create_time', width: 180, render: (row: any) => formatDate(row.create_time) }
]

async function loadArticles() {
  const res = await articleApi.list({ page: 1, page_size: 1000 })
  const list = res.data.data.list || []
  const map: Record<number, string> = {}
  list.forEach((a: any) => { map[a.id] = a.title })
  articleMap.value = map
}

async function loadReads() {
  loading.value = true
  try {
    const params = aid.value ? { aid: aid.value } : {}
    const res = await readApi.list(params)
    reads.value = res.data.data
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadArticles()
  loadReads()
})
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; gap: 8px; align-items: center;">
      <n-input v-model:value="aid" placeholder="文章ID" style="width: 120px" />
      <n-button type="primary" @click="loadReads">搜索</n-button>
    </div>
    <n-data-table :columns="columns" :data="reads" :loading="loading" />
  </div>
</template>
