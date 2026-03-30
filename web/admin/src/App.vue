<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NButton, NDropdown, NAvatar } from 'naive-ui'
import { menuApi, authApi } from './api'

const collapsed = ref(false)
const route = useRoute()
const router = useRouter()

const adminInfo = ref<any>(null)
const menuOptions = ref<any[]>([])

function renderMenuIcon(icon: string) {
  return () => icon
}

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
  try {
    const res = await menuApi.list()
    const menus = res.data.data
    
    const menuMap = new Map<number, any>()
    const rootMenus: any[] = []
    
    menus.forEach((m: any) => {
      if (m.status !== 1) return
      menuMap.set(m.id, {
        label: m.name,
        key: m.url || m.mark,
        icon: renderMenuIcon(m.icon || '📄')
      })
    })
    
    menus.forEach((m: any) => {
      if (m.status !== 1) return
      const menu = menuMap.get(m.id)
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
    
    menuOptions.value = rootMenus
  } catch (e) {
    console.error(e)
  }
}

function handleMenuSelect(key: string) {
  router.push(`/${key}`)
}

const userOptions = [
  { label: '个人设置', key: 'profile' },
  { type: 'divider', key: 'd1' },
  { label: '退出登录', key: 'logout' }
]

function handleUserSelect(key: string) {
  if (key === 'logout') {
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
  <n-layout has-sider style="height: 100vh">
    <n-layout-sider
      v-if="!route.path.includes('/login')"
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="240"
      show-trigger
      @collapse="collapsed = true"
      @expand="collapsed = false"
      :native-scrollbar="false"
    >
      <div style="padding: 16px; font-size: 18px; font-weight: bold; text-align: center">
        {{ collapsed ? 'B' : 'Blog Admin' }}
      </div>
      <n-menu
        :collapsed="collapsed"
        :collapsed-width="64"
        :options="menuOptions"
        @update:value="handleMenuSelect"
        :value="route.name as string"
      />
    </n-layout-sider>
    <n-layout>
      <n-layout-header v-if="!route.path.includes('/login')" bordered style="padding: 12px 24px; display: flex; justify-content: flex-end; align-items: center">
        <n-dropdown :options="userOptions" @select="handleUserSelect">
          <n-button text>
            <n-avatar round size="small" style="margin-right: 8px">
              {{ adminInfo?.name?.charAt(0) || 'A' }}
            </n-avatar>
            {{ adminInfo?.name || '管理员' }}
          </n-button>
        </n-dropdown>
      </n-layout-header>
      <n-layout-content content-style="padding: 24px">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<style>
body {
  margin: 0;
  padding: 0;
}
</style>