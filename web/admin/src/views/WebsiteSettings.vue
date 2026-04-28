<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSpace, useMessage } from 'naive-ui'
import { websiteApi } from '../api'

const message = useMessage()
const websites = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingSetting = ref({ id: 0, key: '', value: '' })

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
  { title: 'Key', key: 'key', width: 200, ellipsis: { tooltip: true } },
  { title: 'Value', key: 'value', ellipsis: { tooltip: true } },
  { title: '创建时间', key: 'create_time', width: 180, render: (row: any) => formatDate(row.create_time) },
  { title: '操作', key: 'actions', width: 150, render: (row: any) => h(NSpace, () => [
    h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑'),
    h(NButton, { size: 'small', type: 'error', onClick: () => deleteSetting(row.id) }, () => '删除')
  ])}
]

async function loadSettings() {
  loading.value = true
  try {
    const res = await websiteApi.list()
    websites.value = res.data.data || []
  } finally {
    loading.value = false
  }
}

function openEdit(row?: any) {
  editingSetting.value = row ? { id: row.id, key: row.key, value: String(row.value) } : { id: 0, key: '', value: '' }
  showModal.value = true
}

async function saveSetting() {
  if (!editingSetting.value.key) {
    message.error('请输入Key')
    return
  }
  const data: Record<string, string> = { [editingSetting.value.key]: editingSetting.value.value }
  await websiteApi.update(data)
  message.success('保存成功')
  showModal.value = false
  loadSettings()
}

async function deleteSetting(id: number) {
  await websiteApi.delete(id)
  message.success('删除成功')
  loadSettings()
}

onMounted(loadSettings)
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: flex-end;">
      <n-button type="primary" @click="openEdit()">新建配置</n-button>
    </div>
    <n-data-table :columns="columns" :data="websites" :loading="loading" />
    <n-modal v-model:show="showModal" preset="card" title="网站设置" style="width: 500px">
      <n-form :model="editingSetting">
        <n-form-item label="Key">
          <n-input v-model:value="editingSetting.key" :disabled="!!editingSetting.id" placeholder="例如 cdn_url, site_title" />
        </n-form-item>
        <n-form-item label="Value">
          <n-input v-model:value="editingSetting.value" placeholder="请输入值" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="saveSetting">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>
