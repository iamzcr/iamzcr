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
            <span class="logo-text">阿呆不是呆哦</span>
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
      <router-view />
    </main>
    
    <footer class="site-footer">
      <div class="footer-content">
        <div class="footer-links">
          <span>© 2024 阿呆不是呆哦</span>
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
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background: #f0f2f5;
  color: #333;
}

.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.navbar {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.nav-content {
  max-width: 100%;
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
  color: #333;
}

.nav-menu {
  display: flex;
  align-items: center;
  gap: 24px;
}

.nav-link {
  font-size: 15px;
  color: #555;
  text-decoration: none;
  padding: 8px 0;
  transition: color 0.2s;
}

.nav-link:hover {
  color: #18a058;
}

.nav-dropdown {
  position: relative;
}

.dropdown-content {
  display: none;
  position: absolute;
  top: 100%;
  left: 0;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  min-width: 160px;
  padding: 8px 0;
}

.nav-dropdown:hover .dropdown-content {
  display: block;
}

.dropdown-item {
  display: block;
  padding: 10px 20px;
  color: #555;
  text-decoration: none;
  transition: all 0.2s;
}

.dropdown-item:hover {
  background: #f5f7fa;
  color: #18a058;
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
  padding: 24px;
}

.site-footer {
  background: #fff;
  border-top: 1px solid #eee;
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
}

.footer-links a {
  color: #666;
  text-decoration: none;
  font-size: 14px;
}

.footer-links a:hover {
  color: #18a058;
}

.footer-links .divider {
  color: #ddd;
}

.footer-copy {
  color: #999;
  font-size: 13px;
}
</style>