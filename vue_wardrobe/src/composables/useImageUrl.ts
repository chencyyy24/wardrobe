const BASE = import.meta.env.VITE_API_BASE || ''

export function useImageUrl() {
  function resolve(path: string | null | undefined): string {
    if (!path) return ''
    if (path.startsWith('http')) return path
    // 相对路径补全为服务器地址（兼容 Capacitor APK）
    if (path.startsWith('/') && BASE) {
      return BASE + path
    }
    return path
  }

  return { resolve }
}
