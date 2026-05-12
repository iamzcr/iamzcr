<script setup lang="ts">
import { ref, onMounted, h, computed } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSelect, NTag, NUpload, NImage, NSpace, useMessage, NCheckbox } from 'naive-ui'
import type { UploadFileInfo } from 'naive-ui'
import { attachApi, attachMediaApi, websiteApi, platformApi } from '../api'

const message = useMessage()
const attaches = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const uploadLoading = ref(false)
const pagination = ref({ page: 1, pageSize: 10, itemCount: 0 })
const editingAttach = ref({ id: 0, name: '', link: '', path: '', status: 1, type: 1, is_thumb: 0 })
const cdnUrl = ref('')
const platforms = ref<any[]>([])
const showSyncModal = ref(false)
const syncPlatformId = ref<number | null>(null)
const syncAttachId = ref(0)

const showDeleteModal = ref(false)
const deleteAttachId = ref(0)
const deleteMediaRecords = ref<any[]>([])
const shouldDeleteMedia = ref(false)

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

const platformOptions = computed(() => platforms.value.map((p: any) => ({ label: p.name, value: p.id })))

function getPlatformName(id: number) {
  const p = platforms.value.find((p: any) => p.id === id)
  return p ? p.name : `ID:${id}`
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
  { title: '封面', key: 'is_thumb', width: 70, render: (row: any) => h(NTag, { type: row.is_thumb === 1 ? 'success' : 'default', size: 'small' }, () => row.is_thumb === 1 ? '是' : '否') },
  { title: '状态', key: 'status', width: 80, render: (row: any) => h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small' }, () => row.status === 1 ? '启用' : '禁用') },
  { title: '创建时间', key: 'create_time', width: 180, render: (row: any) => formatDate(row.create_time) },
  { title: '操作', key: 'actions', width: 250, render: (row: any) => h(NSpace, () => [
    h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑'),
    h(NButton, { size: 'small', type: 'info', onClick: () => openSyncModal(row.id) }, () => '同步到媒体'),
    h(NButton, { size: 'small', type: 'error', onClick: () => openDeleteModal(row.id) }, () => '删除')
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

async function loadPlatforms() {
  try {
    const res = await platformApi.list({ page: 1, page_size: 100 })
    platforms.value = (res.data.data.list || []).filter((p: any) => p.status === 1)
  } catch { /* ignore */ }
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
  editingAttach.value = row ? { ...row } : { id: 0, name: '', link: '', path: '', status: 1, type: 1, is_thumb: 0 }
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

function openSyncModal(id: number) {
  syncAttachId.value = id
  syncPlatformId.value = platforms.value.length > 0 ? platforms.value[0].id : null
  showSyncModal.value = true
}

async function syncNow() {
  if (!syncPlatformId.value) return
  showSyncModal.value = false
  try {
    await attachMediaApi.syncToMedia(syncAttachId.value, syncPlatformId.value)
    message.success('同步成功')
  } catch (e: any) {
    message.error(e.response?.data?.message || '同步失败')
  }
}

async function openDeleteModal(id: number) {
  deleteAttachId.value = id
  shouldDeleteMedia.value = false
  try {
    const res = await attachApi.getMedia(id)
    deleteMediaRecords.value = res.data.data || []
  } catch {
    deleteMediaRecords.value = []
  }
  showDeleteModal.value = true
}

async function confirmDelete() {
  showDeleteModal.value = false
  try {
    await attachApi.delete(deleteAttachId.value, shouldDeleteMedia.value ? { delete_media: 'true' } : undefined)
    message.success('删除成功')
    loadAttaches()
  } catch (e: any) {
    message.error(e.response?.data?.msg || '删除失败')
  }
}

async function initCdn() {
  try {
    const res = await websiteApi.get()
    cdnUrl.value = res.data.data?.cdn_url || ''
  } catch { /* ignore */ }
}

onMounted(() => {
  loadAttaches()
  loadPlatforms()
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
        <n-form-item label="封面图">
          <n-select v-model:value="editingAttach.is_thumb" :options="[{ label: '是', value: 1 }, { label: '否', value: 0 }]" />
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

    <n-modal v-model:show="showSyncModal" preset="card" title="同步到媒体平台" style="width: 400px">
      <n-form>
        <n-form-item label="选择平台">
          <n-select v-model:value="syncPlatformId" :options="platformOptions" placeholder="选择平台" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showSyncModal = false">取消</n-button>
        <n-button type="primary" @click="syncNow">开始同步</n-button>
      </template>
    </n-modal>

    <n-modal v-model:show="showDeleteModal" preset="card" title="确认删除" style="width: 450px">
      <div>
        <p style="margin-bottom: 12px;">确定要删除此附件吗？</p>
        <div v-if="deleteMediaRecords.length > 0" style="margin-bottom: 16px; padding: 12px; background: #fef3c7; border-radius: 4px;">
          <p style="margin-bottom: 8px; font-weight: bold;">已同步到以下平台：</p>
          <div v-for="r in deleteMediaRecords" :key="r.id" style="font-size: 13px; margin-bottom: 4px;">
            {{ getPlatformName(r.platform_id) }} — 状态: {{ r.status === 1 ? '已同步' : r.status === 2 ? '失败' : '待处理' }}
          </div>
          <n-checkbox v-model:checked="shouldDeleteMedia" style="margin-top: 8px;">
            同时从媒体平台删除已同步的附件
          </n-checkbox>
        </div>
        <p v-else style="color: #999;">此附件尚未同步到任何平台</p>
      </div>
      <template #footer>
        <n-button @click="showDeleteModal = false">取消</n-button>
        <n-button type="error" @click="confirmDelete">确认删除</n-button>
      </template>
    </n-modal>
  </div>
</template>