import { apiGet, apiPost, apiPatch, apiPut, apiDelete } from './client'
import type { Outfit, OutfitSummary } from '@/types/outfit'

export const outfitApi = {
  list: () => apiGet<Outfit[]>('/api/outfit'),

  summary: () => apiGet<OutfitSummary[]>('/api/outfit/summary'),

  getById: (id: string) => apiGet<Outfit>(`/api/outfit/${id}`),

  create: (data: Record<string, unknown>) => apiPost<Outfit>('/api/outfit', data),

  updateCard: (id: string, cardImage: string) =>
    apiPatch<void>(`/api/outfit/${id}`, { card_image: cardImage }),

  updateItems: (id: string, items: Record<string, string>) =>
    apiPut<Outfit>(`/api/outfit/${id}/items`, { items }),

  delete: (id: string) => apiDelete(`/api/outfit/${id}`),

  clone: (id: string, name?: string) => apiPost<Outfit>(`/api/outfit/${id}/clone`, { name }),
}
