import { ref } from 'vue'
import { defineStore } from 'pinia'
import { outfitApi } from '@/api/outfit'
import type { Outfit, OutfitSummary } from '@/types/outfit'

export const useOutfitStore = defineStore('outfit', () => {
  const outfits = ref<Outfit[]>([])
  const summaries = ref<OutfitSummary[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchSummaries() {
    loading.value = true
    error.value = null
    try {
      summaries.value = (await outfitApi.summary()) || []
    } catch (e) {
      error.value = '加载搭配列表失败'
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  async function fetchFullList() {
    loading.value = true
    error.value = null
    try {
      outfits.value = await outfitApi.list()
    } catch (e) {
      error.value = '加载搭配列表失败'
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  async function getDetail(id: string): Promise<Outfit> {
    return await outfitApi.getById(id)
  }

  async function createOutfit(data: Record<string, unknown>): Promise<Outfit> {
    const outfit = await outfitApi.create(data)
    outfits.value.unshift(outfit)
    return outfit
  }

  async function cloneOutfit(id: string, name?: string): Promise<Outfit> {
    const outfit = await outfitApi.clone(id, name)
    summaries.value.unshift({
      id: outfit.id,
      name: outfit.name,
      card_image: outfit.card_image,
      created_at: outfit.created_at,
      item_count: outfit.items.length,
    })
    return outfit
  }

  async function deleteOutfit(id: string) {
    await outfitApi.delete(id)
    summaries.value = summaries.value.filter(s => s.id !== id)
    outfits.value = outfits.value.filter(o => o.id !== id)
  }

  async function updateOutfitItems(id: string, items: Record<string, string>): Promise<Outfit> {
    const updated = await outfitApi.updateItems(id, items)
    return updated
  }

  return {
    outfits, summaries, loading, error,
    fetchSummaries, fetchFullList, getDetail,
    createOutfit, cloneOutfit, deleteOutfit, updateOutfitItems,
  }
})
