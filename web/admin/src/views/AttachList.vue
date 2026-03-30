<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSelect, NTag } from 'naive-ui'
import { attachApi } from '../api'

const attaches = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingAttach = ref({ id: 0, name: '', link: '', path: '', status: 1, type: 1 })

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
  { title: '名称', key: 'name' },
  { title: '链接', key: 'link', ellipsis: { tooltip: true } },
  { title: '路径', key: 'path', ellipsis: { tooltip: true } },
  { title: '类型', key: 'type', width: 80, render: (row: any) => h(NTag, { type: row.type === 1 ? 'warning' : 'info', size: 'small' }, () => row.type === 1 ? '图片' : '视频') },
  { title: '状态', key: 'status', width: 80, render: (row: any) => h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small' }, () => row.status === 1 ? '启用' : '禁用') },
  { title: '创建时间', key: 'create_time', width: 180, render: (row: any) => formatDate(row.create_time) },
  { title: '操作', key: 'actions', width: 120, render: (row: any) => h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑') }
]

async function loadAttaches() {
  loading.value = true
  try {
    const res = await attachApi.list()
    attaches.value = res.data.data
  } finally {
    loading.value = false
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

onMounted(loadAttaches)
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: flex-end;">
      <n-button type="primary" @click="openEdit()">新建附件</n-button>
    </div>
    <n-data-table :columns="columns" :data="attaches" :loading="loading" />
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
