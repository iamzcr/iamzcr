import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

export const articleApi = {
  list: (params?: any) => api.get('/articles', { params }),
  get: (id: number) => api.get(`/articles/${id}`),
  create: (data: any) => api.post('/articles', data),
  update: (id: number, data: any) => api.put(`/articles/${id}`, data),
  delete: (id: number) => api.delete(`/articles/${id}`)
}

export const categoryApi = {
  list: () => api.get('/categories'),
  get: (id: number) => api.get(`/categories/${id}`),
  create: (data: any) => api.post('/categories', data),
  update: (id: number, data: any) => api.put(`/categories/${id}`, data),
  delete: (id: number) => api.delete(`/categories/${id}`)
}

export const commentApi = {
  list: () => api.get('/comments'),
  get: (id: number) => api.get(`/comments/${id}`),
  create: (data: any) => api.post('/comments', data),
  update: (id: number, data: any) => api.put(`/comments/${id}`, data),
  delete: (id: number) => api.delete(`/comments/${id}`)
}

export const menuApi = {
  list: () => api.get('/menus'),
  get: (id: number) => api.get(`/menus/${id}`),
  create: (data: any) => api.post('/menus', data),
  update: (id: number, data: any) => api.put(`/menus/${id}`, data),
  delete: (id: number) => api.delete(`/menus/${id}`)
}

export const tagsApi = {
  list: () => api.get('/tags'),
  get: (id: number) => api.get(`/tags/${id}`),
  create: (data: any) => api.post('/tags', data),
  update: (id: number, data: any) => api.put(`/tags/${id}`, data),
  delete: (id: number) => api.delete(`/tags/${id}`)
}

export const directoryApi = {
  list: () => api.get('/directories'),
  get: (id: number) => api.get(`/directories/${id}`),
  create: (data: any) => api.post('/directories', data),
  update: (id: number, data: any) => api.put(`/directories/${id}`, data),
  delete: (id: number) => api.delete(`/directories/${id}`)
}

export const websiteApi = {
  get: () => api.get('/website'),
  update: (data: any) => api.put('/website', data)
}