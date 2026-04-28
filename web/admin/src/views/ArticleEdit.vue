<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, NCard, NSwitch, NInputNumber, NSelect, NSpace, NGrid, NGridItem, NModal, NImage } from 'naive-ui'
import { articleApi, categoryApi, directoryApi, tagsApi, attachApi, websiteApi } from '../api'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const form = ref({
  cid: 0,
  did: 0,
  title: '',
  desc: '',
  keyword: '',
  author: 'nicholas',
  thumb: '',
  summary: '',
  content: '',
  is_hot: 0,
  is_new: 0,
  is_recom: 0,
  weight: 0,
  public_time: 0,
  status: 1,
  month: '',
  tag_ids: [] as number[]
})
const loading = ref(false)

const categories = ref<any[]>([])
const directories = ref<any[]>([])
const allTags = ref<any[]>([])

const showCoverModal = ref(false)
const coverImages = ref<any[]>([])
const coverLoading = ref(false)
const cdnUrl = ref('')

function getFullUrl(path: string) {
  if (!path) return ''
  if (path.startsWith('http://') || path.startsWith('https://')) return path
  if (!cdnUrl.value) return path
  return cdnUrl.value.replace(/\/+$/, '') + '/' + path.replace(/^\/+/, '')
}

async function loadData() {
  const [catRes, dirRes, tagRes, webRes] = await Promise.all([
    categoryApi.list({ page: 1, page_size: 1000 }),
    directoryApi.list({ page: 1, page_size: 1000 }),
    tagsApi.list({ page: 1, page_size: 1000 }),
    websiteApi.get()
  ])
  categories.value = catRes.data.data.list || catRes.data.data || []
  directories.value = dirRes.data.data.list || dirRes.data.data || []
  allTags.value = tagRes.data.data.list || tagRes.data.data || []
  cdnUrl.value = webRes.data.data?.cdn_url || ''
}

async function loadArticle() {
  if (route.params.id) {
    const res = await articleApi.get(Number(route.params.id))
    const data = res.data.data
    form.value = {
      cid: data.article?.cid || data.cid || 0,
      did: data.article?.did || data.did || 0,
      title: data.article?.title || data.title || '',
      desc: data.article?.desc || data.desc || '',
      keyword: data.article?.keyword || data.keyword || '',
      author: data.article?.author || data.author || 'nicholas',
      thumb: data.article?.thumb || data.thumb || '',
      summary: data.article?.summary || data.summary || '',
      content: data.article?.content || data.content || '',
      is_hot: data.article?.is_hot || data.is_hot || 0,
      is_new: data.article?.is_new || data.is_new || 0,
      is_recom: data.article?.is_recom || data.is_recom || 0,
      weight: data.article?.weight || data.weight || 0,
      public_time: data.article?.public_time || data.public_time || 0,
      status: data.article?.status || data.status || 1,
      month: data.article?.month || data.month || '',
      tag_ids: data.tags ? data.tags.map((t: any) => t.id) : []
    }
  }
}

async function openCoverModal() {
  showCoverModal.value = true
  coverLoading.value = true
  try {
    const res = await attachApi.list({ page: 1, page_size: 1000, type: 1 })
    coverImages.value = res.data.data.list || res.data.data || []
  } finally {
    coverLoading.value = false
  }
}

function selectCover(img: any) {
  form.value.thumb = img.link
  showCoverModal.value = false
}

function clearCover() {
  form.value.thumb = ''
}

async function save() {
  loading.value = true
  try {
    if (isEdit.value) {
      await articleApi.update(Number(route.params.id), form.value)
    } else {
      await articleApi.create(form.value)
    }
    router.push('/articles')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
  loadArticle()
})
</script>

