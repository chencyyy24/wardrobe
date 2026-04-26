import { apiPost, apiGet, apiPut } from './client'
import type { ApiResponse } from '@/types/api'

export interface AuthUser {
  id: number
  username: string
  nickname: string
  avatar?: string | null
  created_at?: string
}

export interface AuthResult {
  token: string
  user: AuthUser
}

export const authApi = {
  login: (username: string, password: string) =>
    apiPost<AuthResult>('/api/auth/login', { username, password }),

  register: (username: string, password: string, nickname?: string) =>
    apiPost<AuthResult>('/api/auth/register', { username, password, nickname }),

  getMe: () => apiGet<AuthUser>('/api/auth/me'),

  updateProfile: (nickname: string) =>
    apiPut<{ nickname: string }>('/api/auth/profile', { nickname }),

  changePassword: (oldPassword: string, newPassword: string) =>
    apiPut<void>('/api/auth/password', { old_password: oldPassword, new_password: newPassword }),

  uploadAvatar: async (file: File): Promise<{ avatar: string }> => {
    const token = localStorage.getItem('wardrobe_token')
    const formData = new FormData()
    formData.append('avatar', file)
    const res = await fetch('/api/auth/avatar', {
      method: 'POST',
      headers: token ? { 'Authorization': `Bearer ${token}` } : {},
      body: formData,
    })
    if (!res.ok) throw new Error('上传头像失败')
    const json = await res.json()
    return json.data
  },
}
