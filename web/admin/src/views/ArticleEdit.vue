<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, NCard, NSwitch, NInputNumber, NSelect, NSpace, NGrid, NGridItem, NModal, NImage, NRadioGroup, NRadioButton } from 'naive-ui'
import { Editor } from '@tiptap/core'
import StarterKit from '@tiptap/starter-kit'
import Underline from '@tiptap/extension-underline'
import { articleApi, categoryApi, directoryApi, tagsApi, attachApi, websiteApi, platformApi } from '../api'
import { markdownToHtml, htmlToMarkdown } from '../utils/markdown'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const platforms = ref<any[]>([])
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
  tag_ids: [] as number[],
  publish_platform_ids: [] as number[]
})
const loading = ref(false)

const categories = ref<any[]>([])
const directories = ref<any[]>([])
const allTags = ref<any[]>([])

const showCoverModal = ref(false)
const coverImages = ref<any[]>([])
const coverLoading = ref(false)
const cdnUrl = ref('')

const editorMode = ref<'markdown' | 'richtext'>('markdown')
const richEditor = ref<Editor | null>(null)
const richEditorDom = ref<any>(null)

const platformOptions = computed(() => platforms.value.map((p: any) => ({ label: p.name, value: p.id })))

function initRichEditor() {
  if (richEditor.value) {
    richEditor.value.destroy()
    richEditor.value = null
  }
  nextTick(() => {
    richEditor.value = new Editor({
      element: richEditorDom.value,
      extensions: [StarterKit, Underline],
      content: markdownToHtml(form.value.content),
      onUpdate: ({ editor }) => {
        form.value.content = htmlToMarkdown(editor.getHTML())
      },
    })
  })
}

function destroyRichEditor() {
  if (richEditor.value) {
    richEditor.value.destroy()
    richEditor.value = null
  }
}

function switchEditorMode() {
  if (editorMode.value === 'richtext') {
    initRichEditor()
  } else {
    destroyRichEditor()
  }
}

watch(editorMode, () => {
  switchEditorMode()
})

function getFullUrl(path: string) {
  if (!path) return ''
  if (path.startsWith('http://') || path.startsWith('https://')) return path
  if (!cdnUrl.value) return path
  return cdnUrl.value.replace(/\/+$/, '') + '/' + path.replace(/^\/+/, '')
}

