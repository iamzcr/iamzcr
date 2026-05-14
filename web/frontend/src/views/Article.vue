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
const readCount = ref(0)
const prevArticle = ref<any>(null)
const nextArticle = ref<any>(null)
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
    readCount.value = data.read_count || 0
    prevArticle.value = data.prev_article || null
    nextArticle.value = data.next_article || null
  } catch {
    article.value = null
  } finally {
    loading.value = false
    await nextTick()
    enhanceCodeBlocks()
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

function enhanceCodeBlocks() {
  document.querySelectorAll('.hljs-code').forEach((pre) => {
    if (pre.closest('.code-block')) return

    const codeEl = pre.querySelector('code')
    const langClass = codeEl?.className.match(/language-(\w+)/)?.[1]
    const langName = langClass === 'plaintext' ? 'Plain Text' : (langClass || 'Code')

    const codeText = codeEl?.textContent || ''
    const lineCount = codeText.split('\n').length
    const isLong = lineCount > 15

    const block = document.createElement('div')
    block.className = 'code-block'

    const header = document.createElement('div')
    header.className = 'code-header'

    const lang = document.createElement('span')
    lang.className = 'code-lang'
    const langMap: Record<string, string> = {
      javascript: 'JavaScript', typescript: 'TypeScript', python: 'Python',
      go: 'Go', rust: 'Rust', java: 'Java', cpp: 'C++', c: 'C',
      html: 'HTML', css: 'CSS', scss: 'SCSS', json: 'JSON',
      xml: 'XML', yaml: 'YAML', markdown: 'Markdown', shell: 'Shell',
      bash: 'Bash', sh: 'Shell', zsh: 'Zsh', powershell: 'PowerShell',
      sql: 'SQL', php: 'PHP', ruby: 'Ruby', swift: 'Swift',
      kotlin: 'Kotlin', dart: 'Dart', lua: 'Lua', r: 'R',
      dockerfile: 'Dockerfile', nginx: 'Nginx', plaintext: 'Plain Text',
    }
    lang.textContent = langMap[langName.toLowerCase()] || langName

    const headerRight = document.createElement('div')
    headerRight.className = 'code-header-right'

    // Copy button
    const copyBtn = document.createElement('button')
    copyBtn.className = 'code-btn copy-btn'
    copyBtn.appendChild(makeCopyIcon())
    const copyLabel = document.createElement('span')
    copyLabel.textContent = '复制'
    copyBtn.appendChild(copyLabel)
    copyBtn.title = '复制代码'
    copyBtn.onclick = (e) => {
      e.stopPropagation()
      const code = pre.querySelector('code')?.textContent || ''
      copyToClipboard(code)
      copyBtn.classList.add('copied')
      copyLabel.textContent = '已复制'
      setTimeout(() => {
        copyBtn.classList.remove('copied')
        copyLabel.textContent = '复制'
      }, 2000)
    }

    // Toggle button
    const toggleBtn = document.createElement('button')
    toggleBtn.className = 'code-btn toggle-btn'
    const toggleArrow = makeChevronIcon()
    toggleBtn.appendChild(toggleArrow)
    toggleBtn.title = '折叠代码'
    let collapsed = false

    const body = document.createElement('div')
    body.className = 'code-body'

    toggleBtn.onclick = (e) => {
      e.stopPropagation()
      collapsed = !collapsed
      if (collapsed) {
        block.classList.add('collapsed')
        toggleBtn.title = '展开代码'
      } else {
        block.classList.remove('collapsed')
        toggleBtn.title = '折叠代码'
      }
    }

    if (isLong) {
      const linesBadge = document.createElement('span')
      linesBadge.className = 'code-lines'
      linesBadge.textContent = `${lineCount} 行`
      headerRight.appendChild(linesBadge)
    }

    headerRight.appendChild(copyBtn)
    headerRight.appendChild(toggleBtn)

    header.appendChild(lang)
    header.appendChild(headerRight)

    header.onclick = () => toggleBtn.click()

    block.appendChild(header)
    block.appendChild(body)
    pre.parentElement?.insertBefore(block, pre)
    body.appendChild(pre)
  })
}

function makeCopyIcon(): SVGElement {
  const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg')
  svg.setAttribute('class', 'btn-icon')
  svg.setAttribute('width', '13')
  svg.setAttribute('height', '13')
  svg.setAttribute('viewBox', '0 0 24 24')
  svg.setAttribute('fill', 'none')
  svg.setAttribute('stroke', 'currentColor')
  svg.setAttribute('stroke-width', '2')
  svg.setAttribute('stroke-linecap', 'round')
  svg.setAttribute('stroke-linejoin', 'round')

  const r1 = document.createElementNS('http://www.w3.org/2000/svg', 'rect')
  r1.setAttribute('x', '9'); r1.setAttribute('y', '9')
  r1.setAttribute('width', '13'); r1.setAttribute('height', '13')
  r1.setAttribute('rx', '2'); svg.appendChild(r1)

  const p = document.createElementNS('http://www.w3.org/2000/svg', 'path')
  p.setAttribute('d', 'M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1')
  svg.appendChild(p)

  return svg
}

function makeChevronIcon(): SVGElement {
  const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg')
  svg.setAttribute('class', 'btn-icon toggle-arrow')
  svg.setAttribute('width', '13')
  svg.setAttribute('height', '13')
  svg.setAttribute('viewBox', '0 0 24 24')
  svg.setAttribute('fill', 'none')
  svg.setAttribute('stroke', 'currentColor')
  svg.setAttribute('stroke-width', '2.5')
  svg.setAttribute('stroke-linecap', 'round')
  svg.setAttribute('stroke-linejoin', 'round')

  const poly = document.createElementNS('http://www.w3.org/2000/svg', 'polyline')
  poly.setAttribute('points', '6 9 12 15 18 9')
  svg.appendChild(poly)

  return svg
}

function copyToClipboard(text: string) {
  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard.writeText(text).catch(() => {
      fallbackCopy(text)
    })
  } else {
    fallbackCopy(text)
  }
}

