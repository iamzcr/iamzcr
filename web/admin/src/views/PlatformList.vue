<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSwitch, NTag, NSpace, useMessage } from 'naive-ui'
import { platformApi } from '../api'

const platforms = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const pagination = ref({ page: 1, pageSize: 10, itemCount: 0 })
const editingPlatform = ref({ id: 0, mark: '', name: '', status: 1 })
const message = useMessage()

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '名称', key: 'name' },
  { title: '标识', key: 'mark' },
  { title: '状态', key: 'status', width: 80, render: (row: any) => h(NTag, { type: row.status === 1 ? 'success' : 'default', size: 'small' }, () => row.status === 1 ? '启用' : '禁用') },
  { title: '操作', key: 'actions', width: 150, render: (row: any) => h(NSpace, () => [
    h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑'),
    h(NButton, { size: 'small', type: 'error', onClick: () => deletePlatform(row.id) }, () => '删除')
  ])}
]

async function loadPlatforms() {
  loading.value = true
  try {
    const res = await platformApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    platforms.value = res.data.data.list || res.data.data
    pagination.value.itemCount = res.data.data.total || 0
  } finally {
    loading.value = false
  }
}

function openEdit(row?: any) {
  editingPlatform.value = row ? { ...row } : { id: 0, mark: '', name: '', status: 1 }
  showModal.value = true
}

async function savePlatform() {
  if (editingPlatform.value.id) {
    await platformApi.update(editingPlatform.value.id, editingPlatform.value)
  } else {
    await platformApi.create(editingPlatform.value)
  }
  showModal.value = false
  loadPlatforms()
}

async function deletePlatform(id: number) {
  await platformApi.delete(id)
  message.success('删除成功')
  loadPlatforms()
}

onMounted(loadPlatforms)
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: flex-end;">
      <n-button type="primary" @click="openEdit()">新建平台</n-button>
    </div>
    <n-data-table 
      :columns="columns" 
      :data="platforms" 
      :loading="loading"
      remote
      :pagination="pagination"
      @update:page="pagination.page = $event; loadPlatforms()"
    />
    <n-modal v-model:show="showModal" preset="card" title="平台管理" style="width: 500px">
      <n-form :model="editingPlatform">
        <n-form-item label="名称">
          <n-input v-model:value="editingPlatform.name" />
        </n-form-item>
        <n-form-item label="标识">
          <n-input v-model:value="editingPlatform.mark" placeholder="如: wechat, bilibili" />
        </n-form-item>
        <n-form-item label="状态">
          <n-switch v-model:value="editingPlatform.status" :checked-value="1" :unchecked-value="0" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="savePlatform">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>
