<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NMenu, NDropdown } from 'naive-ui'
import { menuApi, authApi } from './api'

function getCollapsed(): boolean {
  const stored = localStorage.getItem('sidebar_collapsed')
  return stored !== null ? stored === 'true' : false
}

function setCollapsed(val: boolean) {
  localStorage.setItem('sidebar_collapsed', String(val))
}

const collapsed = ref(getCollapsed())
const route = useRoute()
const router = useRouter()

const adminInfo = ref<any>(null)
const menuOptions = ref<any[]>([])
const expandedKeys = ref<string[]>([])
const isLoginPage = computed(() => route.path === '/login')
const selectedMenuKey = computed(() => {
  if (route.path.startsWith('/articles/edit/')) return '/articles'
  return route.path
})

function toggleCollapsed() {
  collapsed.value = !collapsed.value
  setCollapsed(collapsed.value)
}

async function loadAdminInfo() {
  const token = localStorage.getItem('admin_token')
  if (!token) {
    adminInfo.value = null
    return
  }
  
  try {
    const info = localStorage.getItem('admin_info')
    if (info) {
      adminInfo.value = JSON.parse(info)
    } else {
      const res = await authApi.getInfo()
      adminInfo.value = res.data.data
      localStorage.setItem('admin_info', JSON.stringify(res.data.data))
    }
  } catch (e) {
    console.error(e)
  }
}

async function loadMenus() {
  const token = localStorage.getItem('admin_token')
  if (!token) {
    menuOptions.value = []
    expandedKeys.value = []
    return
  }
  
  try {
    const res = await menuApi.list({ page: 1, page_size: 1000 })
    const menus = res.data.data.list || res.data.data
    
    if (!menus || menus.length === 0) {
      menuOptions.value = getDefaultMenus()
      expandedKeys.value = ['menu-articles']
      return
    }

    const normalizedMenus = menus
      .filter((menu: any) => menu.status === 1)
      .sort((a: any, b: any) => (b.weight || 0) - (a.weight || 0))

    const menuMap = new Map<number, any>()
    const childrenMap = new Map<number, any[]>()

    normalizedMenus.forEach((menu: any) => {
      menuMap.set(menu.id, menu)
      const parent = Number(menu.parent || 0)
      const siblings = childrenMap.get(parent) || []
      siblings.push(menu)
      childrenMap.set(parent, siblings)
    })

    function toMenuOption(menu: any): any {
      const childMenus = (childrenMap.get(menu.id) || []).map(toMenuOption)
      const path = mapMenuPath(menu.url || String(menu.id))
      return childMenus.length > 0
        ? { label: menu.name, key: `menu-${menu.id}`, children: childMenus }
        : { label: menu.name, key: path }
    }

    const roots = (childrenMap.get(0) || normalizedMenus.filter((menu: any) => !menuMap.has(Number(menu.parent))))
      .map(toMenuOption)

    menuOptions.value = roots.length > 0 ? roots : getDefaultMenus()
    expandedKeys.value = ['menu-articles']
  } catch (e) {
    console.error('Load menus error:', e)
    menuOptions.value = getDefaultMenus()
    expandedKeys.value = ['menu-articles']
  }
}

function getDefaultMenus() {
  return [
    { label: '概况', key: '/' },
    { label: '文章管理', key: 'menu-articles', children: [
      { label: '文章列表', key: '/articles' },
      { label: '新建文章', key: '/articles/new' }
    ]},
    { label: '分类管理', key: '/categories' },
    { label: '目录管理', key: '/directories' },
    { label: '标签管理', key: '/tags' },
    { label: '评论管理', key: '/comments' },
    { label: '菜单管理', key: '/menus' },
    { label: '网站设置', key: '/website' },
    { label: '平台管理', key: '/platforms' },
    { label: '附件同步', key: '/attach_media' }
  ]
}

function mapMenuPath(key: string) {
  const pathMap: Record<string, string> = {
    '/admin/index': '/',
    '/admin/article/list': '/articles',
    '/admin/article/add': '/articles/new',
    '/admin/category/list': '/categories',
    '/admin/directory/list': '/directories',
    '/admin/tags/list': '/tags',
    '/admin/comment/list': '/comments',
    '/admin/menu/list': '/menus',
    '/admin/website/list': '/website',
    '/admin/admin/list': '/admins',
    '/admin/admin_group/list': '/admin_groups',
    '/admin/admin/password': '/password',
    '/admin/attach/list': '/attaches',
    '/admin/lang/list': '/langs',
    '/admin/log/list': '/logs',
    '/admin/message/list': '/messages',
    '/admin/permit/list': '/permits',
    '/admin/read/list': '/reads',
    '/admin/platform/list': '/platforms',
    '/admin/attach_media/list': '/attach_media',
  }

  let path = pathMap[key] || key
  if (!path.startsWith('/')) {
    path = '/' + path
  }
  return path
}

function handleMenuSelect(key: string) {
  if (!key || !key.startsWith('/')) return
  router.push(mapMenuPath(key))
}

const userOptions = [
  { label: '修改密码', key: 'password' },
  { type: 'divider' as const, key: 'd1' },
  { label: '退出登录', key: 'logout' }
]

