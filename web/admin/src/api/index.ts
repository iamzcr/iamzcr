import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('admin_token')
  if (token && config.url !== '/login') {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401 && error.config?.url !== '/login') {
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_info')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const authApi = {
  login: (data: any) => api.post('/login', data),
  logout: () => api.post('/logout'),
  getInfo: () => api.get('/admin/info')
}

export const articleApi = {
  list: (params?: any) => api.get('/articles', { params }),
  get: (id: number) => api.get(`/articles/${id}`),
  create: (data: any) => api.post('/articles', data),
  update: (id: number, data: any) => api.put(`/articles/${id}`, data),
  delete: (id: number) => api.delete(`/articles/${id}`)
}

export const categoryApi = {
  list: (params?: any) => api.get('/categories', { params }),
  get: (id: number) => api.get(`/categories/${id}`),
  create: (data: any) => api.post('/categories', data),
  update: (id: number, data: any) => api.put(`/categories/${id}`, data),
  delete: (id: number) => api.delete(`/categories/${id}`)
}

export const commentApi = {
  list: (params?: any) => api.get('/comments', { params }),
  get: (id: number) => api.get(`/comments/${id}`),
  create: (data: any) => api.post('/comments', data),
  update: (id: number, data: any) => api.put(`/comments/${id}`, data),
  delete: (id: number) => api.delete(`/comments/${id}`)
}

export const menuApi = {
  list: (params?: any) => api.get('/menus', { params }),
  get: (id: number) => api.get(`/menus/${id}`),
  create: (data: any) => api.post('/menus', data),
  update: (id: number, data: any) => api.put(`/menus/${id}`, data),
  delete: (id: number) => api.delete(`/menus/${id}`)
}

export const tagsApi = {
  list: (params?: any) => api.get('/tags', { params }),
  get: (id: number) => api.get(`/tags/${id}`),
  create: (data: any) => api.post('/tags', data),
  update: (id: number, data: any) => api.put(`/tags/${id}`, data),
  delete: (id: number) => api.delete(`/tags/${id}`)
}

export const directoryApi = {
  list: (params?: any) => api.get('/directories', { params }),
  get: (id: number) => api.get(`/directories/${id}`),
  create: (data: any) => api.post('/directories', data),
  update: (id: number, data: any) => api.put(`/directories/${id}`, data),
  delete: (id: number) => api.delete(`/directories/${id}`)
}

export const websiteApi = {
  get: () => api.get('/website'),
  update: (data: any) => api.put('/website', data)
}

export const attachApi = {
  list: (params?: any) => api.get('/attaches', { params }),
  get: (id: number) => api.get(`/attaches/${id}`),
  create: (data: any) => api.post('/attaches', data),
  update: (id: number, data: any) => api.put(`/attaches/${id}`, data),
  delete: (id: number) => api.delete(`/attaches/${id}`)
}

export const langApi = {
  list: (params?: any) => api.get('/langs', { params }),
  get: (id: number) => api.get(`/langs/${id}`),
  create: (data: any) => api.post('/langs', data),
  update: (id: number, data: any) => api.put(`/langs/${id}`, data),
  delete: (id: number) => api.delete(`/langs/${id}`)
}

export const logApi = {
  list: (params?: any) => api.get('/logs', { params }),
  get: (id: number) => api.get(`/logs/${id}`),
  create: (data: any) => api.post('/logs', data),
  delete: (id: number) => api.delete(`/logs/${id}`)
}

export const messageApi = {
  list: (params?: any) => api.get('/messages', { params }),
  get: (id: number) => api.get(`/messages/${id}`),
  create: (data: any) => api.post('/messages', data),
  update: (id: number, data: any) => api.put(`/messages/${id}`, data),
  delete: (id: number) => api.delete(`/messages/${id}`)
}

export const permitApi = {
  list: (params?: any) => api.get('/permits', { params }),
  get: (id: number) => api.get(`/permits/${id}`),
  create: (data: any) => api.post('/permits', data),
  update: (id: number, data: any) => api.put(`/permits/${id}`, data),
  delete: (id: number) => api.delete(`/permits/${id}`)
}

export const readApi = {
  list: (params?: any) => api.get('/reads', { params }),
  get: (id: number) => api.get(`/reads/${id}`),
  create: (data: any) => api.post('/reads', data),
  delete: (id: number) => api.delete(`/reads/${id}`)
}

export const adminApi = {
  list: () => api.get('/admins'),
  get: (id: number) => api.get(`/admins/${id}`),
  create: (data: any) => api.post('/admins', data),
  update: (id: number, data: any) => api.put(`/admins/${id}`, data),
  delete: (id: number) => api.delete(`/admins/${id}`),
  changePassword: (data: any) => api.post('/admin/password', data)
}

export const adminGroupApi = {
  list: () => api.get('/admin_groups'),
  get: (id: number) => api.get(`/admin_groups/${id}`),
  create: (data: any) => api.post('/admin_groups', data),
  update: (id: number, data: any) => api.put(`/admin_groups/${id}`, data),
  delete: (id: number) => api.delete(`/admin_groups/${id}`)
}
