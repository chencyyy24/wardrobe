export function useImageUrl() {
  function resolve(path: string | null | undefined): string {
    if (!path) return ''
    if (path.startsWith('http')) return path
    return path
  }

  return { resolve }
}
