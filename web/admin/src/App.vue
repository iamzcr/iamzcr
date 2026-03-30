<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NMenu, NDropdown, NButton } from 'naive-ui'
import { menuApi, authApi } from './api'

const collapsed = ref(false)
const route = useRoute()
const router = useRouter()

const adminInfo = ref<any>(null)
const menuOptions = ref<any[]>([])

async function loadAdminInfo() {
  const token = localStorage.getItem('admin_token')
  if (!token) return
  
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
  if (!token) return
  
  try {
    const res = await menuApi.list()
    const menus = res.data.data
    
    if (!menus || menus.length === 0) {
      menuOptions.value = getDefaultMenus()
      return
    }
    
    const menuMap = new Map<number, any>()
    const rootMenus: any[] = []
    
    menus.forEach((m: any) => {
      if (m.status !== 1) return
      menuMap.set(m.id, {
        label: m.name,
        key: m.url || String(m.id)
      })
    })
    
    menus.forEach((m: any) => {
      if (m.status !== 1) return
      const menu = menuMap.get(m.id)
      if (!menu) return
      
      if (m.parent === 0) {
        rootMenus.push(menu)
      } else {
        const parent = menuMap.get(m.parent)
        if (parent) {
          if (!parent.children) {
            parent.children = []
          }
          parent.children.push(menu)
        }
      }
    })
    
    menuOptions.value = rootMenus.length > 0 ? rootMenus : getDefaultMenus()
  } catch (e) {
    console.error('Load menus error:', e)
    menuOptions.value = getDefaultMenus()
  }
}

function getDefaultMenus() {
  return [
    { label: '仪表盘', key: '/' },
    { label: '文章管理', key: 'articles', children: [
      { label: '文章列表', key: 'articles' },
      { label: '新建文章', key: 'articles/new' }
    ]},
    { label: '分类管理', key: 'categories' },
    { label: '目录管理', key: 'directories' },
    { label: '标签管理', key: 'tags' },
    { label: '评论管理', key: 'comments' },
    { label: '菜单管理', key: 'menus' },
    { label: '网站设置', key: 'website' }
  ]
}

function handleMenuSelect(key: string) {
  console.log('Menu selected:', key)
  
  const pathMap: Record<string, string> = {
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
  router.push(path)
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

onMounted(() => {
  loadAdminInfo()
  loadMenus()
})
</script>

<template>
  <n-message-provider>
    <div id="app-container">
      <div class="sidebar" v-if="route.path !== '/login'">
        <div class="logo">{{ collapsed ? 'B' : 'Blog Admin' }}</div>
        <n-menu
          :collapsed="collapsed"
          :options="menuOptions"
          @update:value="handleMenuSelect"
        />
      </div>
      <div class="main-content" v-if="route.path !== '/login'">
        <div class="header">
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
      <div v-if="route.path === '/login'" class="login-page">
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
  background: #1a365d;
  border-right: 1px solid #2c5282;
  overflow-y: auto;
  flex-shrink: 0;
}

.sidebar :deep(.n-menu) {
  background: transparent;
  color: #e2e8f0;
}

.sidebar :deep(.n-menu-item-content) {
  background: transparent;
}

.sidebar :deep(.n-menu-item-content:hover) {
  background: #2c5282;
}

.sidebar :deep(.n-menu-item-content--selected) {
  background: #2c5282;
}

.sidebar :deep(.n-menu-item-content--selected::before) {
  background: #3182ce;
}

.sidebar .logo {
  background: #1a365d;
  color: #fff;
  border-bottom: 1px solid #2c5282;
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

.sidebar .logo {
  padding: 16px;
  font-size: 18px;
  font-weight: bold;
  text-align: center;
  border-bottom: 1px solid #eee;
}

.header {
  height: 60px;
  background: #2c5282;
  border-bottom: 1px solid #1a365d;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 0 24px;
  flex-shrink: 0;
}

.header .username {
  color: #e2e8f0;
}

.user-info {
  cursor: pointer;
}

.username {
  color: #333;
}

.content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  background: #f5f7fa;
}
</style>
