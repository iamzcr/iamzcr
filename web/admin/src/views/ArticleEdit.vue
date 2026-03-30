<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, NCard, NTabs, NTabPane, NSwitch, NInputNumber, NSelect, NSpace, NGrid, NGridItem } from 'naive-ui'
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
        <n-input v-model:value="form.thumb" placeholder="请输入封面图URL" />
      </n-form-item>
      
      <n-form-item label="属性">
        <n-space>
          <n-form-item label="热门" path="is_hot" label-placement="left" :show-label="false">
            <n-switch v-model:value="form.is_hot" :checked-value="1" :unchecked-value="0">
              <template #checked>是</template>
              <template #unchecked>否</template>
            </n-switch>
          </n-form-item>
          <n-form-item label="最新" path="is_new" label-placement="left" :show-label="false">
            <n-switch v-model:value="form.is_new" :checked-value="1" :unchecked-value="0">
              <template #checked>是</template>
              <template #unchecked>否</template>
            </n-switch>
          </n-form-item>
          <n-form-item label="推荐" path="is_recom" label-placement="left" :show-label="false">
            <n-switch v-model:value="form.is_recom" :checked-value="1" :unchecked-value="0">
              <template #checked>是</template>
              <template #unchecked>否</template>
            </n-switch>
          </n-form-item>
          <n-form-item label="状态" path="status" label-placement="left" :show-label="false">
            <n-switch v-model:value="form.status" :checked-value="1" :unchecked-value="0">
              <template #checked>已发布</template>
              <template #unchecked>草稿</template>
            </n-switch>
          </n-form-item>
          <n-form-item label="权重" path="weight" label-placement="left" :show-label="false">
            <n-input-number v-model:value="form.weight" :min="0" :show-button="false" placeholder="权重" style="width: 80px" />
          </n-form-item>
        </n-space>
      </n-form-item>
      
      <n-form-item label="内容" path="content">
        <n-tabs v-model:value="activeTab" type="line">
          <n-tab-pane name="write" tab="编辑">
            <n-input v-model:value="form.content" type="textarea" :rows="20" placeholder="请输入Markdown内容" style="text-align: left;" />
          </n-tab-pane>
          <n-tab-pane name="preview" tab="预览">
            <div class="markdown-preview" v-html="renderedContent"></div>
          </n-tab-pane>
          <n-tab-pane name="split" tab="分屏">
            <div class="split-view">
              <div class="split-editor">
                <n-input v-model:value="form.content" type="textarea" :rows="20" placeholder="请输入Markdown内容" />
              </div>
              <div class="split-preview">
                <div class="markdown-preview" v-html="renderedContent"></div>
              </div>
            </div>
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
.form-container {
  max-width: 800px;
}
.split-view {
  display: flex;
  gap: 16px;
  height: 500px;
}
.split-editor {
  flex: 1;
  min-width: 0;
}
.split-editor :deep(.n-input__textarea-el) {
  text-align: left;
}
.split-preview {
  flex: 1;
  min-width: 0;
  overflow: auto;
}
.markdown-preview {
  padding: 16px;
  border: 1px solid #eee;
  min-height: 400px;
  background: #fafafa;
  text-align: left;
}
.markdown-preview :deep(h1) { font-size: 2em; border-bottom: 1px solid #eee; padding-bottom: 0.3em; }
.markdown-preview :deep(h2) { font-size: 1.5em; border-bottom: 1px solid #eee; padding-bottom: 0.3em; }
.markdown-preview :deep(pre) { background: #f6f8fa; padding: 16px; border-radius: 6px; overflow-x: auto; }
.markdown-preview :deep(code) { background: #f6f8fa; padding: 2px 6px; border-radius: 3px; }
.markdown-preview :deep(blockquote) { border-left: 4px solid #ddd; padding-left: 16px; color: #666; }
.markdown-preview :deep(img) { max-width: 100; }
</style>