<template>
  <n-card>
    <n-form :model="form" label-placement="top" class="form-container">
      <div class="form-meta">
      <n-form-item label="标题" path="title">
        <n-input v-model:value="form.title" placeholder="请输入标题" />
      </n-form-item>
      
      <n-grid :cols="2" :x-gap="16">
        <n-grid-item>
          <n-form-item label="分类" path="cid">
            <n-select v-model:value="form.cid" :options="categories.map(c => ({ label: c.name, value: c.id }))" placeholder="选择分类" clearable />
          </n-form-item>
        </n-grid-item>
        <n-grid-item>
          <n-form-item label="目录" path="did">
            <n-select v-model:value="form.did" :options="directories.filter(d => d.cid === form.cid || form.cid === 0).map(d => ({ label: d.name, value: d.id }))" placeholder="选择目录" clearable />
          </n-form-item>
        </n-grid-item>
      </n-grid>
      
      <n-form-item label="标签" path="tag_ids">
        <n-select v-model:value="form.tag_ids" :options="allTags.map(t => ({ label: t.name, value: t.id }))" multiple placeholder="选择标签" />
      </n-form-item>
      
      <n-form-item label="描述/摘要" path="desc">
        <n-input v-model:value="form.desc" type="textarea" :rows="3" placeholder="请输入描述（同时作为摘要）" />
      </n-form-item>
      
      <n-grid :cols="2" :x-gap="16">
        <n-grid-item>
          <n-form-item label="关键词" path="keyword">
            <n-input v-model:value="form.keyword" placeholder="请输入关键词" />
          </n-form-item>
        </n-grid-item>
        <n-grid-item>
          <n-form-item label="作者" path="author">
            <n-input v-model:value="form.author" placeholder="请输入作者" />
          </n-form-item>
        </n-grid-item>
      </n-grid>
      
      <n-form-item label="封面图" path="thumb">
        <div style="display: flex; gap: 12px; align-items: flex-end;">
          <div style="flex: 1;">
            <n-input v-model:value="form.thumb" placeholder="手动输入URL或点击右侧按钮选择" />
          </div>
          <n-button @click="openCoverModal">选择封面</n-button>
          <n-button v-if="form.thumb" @click="clearCover" secondary>清除</n-button>
        </div>
        <div v-if="form.thumb" style="margin-top: 8px;">
          <n-image width="200" :src="getFullUrl(form.thumb)" style="border-radius: 4px; border: 1px solid #eee;" />
        </div>
      </n-form-item>
      
      <n-form-item label="属性">
        <n-space>
          <n-form-item label="热门" path="is_hot" label-placement="left">
            <n-switch v-model:value="form.is_hot" :checked-value="1" :unchecked-value="0">
              <template #checked>是</template>
              <template #unchecked>否</template>
            </n-switch>
          </n-form-item>
          <n-form-item label="最新" path="is_new" label-placement="left">
            <n-switch v-model:value="form.is_new" :checked-value="1" :unchecked-value="0">
              <template #checked>是</template>
              <template #unchecked>否</template>
            </n-switch>
          </n-form-item>
          <n-form-item label="推荐" path="is_recom" label-placement="left">
            <n-switch v-model:value="form.is_recom" :checked-value="1" :unchecked-value="0">
              <template #checked>是</template>
              <template #unchecked>否</template>
            </n-switch>
          </n-form-item>
          <n-form-item label="状态" path="status" label-placement="left">
            <n-switch v-model:value="form.status" :checked-value="1" :unchecked-value="0">
              <template #checked>已发布</template>
              <template #unchecked>草稿</template>
            </n-switch>
          </n-form-item>
          <n-form-item label="权重" path="weight" label-placement="left">
            <n-input-number v-model:value="form.weight" :min="0" :show-button="false" placeholder="权重" style="width: 80px" />
          </n-form-item>
        </n-space>
      </n-form-item>
      </div>
      
      <n-form-item label="内容" path="content" class="content-form-item">
        <v-md-editor v-model="form.content" height="68vh"></v-md-editor>
      </n-form-item>
      
      <n-form-item>
        <n-space>
          <n-button type="primary" @click="save" :loading="loading">保存</n-button>
          <n-button @click="router.push('/articles')">取消</n-button>
        </n-space>
      </n-form-item>
    </n-form>

    <n-modal v-model:show="showCoverModal" preset="card" title="选择封面图" style="width: 720px">
      <div style="display: flex; flex-wrap: wrap; gap: 12px; max-height: 500px; overflow-y: auto;">
        <div
          v-for="img in coverImages"
          :key="img.id"
          style="width: 160px; cursor: pointer; border: 2px solid transparent; border-radius: 4px; padding: 4px;"
          :style="{ borderColor: form.thumb === img.link ? '#2080f0' : 'transparent' }"
          @click="selectCover(img)"
        >
          <n-image width="150" height="150" :src="getFullUrl(img.link)" style="object-fit: cover; border-radius: 4px;" />
          <div style="font-size: 12px; text-align: center; margin-top: 4px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ img.name }}</div>
        </div>
        <div v-if="coverImages.length === 0 && !coverLoading" style="padding: 24px; color: #999; text-align: center; width: 100%;">
          暂无上传图片，请先在附件管理中上传
        </div>
      </div>
    </n-modal>
  </n-card>
</template>

<style scoped>
.form-container {
  max-width: 800px;
}

.form-meta {
  max-width: 800px;
}

.content-form-item {
  width: 100%;
}

.content-form-item :deep(.n-form-item-blank) {
  display: block;
  width: 100%;
}
</style>
