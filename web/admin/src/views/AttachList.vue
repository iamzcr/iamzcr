<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSelect, NTag, NUpload, NImage, NSpace, useMessage } from 'naive-ui'
import type { UploadFileInfo } from 'naive-ui'
import { attachApi, websiteApi } from '../api'

const message = useMessage()
const attaches = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const uploadLoading = ref(false)
const pagination = ref({ page: 1, pageSize: 10, itemCount: 0 })
const editingAttach = ref({ id: 0, name: '', link: '', path: '', status: 1, type: 1 })
const cdnUrl = ref('')

function getFullUrl(path: string) {
  if (!path) return ''
  if (path.startsWith('http://') || path.startsWith('https://')) return path
  if (!cdnUrl.value) return path
  return cdnUrl.value.replace(/\/+$/, '') + '/' + path.replace(/^\/+/, '')
}

const typeOptions = [
  { label: '图片', value: 1 },
  { label: '视频', value: 2 }
]

const statusOptions = [
  { label: '禁用', value: 0 },
  { label: '启用', value: 1 }
]

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
  { title: '预览', key: 'link', width: 80, render: (row: any) => {
    if (row.type === 1 && row.link) {
      return h(NImage, { width: 48, height: 48, src: getFullUrl(row.link), style: { objectFit: 'cover', borderRadius: '4px' } })
    }
    return h('span', '-')
  }},
  { title: '名称', key: 'name', ellipsis: { tooltip: true } },
  { title: '链接', key: 'link', ellipsis: { tooltip: true } },
  { title: '类型', key: 'type', width: 80, render: (row: any) => h(NTag, { type: row.type === 1 ? 'warning' : 'info', size: 'small' }, () => row.type === 1 ? '图片' : '视频') },
  { title: '状态', key: 'status', width: 80, render: (row: any) => h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small' }, () => row.status === 1 ? '启用' : '禁用') },
  { title: '创建时间', key: 'create_time', width: 180, render: (row: any) => formatDate(row.create_time) },
  { title: '操作', key: 'actions', width: 150, render: (row: any) => h(NSpace, () => [
    h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑'),
    h(NButton, { size: 'small', type: 'error', onClick: () => deleteAttach(row.id) }, () => '删除')
  ])}
]

async function loadAttaches() {
  loading.value = true
  try {
    const res = await attachApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    attaches.value = res.data.data.list || res.data.data
    pagination.value.itemCount = res.data.data.total || 0
  } finally {
    loading.value = false
  }
}

async function handleUpload(options: { file: UploadFileInfo; onFinish: () => void; onError: () => void }) {
  uploadLoading.value = true
  try {
    const rawFile = options.file.file
    if (!rawFile) {
      message.error('文件无效')
      options.onError()
      return
    }
    await attachApi.upload(rawFile)
    message.success('上传成功')
    options.onFinish()
    loadAttaches()
  } catch (e: any) {
    message.error(e.response?.data?.msg || '上传失败')
    options.onError()
  } finally {
    uploadLoading.value = false
  }
}

function openEdit(row?: any) {
  editingAttach.value = row ? { ...row } : { id: 0, name: '', link: '', path: '', status: 1, type: 1 }
  showModal.value = true
}

async function saveAttach() {
  if (editingAttach.value.id) {
    await attachApi.update(editingAttach.value.id, editingAttach.value)
  } else {
    await attachApi.create(editingAttach.value)
  }
  showModal.value = false
  loadAttaches()
}

async function deleteAttach(id: number) {
  await attachApi.delete(id)
  message.success('删除成功')
  loadAttaches()
}

async function initCdn() {
  try {
    const res = await websiteApi.get()
    cdnUrl.value = res.data.data?.cdn_url || ''
  } catch { /* ignore */ }
}

onMounted(() => {
  loadAttaches()
  initCdn()
})
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: space-between; align-items: center;">
      <n-upload
        accept="image/*"
        :show-file-list="false"
        :custom-request="handleUpload"
      >
        <n-button type="primary" :loading="uploadLoading">上传图片</n-button>
      </n-upload>
    </div>
    <n-data-table :columns="columns" :data="attaches" :loading="loading" remote :pagination="pagination" @update:page="pagination.page = $event; loadAttaches()" />
    <n-modal v-model:show="showModal" preset="card" title="附件管理" style="width: 500px">
      <n-form :model="editingAttach">
        <n-form-item label="名称">
          <n-input v-model:value="editingAttach.name" />
        </n-form-item>
        <n-form-item label="链接">
          <n-input v-model:value="editingAttach.link" />
        </n-form-item>
        <n-form-item label="路径">
          <n-input v-model:value="editingAttach.path" />
        </n-form-item>
        <n-form-item label="类型">
          <n-select v-model:value="editingAttach.type" :options="typeOptions" />
        </n-form-item>
        <n-form-item label="状态">
          <n-select v-model:value="editingAttach.status" :options="statusOptions" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="saveAttach">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>
