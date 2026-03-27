import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

export const articleApi = {
  list: (params?: any) => api.get('/articles', { params }),
  get: (id: number) => api.get(`/articles/${id}`)
}

export const categoryApi = {
  list: () => api.get('/categories')
}

export const directoryApi = {
  list: () => api.get('/directories')
}

export const tagsApi = {
  list: () => api.get('/tags')
}