import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import naive from 'naive-ui'
import VMdEditor from '@kangc/v-md-editor'
import '@kangc/v-md-editor/lib/style/base-editor.css'
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js'
import '@kangc/v-md-editor/lib/theme/style/github.css'
import App from './App.vue'
import './style.css'
import Login from './views/Login.vue'
import Index from './views/Index.vue'
import ArticleList from './views/ArticleList.vue'
import ArticleEdit from './views/ArticleEdit.vue'
import CategoryList from './views/CategoryList.vue'
import DirectoryList from './views/DirectoryList.vue'
import CommentList from './views/CommentList.vue'
import TagsList from './views/TagsList.vue'
import MenuList from './views/MenuList.vue'
import WebsiteSettings from './views/WebsiteSettings.vue'
import AttachList from './views/AttachList.vue'
import LangList from './views/LangList.vue'
import LogList from './views/LogList.vue'
import MessageList from './views/MessageList.vue'
import PermitList from './views/PermitList.vue'
import ReadList from './views/ReadList.vue'
import AdminList from './views/AdminList.vue'
import AdminGroupList from './views/AdminGroupList.vue'
import PasswordChange from './views/PasswordChange.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', component: Login },
    { path: '/', component: Index },
    { path: '/articles', component: ArticleList },
    { path: '/articles/new', component: ArticleEdit },
    { path: '/articles/edit/:id', component: ArticleEdit },
    { path: '/categories', component: CategoryList },
    { path: '/directories', component: DirectoryList },
    { path: '/comments', component: CommentList },
    { path: '/tags', component: TagsList },
    { path: '/menus', component: MenuList },
    { path: '/website', component: WebsiteSettings },
    { path: '/attaches', component: AttachList },
    { path: '/langs', component: LangList },
    { path: '/logs', component: LogList },
    { path: '/messages', component: MessageList },
    { path: '/permits', component: PermitList },
    { path: '/reads', component: ReadList },
    { path: '/admins', component: AdminList },
    { path: '/admin_groups', component: AdminGroupList },
    { path: '/password', component: PasswordChange },
  ]
})

router.beforeEach((to) => {
  console.log('Navigating to:', to.path)
  const token = localStorage.getItem('admin_token')
  if (!token && to.path !== '/login') {
    return '/login'
  }
  if (token && to.path === '/login') {
    return '/'
  }
})

VMdEditor.use(githubTheme)

const app = createApp(App)
app.use(router)
app.use(naive)
app.use(VMdEditor)
app.mount('#app')