function fallbackCopy(text: string) {
  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.style.position = 'fixed'
  textarea.style.left = '-9999px'
  textarea.style.top = '-9999px'
  document.body.appendChild(textarea)
  textarea.focus()
  textarea.select()
  try {
    document.execCommand('copy')
  } catch { /* ignore */ }
  document.body.removeChild(textarea)
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
              <span>阅读: {{ readCount }}</span>
            </div>
          </div>
          <div class="markdown-content" v-html="renderedContent"></div>
          <div class="article-footer">
            <n-button @click="router.push('/')">返回首页</n-button>
          </div>
          <div class="article-nav" v-if="prevArticle || nextArticle">
            <div class="nav-item" v-if="prevArticle" @click="goToArticle(prevArticle.id)">
              <span class="nav-label">上一篇</span>
              <span class="nav-title">{{ prevArticle.title }}</span>
            </div>
            <div class="nav-spacer"></div>
            <div class="nav-item nav-item-next" v-if="nextArticle" @click="goToArticle(nextArticle.id)">
              <span class="nav-label">下一篇</span>
              <span class="nav-title">{{ nextArticle.title }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </n-spin>
</template>

<style scoped>
.article-layout {
  display: flex;
  gap: 32px;
  width: 100%;
}

.article-main {
  flex: 1;
  min-width: 0;
}

.article-container {
  background: var(--card-bg);
  padding: 32px;
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border);
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
  font-size: 30px;
  font-weight: 700;
  color: var(--text-h);
  margin: 0 0 16px 0;
  line-height: 1.3;
  letter-spacing: -0.01em;
}
.article-info {
  display: flex;
  gap: 24px;
  color: var(--text-muted);
  font-size: 14px;
}
.article-cover {
  margin-bottom: 24px;
  border-radius: var(--radius);
  overflow: hidden;
}
.article-cover img {
  width: 100%;
  height: auto;
}
.markdown-content {
  background: transparent;
  padding: 0;
  line-height: 1.8;
  text-align: left;
}
.markdown-content :deep(h1) { font-size: 2em; border-bottom: 1px solid var(--border); padding-bottom: 0.3em; margin-top: 1.5em; color: var(--text-h); }
.markdown-content :deep(h2) { font-size: 1.5em; border-bottom: 1px solid var(--border); padding-bottom: 0.3em; margin-top: 1.5em; color: var(--text-h); }
.markdown-content :deep(pre) { 
  background: var(--code-bg); 
  padding: 16px; 
  border-radius: var(--radius-sm); 
  overflow-x: auto; 
  position: relative; 
  white-space: pre;
}
.markdown-content :deep(.hljs-code) { 
  background: var(--code-bg); 
  padding: 16px; 
  border-radius: 0 0 var(--radius-sm) var(--radius-sm); 
  overflow-x: auto; 
  position: relative;
  margin: 0;
  white-space: pre;
  display: block;
}
.markdown-content :deep(.hljs-code code) { 
  background: none; 
  padding: 0;
  white-space: pre;
  font-family: var(--mono);
  display: block;
  font-size: 14px;
}
.markdown-content :deep(.code-block) {
  margin: 1.2em 0;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  overflow: hidden;
  background: var(--code-bg);
}
.markdown-content :deep(.code-header) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: var(--card-bg);
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  user-select: none;
  transition: background 0.15s;
}
.markdown-content :deep(.code-header:hover) {
  background: var(--border);
}
.markdown-content :deep(.code-lang) {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.markdown-content :deep(.code-lines) {
  font-size: 11px;
  color: var(--text-muted);
  margin-right: 4px;
}
.markdown-content :deep(.code-header-right) {
  display: flex;
  align-items: center;
  gap: 8px;
}
.markdown-content :deep(.code-btn) {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 4px 12px;
  font-size: 11px;
  font-weight: 500;
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: var(--text-muted);
  line-height: 1.5;
}
.markdown-content :deep(.code-btn .btn-icon) {
  flex-shrink: 0;
}
.markdown-content :deep(.code-btn:hover) {
  background: var(--accent-bg);
  border-color: var(--accent-border);
  color: var(--accent);
}
.markdown-content :deep(.code-btn:active) {
  transform: scale(0.96);
}
/* Copy button — copied state */
.markdown-content :deep(.copy-btn.copied) {
  background: var(--accent-bg);
  border-color: var(--accent);
  color: var(--accent);
}
/* Toggle button */
.markdown-content :deep(.toggle-btn) {
  padding: 4px 8px;
}
.markdown-content :deep(.toggle-arrow) {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.markdown-content :deep(.code-block.collapsed .toggle-arrow) {
  transform: rotate(-90deg);
}
.markdown-content :deep(.code-body) {
  overflow: hidden;
  transition: max-height 0.35s ease, opacity 0.25s ease, padding 0.25s ease;
  max-height: 6000px;
  opacity: 1;
}
.markdown-content :deep(.code-block.collapsed .code-body) {
  max-height: 0;
  opacity: 0;
  padding-top: 0;
  padding-bottom: 0;
}
.markdown-content :deep(blockquote) { border-left: 4px solid var(--border); padding-left: 16px; color: var(--text-muted); margin: 1em 0; }
.markdown-content :deep(img) { max-width: 100%; border-radius: 8px; }
.markdown-content :deep(p) { margin: 1em 0; }
.markdown-content :deep(ul), .markdown-content :deep(ol) { padding-left: 2em; margin: 1em 0; }
.markdown-content :deep(table) { border-collapse: collapse; width: 100%; margin: 1em 0; }
.markdown-content :deep(th), .markdown-content :deep(td) { border: 1px solid var(--border); padding: 8px; }
.markdown-content :deep(th) { background: var(--code-bg); }
.article-footer {
  margin-top: 32px;
  text-align: center;
}

.sidebar {
  width: 320px;
  flex-shrink: 0;
}

.sidebar-section {
  background: var(--card-bg);
  border-radius: var(--radius);
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border);
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-h);
  margin: 0 0 16px 0;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border);
  letter-spacing: 0.01em;
}

