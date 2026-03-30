<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NInputNumber } from 'naive-ui'
import { menuApi } from '../api'

const menus = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const pagination = ref({ page: 1, pageSize: 10, total: 0 })
const editingMenu = ref({ id: 0, type: 0, mark: '', author: '', name: '', url: '', parent: 0, icon: '', weight: 0, status: 1 })

const columns = [
  { title: 'ID', key: 'id', width: 60 },
  { title: '名称', key: 'name' },
  { title: '链接', key: 'url', ellipsis: { tooltip: true } },
  { title: '类型', key: 'type', width: 60 },
  { title: '权重', key: 'weight', width: 60 },
  { title: '操作', key: 'actions', width: 100, render: (row: any) => h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑') }
]

async function loadMenus() {
  loading.value = true
  try {
    const res = await menuApi.list({ page: pagination.value.page, page_size: pagination.value.pageSize })
    menus.value = res.data.data.list || res.data.data
    pagination.value.total = res.data.data.total || 0
  } finally {
    loading.value = false
  }
}

function openEdit(row?: any) {
  editingMenu.value = row ? { ...row } : { id: 0, type: 0, mark: '', author: '', name: '', url: '', parent: 0, icon: '', weight: 0, status: 1 }
  showModal.value = true
}

async function saveMenu() {
  if (editingMenu.value.id) {
    await menuApi.update(editingMenu.value.id, editingMenu.value)
  } else {
    await menuApi.create(editingMenu.value)
  }
  showModal.value = false
  loadMenus()
}

onMounted(loadMenus)
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: flex-end;">
      <n-button type="primary" @click="openEdit()">新建菜单</n-button>
    </div>
    <n-data-table :columns="columns" :data="menus" :loading="loading" :pagination="pagination" @update:page="pagination.page = $event; loadMenus()" />
    <n-modal v-model:show="showModal" preset="card" title="菜单管理" style="width: 500px">
      <n-form :model="editingMenu">
        <n-form-item label="名称">
          <n-input v-model:value="editingMenu.name" />
        </n-form-item>
        <n-form-item label="链接">
          <n-input v-model:value="editingMenu.url" />
        </n-form-item>
        <n-form-item label="标识">
          <n-input v-model:value="editingMenu.mark" />
        </n-form-item>
        <n-form-item label="图标">
          <n-input v-model:value="editingMenu.icon" />
        </n-form-item>
        <n-form-item label="父级">
          <n-input-number v-model:value="editingMenu.parent" :min="0" />
        </n-form-item>
        <n-form-item label="类型">
          <n-input-number v-model:value="editingMenu.type" :min="0" />
        </n-form-item>
        <n-form-item label="权重">
          <n-input-number v-model:value="editingMenu.weight" :min="0" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="saveMenu">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>
