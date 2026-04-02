<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NEmpty, NTag, NSpin, NPagination } from 'naive-ui'
import { articleApi, categoryApi, directoryApi, tagsApi } from '../api'

const router = useRouter()
const route = useRoute()

const articles = ref<any[]>([])
const latestArticles = ref<any[]>([])
const categories = ref<any[]>([])
const directories = ref<any[]>([])
const allTags = ref<any[]>([])
const loading = ref(false)

const currentType = ref('')
const currentId = ref(0)
const page = ref(1)
const pageSize = ref(10)
const totalArticles = ref(0)

async function loadData() {
  loading.value = true
  try {
    updateCurrentFilter()

    const articleParams: any = { page: page.value, page_size: pageSize.value }
    if (currentType.value === 'category' && currentId.value) articleParams.cid = currentId.value
    if (currentType.value === 'directory' && currentId.value) articleParams.did = currentId.value
    if (currentType.value === 'tag' && currentId.value) articleParams.tid = currentId.value

    const [articleRes, latestRes, catRes, dirRes, tagRes] = await Promise.all([
      articleApi.list(articleParams),
      articleApi.list({ page: 1, page_size: 5 }),
      categoryApi.list({ page: 1, page_size: 1000 }),
      directoryApi.list({ page: 1, page_size: 1000 }),
      tagsApi.list({ page: 1, page_size: 1000 })
    ])

    articles.value = articleRes.data.data.list || []
    totalArticles.value = articleRes.data.data.total || 0
    latestArticles.value = latestRes.data.data.list || []
    categories.value = catRes.data.data.list || catRes.data.data || []
    directories.value = dirRes.data.data.list || dirRes.data.data || []
    allTags.value = tagRes.data.data.list || tagRes.data.data || []
  } finally {
    loading.value = false
  }
}

function updateCurrentFilter() {
  currentType.value = route.name === 'category' ? 'category' : 
                     route.name === 'tag' ? 'tag' :
                     route.name === 'directory' ? 'directory' : ''
  currentId.value = Number(route.params.id) || 0
}

function getCategoryName(cid: number) {
  const cat = categories.value.find(c => c.id === cid)
  return cat ? cat.name : ''
}

function getDirectoryName(did: number) {
  const dir = directories.value.find(d => d.id === did)
  return dir ? dir.name : ''
}

function getTags(article: any) {
  return article.tags || []
}

function formatDate(timestamp: number) {
  if (!timestamp) return '-'
  return new Date(timestamp * 1000).toLocaleDateString('zh-CN')
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

watch(() => route.params.id, () => {
  page.value = 1
  loadData()
})

watch(() => route.name, () => {
  page.value = 1
  loadData()
})

onMounted(loadData)
</script>

<template>
  <div class="home-page">
    <n-spin :show="loading">
      <div class="layout">
        <main class="main-content">
          <n-empty v-if="!loading && articles.length === 0" description="暂无文章" />
          <div class="articles-list">
            <div 
              v-for="article in articles" 
              :key="article.id" 
              class="article-card"
              @click="goToArticle(article.id)"
            >
              <div class="article-cover" v-if="article.thumb">
                <img :src="article.thumb" :alt="article.title" />
              </div>
              <div class="article-body">
                <div class="article-meta">
                  <n-tag v-if="getCategoryName(article.cid)" size="small" type="primary" @click.stop="goToCategory(article.cid)">
                    {{ getCategoryName(article.cid) }}
                  </n-tag>
                  <n-tag v-if="getDirectoryName(article.did)" size="small" type="info">
                    {{ getDirectoryName(article.did) }}
                  </n-tag>
                  <span class="article-date">{{ formatDate(article.create_time) }}</span>
                </div>
                <h2 class="article-title">{{ article.title }}</h2>
                <p class="article-summary">{{ article.summary || article.desc || '暂无摘要' }}</p>
                <div class="article-tags" v-if="getTags(article).length > 0">
                  <n-tag 
                    v-for="tag in getTags(article)" 
                    :key="tag.id" 
                    size="small" 
                    type="warning"
                    class="tag-item"
                    @click.stop="goToTag(tag.id)"
                  >
                    {{ tag.name }}
                  </n-tag>
                </div>
              </div>
            </div>
          </div>
          
            <div class="pagination-wrapper" v-if="totalArticles > pageSize">
              <n-pagination 
                v-model:page="page" 
                :page-size="pageSize" 
                :item-count="totalArticles"
                @update:page="(p: number) => { page = p; loadData() }"
              />
            </div>
        </main>

        <aside class="sidebar">
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
            <h3 class="section-title">分类目录</h3>
            <div class="category-list">
              <div 
                v-for="cat in categories" 
                :key="cat.id" 
                class="category-item"
                @click="goToCategory(cat.id)"
              >
                <span class="category-name">{{ cat.name }}</span>
                <span class="category-count">></span>
              </div>
            </div>
          </div>
        </aside>
      </div>
    </n-spin>
  </div>
</template>

<style scoped>
.home-page {
  min-height: calc(100vh - 120px);
}

.layout {
  display: flex;
  gap: 32px;
  padding: 0 32px;
  width: 100%;
}

.main-content {
  flex: 1;
  min-width: 0;
}

.content-header {
  margin-bottom: 24px;
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

.articles-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.article-card {
  display: flex;
  background: #fff;
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
}

.article-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.1);
}

.article-cover {
  width: 400px;
  min-width: 400px;
  height: 240px;
  overflow: hidden;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.article-card:hover .article-cover img {
  transform: scale(1.05);
}

.article-body {
  flex: 1;
  padding: 24px;
  display: flex;
  flex-direction: column;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.article-date {
  color: #999;
  font-size: 13px;
}

.article-title {
  font-size: 22px;
  font-weight: 600;
  color: #333;
  margin: 0 0 12px 0;
  line-height: 1.4;
}

.article-summary {
  color: #666;
  font-size: 14px;
  line-height: 1.7;
  margin: 0 0 16px 0;
  flex: 1;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tag-item {
  font-size: 12px;
  cursor: pointer;
}

@media (max-width: 1200px) {
  .article-cover {
    width: 320px;
    min-width: 320px;
    height: 200px;
  }
}

@media (max-width: 900px) {
  .layout {
    flex-direction: column-reverse;
  }
  
  .sidebar {
    width: 100%;
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 20px;
  }
  
  .sidebar-section {
    margin-bottom: 0;
  }
  
  .article-card {
    flex-direction: column;
  }
  
  .article-cover {
    width: 100%;
    height: 220px;
  }
}

@media (max-width: 600px) {
  .sidebar {
    grid-template-columns: 1fr;
  }
  
  .layout {
    padding: 0 16px;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 32px 0;
}
</style>
