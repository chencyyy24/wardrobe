export interface ClothingItem {
  id: string
  name: string
  user_id: number
  category_id: number
  subcategory_id: number | null
  original_image: string
  masked_image: string | null
  status: string
  price: number | null
  created_at: string
  updated_at: string
  category?: {
    id: number
    name: string
  }
  subcategory?: {
    id: number
    name: string
  }
}
