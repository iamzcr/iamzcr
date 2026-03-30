<script setup lang="ts">
import { ref, onMounted, computed, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NSpace, NTag, NSpin, NEmpty } from 'naive-ui'
import { articleApi, categoryApi, directoryApi, tagsApi } from '../api'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

const route = useRoute()
const router = useRouter()
const md = new MarkdownIt({
  breaks: true,
  highlight: (code, lang) => {
    const language = hljs.getLanguage(lang) ? lang : 'plaintext'
    const highlighted = hljs.highlight(code, { language }).value
    return `<pre class="hljs-code"><code class="language-${language}">${highlighted}</code></pre>`
  }
})

const article = ref<any>(null)
const category = ref<any>(null)
const directory = ref<any>(null)
const tags = ref<any[]>([])
const loading = ref(false)

const latestArticles = ref<any[]>([])
const categories = ref<any[]>([])
const directories = ref<any[]>([])
const allTags = ref<any[]>([])
const articlesByDirectory = ref<Record<number, any[]>>({})
const expandedDirs = ref<Set<number>>(new Set())

const renderedContent = computed(() => {
  const content = article.value?.article?.content || article.value?.content || ''
  return content ? md.render(content) : ''
})

async function loadArticle() {
  loading.value = true
  try {
    const res = await articleApi.get(Number(route.params.id))
    const data = res.data.data
    article.value = data.article || data
    category.value = data.category || null
    directory.value = data.directory || null
    tags.value = data.tags || []
  } catch {
    article.value = null
  } finally {
    loading.value = false
    await nextTick()
    addCopyButtons()
  }
}

async function loadSidebarData() {
  const cid = category.value?.id
  const [articleRes, catRes, tagRes, dirRes] = await Promise.all([
    articleApi.list({ page: 1, page_size: 100 }),
    categoryApi.list(),
    tagsApi.list(),
    directoryApi.list({ cid })
  ])
  const allArticles = articleRes.data.data.list
  latestArticles.value = allArticles.slice(0, 5)
  categories.value = catRes.data.data
  allTags.value = tagRes.data.data
  directories.value = dirRes.data.data
  
  const byDir: Record<number, any[]> = {}
  allArticles.forEach((a: any) => {
    if (a.did) {
      if (!byDir[a.did]) byDir[a.did] = []
      byDir[a.did].push(a)
    }
  })
  articlesByDirectory.value = byDir
}

function getDirectoriesByCategory() {
  if (!category.value) return []
  return directories.value.filter((d: any) => d.cid === category.value.id)
}

function toggleDir(id: number) {
  if (expandedDirs.value.has(id)) {
    expandedDirs.value.delete(id)
  } else {
    expandedDirs.value.add(id)
  }
}

function goToArticle(id: number) {
  router.push(`/article/${id}`)
}

function goToCategory(id: number) {
  router.push(`/category/${id}`)
}

function goToTag(id: number) {
  router.push(`/tag/${id}`)
}

function formatDate(timestamp: number) {
  if (!timestamp) return '-'
  return new Date(timestamp * 1000).toLocaleDateString('zh-CN')
}

function addCopyButtons() {
  document.querySelectorAll('.hljs-code').forEach((pre) => {
    if (pre.parentElement?.classList.contains('code-wrapper')) return
    
    const wrapper = document.createElement('div')
    wrapper.className = 'code-wrapper'
    
    const copyBtn = document.createElement('button')
    copyBtn.className = 'copy-btn'
    copyBtn.textContent = '复制'
    copyBtn.onclick = () => {
      const code = pre.querySelector('code')?.textContent || ''
      navigator.clipboard.writeText(code)
      copyBtn.textContent = '已复制'
      setTimeout(() => { copyBtn.textContent = '复制' }, 2000)
    }
    
    pre.parentElement?.insertBefore(wrapper, pre)
    wrapper.appendChild(pre)
    wrapper.appendChild(copyBtn)
  })
}

onMounted(() => {
  loadArticle()
  loadSidebarData()
})

watch(() => route.params.id, () => {
  loadArticle()
  window.scrollTo(0, 0)
})
</script>

