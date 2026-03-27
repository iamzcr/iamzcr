import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import naive from 'naive-ui'
import App from './App.vue'
import './style.css'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: () => import('./views/Home.vue') },
    { path: '/article/:id', name: 'article', component: () => import('./views/Article.vue') },
    { path: '/category/:id', name: 'category', component: () => import('./views/Home.vue') },
    { path: '/tag/:id', name: 'tag', component: () => import('./views/Home.vue') }
  ]
})

const app = createApp(App)
app.use(router)
app.use(naive)
app.mount('#app')