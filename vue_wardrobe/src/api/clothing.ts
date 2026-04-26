import { apiGet, apiPost, apiDelete } from './client'
import type { ClothingItem } from '@/types/clothing'

export const clothingApi = {
  list: (categoryId?: number) => {
    const params: Record<string, string> = {}
    if (categoryId) params.category_id = String(categoryId)
    return apiGet<ClothingItem[]>('/api/clothing', params)
  },

  getById: (id: string) => apiGet<ClothingItem>(`/api/clothing/${id}`),

  upload: (formData: FormData) => apiPost<ClothingItem>('/api/clothing', formData),

  delete: (id: string) => apiDelete(`/api/clothing/${id}`),
}
