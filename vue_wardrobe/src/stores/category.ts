import { ref } from 'vue'
import { defineStore } from 'pinia'
import { categoryApi } from '@/api/category'
import type { Category } from '@/types/category'

export const useCategoryStore = defineStore('category', () => {
  const categories = ref<Category[]>([])
  const loaded = ref(false)
  const loading = ref(false)

  async function load() {
    if (loaded.value || loading.value) return
    loading.value = true
    try {
      categories.value = await categoryApi.getAll()
      loaded.value = true
    } catch (e) {
      console.error('Failed to load categories:', e)
    } finally {
      loading.value = false
    }
  }

  function getById(id: number): Category | undefined {
    return categories.value.find(c => c.id === id)
  }

  return { categories, loaded, loading, load, getById }
})
