<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NInputNumber, NSwitch, NTag, NSpace, useMessage } from 'naive-ui'
import { categoryApi } from '../api'

const categories = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const pagination = ref({ page: 1, pageSize: 10, itemCount: 0 })
const editingCategory = ref({ id: 0, type: '', parent: '', mark: '', author: '', name: '', weight: 0, status: 1 })
const message = useMessage()

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '名称', key: 'name' },
  { title: '标识', key: 'mark' },
  { title: '类型', key: 'type' },
  { title: '权重', key: 'weight', width: 60 },
  { title: '状态', key: 'status', width: 60, render: (row: any) => h(NTag, { type: row.status === 1 ? 'success' : 'default', size: 'small' }, () => row.status === 1 ? '启用' : '禁用') },
  { title: '操作', key: 'actions', width: 150, render: (row: any) => h(NSpace, () => [
    h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑'),
    h(NButton, { size: 'small', type: 'error', onClick: () => deleteCategory(row.id) }, () => '删除')
  ])}
]

async function loadCategories() {
  loading.value = true
  try {
    const res = await categoryApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    categories.value = res.data.data.list || res.data.data
    pagination.value.itemCount = res.data.data.total || 0
  } finally {
    loading.value = false
  }
}

function openEdit(row?: any) {
  editingCategory.value = row ? { ...row } : { id: 0, type: '', parent: '', mark: '', author: '', name: '', weight: 0, status: 1 }
  showModal.value = true
}

async function saveCategory() {
  if (editingCategory.value.id) {
    await categoryApi.update(editingCategory.value.id, editingCategory.value)
  } else {
    await categoryApi.create(editingCategory.value)
  }
  showModal.value = false
  loadCategories()
}

async function deleteCategory(id: number) {
  await categoryApi.delete(id)
  message.success('删除成功')
  loadCategories()
}

onMounted(loadCategories)
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: flex-end;">
      <n-button type="primary" @click="openEdit()">新建分类</n-button>
    </div>
    <n-data-table 
      :columns="columns" 
      :data="categories" 
      :loading="loading"
      remote
      :pagination="pagination"
      @update:page="pagination.page = $event; loadCategories()"
    />
    <n-modal v-model:show="showModal" preset="card" title="分类管理" style="width: 500px">
      <n-form :model="editingCategory">
        <n-form-item label="名称">
          <n-input v-model:value="editingCategory.name" />
        </n-form-item>
        <n-form-item label="标识">
          <n-input v-model:value="editingCategory.mark" />
        </n-form-item>
        <n-form-item label="类型">
          <n-input v-model:value="editingCategory.type" />
        </n-form-item>
        <n-form-item label="父级">
          <n-input v-model:value="editingCategory.parent" />
        </n-form-item>
        <n-form-item label="权重">
          <n-input-number v-model:value="editingCategory.weight" :min="0" />
        </n-form-item>
        <n-form-item label="状态">
          <n-switch v-model:value="editingCategory.status" :checked-value="1" :unchecked-value="0" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="saveCategory">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>
