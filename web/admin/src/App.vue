<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NMenu, NDropdown, NButton } from 'naive-ui'
import { menuApi, authApi } from './api'

const collapsed = ref(true)
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
      expandedKeys.value = []
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
    expandedKeys.value = []
  } catch (e) {
    console.error('Load menus error:', e)
    menuOptions.value = getDefaultMenus()
    expandedKeys.value = []
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
    { label: '网站设置', key: '/website' }
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
      <div class="sidebar" :class="{ collapsed }" v-if="!isLoginPage">
        <div class="logo">{{ collapsed ? 'B' : 'Blog Admin' }}</div>
        <n-menu
          :collapsed="collapsed"
          :collapsed-width="64"
          :collapsed-icon-size="18"
          :options="menuOptions"
          :value="selectedMenuKey"
          :expanded-keys="expandedKeys"
          @update:expanded-keys="expandedKeys = $event"
          @update:value="handleMenuSelect"
        />
      </div>
      <div class="main-content" v-if="!isLoginPage">
        <div class="header">
          <n-button text class="collapse-trigger" @click="collapsed = !collapsed" :title="collapsed ? '展开菜单' : '收起菜单'">
            <span class="collapse-icon" aria-hidden="true">{{ collapsed ? '>>' : '<<' }}</span>
          </n-button>
          <div class="user-info">
            <n-dropdown :options="userOptions" @select="handleUserSelect">
              <n-button text>
                <span class="username">{{ adminInfo?.name || '管理员' }} ▼</span>
              </n-button>
            </n-dropdown>
          </div>
        </div>
        <div class="content">
          <router-view :key="route.fullPath" />
        </div>
      </div>
      <div v-if="isLoginPage" class="login-page">
        <router-view />
      </div>
    </div>
  </n-message-provider>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body, html {
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

#app-container {
  height: 100vh;
  display: flex;
  flex-direction: row;
}

#app-container > .sidebar {
  width: 240px;
  background: #ffffff;
  border-right: 1px solid #e2e8f0;
  overflow-y: auto;
  flex-shrink: 0;
  transition: width 0.2s ease, box-shadow 0.2s ease;
}

#app-container > .sidebar.collapsed {
  width: 64px;
}

.sidebar :deep(.n-menu) {
  background: transparent;
  color: #334155;
}

.sidebar :deep(.n-menu-item-content) {
  background: transparent;
  margin: 4px 10px;
  border-radius: 10px;
}

.sidebar :deep(.n-menu-item-content:hover) {
  background: #f1f5f9;
}

.sidebar :deep(.n-menu-item-content--selected) {
  background: #e0ecff;
}

.sidebar :deep(.n-menu-item-content--selected::before) {
  background: #2563eb;
}

.sidebar .logo {
  background: #ffffff;
  color: #0f172a;
  border-bottom: 1px solid #e2e8f0;
  padding: 16px;
  font-size: 18px;
  font-weight: bold;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  transition: padding 0.2s ease, font-size 0.2s ease;
}

.sidebar.collapsed .logo {
  padding: 16px 0;
  font-size: 16px;
}

#app-container > .main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

#app-container > .login-page {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.header {
  height: 60px;
  background: #ffffff;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  flex-shrink: 0;
}

.collapse-trigger {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  color: #475569;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
}

.collapse-trigger:hover {
  background: #eef2ff;
  color: #1d4ed8;
}

.collapse-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  font-size: 14px;
  font-weight: 700;
  letter-spacing: -1px;
}

.header .username {
  color: #334155;
}

.user-info {
  cursor: pointer;
}

.username {
  color: #334155;
}

.content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  background: #f5f7fa;
}
</style>