<template>
  <n-spin :show="loading">
    <n-empty v-if="!loading && !article" description="文章不存在" />
    <div v-else-if="article" class="article-layout">
      <div class="article-main">
        <div class="article-container">
          <div class="article-header">
            <div class="article-meta">
              <n-space>
                <n-tag v-if="category" type="primary" @click="goToCategory(category.id)" class="clickable-tag">{{ category.name }}</n-tag>
                <n-tag v-if="directory" type="info">{{ directory.name }}</n-tag>
                <n-tag v-for="tag in tags" :key="tag.id" type="warning" @click="goToTag(tag.id)" class="clickable-tag">{{ tag.name }}</n-tag>
              </n-space>
            </div>
            <h1 class="article-title">{{ article.title }}</h1>
            <div class="article-info">
              <span v-if="article.author">作者: {{ article.author }}</span>
              <span>发布时间: {{ formatDate(article.create_time) }}</span>
              <span>阅读: {{ article.weight }}</span>
            </div>
          </div>
          <div class="article-cover" v-if="article.thumb">
            <img :src="article.thumb" :alt="article.title" />
          </div>
          <div class="markdown-content" v-html="renderedContent"></div>
          <div class="article-footer">
            <n-button @click="router.push('/')">返回首页</n-button>
          </div>
        </div>
      </div>
      
      <aside class="sidebar">
        <div class="sidebar-section" v-if="category">
          <h3 class="section-title">{{ category.name }} - 目录</h3>
          <div class="directory-list">
            <div 
              v-for="dir in getDirectoriesByCategory()" 
              :key="dir.id" 
              class="directory-item"
            >
              <div 
                class="directory-header"
                :class="{ active: dir.id === directory?.id, expanded: expandedDirs.has(dir.id) }"
                @click="toggleDir(dir.id)"
              >
                <span class="expand-icon">{{ expandedDirs.has(dir.id) ? '▼' : '▶' }}</span>
                <span class="directory-name">{{ dir.name }}</span>
                <span class="directory-count">({{ (articlesByDirectory[dir.id] || []).length }})</span>
              </div>
              <div class="directory-articles" v-show="expandedDirs.has(dir.id)">
                <div 
                  v-for="a in articlesByDirectory[dir.id]" 
                  :key="a.id"
                  class="directory-article-item"
                  :class="{ active: a.id === article?.id }"
                  @click="goToArticle(a.id)"
                >
                  {{ a.title }}
                </div>
                <div v-if="!articlesByDirectory[dir.id]?.length" class="no-articles">暂无文章</div>
              </div>
            </div>
          </div>
        </div>

        <div class="sidebar-section">
          <h3 class="section-title">标签云</h3>
          <div class="tag-cloud">
            <n-tag 
              v-for="tag in allTags" 
              :key="tag.id" 
              :type="tag.is_hot ? 'warning' : 'default'"
              class="cloud-tag"
              @click="goToTag(tag.id)"
            >
              {{ tag.name }}
            </n-tag>
          </div>
        </div>

        <div class="sidebar-section">
          <h3 class="section-title">最新文章</h3>
          <div class="latest-articles">
            <div 
              v-for="article in latestArticles" 
              :key="article.id" 
              class="latest-item"
              @click="goToArticle(article.id)"
            >
              <div class="latest-title">{{ article.title }}</div>
              <div class="latest-meta">{{ formatDate(article.create_time) }}</div>
            </div>
          </div>
        </div>
      </aside>
    </div>
  </n-spin>
</template>

<style scoped>
.article-layout {
  display: flex;
  gap: 32px;
  padding: 0 32px;
  width: 100%;
}

.article-main {
  flex: 1;
  min-width: 0;
}

