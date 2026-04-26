export function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  const month = d.getMonth() + 1
  const day = d.getDate()
  return `${month}月${day}日`
}

export function generateOutfitName(): string {
  const d = new Date()
  const styles = ['简约风', '通勤风', '休闲风', '清新风', '优雅风']
  const style = styles[Math.floor(Math.random() * styles.length)]
  return `${d.getMonth() + 1}月${d.getDate()}日${style}`
}
