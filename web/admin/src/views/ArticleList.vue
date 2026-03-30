<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { NDataTable, NButton, NTag } from 'naive-ui'
import { articleApi } from '../api'

const router = useRouter()

const articles = ref<any[]>([])
const loading = ref(false)
const pagination = ref({ page: 1, pageSize: 10, total: 0 })

function formatDate(time: number | string) {
  if (!time) return '-'
  const date = new Date(time)
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const h = String(date.getHours()).padStart(2, '0')
  const min = String(date.getMinutes()).padStart(2, '0')
  const s = String(date.getSeconds()).padStart(2, '0')
  return `${y}-${m}-${d} ${h}:${min}:${s}`
}

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '标题', key: 'title', ellipsis: { tooltip: true } },
  { title: '作者', key: 'author', width: 100 },
  { title: '状态', key: 'status', width: 80, render: (row: any) => h(NTag, { type: row.status === 1 ? 'success' : 'default', size: 'small' }, () => row.status === 1 ? '已发布' : '草稿') },
  { title: '热门', key: 'is_hot', width: 60, render: (row: any) => h(NTag, { type: row.is_hot === 1 ? 'error' : 'default', size: 'small' }, () => row.is_hot === 1 ? '是' : '否') },
  { title: '权重', key: 'weight', width: 60 },
  { title: '创建时间', key: 'create_time', width: 180, render: (row: any) => formatDate(row.create_time) },
  { title: '操作', key: 'actions', width: 120, render: (row: any) => h(NButton, { size: 'small', onClick: () => router.push(`/articles/${row.id}`) }, () => '编辑') }
]

async function loadArticles() {
  loading.value = true
  try {
    const res = await articleApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    articles.value = res.data.data.list
    pagination.value.total = res.data.data.total
  } finally {
    loading.value = false
  }
}

onMounted(loadArticles)
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: flex-end;">
      <n-button type="primary" @click="router.push('/articles/new')">新建文章</n-button>
    </div>
    <n-data-table
      :columns="columns"
      :data="articles"
      :loading="loading"
      :pagination="pagination"
      @update:page="pagination.page = $event; loadArticles()"
    />
  </div>
</template>
