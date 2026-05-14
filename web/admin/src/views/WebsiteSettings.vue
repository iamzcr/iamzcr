<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSpace, NSwitch, NTag, useMessage } from 'naive-ui'
import { websiteApi } from '../api'

const message = useMessage()
const websites = ref<any[]>([])
const loading = ref(false)
const showModal = ref(false)
const editingSetting = ref({ id: 0, key: '', value: '', is_to_frontend: 1 })

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
  { title: '前端可见', key: 'is_to_frontend', width: 90, render: (row: any) => h(NTag, { type: row.is_to_frontend === 1 ? 'success' : 'default', size: 'small' }, () => row.is_to_frontend === 1 ? '是' : '否') },
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
  editingSetting.value = row ? { id: row.id, key: row.key, value: String(row.value), is_to_frontend: row.is_to_frontend ?? 1 } : { id: 0, key: '', value: '', is_to_frontend: 1 }
  showModal.value = true
}

function getKeyPlaceholder(key: string) {
  const hints: Record<string, string> = {
    wechat_app_id: '微信公众号 AppID',
    wechat_app_secret: '微信公众号 AppSecret',
    wechat_token: '微信公众号 Token (服务器配置)',
    wechat_aes_key: '微信公众号 EncodingAESKey',
    wechat_original_id: '微信公众号原始ID',
    cdn_url: 'CDN/资源访问域名',
    site_title: '网站标题'
  }
  return hints[key] || '例如 cdn_url, site_title, wechat_app_id'
}

function getValuePlaceholder(key: string) {
  const hints: Record<string, string> = {
    wechat_app_id: 'wx... 开头的 AppID',
    wechat_app_secret: 'AppSecret 密钥',
    wechat_token: 'Token 需与微信后台配置一致',
    wechat_aes_key: '43位随机字符串',
    wechat_original_id: 'gh_ 开头的原始ID',
    cdn_url: 'https://cdn.example.com',
    site_title: '堆栈人生'
  }
  return hints[key] || '请输入值'
}

async function saveSetting() {
  if (!editingSetting.value.key) {
    message.error('请输入Key')
    return
  }
  const data = [{ key: editingSetting.value.key, value: editingSetting.value.value, is_to_frontend: editingSetting.value.is_to_frontend }]
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
  <div class="page-wrap">
    <div class="page-toolbar" style="justify-content: space-between; align-items: center;">
      <div style="color: #666; font-size: 13px;">
        微信公众号配置 Key: wechat_app_id, wechat_app_secret, wechat_token, wechat_aes_key, wechat_original_id
      </div>
      <n-button type="primary" @click="openEdit()">新建配置</n-button>
    </div>
    <n-data-table :columns="columns" :data="websites" :loading="loading" />
    <n-modal v-model:show="showModal" preset="card" title="网站设置" style="width: 500px">
      <n-form :model="editingSetting">
        <n-form-item label="Key">
          <n-input v-model:value="editingSetting.key" :disabled="!!editingSetting.id" :placeholder="getKeyPlaceholder(editingSetting.key)" />
        </n-form-item>
        <n-form-item label="Value">
          <n-input v-model:value="editingSetting.value" :placeholder="getValuePlaceholder(editingSetting.key)" />
        </n-form-item>
        <n-form-item label="前端可见">
          <n-switch v-model:value="editingSetting.is_to_frontend" :checked-value="1" :unchecked-value="0">
            <template #checked>是</template>
            <template #unchecked>否</template>
          </n-switch>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="saveSetting">保存</n-button>
      </template>
    </n-modal>
  </div>
</template>
