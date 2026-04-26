import type { SlotType } from './category'
import type { ClothingItem } from './clothing'

export interface OutfitItem {
  id: number
  outfit_id: string
  clothing_id: string
  slot: SlotType
  clothing: ClothingItem
}

export interface Outfit {
  id: string
  user_id: number
  name: string
  card_image: string | null
  created_at: string
  updated_at: string
  items: OutfitItem[]
}

export interface OutfitSummary {
  id: string
  name: string
  card_image: string | null
  created_at: string
  item_count: number
}

export interface OutfitSlotSelection {
  clothing: ClothingItem | null
}

export interface OutfitEditorState {
  name: string
  slots: Record<SlotType, OutfitSlotSelection>
}

export interface CreateOutfitPayload {
  name: string
  items: Partial<Record<SlotType, string>>
  card_image?: string
}
