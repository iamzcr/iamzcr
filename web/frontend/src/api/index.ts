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
	list: (params?: any) => api.get('/categories', { params })
}

export const directoryApi = {
  list: (params?: any) => api.get('/directories', { params })
}

export const tagsApi = {
	list: (params?: any) => api.get('/tags', { params })
}
