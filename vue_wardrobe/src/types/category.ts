export interface Subcategory {
  id: number
  name: string
}

export interface Category {
  id: number
  name: string
  subcategories: Subcategory[]
}

export type SlotType = 'outer' | 'top' | 'bottom' | 'skirt' | 'shoes' | 'accessory'

export const SLOT_CATEGORY_MAP: Record<SlotType, number> = {
  outer: 1,
  top: 2,
  bottom: 3,
  skirt: 4,
  shoes: 5,
  accessory: 6,
}

export const SLOT_LABELS: Record<SlotType, string> = {
  outer: '外套',
  top: '上衣',
  bottom: '裤子',
  skirt: '裙子',
  shoes: '鞋子',
  accessory: '配饰',
}

export const ALL_SLOTS: SlotType[] = ['outer', 'top', 'bottom', 'skirt', 'shoes', 'accessory']
