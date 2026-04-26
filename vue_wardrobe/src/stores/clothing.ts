import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { clothingApi } from '@/api/clothing'
import type { ClothingItem } from '@/types/clothing'

export const useClothingStore = defineStore('clothing', () => {
  const items = ref<ClothingItem[]>([])
  const activeCategoryId = ref<number | undefined>(undefined)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const filteredItems = computed(() => {
    if (!activeCategoryId.value) return items.value
    return items.value.filter(i => i.category_id === activeCategoryId.value)
  })

  async function fetchList(categoryId?: number) {
    loading.value = true
    error.value = null
    try {
      items.value = await clothingApi.list(categoryId)
      activeCategoryId.value = categoryId
    } catch (e) {
      error.value = '加载衣物列表失败'
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  function addItem(item: ClothingItem) {
    items.value.unshift(item)
  }

  async function deleteItem(id: string) {
    await clothingApi.delete(id)
    items.value = items.value.filter(i => i.id !== id)
  }

  function getItemsByCategory(categoryId: number): ClothingItem[] {
    return items.value.filter(i => i.category_id === categoryId)
  }

  async function fetchAll() {
    loading.value = true
    error.value = null
    try {
      items.value = await clothingApi.list()
      activeCategoryId.value = undefined
    } catch (e) {
      error.value = '加载衣物列表失败'
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  return {
    items, filteredItems, activeCategoryId, loading, error,
    fetchList, fetchAll, addItem, deleteItem, getItemsByCategory,
  }
})
