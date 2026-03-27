<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const collapsed = ref(false)
const route = useRoute()
const router = useRouter()

const menuOptions = [
  { label: '文章管理', key: 'articles', icon: '📝' },
  { label: '分类管理', key: 'categories', icon: '📁' },
  { label: '目录管理', key: 'directories', icon: '📂' },
  { label: '评论管理', key: 'comments', icon: '💬' },
  { label: '标签管理', key: 'tags', icon: '🏷️' },
  { label: '菜单管理', key: 'menus', icon: '🔗' },
  { label: '网站设置', key: 'website', icon: '⚙️' }
]

function handleMenuSelect(key: string) {
  router.push(`/${key}`)
}
</script>

<template>
  <n-layout has-sider style="height: 100vh">
    <n-layout-sider
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
      <n-layout-header bordered style="padding: 16px">
        <h2>个人博客管理系统</h2>
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