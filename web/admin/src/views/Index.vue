<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { articleApi, categoryApi, tagsApi } from '../api'

const loading = ref(false)
const stats = ref([
  { label: '文章数量', value: 0, accent: '#2563eb' },
  { label: '分类数量', value: 0, accent: '#0f766e' },
  { label: '标签数量', value: 0, accent: '#7c3aed' }
])

function getTotal(data: any) {
  if (typeof data?.total === 'number') return data.total
  if (Array.isArray(data?.list)) return data.list.length
  if (Array.isArray(data)) return data.length
  return 0
}

async function loadStats() {
  loading.value = true
  try {
    const [articleRes, categoryRes, tagsRes] = await Promise.all([
      articleApi.list({ page: 1, page_size: 1 }),
      categoryApi.list({ page: 1, page_size: 1 }),
      tagsApi.list({ page: 1, page_size: 1 })
    ])

    stats.value = [
      { label: '文章数量', value: getTotal(articleRes.data.data), accent: '#2563eb' },
      { label: '分类数量', value: getTotal(categoryRes.data.data), accent: '#0f766e' },
      { label: '标签数量', value: getTotal(tagsRes.data.data), accent: '#7c3aed' }
    ]
  } finally {
    loading.value = false
  }
}

onMounted(loadStats)
</script>

<template>
  <div class="dashboard-page">
    <div class="dashboard-head">
      <div>
        <h1>概况</h1>
        <p>查看当前博客的基础数据。</p>
      </div>
    </div>

    <div class="stats-grid" v-if="!loading">
      <div v-for="item in stats" :key="item.label" class="stat-card">
        <span class="stat-bar" :style="{ background: item.accent }"></span>
        <div class="stat-label">{{ item.label }}</div>
        <div class="stat-value">{{ item.value }}</div>
      </div>
    </div>

    <div v-else class="dashboard-loading">加载中...</div>
  </div>
</template>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.dashboard-head h1 {
  margin-bottom: 8px;
  color: #0f172a;
  font-size: 28px;
}

.dashboard-head p {
  color: #64748b;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 20px;
}

.stat-card {
  position: relative;
  overflow: hidden;
  padding: 24px;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  background: #ffffff;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.06);
}

.stat-bar {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 4px;
}

.stat-label {
  margin-bottom: 14px;
  color: #64748b;
  font-size: 14px;
}

.stat-value {
  color: #0f172a;
  font-size: 34px;
  font-weight: 700;
}

.dashboard-loading {
  color: #64748b;
}
</style>
