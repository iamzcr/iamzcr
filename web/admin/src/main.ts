import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import naive from 'naive-ui'
import App from './App.vue'
import './style.css'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/articles' },
    { path: '/articles', component: () => import('./views/ArticleList.vue') },
    { path: '/articles/new', component: () => import('./views/ArticleEdit.vue') },
    { path: '/articles/edit/:id', component: () => import('./views/ArticleEdit.vue') },
    { path: '/categories', component: () => import('./views/CategoryList.vue') },
    { path: '/directories', component: () => import('./views/DirectoryList.vue') },
    { path: '/comments', component: () => import('./views/CommentList.vue') },
    { path: '/tags', component: () => import('./views/TagsList.vue') },
    { path: '/menus', component: () => import('./views/MenuList.vue') },
    { path: '/website', component: () => import('./views/WebsiteSettings.vue') },
  ]
})

const app = createApp(App)
app.use(router)
app.use(naive)
app.mount('#app')