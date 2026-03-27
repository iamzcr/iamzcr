<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NPopconfirm, NTag } from 'naive-ui'
import { commentApi } from '../api'

const comments = ref<any[]>([])
const loading = ref(false)

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '文章ID', key: 'article_id', width: 80 },
  { title: '作者', key: 'author', width: 100 },
  { title: '邮箱', key: 'email', width: 180 },
  { title: '内容', key: 'content', ellipsis: { tooltip: true } },
  { title: '状态', key: 'status', width: 80, render: (row: any) => 
    h(NTag, { type: row.status === 1 ? 'success' : 'warning', size: 'small' }, 
      () => row.status === 1 ? '已审核' : '待审核') 
  },
  { title: '时间', key: 'created_at', width: 180 },
  { title: '操作', key: 'actions', width: 120, render: (row: any) => 
    h(NPopconfirm, { onPositiveClick: () => deleteComment(row.id) }, 
      () => h(NButton, { size: 'small', type: 'error' }, () => '删除'))
  }
]

async function loadComments() {
  loading.value = true
  try {
    const res = await commentApi.list()
    comments.value = res.data.data
  } finally {
    loading.value = false
  }
}

async function deleteComment(id: number) {
  await commentApi.delete(id)
  loadComments()
}

onMounted(loadComments)
</script>

<template>
  <n-data-table :columns="columns" :data="comments" :loading="loading" />
</template>