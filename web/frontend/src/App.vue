<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NInput, NButton } from 'naive-ui'
import { categoryApi, tagsApi } from './api'

const categories = ref<any[]>([])
const tags = ref<any[]>([])
const searchKeyword = ref('')

const menuOptions = [
  { label: '首页', key: '/' },
]

function openAdmin() {
  window.open('http://localhost:3001', '_blank')
}

async function loadMenuData() {
  const [catRes, tagRes] = await Promise.all([
    categoryApi.list(),
    tagsApi.list()
  ])
  categories.value = catRes.data.data
  tags.value = tagRes.data.data
  
  // Add categories to menu
  categories.value.forEach(cat => {
    menuOptions.push({ label: cat.name, key: `/category/${cat.id}` })
  })
}

onMounted(loadMenuData)
</script>

<template>
  <div class="app-container">
    <nav class="navbar">
      <div class="nav-content">
        <div class="nav-left">
          <router-link to="/" class="logo">
            <span class="logo-icon">✍️</span>
            <span class="logo-text">堆栈人生</span>
          </router-link>
          <div class="nav-menu">
            <router-link to="/" class="nav-link">首页</router-link>
            <div class="nav-dropdown">
              <span class="nav-link">分类 ▾</span>
              <div class="dropdown-content">
                <router-link 
                  v-for="cat in categories" 
                  :key="cat.id" 
                  :to="`/category/${cat.id}`"
                  class="dropdown-item"
                >
                  {{ cat.name }}
                </router-link>
              </div>
            </div>
          </div>
        </div>
        <div class="nav-right">
          <n-input 
            v-model:value="searchKeyword" 
            placeholder="搜索文章..." 
            class="search-input"
            clearable
          />
          <n-button type="primary" @click="$router.push('/')">搜索</n-button>
        </div>
      </div>
    </nav>
    
    <main class="main-wrapper">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" :key="$route.fullPath" />
        </transition>
      </router-view>
    </main>
    
    <footer class="site-footer">
      <div class="footer-content">
          <div class="footer-links">
          <span>© 2024 堆栈人生</span>
          <a href="javascript:void(0)" @click="$router.push('/')">首页</a>
          <span class="divider">|</span>
          <a href="javascript:void(0)" @click="openAdmin">管理后台</a>
        </div>
        <p class="footer-copy">Powered by Go + Gin + Vue</p>
      </div>
    </footer>
  </div>
</template>

<style>
*,
*::before,
*::after {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: var(--sans);
  background: var(--bg);
  color: var(--text);
}

.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.navbar {
  background: var(--card-bg);
  box-shadow: var(--shadow-sm);
  position: sticky;
  top: 0;
  z-index: 1000;
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border);
}

.nav-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 32px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-left {
  display: flex;
  align-items: center;
  gap: 32px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
}

.logo-icon {
  font-size: 28px;
}

.logo-text {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-h);
  font-family: var(--serif);
  letter-spacing: 0.02em;
}

.nav-menu {
  display: flex;
  align-items: center;
  gap: 24px;
}

.nav-link {
  font-size: 15px;
  color: var(--text);
  text-decoration: none;
  padding: 8px 0;
  transition: color 0.2s;
}

.nav-link:hover {
  color: var(--accent);
}

.nav-dropdown {
  position: relative;
}

.dropdown-content {
  display: none;
  position: absolute;
  top: 100%;
  left: 0;
  background: var(--card-bg);
  border-radius: var(--radius);
  box-shadow: var(--shadow-lg);
  border: 1px solid var(--border);
  min-width: 160px;
  padding: 6px 0;
  opacity: 0;
  transform: translateY(-4px);
  transition: opacity 0.2s, transform 0.2s;
}

.nav-dropdown:hover .dropdown-content {
  display: block;
  opacity: 1;
  transform: translateY(0);
}

.dropdown-item {
  display: block;
  padding: 10px 20px;
  color: var(--text);
  text-decoration: none;
  transition: all 0.2s;
  font-size: 14px;
}

.dropdown-item:hover {
  background: var(--accent-bg);
  color: var(--accent);
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search-input {
  width: 240px;
}

.main-wrapper {
  flex: 1;
  max-width: 1440px;
  width: 100%;
  margin: 0 auto;
  padding: 28px 32px;
}

@media (max-width: 600px) {
  .main-wrapper {
    padding: 20px 16px;
  }
  
  .nav-content {
    padding: 0 16px;
  }
}

.site-footer {
  background: var(--card-bg);
  border-top: 1px solid var(--border);
  padding: 24px;
}

.footer-content {
  max-width: 1400px;
  margin: 0 auto;
  text-align: center;
}

.footer-links {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-bottom: 8px;
  color: var(--text-muted);
  font-size: 14px;
}

.footer-links a {
  color: var(--text-muted);
  text-decoration: none;
  font-size: 14px;
  transition: color 0.2s;
}

.footer-links a:hover {
  color: var(--accent);
}

.footer-links .divider {
  color: var(--border);
}

.footer-copy {
  color: var(--text-muted);
  font-size: 13px;
}

/* Page transition */
.page-enter-active,
.page-leave-active {
  transition: opacity 0.28s ease, transform 0.28s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(12px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>