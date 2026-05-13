<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { articleApi, categoryApi, tagsApi } from '../api'

const loading = ref(false)
const stats = ref([
  { label: '文章数量', value: 0, accent: '#18a058' },
  { label: '分类数量', value: 0, accent: '#2563eb' },
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
      { label: '文章数量', value: getTotal(articleRes.data.data), accent: '#18a058' },
      { label: '分类数量', value: getTotal(categoryRes.data.data), accent: '#2563eb' },
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
    <div class="dash-head">
      <h1>概况</h1>
      <p>查看当前博客的基础数据</p>
    </div>

    <div class="dash-grid" v-if="!loading">
      <div v-for="(item, i) in stats" :key="item.label" class="dash-card" :style="{ '--i': i }">
        <span class="dash-bar" :style="{ background: item.accent }"></span>
        <div class="dash-label">{{ item.label }}</div>
        <div class="dash-value">{{ item.value.toLocaleString() }}</div>
      </div>
    </div>

    <div v-else class="dash-loading">加载中...</div>
  </div>
</template>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.dash-head h1 {
  margin-bottom: 4px;
}

.dash-head p {
  color: var(--text-muted);
  font-size: 14px;
}

.dash-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 20px;
}

.dash-card {
  position: relative;
  overflow: hidden;
  padding: 24px;
  border-radius: 14px;
  border: 1px solid var(--card-border);
  background: var(--card-bg);
  box-shadow: var(--card-shadow);
  animation: cardIn 0.4s ease both;
  animation-delay: calc(var(--i) * 80ms + 0.05s);
}

@keyframes cardIn {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dash-bar {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 3px;
}

.dash-label {
  margin-bottom: 12px;
  color: var(--text-muted);
  font-size: 13px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.dash-value {
  color: var(--text-h);
  font-size: 36px;
  font-weight: 700;
  letter-spacing: -0.02em;
  line-height: 1;
}

.dash-loading {
  color: var(--text-muted);
}
</style>