.latest-articles {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.latest-item {
  cursor: pointer;
  padding: 8px 0;
  border-bottom: 1px solid var(--border);
  transition: all 0.2s;
}

.latest-item:last-child {
  border-bottom: none;
}

.latest-item:hover .latest-title {
  color: var(--accent);
}

.latest-title {
  font-size: 14px;
  color: var(--text-h);
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.latest-meta {
  font-size: 12px;
  color: var(--text-muted);
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
  gap: 6px;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
}

.category-item:hover {
  background: var(--accent-bg);
}

.category-item.active {
  background: var(--accent-bg);
  color: var(--accent);
}

.directory-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.directory-item {
  border-radius: var(--radius-sm);
  overflow: hidden;
}

.directory-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  cursor: pointer;
  border-radius: var(--radius-sm);
  transition: all 0.2s;
}

.directory-header:hover {
  background: var(--accent-bg);
}

.directory-header.active {
  background: var(--accent-bg);
  color: var(--accent);
}

.directory-header.expanded {
  background: var(--accent-bg);
}

.expand-icon {
  font-size: 10px;
  color: var(--text-muted);
  width: 12px;
}

.directory-name {
  font-size: 14px;
  color: var(--text-h);
  flex: 1;
}

.directory-count {
  font-size: 12px;
  color: var(--text-muted);
}

.directory-articles {
  padding-left: 24px;
}

.directory-article-item {
  padding: 8px 12px;
  font-size: 13px;
  color: var(--text);
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.directory-article-item:hover {
  background: var(--accent-bg);
  color: var(--accent);
}

.directory-article-item.active {
  color: var(--accent);
  font-weight: 500;
}

.no-articles {
  padding: 8px 12px;
  font-size: 12px;
  color: var(--text-muted);
}

.category-name {
  font-size: 14px;
  color: var(--text);
}

.category-count {
  font-size: 12px;
  color: var(--text-muted);
  background: var(--border);
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

.article-nav {
  display: flex;
  align-items: stretch;
  margin-top: 32px;
  border-top: 1px solid var(--border);
  padding-top: 20px;
}

.nav-item {
  flex: 1;
  cursor: pointer;
  padding: 12px 16px;
  border-radius: 8px;
  transition: background 0.2s;
}

.nav-item:hover {
  background: var(--bg-secondary);
}

.nav-item-next {
  text-align: right;
}

.nav-label {
  display: block;
  font-size: 12px;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.nav-title {
  display: block;
  font-size: 14px;
  color: var(--text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.nav-spacer {
  width: 1px;
  background: var(--border);
  margin: 0 16px;
}
</style>