.article-container {
  background: #fff;
  padding: 32px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.article-header {
  margin-bottom: 24px;
}
.article-meta {
  margin-bottom: 16px;
}
.clickable-tag {
  cursor: pointer;
}
.clickable-tag:hover {
  opacity: 0.8;
}
.article-title {
  font-size: 32px;
  font-weight: 700;
  color: #333;
  margin: 0 0 16px 0;
  line-height: 1.3;
}
.article-info {
  display: flex;
  gap: 24px;
  color: #999;
  font-size: 14px;
}
.article-cover {
  margin-bottom: 24px;
  border-radius: 12px;
  overflow: hidden;
}
.article-cover img {
  width: 100%;
  height: auto;
}
.markdown-content {
  background: #fff;
  padding: 0;
  line-height: 1.8;
  text-align: left;
}
.markdown-content :deep(h1) { font-size: 2em; border-bottom: 1px solid #eee; padding-bottom: 0.3em; margin-top: 1.5em; }
.markdown-content :deep(h2) { font-size: 1.5em; border-bottom: 1px solid #eee; padding-bottom: 0.3em; margin-top: 1.5em; }
.markdown-content :deep(pre) { 
  background: #f6f8fa; 
  padding: 16px; 
  border-radius: 6px; 
  overflow-x: auto; 
  position: relative; 
  white-space: pre;
}
.markdown-content :deep(.hljs-code) { 
  background: #f6f8fa; 
  padding: 16px; 
  border-radius: 6px; 
  overflow-x: auto; 
  position: relative;
  margin: 1em 0;
  white-space: pre;
  display: block;
}
.markdown-content :deep(.hljs-code code) { 
  background: none; 
  padding: 0;
  white-space: pre;
  font-family: var(--mono);
  display: block;
}
.markdown-content :deep(.code-wrapper) {
  position: relative;
  margin: 1em 0;
}
.markdown-content :deep(.copy-btn) {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 4px 12px;
  font-size: 12px;
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}
.markdown-content :deep(.copy-btn:hover) {
  background: #f0f0f0;
}
.markdown-content :deep(blockquote) { border-left: 4px solid #ddd; padding-left: 16px; color: #666; margin: 1em 0; }
.markdown-content :deep(img) { max-width: 100%; border-radius: 8px; }
.markdown-content :deep(p) { margin: 1em 0; }
.markdown-content :deep(ul), .markdown-content :deep(ol) { padding-left: 2em; margin: 1em 0; }
.markdown-content :deep(table) { border-collapse: collapse; width: 100%; margin: 1em 0; }
.markdown-content :deep(th), .markdown-content :deep(td) { border: 1px solid #ddd; padding: 8px; }
.markdown-content :deep(th) { background: #f6f8fa; }
.article-footer {
  margin-top: 32px;
  text-align: center;
}

.sidebar {
  width: 320px;
  flex-shrink: 0;
}

.sidebar-section {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0 0 16px 0;
  padding-bottom: 12px;
  border-bottom: 1px solid #eee;
}

.latest-articles {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.latest-item {
  cursor: pointer;
  padding: 8px 0;
  border-bottom: 1px solid #f5f5f5;
  transition: all 0.2s;
}

.latest-item:last-child {
  border-bottom: none;
}

.latest-item:hover .latest-title {
  color: #18a058;
}

.latest-title {
  font-size: 14px;
  color: #333;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.latest-meta {
  font-size: 12px;
  color: #999;
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.cloud-tag {
  cursor: pointer;
  transition: all 0.2s;
}

.cloud-tag:hover {
  transform: scale(1.1);
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.category-item:hover {
  background: #f5f7fa;
}

.category-item.active {
  background: #e6f7ed;
  color: #18a058;
}

.directory-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.directory-item {
  border-radius: 6px;
  overflow: hidden;
}

.directory-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s;
}

.directory-header:hover {
  background: #f5f7fa;
}

.directory-header.active {
  background: #e6f7ed;
  color: #18a058;
}

.directory-header.expanded {
  background: #f0f5ff;
}

.expand-icon {
  font-size: 10px;
  color: #999;
  width: 12px;
}

.directory-name {
  font-size: 14px;
  color: #333;
  flex: 1;
}

.directory-count {
  font-size: 12px;
  color: #999;
}

.directory-articles {
  padding-left: 24px;
}

.directory-article-item {
  padding: 8px 12px;
  font-size: 13px;
  color: #666;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.directory-article-item:hover {
  background: #f5f7fa;
  color: #18a058;
}

.directory-article-item.active {
  color: #18a058;
  font-weight: 500;
}

.no-articles {
  padding: 8px 12px;
  font-size: 12px;
  color: #999;
}

.category-name {
  font-size: 14px;
  color: #333;
}

.category-count {
  font-size: 12px;
  color: #999;
  background: #f0f0f0;
  padding: 2px 8px;
  border-radius: 10px;
}

@media (max-width: 900px) {
  .article-layout {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
  }
}
</style>