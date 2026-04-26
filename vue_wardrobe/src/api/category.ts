import { apiGet } from './client'
import type { Category } from '@/types/category'

export const categoryApi = {
  getAll: () => apiGet<Category[]>('/api/categories'),
}
