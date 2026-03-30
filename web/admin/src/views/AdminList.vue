<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSelect, NSwitch } from 'naive-ui'
import { adminApi, adminGroupApi } from '../api'

const admins = ref<any[]>([])
const groups = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingAdmin = ref({ id: 0, username: '', password: '', name: '', group_id: 0, status: 1 })

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
  { title: '用户名', key: 'username', width: 120 },
  { title: '姓名', key: 'name', width: 100 },
  { title: '用户组', key: 'group_id', width: 100, render: (row: any) => {
    const g = groups.value.find((x: any) => x.id === row.group_id)
    return g ? g.name : '-'
  }},
  { title: '状态', key: 'status', width: 80, render: (row: any) => h('span', row.status === 1 ? '正常' : '禁用') },
  { title: '登录次数', key: 'login_num', width: 80 },
  { title: '最后登录', key: 'last_login_time', width: 180, render: (row: any) => formatDate(row.last_login_time) },
  { title: '操作', key: 'actions', width: 120, render: (row: any) => h(NButton, { size: 'small', onClick: () => openEdit(row) }, () => '编辑') }
]

async function loadAdmins() {
  loading.value = true
  try {
    const res = await adminApi.list()
    admins.value = res.data.data
  } finally {
    loading.value = false
  }
}

async function loadGroups() {
  const res = await adminGroupApi.list()
  groups.value = res.data.data
}

function openEdit(row?: any) {
  editingAdmin.value = row ? { ...row, password: '' } : { id: 0, username: '', password: '', name: '', group_id: 0, status: 1 }
  showModal.value = true
}

async function saveAdmin() {
  if (editingAdmin.value.id) {
    await adminApi.update(editingAdmin.value.id, editingAdmin.value)
  } else {
    await adminApi.create(editingAdmin.value)
  }
  showModal.value = false
  loadAdmins()
}

onMounted(() => {
  loadAdmins()
  loadGroups()
})
</script>

<template>
  <div>
    <div style="margin-bottom: 16px; display: flex; justify-content: flex-end;">
      <n-button type="primary" @click="openEdit()">新建管理员</n-button>
    </div>
    <n-data-table :columns="columns" :data="admins" :loading="loading" />
    <n-modal v-model:show="showModal" preset="card" title="管理员管理" style="width: 500px">
      <n-form :model="editingAdmin">
        <n-form-item label="用户名">
          <n-input v-model:value="editingAdmin.username" :disabled="!!editingAdmin.id" />
        </n-form-item>
        <n-form-item v-if="!editingAdmin.id" label="密码">
          <n-input v-model:value="editingAdmin.password" type="password" />
        </n-form-item>
        <n-form-item label="姓名">
          <n-input v-model:value="editingAdmin.name" />
        </n-form-item>
        <n-form-item label="用户组">
          <n-select v-model:value="editingAdmin.group_id" :options="groups.map((g: any) => ({ label: g.name, value: g.id }))" />
        </n-form-item>
        <n-form-item label="状态">
          <n-switch v-model:value="editingAdmin.status" :checked-value="1" :unchecked-value="0">
            <template #checked>启用</template>
            <template #unchecked>禁用</template>
          </n-switch>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="saveAdmin">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>
