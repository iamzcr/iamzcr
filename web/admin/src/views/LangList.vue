<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NInputNumber, NSwitch } from 'naive-ui'
import { langApi } from '../api'

const langs = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const pagination = ref({ page: 1, pageSize: 10, itemCount: 0 })
const editingLang = ref({ id: 0, name: '', lang: '', default: 0, status: 1, author: '', weight: 0 })

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
  { title: '语言代码', key: 'lang', width: 100 },
  { title: '默认', key: 'default', width: 80, render: (row: any) => h('span', row.default === 1 ? '是' : '否') },
  { title: '权重', key: 'weight', width: 80 },
  { title: '状态', key: 'status', width: 80, render: (row: any) => h('span', row.status === 1 ? '启用' : '禁用') },
  { title: '创建时间', key: 'create_time', width: 180, render: (row: any) => formatDate(row.create_time) },
  { title: '操作', key: 'actions', width: 120, render: (row: any) => h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑') }
]

async function loadLangs() {
  loading.value = true
  try {
    const res = await langApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    langs.value = res.data.data.list || res.data.data
    pagination.value.itemCount = res.data.data.total || 0
  } finally {
    loading.value = false
  }
}

function openEdit(row?: any) {
  editingLang.value = row ? { ...row } : { id: 0, name: '', lang: '', default: 0, status: 1, author: '', weight: 0 }
  showModal.value = true
}

async function saveLang() {
  if (editingLang.value.id) {
    await langApi.update(editingLang.value.id, editingLang.value)
  } else {
    await langApi.create(editingLang.value)
  }
  showModal.value = false
  loadLangs()
}

onMounted(loadLangs)
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: flex-end;">
      <n-button type="primary" @click="openEdit()">新建语言</n-button>
    </div>
    <n-data-table :columns="columns" :data="langs" :loading="loading" remote :pagination="pagination" @update:page="pagination.page = $event; loadLangs()" />
    <n-modal v-model:show="showModal" preset="card" title="语言管理" style="width: 500px">
      <n-form :model="editingLang">
        <n-form-item label="名称">
          <n-input v-model:value="editingLang.name" placeholder="如: 简体中文" />
        </n-form-item>
        <n-form-item label="语言代码">
          <n-input v-model:value="editingLang.lang" placeholder="如: zh, en" />
        </n-form-item>
        <n-form-item label="默认">
          <n-switch v-model:value="editingLang.default" :checked-value="1" :unchecked-value="0">
            <template #checked>是</template>
            <template #unchecked>否</template>
          </n-switch>
        </n-form-item>
        <n-form-item label="权重">
          <n-input-number v-model:value="editingLang.weight" :min="0" />
        </n-form-item>
        <n-form-item label="状态">
          <n-switch v-model:value="editingLang.status" :checked-value="1" :unchecked-value="0">
            <template #checked>启用</template>
            <template #unchecked>禁用</template>
          </n-switch>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="saveLang">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>