async function loadData() {
  const [catRes, dirRes, tagRes, webRes, platRes] = await Promise.all([
    categoryApi.list({ page: 1, page_size: 1000 }),
    directoryApi.list({ page: 1, page_size: 1000 }),
    tagsApi.list({ page: 1, page_size: 1000 }),
    websiteApi.get(),
    platformApi.list({ page: 1, page_size: 100 })
  ])
  categories.value = catRes.data.data.list || catRes.data.data || []
  directories.value = dirRes.data.data.list || dirRes.data.data || []
  allTags.value = tagRes.data.data.list || tagRes.data.data || []
  cdnUrl.value = webRes.data.data?.cdn_url || ''
  platforms.value = (platRes.data.data.list || []).filter((p: any) => p.status === 1)
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
      tag_ids: data.tags ? data.tags.map((t: any) => t.id) : [],
      publish_platform_ids: [] as number[]
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
  showCoverModal.value = false
  form.value.thumb = img.link
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

onBeforeUnmount(() => {
  destroyRichEditor()
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
      
      <n-form-item>
        <n-space align="center">
          <span class="prop-label">热门</span>
          <n-switch v-model:value="form.is_hot" :checked-value="1" :unchecked-value="0">
            <template #checked>是</template>
            <template #unchecked>否</template>
          </n-switch>
          <span class="prop-label">最新</span>
          <n-switch v-model:value="form.is_new" :checked-value="1" :unchecked-value="0">
            <template #checked>是</template>
            <template #unchecked>否</template>
          </n-switch>
          <span class="prop-label">推荐</span>
          <n-switch v-model:value="form.is_recom" :checked-value="1" :unchecked-value="0">
            <template #checked>是</template>
            <template #unchecked>否</template>
          </n-switch>
          <span class="prop-label">状态</span>
          <n-switch v-model:value="form.status" :checked-value="1" :unchecked-value="0">
            <template #checked>已发布</template>
            <template #unchecked>草稿</template>
          </n-switch>
          <span class="prop-label">权重</span>
          <n-input-number v-model:value="form.weight" :min="0" :show-button="false" placeholder="权重" style="width: 80px" />
        </n-space>
      </n-form-item>

      <n-form-item label="发布到平台" path="publish_platform_ids">
        <n-select v-model:value="form.publish_platform_ids" :options="platformOptions" multiple placeholder="选择要发布到的平台（可多选）" clearable />
      </n-form-item>
      </div>
      
      <n-form-item label="编辑模式" path="editorMode">
        <n-radio-group v-model:value="editorMode">
          <n-radio-button value="markdown">Markdown</n-radio-button>
          <n-radio-button value="richtext">富文本</n-radio-button>
        </n-radio-group>
      </n-form-item>

      <n-form-item label="内容" path="content" class="content-form-item">
        <v-md-editor v-if="editorMode === 'markdown'" v-model="form.content" height="68vh"></v-md-editor>
        <div v-else class="richtext-editor">
          <div class="richtext-toolbar">
            <button class="toolbar-btn" title="加粗" @click.prevent="richEditor?.chain().focus().toggleBold().run()" :class="{ active: richEditor?.isActive('bold') }"><b>B</b></button>
            <button class="toolbar-btn" title="斜体" @click.prevent="richEditor?.chain().focus().toggleItalic().run()" :class="{ active: richEditor?.isActive('italic') }"><i>I</i></button>
            <button class="toolbar-btn" title="下划线" @click.prevent="richEditor?.chain().focus().toggleUnderline().run()" :class="{ active: richEditor?.isActive('underline') }"><u>U</u></button>
            <button class="toolbar-btn" title="删除线" @click.prevent="richEditor?.chain().focus().toggleStrike().run()" :class="{ active: richEditor?.isActive('strike') }"><s>S</s></button>
            <span class="toolbar-divider"></span>
            <button class="toolbar-btn" title="标题1" @click.prevent="richEditor?.chain().focus().toggleHeading({ level: 1 }).run()" :class="{ active: richEditor?.isActive('heading', { level: 1 }) }">H1</button>
            <button class="toolbar-btn" title="标题2" @click.prevent="richEditor?.chain().focus().toggleHeading({ level: 2 }).run()" :class="{ active: richEditor?.isActive('heading', { level: 2 }) }">H2</button>
            <button class="toolbar-btn" title="标题3" @click.prevent="richEditor?.chain().focus().toggleHeading({ level: 3 }).run()" :class="{ active: richEditor?.isActive('heading', { level: 3 }) }">H3</button>
            <span class="toolbar-divider"></span>
            <button class="toolbar-btn" title="无序列表" @click.prevent="richEditor?.chain().focus().toggleBulletList().run()" :class="{ active: richEditor?.isActive('bulletList') }">•</button>
            <button class="toolbar-btn" title="有序列表" @click.prevent="richEditor?.chain().focus().toggleOrderedList().run()" :class="{ active: richEditor?.isActive('orderedList') }">1.</button>
            <button class="toolbar-btn" title="引用" @click.prevent="richEditor?.chain().focus().toggleBlockquote().run()" :class="{ active: richEditor?.isActive('blockquote') }">"</button>
            <button class="toolbar-btn" title="代码块" @click.prevent="richEditor?.chain().focus().toggleCodeBlock().run()" :class="{ active: richEditor?.isActive('codeBlock') }">&lt;/&gt;</button>
            <button class="toolbar-btn" title="分割线" @click.prevent="richEditor?.chain().focus().setHorizontalRule().run()">—</button>
            <span class="toolbar-divider"></span>
            <button class="toolbar-btn" title="撤销" @click.prevent="richEditor?.chain().focus().undo().run()">↩</button>
            <button class="toolbar-btn" title="重做" @click.prevent="richEditor?.chain().focus().redo().run()">↪</button>
          </div>
          <div ref="richEditorDom" class="richtext-content"></div>
        </div>
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
  width: 100%;
}

.form-meta {
  max-width: 800px;
}

.prop-label {
  font-size: 13px;
  color: #666;
  margin-left: 8px;
}
.prop-label:first-child {
  margin-left: 0;
}

.content-form-item {
  width: 100%;
}

.content-form-item :deep(.n-form-item-blank) {
  display: block;
  width: 100%;
}

.richtext-editor {
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  overflow: hidden;
}

.richtext-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
  padding: 6px 8px;
  background: #fafafa;
  border-bottom: 1px solid #e0e0e0;
}

.richtext-toolbar .toolbar-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #333;
  transition: background 0.15s;
}

.richtext-toolbar .toolbar-btn:hover {
  background: #e8e8e8;
}

.richtext-toolbar .toolbar-btn.active {
  background: #d0e0ff;
  color: #2080f0;
}

.richtext-toolbar .toolbar-divider {
  width: 1px;
  background: #ddd;
  margin: 0 4px;
}

.richtext-content {
  min-height: 68vh;
  padding: 12px 16px;
  outline: none;
}

.richtext-content :deep(.ProseMirror) {
  min-height: 66vh;
  outline: none;
  font-size: 15px;
  line-height: 1.75;
}

.richtext-content :deep(.ProseMirror p) {
  margin: 0 0 8px 0;
}

.richtext-content :deep(.ProseMirror h1) {
  font-size: 24px;
  margin: 16px 0 8px 0;
}

.richtext-content :deep(.ProseMirror h2) {
  font-size: 20px;
  margin: 14px 0 8px 0;
}

.richtext-content :deep(.ProseMirror h3) {
  font-size: 17px;
  margin: 12px 0 6px 0;
}

.richtext-content :deep(.ProseMirror blockquote) {
  border-left: 3px solid #ddd;
  padding-left: 12px;
  margin: 8px 0;
  color: #666;
}

.richtext-content :deep(.ProseMirror ul),
.richtext-content :deep(.ProseMirror ol) {
  padding-left: 24px;
  margin: 4px 0;
}

.richtext-content :deep(.ProseMirror code) {
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 13px;
}

.richtext-content :deep(.ProseMirror pre) {
  background: #282c34;
  color: #abb2bf;
  padding: 12px 16px;
  border-radius: 6px;
  overflow-x: auto;
}

.richtext-content :deep(.ProseMirror pre code) {
  background: none;
  padding: 0;
  color: inherit;
}

.richtext-content :deep(.ProseMirror hr) {
  border: none;
  border-top: 1px solid #ddd;
  margin: 16px 0;
}
</style>