<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NInputNumber, NSwitch, NTag } from 'naive-ui'
import { tagsApi } from '../api'

const tags = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const pagination = ref({ page: 1, pageSize: 10, total: 0 })
const editingTag = ref({ id: 0, type: '', mark: '', author: '', name: '', weight: 0, status: 1, is_hot: 0 })

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '名称', key: 'name' },
  { title: '标识', key: 'mark' },
  { title: '类型', key: 'type' },
  { title: '权重', key: 'weight', width: 60 },
  { title: '热门', key: 'is_hot', width: 60, render: (row: any) => h(NTag, { type: row.is_hot === 1 ? 'error' : 'default', size: 'small' }, () => row.is_hot === 1 ? '是' : '否') },
  { title: '操作', key: 'actions', width: 100, render: (row: any) => h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑') }
]

async function loadTags() {
  loading.value = true
  try {
    const res = await tagsApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    tags.value = res.data.data.list || res.data.data
    pagination.value.total = res.data.data.total || 0
  } finally {
    loading.value = false
  }
}

function openEdit(row?: any) {
  editingTag.value = row ? { ...row } : { id: 0, type: '', mark: '', author: '', name: '', weight: 0, status: 1, is_hot: 0 }
  showModal.value = true
}

async function saveTag() {
  if (editingTag.value.id) {
    await tagsApi.update(editingTag.value.id, editingTag.value)
  } else {
    await tagsApi.create(editingTag.value)
  }
  showModal.value = false
  loadTags()
}

onMounted(loadTags)
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: flex-end;">
      <n-button type="primary" @click="openEdit()">新建标签</n-button>
    </div>
    <n-data-table :columns="columns" :data="tags" :loading="loading" :pagination="pagination" @update:page="pagination.page = $event; loadTags()" />
    <n-modal v-model:show="showModal" preset="card" title="标签管理" style="width: 500px">
      <n-form :model="editingTag">
        <n-form-item label="名称">
          <n-input v-model:value="editingTag.name" />
        </n-form-item>
        <n-form-item label="标识">
          <n-input v-model:value="editingTag.mark" />
        </n-form-item>
        <n-form-item label="类型">
          <n-input v-model:value="editingTag.type" />
        </n-form-item>
        <n-form-item label="权重">
          <n-input-number v-model:value="editingTag.weight" :min="0" />
        </n-form-item>
        <n-form-item label="热门">
          <n-switch v-model:value="editingTag.is_hot" :checked-value="1" :unchecked-value="0" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="saveTag">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>
