import type { ApiResponse } from '@/types/api'

const BASE_URL = import.meta.env.VITE_API_BASE || ''

function authHeaders(): Record<string, string> {
  const token = localStorage.getItem('wardrobe_token')
  if (token) return { 'Authorization': `Bearer ${token}` }
  return {}
}

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

export async function apiGet<T>(path: string, params?: Record<string, string>): Promise<T> {
  const url = new URL(path, window.location.origin)
  if (params) {
    Object.entries(params).forEach(([k, v]) => url.searchParams.set(k, v))
  }
  const res = await fetch(`${BASE_URL}${url.pathname}${url.search}`, {
    headers: { ...authHeaders() },
  })
  if (!res.ok) throw new ApiError(res.status, await res.text())
  const json: ApiResponse<T> = await res.json()
  return json.data
}

export async function apiPost<T>(path: string, body?: FormData | Record<string, unknown>): Promise<T> {
  const isFormData = body instanceof FormData
  const res = await fetch(`${BASE_URL}${path}`, {
    method: 'POST',
    headers: {
      ...authHeaders(),
      ...(isFormData ? {} : { 'Content-Type': 'application/json' }),
    },
    body: isFormData ? body : body ? JSON.stringify(body) : undefined,
  })
  if (!res.ok) throw new ApiError(res.status, await res.text())
  const json: ApiResponse<T> = await res.json()
  return json.data
}

export async function apiDelete(path: string): Promise<void> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method: 'DELETE',
    headers: { ...authHeaders() },
  })
  if (!res.ok) throw new ApiError(res.status, await res.text())
}

export { ApiError }

export async function apiPut<T>(path: string, body?: Record<string, unknown>): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method: 'PUT',
    headers: {
      ...authHeaders(),
      'Content-Type': 'application/json',
    },
    body: body ? JSON.stringify(body) : undefined,
  })
  if (!res.ok) throw new ApiError(res.status, await res.text())
  const json: ApiResponse<T> = await res.json()
  return json.data
}

export async function apiPatch<T>(path: string, body?: Record<string, unknown>): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method: 'PATCH',
    headers: {
      ...authHeaders(),
      'Content-Type': 'application/json',
    },
    body: body ? JSON.stringify(body) : undefined,
  })
  if (!res.ok) throw new ApiError(res.status, await res.text())
  const json: ApiResponse<T> = await res.json()
  return json.data
}
