<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { NInput, NButton } from 'naive-ui'
import { categoryApi, tagsApi } from './api'

const router = useRouter()

const categories = ref<any[]>([])
const tags = ref<any[]>([])
const searchKeyword = ref('')
const mobileMenuOpen = ref(false)

function toggleMobileMenu() {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

function closeMobileMenu() {
  mobileMenuOpen.value = false
}

function onSearch() {
  closeMobileMenu()
  router.push('/')
}

// Close mobile menu on route change
watch(() => router.currentRoute.value.fullPath, () => {
  closeMobileMenu()
})

// Close mobile menu on Escape
function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && mobileMenuOpen.value) {
    closeMobileMenu()
  }
}

async function loadMenuData() {
  const [catRes, tagRes] = await Promise.all([
    categoryApi.list(),
    tagsApi.list()
  ])
  categories.value = catRes.data.data
  tags.value = tagRes.data.data
}

onMounted(() => {
  loadMenuData()
  document.addEventListener('keydown', onKeydown)
})
</script>

<template>
  <div class="app-container">
    <nav class="navbar">
      <div class="nav-content">
        <div class="nav-left">
          <router-link to="/" class="logo" @click="closeMobileMenu">
            <span class="logo-icon">✍️</span>
            <span class="logo-text">堆栈人生</span>
          </router-link>
          <div class="nav-menu">
            <router-link to="/" class="nav-link" @click="closeMobileMenu">首页</router-link>
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
            @keyup.enter="onSearch"
          />
          <n-button type="primary" @click="onSearch" class="search-btn">搜索</n-button>
        </div>
        <button 
          class="hamburger" 
          :class="{ active: mobileMenuOpen }"
          @click="toggleMobileMenu"
          aria-label="菜单"
        >
          <span></span>
          <span></span>
          <span></span>
        </button>
      </div>

      <!-- Mobile menu overlay -->
      <transition name="mobile-menu">
        <div v-if="mobileMenuOpen" class="mobile-overlay" @click.self="closeMobileMenu">
          <div class="mobile-menu" @click.stop>
            <div class="mobile-search">
              <n-input 
                v-model:value="searchKeyword" 
                placeholder="搜索文章..." 
                size="large"
                clearable
                @keyup.enter="onSearch"
              />
              <n-button type="primary" size="large" block @click="onSearch" class="mobile-search-btn">
                搜索
              </n-button>
            </div>

            <div class="mobile-links">
              <router-link 
                to="/" 
                class="mobile-link" 
                :style="{ '--i': 0 }"
                @click="closeMobileMenu"
              >
                <span class="mobile-link-icon">⌂</span>
                首页
              </router-link>
              <div class="mobile-section-label" :style="{ '--i': 1 }">分类浏览</div>
              <router-link
                v-for="(cat, idx) in categories"
                :key="cat.id"
                :to="`/category/${cat.id}`"
                class="mobile-link mobile-link-sub"
                :style="{ '--i': idx + 2 }"
                @click="closeMobileMenu"
              >
                {{ cat.name }}
              </router-link>
            </div>

            <div class="mobile-footer">
              <span>© 2024 堆栈人生</span>
            </div>
          </div>
        </div>
      </transition>
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

/* Hamburger */
.hamburger {
  display: none;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 5px;
  width: 40px;
  height: 40px;
  padding: 0;
  background: transparent;
  border: none;
  cursor: pointer;
  z-index: 1001;
  position: relative;
}

.hamburger span {
  display: block;
  width: 22px;
  height: 2px;
  background: var(--text);
  border-radius: 2px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: center;
}

.hamburger.active span:nth-child(1) {
  transform: translateY(7px) rotate(45deg);
}

.hamburger.active span:nth-child(2) {
  opacity: 0;
  transform: scaleX(0);
}

.hamburger.active span:nth-child(3) {
  transform: translateY(-7px) rotate(-45deg);
}

/* Mobile overlay */
.mobile-overlay {
  display: none;
}

/* Mobile menu */
.mobile-menu {
  background: var(--card-bg);
  border-bottom: 1px solid var(--border);
  box-shadow: var(--shadow-lg);
  max-height: calc(100vh - 64px);
  overflow-y: auto;
}

.mobile-search {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  border-bottom: 1px solid var(--border);
}

.mobile-search-btn {
  margin-top: 4px;
}

.mobile-links {
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mobile-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  color: var(--text);
  text-decoration: none;
  font-size: 15px;
  border-radius: var(--radius-sm);
  transition: all 0.2s;
  animation: mobileLinkIn 0.35s ease both;
  animation-delay: calc(var(--i) * 40ms + 0.05s);
}

@keyframes mobileLinkIn {
  from {
    opacity: 0;
    transform: translateX(-12px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.mobile-link:hover {
  background: var(--accent-bg);
  color: var(--accent);
}

.mobile-link-icon {
  font-size: 16px;
  opacity: 0.5;
}

.mobile-link-sub {
  padding-left: 32px;
  font-size: 14px;
}

.mobile-section-label {
  padding: 16px 16px 6px;
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  animation: mobileLinkIn 0.35s ease both;
  animation-delay: calc(var(--i) * 40ms + 0.05s);
}

.mobile-footer {
  padding: 20px;
  display: flex;
  justify-content: center;
  gap: 16px;
  border-top: 1px solid var(--border);
  font-size: 13px;
  color: var(--text-muted);
}

.mobile-footer a {
  color: var(--text-muted);
  text-decoration: none;
  transition: color 0.2s;
}

.mobile-footer a:hover {
  color: var(--accent);
}

/* Mobile menu transition */
.mobile-menu-enter-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.mobile-menu-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.mobile-menu-enter-from,
.mobile-menu-leave-to {
  opacity: 0;
}

.mobile-menu-enter-from .mobile-menu,
.mobile-menu-leave-to .mobile-menu {
  transform: translateY(-12px);
}

/* --- Responsive --- */

@media (max-width: 768px) {
  .nav-content {
    padding: 0 20px;
  }

  .nav-menu,
  .nav-right {
    display: none;
  }

  .hamburger {
    display: flex;
  }

  .mobile-overlay {
    display: block;
    position: fixed;
    inset: 64px 0 0 0;
    background: rgba(0, 0, 0, 0.3);
    z-index: 999;
    backdrop-filter: blur(2px);
  }

  .mobile-menu {
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  }

  .logo-text {
    font-size: 20px;
  }
}

@media (max-width: 1024px) {
  .nav-content {
    padding: 0 24px;
  }

  .search-input {
    width: 180px;
  }

  .nav-menu {
    gap: 16px;
  }

  .nav-left {
    gap: 20px;
  }
}

@media (max-width: 600px) {
  .main-wrapper {
    padding: 20px 16px;
  }
  
  .nav-content {
    padding: 0 16px;
    height: 56px;
  }

  .logo-icon {
    font-size: 24px;
  }

  .logo-text {
    font-size: 18px;
  }

  .mobile-overlay {
    inset: 56px 0 0 0;
  }

  .mobile-menu {
    max-height: calc(100vh - 56px);
  }
}

.main-wrapper {
  flex: 1;
  max-width: 1440px;
  width: 100%;
  margin: 0 auto;
  padding: 28px 32px;
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
