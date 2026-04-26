import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { authApi } from '@/api/auth'
import type { AuthUser } from '@/api/auth'
import { showToast } from 'vant'

const TOKEN_KEY = 'wardrobe_token'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(TOKEN_KEY))
  const user = ref<AuthUser | null>(null)
  const loading = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  async function login(username: string, password: string) {
    loading.value = true
    try {
      const result = await authApi.login(username, password)
      token.value = result.token
      user.value = result.user
      localStorage.setItem(TOKEN_KEY, result.token)
      return true
    } catch (e: any) {
      const msg = e.message || '登录失败'
      showToast(msg)
      return false
    } finally {
      loading.value = false
    }
  }

  async function register(username: string, password: string, nickname?: string) {
    loading.value = true
    try {
      const result = await authApi.register(username, password, nickname)
      token.value = result.token
      user.value = result.user
      localStorage.setItem(TOKEN_KEY, result.token)
      return true
    } catch (e: any) {
      const msg = e.message || '注册失败'
      showToast(msg)
      return false
    } finally {
      loading.value = false
    }
  }

  async function fetchUser() {
    if (!token.value) return
    try {
      user.value = await authApi.getMe()
    } catch {
      // Token invalid, clear
      logout()
    }
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem(TOKEN_KEY)
  }

  function restoreSession() {
    const saved = localStorage.getItem(TOKEN_KEY)
    if (saved) {
      token.value = saved
      fetchUser()
    }
  }

  return { token, user, loading, isLoggedIn, login, register, logout, fetchUser, restoreSession }
})