function handleUserSelect(key: string) {
  if (key === 'password') {
    router.push('/password')
  } else if (key === 'logout') {
    authApi.logout().then(() => {
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_info')
      router.push('/login')
    })
  }
}

watch(
  isLoginPage,
  (loginPage) => {
    if (loginPage) {
      adminInfo.value = null
      menuOptions.value = []
      expandedKeys.value = []
      return
    }

    loadAdminInfo()
    loadMenus()
  },
  { immediate: true }
)
</script>

<template>
  <n-message-provider>
    <div id="app-container">
      <aside class="sidebar" :class="{ collapsed }" v-if="!isLoginPage">
        <div class="sb-logo">
          <div class="sb-logo-icon">☰</div>
          <span v-if="!collapsed" class="sb-logo-text">Blog Admin</span>
        </div>
        <div class="sb-menu">
          <n-menu
            :collapsed="collapsed"
            :collapsed-width="64"
            :collapsed-icon-size="18"
            :options="menuOptions"
            :value="selectedMenuKey"
            :expanded-keys="expandedKeys"
            inverted
            @update:expanded-keys="expandedKeys = $event"
            @update:value="handleMenuSelect"
          />
        </div>
        <div v-if="!collapsed" class="sb-footer">
          <span class="sb-version">v1.0</span>
        </div>
      </aside>

      <div class="main-area" v-if="!isLoginPage">
        <header class="topbar">
          <button class="topbar-trigger" @click="toggleCollapsed" :title="collapsed ? '展开' : '收起'">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="3" y1="6" x2="21" y2="6"></line>
              <line x1="3" y1="12" x2="21" y2="12"></line>
              <line x1="3" y1="18" x2="21" y2="18"></line>
            </svg>
          </button>

          <div class="topbar-right">
            <n-dropdown :options="userOptions" @select="handleUserSelect" placement="bottom-end">
              <button class="topbar-user">
                <span class="topbar-avatar">{{ (adminInfo?.name || 'A')[0] }}</span>
                <span class="topbar-username">{{ adminInfo?.name || '管理员' }}</span>
                <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
                  <polyline points="6 9 12 15 18 9"></polyline>
                </svg>
              </button>
            </n-dropdown>
          </div>
        </header>

        <main class="page-content">
          <router-view v-slot="{ Component }">
            <transition name="page" mode="out-in">
              <component :is="Component" :key="$route.fullPath" />
            </transition>
          </router-view>
        </main>
      </div>

      <div v-if="isLoginPage" class="login-wrapper">
        <router-view />
      </div>
    </div>
  </n-message-provider>
</template>

<style>
*,
*::before,
*::after {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body, html {
  height: 100%;
}

#app-container {
  height: 100vh;
  display: flex;
  flex-direction: row;
  background: var(--page-bg);
}

/* Sidebar */
.sidebar {
  width: 240px;
  flex-shrink: 0;
  background: var(--sidebar-bg);
  display: flex;
  flex-direction: column;
  transition: width 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.sidebar.collapsed {
  width: 64px;
}

.sb-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 56px;
  padding: 0 18px;
  border-bottom: 1px solid var(--sidebar-border);
  flex-shrink: 0;
  overflow: hidden;
  white-space: nowrap;
}

.sidebar.collapsed .sb-logo {
  padding: 0;
  justify-content: center;
}

.sb-logo-icon {
  font-size: 18px;
  color: var(--accent);
  flex-shrink: 0;
  width: 24px;
  text-align: center;
}

.sb-logo-text {
  font-size: 16px;
  font-weight: 700;
  color: var(--sidebar-text-active);
  letter-spacing: 0.02em;
}

.sb-menu {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.sb-footer {
  padding: 10px 18px;
  border-top: 1px solid var(--sidebar-border);
  flex-shrink: 0;
}

.sb-version {
  font-size: 11px;
  color: var(--sidebar-text);
  opacity: 0.5;
}

/* Main area */
.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

/* Top bar */
.topbar {
  height: 56px;
  background: var(--header-bg);
  border-bottom: 1px solid var(--header-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  flex-shrink: 0;
}

.topbar-trigger {
  width: 34px;
  height: 34px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  background: var(--card-bg);
  color: var(--text);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.topbar-trigger:hover {
  background: var(--page-bg);
  color: var(--accent);
  border-color: var(--accent);
}

.topbar-right {
  display: flex;
  align-items: center;
}

.topbar-user {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px 6px 6px;
  border-radius: 20px;
  border: 1px solid transparent;
  background: transparent;
  cursor: pointer;
  transition: all 0.15s;
  font-size: 13px;
  color: var(--text);
}

.topbar-user:hover {
  background: var(--page-bg);
  border-color: var(--border);
}

.topbar-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--accent);
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.topbar-username {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-h);
}

/* Page content */
.page-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

/* Login wrapper */
.login-wrapper {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background: var(--page-bg);
}

/* Page transition */
.page-enter-active,
.page-leave-active {
  transition: opacity 0.22s ease, transform 0.22s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

.page-leave-to {
  opacity: 0;
}

/* Shared page layout for table/list views */
.page-wrap {
  background: var(--card-bg);
  border: 1px solid var(--card-border);
  border-radius: 12px;
  box-shadow: var(--card-shadow);
  padding: 20px;
}

.page-toolbar {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.page-toolbar:empty {
  display: none;
}
</style>
