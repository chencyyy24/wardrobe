import { apiPost } from './client'

export const uploadApi = {
  uploadCard: (formData: FormData) => apiPost<{ url: string }>('/api/upload/card', formData),
}
