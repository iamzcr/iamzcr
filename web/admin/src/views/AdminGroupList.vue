<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSwitch, NSpace, useMessage } from 'naive-ui'
import { adminGroupApi } from '../api'

const groups = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingGroup = ref({ id: 0, mark: '', name: '', description: '', menu_permit: '', menu_modules: '', allow_ip: '', status: 1 })
const message = useMessage()

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
  { title: '标识', key: 'mark', width: 100 },
  { title: '名称', key: 'name' },
  { title: '描述', key: 'description', ellipsis: { tooltip: true } },
  { title: '状态', key: 'status', width: 80, render: (row: any) => h('span', row.status === 1 ? '启用' : '禁用') },
  { title: '创建时间', key: 'create_time', width: 180, render: (row: any) => formatDate(row.create_time) },
  { title: '操作', key: 'actions', width: 150, render: (row: any) => h(NSpace, () => [h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑'), h(NButton, { size: 'small', type: 'error', onClick: () => deleteGroup(row.id) }, () => '删除')]) }
]

async function loadGroups() {
  loading.value = true
  try {
    const res = await adminGroupApi.list()
    groups.value = res.data.data
  } finally {
    loading.value = false
  }
}

function openEdit(row?: any) {
  editingGroup.value = row ? { ...row } : { id: 0, mark: '', name: '', description: '', menu_permit: '', menu_modules: '', allow_ip: '', status: 1 }
  showModal.value = true
}

async function saveGroup() {
  if (editingGroup.value.id) {
    await adminGroupApi.update(editingGroup.value.id, editingGroup.value)
  } else {
    await adminGroupApi.create(editingGroup.value)
  }
  showModal.value = false
  loadGroups()
}

async function deleteGroup(id: number) {
  await adminGroupApi.delete(id)
  message.success('删除成功')
  loadGroups()
}

onMounted(loadGroups)
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: flex-end;">
      <n-button type="primary" @click="openEdit()">新建用户组</n-button>
    </div>
    <n-data-table :columns="columns" :data="groups" :loading="loading" />
    <n-modal v-model:show="showModal" preset="card" title="用户组管理" style="width: 600px">
      <n-form :model="editingGroup">
        <n-form-item label="标识">
          <n-input v-model:value="editingGroup.mark" />
        </n-form-item>
        <n-form-item label="名称">
          <n-input v-model:value="editingGroup.name" />
        </n-form-item>
        <n-form-item label="描述">
          <n-input v-model:value="editingGroup.description" type="textarea" :rows="2" />
        </n-form-item>
        <n-form-item label="菜单权限">
          <n-input v-model:value="editingGroup.menu_permit" type="textarea" :rows="3" />
        </n-form-item>
        <n-form-item label="模块权限">
          <n-input v-model:value="editingGroup.menu_modules" type="textarea" :rows="3" />
        </n-form-item>
        <n-form-item label="允许IP">
          <n-input v-model:value="editingGroup.allow_ip" placeholder="多个用逗号分隔" />
        </n-form-item>
        <n-form-item label="状态">
          <n-switch v-model:value="editingGroup.status" :checked-value="1" :unchecked-value="0">
            <template #checked>启用</template>
            <template #unchecked>禁用</template>
          </n-switch>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="saveGroup">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>
