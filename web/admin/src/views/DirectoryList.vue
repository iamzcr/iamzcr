<script setup lang="ts">
import { ref, onMounted, h, computed } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NInputNumber, NSelect } from 'naive-ui'
import { directoryApi, categoryApi } from '../api'

const directories = ref<any[]>([])
const categories = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const pagination = ref({ page: 1, pageSize: 10, itemCount: 0 })
const editingDirectory = ref({ id: 0, cid: null as number | null, type: '', parent: '', mark: '', author: '', name: '', weight: 0, status: 1 })

const categoryMap = computed(() => {
  const map: Record<number, string> = {}
  categories.value.forEach(c => { map[c.id] = c.name })
  return map
})

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '名称', key: 'name' },
  { title: '分类', key: 'cid', width: 100, render: (row: any) => row.cid ? (categoryMap.value[row.cid] || row.cid) : '-' },
  { title: '标识', key: 'mark' },
  { title: '类型', key: 'type' },
  { title: '权重', key: 'weight', width: 60 },
  { title: '操作', key: 'actions', width: 100, render: (row: any) => h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑') }
]

async function loadDirectories() {
  loading.value = true
  try {
    const res = await directoryApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    directories.value = res.data.data.list || res.data.data
    pagination.value.itemCount = res.data.data.total || 0
  } finally {
    loading.value = false
  }
}

async function loadCategories() {
  const res = await categoryApi.list({ page: 1, page_size: 1000 })
  categories.value = res.data.data.list || res.data.data
}

function openEdit(row?: any) {
  editingDirectory.value = row ? { ...row } : { id: 0, cid: null, type: '', parent: '', mark: '', author: '', name: '', weight: 0, status: 1 }
  showModal.value = true
}

async function saveDirectory() {
  const data = { ...editingDirectory.value }
  if (data.cid === null) delete (data as any).cid
  if (editingDirectory.value.id) {
    await directoryApi.update(editingDirectory.value.id, data)
  } else {
    await directoryApi.create(data)
  }
  showModal.value = false
  loadDirectories()
}

onMounted(() => {
  loadDirectories()
  loadCategories()
})
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: flex-end;">
      <n-button type="primary" @click="openEdit()">新建目录</n-button>
    </div>
    <n-data-table :columns="columns" :data="directories" :loading="loading" remote :pagination="pagination" @update:page="pagination.page = $event; loadDirectories()" />
    <n-modal v-model:show="showModal" preset="card" title="目录管理" style="width: 500px">
      <n-form :model="editingDirectory">
        <n-form-item label="所属分类">
          <n-select v-model:value="editingDirectory.cid" :options="categories.map(c => ({ label: c.name, value: c.id }))" placeholder="选择分类（可选）" clearable />
        </n-form-item>
        <n-form-item label="名称">
          <n-input v-model:value="editingDirectory.name" />
        </n-form-item>
        <n-form-item label="标识">
          <n-input v-model:value="editingDirectory.mark" />
        </n-form-item>
        <n-form-item label="类型">
          <n-input v-model:value="editingDirectory.type" />
        </n-form-item>
        <n-form-item label="父级">
          <n-input v-model:value="editingDirectory.parent" />
        </n-form-item>
        <n-form-item label="权重">
          <n-input-number v-model:value="editingDirectory.weight" :min="0" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="saveDirectory">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>
