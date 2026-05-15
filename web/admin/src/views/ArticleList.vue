<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { NDataTable, NButton, NTag, NSpace, NModal, NCheckboxGroup, NCheckbox, NInput, useMessage } from 'naive-ui'
import { articleApi, platformApi } from '../api'

const router = useRouter()
const message = useMessage()

const articles = ref<any[]>([])
const loading = ref(false)
const platforms = ref<any[]>([])
const wechatStatusMap = ref<Record<number, boolean>>({})
const publishingMap = ref<Record<number, boolean>>({})
const pagination = ref({ page: 1, pageSize: 10, itemCount: 0 })
const searchKeyword = ref('')

const showPublishModal = ref(false)
const selectedPlatformIds = ref<number[]>([])
const publishTargetId = ref(0)

function formatDate(time: number | string) {
  if (!time) return '-'
  const ts = typeof time === 'string' ? parseInt(time) : time
  const date = new Date(ts * 1000)
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
  { title: '微信', key: 'wechat', width: 90, render: (row: any) => h(NTag, { type: wechatStatusMap.value[row.id] ? 'success' : 'default', size: 'small' }, () => wechatStatusMap.value[row.id] ? '已发布' : '未发布') },
  { title: '热门', key: 'is_hot', width: 60, render: (row: any) => h(NTag, { type: row.is_hot === 1 ? 'error' : 'default', size: 'small' }, () => row.is_hot === 1 ? '是' : '否') },
  { title: '权重', key: 'weight', width: 60 },
  { title: '发布时间', key: 'public_time', width: 180, render: (row: any) => formatDate(row.public_time) },
  { title: '创建时间', key: 'create_time', width: 180, render: (row: any) => formatDate(row.create_time) },
  { title: '操作', key: 'actions', width: 220, render: (row: any) => h(NSpace, () => [
    h(NButton, { size: 'small', onClick: () => router.push(`/articles/edit/${row.id}`) }, () => '编辑'),
    h(NButton, { size: 'small', type: 'info', loading: publishingMap.value[row.id], onClick: () => openPublishModal(row.id) }, () => '发布'),
    h(NButton, { size: 'small', type: 'error', onClick: () => deleteArticle(row.id) }, () => '删除')
  ])}
]

async function loadArticles() {
  loading.value = true
  try {
    const res = await articleApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize, keyword: searchKeyword.value })
    articles.value = res.data.data.list
    pagination.value.itemCount = res.data.data.total
    loadWechatStatus()
  } finally {
    loading.value = false
  }
}

async function loadPlatforms() {
  try {
    const res = await platformApi.list({ page: 1, page_size: 100 })
    platforms.value = (res.data.data.list || []).filter((p: any) => p.status === 1)
  } catch { /* ignore */ }
}

async function loadWechatStatus() {
  const wechatPlatform = platforms.value.find((p: any) => p.mark === 'wechat')
  for (const article of articles.value) {
    try {
      const res = await articleApi.getMedia(article.id)
      const records = res.data.data || []
      wechatStatusMap.value[article.id] = records.some((r: any) => r.platform_id === wechatPlatform?.id && r.status === 1)
    } catch {
      wechatStatusMap.value[article.id] = false
    }
  }
}

function openPublishModal(id: number) {
  publishTargetId.value = id
  selectedPlatformIds.value = platforms.value.filter((p: any) => p.status === 1).map((p: any) => p.id)
  showPublishModal.value = true
}

async function publishNow() {
  if (selectedPlatformIds.value.length === 0) return
  const id = publishTargetId.value
  publishingMap.value[id] = true
  showPublishModal.value = false
  try {
    await articleApi.publishToMedia(id, selectedPlatformIds.value)
    message.success('发布成功')
    wechatStatusMap.value[id] = true
  } catch (e: any) {
    message.error(e.response?.data?.message || '发布失败')
  } finally {
    publishingMap.value[id] = false
  }
}

async function deleteArticle(id: number) {
  await articleApi.delete(id)
  message.success('删除成功')
  loadArticles()
}

onMounted(() => {
  loadPlatforms().then(() => loadArticles())
})
</script>

<template>
  <div class="page-wrap">
    <div class="page-toolbar">
      <div style="display: flex; gap: 12px; align-items: center;">
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索文章标题..."
          clearable
          style="width: 240px"
          @keyup.enter="pagination.page = 1; loadArticles()"
          @clear="pagination.page = 1; loadArticles()"
        />
        <n-button type="primary" @click="router.push('/articles/new')">新建文章</n-button>
      </div>
    </div>
    <n-data-table
      :columns="columns"
      :data="articles"
      :loading="loading"
      remote
      :pagination="pagination"
      @update:page="pagination.page = $event; loadArticles()"
    />
    <n-modal v-model:show="showPublishModal" preset="card" title="选择发布平台" style="width: 400px">
      <n-checkbox-group v-model:value="selectedPlatformIds">
        <n-space vertical>
          <n-checkbox v-for="p in platforms" :key="p.id" :value="p.id">{{ p.name }}</n-checkbox>
        </n-space>
      </n-checkbox-group>
      <template #footer>
        <n-button @click="showPublishModal = false">取消</n-button>
        <n-button type="primary" @click="publishNow">确定发布</n-button>
      </template>
    </n-modal>
  </div>
</template>