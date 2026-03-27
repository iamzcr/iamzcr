<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, NCard, NTabs, NTabPane, NSwitch, NInputNumber, NSelect } from 'naive-ui'
import { articleApi, categoryApi, directoryApi, tagsApi } from '../api'
import MarkdownIt from 'markdown-it'

const route = useRoute()
const router = useRouter()
const md = new MarkdownIt()

const isEdit = computed(() => !!route.params.id)
const form = ref({
  cid: 0,
  did: 0,
  title: '',
  desc: '',
  keyword: '',
  author: '',
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
const activeTab = ref('write')

const categories = ref<any[]>([])
const directories = ref<any[]>([])
const allTags = ref<any[]>([])

const renderedContent = computed(() => md.render(form.value.content))

async function loadData() {
  const [catRes, dirRes, tagRes] = await Promise.all([
    categoryApi.list(),
    directoryApi.list(),
    tagsApi.list()
  ])
  categories.value = catRes.data.data
  directories.value = dirRes.data.data
  allTags.value = tagRes.data.data
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
      author: data.article?.author || data.author || '',
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
    <n-form :model="form" label-placement="left">
      <n-form-item label="标题">
        <n-input v-model:value="form.title" placeholder="请输入标题" />
      </n-form-item>
      <n-form-item label="分类">
        <n-select v-model:value="form.cid" :options="categories.map(c => ({ label: c.name, value: c.id }))" placeholder="选择分类" clearable />
      </n-form-item>
      <n-form-item label="目录">
        <n-select v-model:value="form.did" :options="directories.filter(d => d.cid === form.cid || form.cid === 0).map(d => ({ label: d.name, value: d.id }))" placeholder="选择目录" clearable />
      </n-form-item>
      <n-form-item label="标签">
        <n-select v-model:value="form.tag_ids" :options="allTags.map(t => ({ label: t.name, value: t.id }))" multiple placeholder="选择标签" />
      </n-form-item>
      <n-form-item label="描述">
        <n-input v-model:value="form.desc" placeholder="请输入描述" />
      </n-form-item>
      <n-form-item label="关键词">
        <n-input v-model:value="form.keyword" placeholder="请输入关键词" />
      </n-form-item>
      <n-form-item label="作者">
        <n-input v-model:value="form.author" placeholder="请输入作者" />
      </n-form-item>
      <n-form-item label="封面图">
        <n-input v-model:value="form.thumb" placeholder="请输入封面图URL" />
      </n-form-item>
      <n-form-item label="摘要">
        <n-input v-model:value="form.summary" type="textarea" :rows="2" placeholder="请输入摘要" />
      </n-form-item>
      <n-form-item label="权重">
        <n-input-number v-model:value="form.weight" :min="0" />
      </n-form-item>
      <n-form-item label="热门">
        <n-switch v-model:value="form.is_hot" :checked-value="1" :unchecked-value="0" />
      </n-form-item>
      <n-form-item label="最新">
        <n-switch v-model:value="form.is_new" :checked-value="1" :unchecked-value="0" />
      </n-form-item>
      <n-form-item label="推荐">
        <n-switch v-model:value="form.is_recom" :checked-value="1" :unchecked-value="0" />
      </n-form-item>
      <n-form-item label="状态">
        <n-switch v-model:value="form.status" :checked-value="1" :unchecked-value="0">
          <template #checked>已发布</template>
          <template #unchecked>草稿</template>
        </n-switch>
      </n-form-item>
      <n-form-item label="内容">
        <n-tabs v-model:value="activeTab" type="line">
          <n-tab-pane name="write" tab="编辑">
            <n-input v-model:value="form.content" type="textarea" :rows="20" placeholder="请输入Markdown内容" />
          </n-tab-pane>
          <n-tab-pane name="preview" tab="预览">
            <div class="markdown-preview" v-html="renderedContent"></div>
          </n-tab-pane>
        </n-tabs>
      </n-form-item>
      <n-form-item>
        <n-space>
          <n-button type="primary" @click="save" :loading="loading">保存</n-button>
          <n-button @click="router.push('/articles')">取消</n-button>
        </n-space>
      </n-form-item>
    </n-form>
  </n-card>
</template>

<style scoped>
.markdown-preview {
  padding: 16px;
  border: 1px solid #eee;
  min-height: 400px;
  background: #fafafa;
}
.markdown-preview :deep(h1) { font-size: 2em; border-bottom: 1px solid #eee; padding-bottom: 0.3em; }
.markdown-preview :deep(h2) { font-size: 1.5em; border-bottom: 1px solid #eee; padding-bottom: 0.3em; }
.markdown-preview :deep(pre) { background: #f6f8fa; padding: 16px; border-radius: 6px; overflow-x: auto; }
.markdown-preview :deep(code) { background: #f6f8fa; padding: 2px 6px; border-radius: 3px; }
.markdown-preview :deep(blockquote) { border-left: 4px solid #ddd; padding-left: 16px; color: #666; }
.markdown-preview :deep(img) { max-width: 100%; }
</style>