<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NInputNumber, NSelect } from 'naive-ui'
import { permitApi } from '../api'

const permits = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const pagination = ref({ page: 1, pageSize: 10, total: 0 })
const editingPermit = ref({ id: 0, type: 1, mark: '', parent: 0, name: '', modules: '', author: '', status: 1, weight: 0 })

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
  { title: '标识', key: 'mark', width: 120 },
  { title: '名称', key: 'name' },
  { title: '模块', key: 'modules', ellipsis: { tooltip: true } },
  { title: '权重', key: 'weight', width: 80 },
  { title: '状态', key: 'status', width: 80, render: (row: any) => h('span', row.status === 1 ? '启用' : '禁用') },
  { title: '创建时间', key: 'create_time', width: 180, render: (row: any) => formatDate(row.create_time) },
  { title: '操作', key: 'actions', width: 120, render: (row: any) => h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑') }
]

async function loadPermits() {
  loading.value = true
  try {
    const res = await permitApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    permits.value = res.data.data.list || res.data.data
    pagination.value.total = res.data.data.total || 0
  } finally {
    loading.value = false
  }
}

function openEdit(row?: any) {
  editingPermit.value = row ? { ...row } : { id: 0, type: 1, mark: '', parent: 0, name: '', modules: '', author: '', status: 1, weight: 0 }
  showModal.value = true
}

async function savePermit() {
  if (editingPermit.value.id) {
    await permitApi.update(editingPermit.value.id, editingPermit.value)
  } else {
    await permitApi.create(editingPermit.value)
  }
  showModal.value = false
  loadPermits()
}

onMounted(loadPermits)
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: flex-end;">
      <n-button type="primary" @click="openEdit()">新建权限</n-button>
    </div>
    <n-data-table :columns="columns" :data="permits" :loading="loading" :pagination="pagination" @update:page="pagination.page = $event; loadPermits()" />
    <n-modal v-model:show="showModal" preset="card" title="权限管理" style="width: 500px">
      <n-form :model="editingPermit">
        <n-form-item label="标识">
          <n-input v-model:value="editingPermit.mark" />
        </n-form-item>
        <n-form-item label="名称">
          <n-input v-model:value="editingPermit.name" />
        </n-form-item>
        <n-form-item label="模块">
          <n-input v-model:value="editingPermit.modules" type="textarea" :rows="3" />
        </n-form-item>
        <n-form-item label="父级">
          <n-input-number v-model:value="editingPermit.parent" :min="0" />
        </n-form-item>
        <n-form-item label="类型">
          <n-input-number v-model:value="editingPermit.type" :min="1" />
        </n-form-item>
        <n-form-item label="权重">
          <n-input-number v-model:value="editingPermit.weight" :min="0" />
        </n-form-item>
        <n-form-item label="状态">
          <n-select v-model:value="editingPermit.status" :options="[{label: '启用', value: 1}, {label: '禁用', value: 0}]" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="savePermit">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>
