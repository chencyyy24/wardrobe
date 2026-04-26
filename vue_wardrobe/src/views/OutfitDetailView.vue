<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useOutfitStore } from '@/stores/outfit'
import { useClothingStore } from '@/stores/clothing'
import { useImageUrl } from '@/composables/useImageUrl'
import { SLOT_LABELS, SLOT_CATEGORY_MAP, ALL_SLOTS } from '@/types/category'
import { outfitApi } from '@/api/outfit'
import type { Outfit } from '@/types/outfit'
import type { ClothingItem } from '@/types/clothing'
import type { SlotType } from '@/types/category'
import { showToast, showConfirmDialog, showLoadingToast, closeToast } from 'vant'

const route = useRoute()
const router = useRouter()
const outfitStore = useOutfitStore()
const clothingStore = useClothingStore()
const { resolve } = useImageUrl()

const outfit = ref<Outfit | null>(null)
const loading = ref(true)
const coverImage = ref('')

// Edit mode
const editing = ref(false)
const saving = ref(false)
const editSlots = ref<Record<string, ClothingItem | null>>({})
const showSelector = ref(false)
const activeSlot = ref<SlotType | null>(null)

const groupedItems = computed(() => {
  if (!outfit.value) return []
  return outfit.value.items.map(item => ({
    slot: item.slot,
    label: SLOT_LABELS[item.slot as keyof typeof SLOT_LABELS] || item.slot,
    clothing: item.clothing,
  }))
})

const slotItems = computed<ClothingItem[]>(() => {
  if (!activeSlot.value) return []
  const catId = SLOT_CATEGORY_MAP[activeSlot.value]
  return clothingStore.getItemsByCategory(catId)
})

const selectedSlotId = computed(() => {
  if (!activeSlot.value) return null
  return editSlots.value[activeSlot.value]?.id || null
})

onMounted(async () => {
  const id = route.params.id as string
  try {
    outfit.value = await outfitStore.getDetail(id)
    coverImage.value = outfit.value.card_image || ''
    // Init edit slots from current outfit items
    for (const item of outfit.value.items) {
      editSlots.value[item.slot] = item.clothing
    }
  } catch (e) {
    showToast('搭配不存在')
    router.push('/outfit')
  } finally {
    loading.value = false
  }
})

function startEdit() {
  editing.value = true
  // Load all clothing items for selectors
  clothingStore.fetchAll()
}

function cancelEdit() {
  editing.value = false
  // Reset edit slots to current outfit items
  editSlots.value = {}
  if (outfit.value) {
    for (const item of outfit.value.items) {
      editSlots.value[item.slot] = item.clothing
    }
  }
}

function openSelector(slot: SlotType) {
  activeSlot.value = slot
  showSelector.value = true
}

function selectItem(item: ClothingItem) {
  if (!activeSlot.value) return
  editSlots.value[activeSlot.value] = item
  showSelector.value = false
}

async function saveEdits() {
  if (!outfit.value) return

  // Validate: need top + (bottom or skirt)
  const hasTop = !!editSlots.value['top']
  const hasBottom = !!editSlots.value['bottom']
  const hasSkirt = !!editSlots.value['skirt']
  if (!hasTop) { showToast('请选择上衣'); return }
  if (hasBottom && hasSkirt) { showToast('裤子和裙子不能同时选择'); return }
  if (!hasBottom && !hasSkirt) { showToast('请选择裤子或裙子'); return }

  saving.value = true
  showLoadingToast({ message: '保存中...', forbidClick: true })

  try {
    const items: Record<string, string> = {}
    for (const slot of ALL_SLOTS) {
      const clothing = editSlots.value[slot]
      if (clothing) {
        items[slot] = clothing.id
      }
    }

    const updated = await outfitStore.updateOutfitItems(outfit.value.id, items)
    outfit.value = updated
    coverImage.value = updated.card_image || ''
    editing.value = false
    closeToast()
    showToast('更新成功')
  } catch (e) {
    closeToast()
    showToast('保存失败')
    console.error(e)
  } finally {
    saving.value = false
  }
}

async function onUploadCover() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.onchange = async () => {
    const file = input.files?.[0]
    if (!file || !outfit.value) return

    showLoadingToast({ message: '上传封面中...', forbidClick: true })
    try {
      const formData = new FormData()
      formData.append('image', file)
      const token = localStorage.getItem('wardrobe_token')
      const res = await fetch('/api/upload/card', {
        method: 'POST',
        headers: token ? { 'Authorization': `Bearer ${token}` } : {},
        body: formData,
      })
      const json = await res.json()
      const url = json.url || json.data?.url
      if (!url) throw new Error('上传失败')

      await outfitApi.updateCard(outfit.value.id, url)
      coverImage.value = url
      outfit.value.card_image = url
      closeToast()
      showToast('封面已更新')
    } catch (e) {
      closeToast()
      showToast('封面上传失败')
      console.error(e)
    }
  }
  input.click()
}

async function cloneOutfit() {
  if (!outfit.value) return
  showLoadingToast({ message: '复制中...', forbidClick: true })
  try {
    await outfitStore.cloneOutfit(outfit.value.id)
    closeToast()
    showToast('复制成功')
    router.push('/outfit')
  } catch {
    closeToast()
    showToast('复制失败')
  }
}

async function deleteOutfit() {
  if (!outfit.value) return
  showConfirmDialog({
    title: '删除搭配',
    message: `确定要删除「${outfit.value.name}」吗？`,
    confirmButtonColor: '#ee0a24',
  }).then(async () => {
    if (!outfit.value) return
    try {
      await outfitStore.deleteOutfit(outfit.value.id)
      showToast('删除成功')
      router.push('/outfit')
    } catch {
      showToast('删除失败')
    }
  }).catch(() => {})
}
</script>

<template>
  <div v-if="loading" class="loading-page">
    <van-loading size="24">加载中...</van-loading>
  </div>

  <div v-else-if="outfit" class="detail-page">
    <!-- Card image - tappable to upload -->
    <div class="card-banner" @click="onUploadCover">
      <div v-if="coverImage" class="cover-wrapper">
        <van-image :src="resolve(coverImage)" fit="cover" class="cover-img" />
        <div class="cover-overlay">
          <van-icon name="photograph" size="24" color="#fff" />
          <span>点击更换封面</span>
        </div>
      </div>
      <div v-else class="cover-upload">
        <van-icon name="photograph" size="48" color="#c8c9cc" />
        <p>点击上传穿搭封面</p>
        <span class="cover-hint">上传后可在历史记录中快速识别</span>
      </div>
    </div>

    <!-- Outfit name -->
    <div class="name-section">
      <h2>{{ outfit.name || '未命名搭配' }}</h2>
      <span class="date">{{ new Date(outfit.created_at).toLocaleDateString('zh-CN') }}</span>
    </div>

    <!-- Items list -->
    <div class="items-section">
      <div class="section-title">{{ editing ? '点击单品可替换' : '搭配单品' }}</div>
      <div v-if="groupedItems.length === 0 && !editing" class="empty-items">
        <van-empty description="暂无单品信息" />
      </div>
      <div
        v-for="group in groupedItems"
        :key="group.slot"
        class="item-row"
        :class="{ clickable: editing }"
        @click="editing ? openSelector(group.slot as SlotType) : undefined"
      >
        <div class="item-slot">{{ group.label }}</div>
        <div class="item-clothing">
          <van-image
            :src="resolve(group.clothing.masked_image || group.clothing.original_image)"
            fit="cover"
            class="item-thumb"
          />
          <div class="item-detail">
            <div class="item-name">{{ group.clothing.name }}</div>
            <div class="item-category">
              {{ group.clothing.category?.name }}
              <template v-if="group.clothing.subcategory"> / {{ group.clothing.subcategory.name }}</template>
            </div>
          </div>
        </div>
        <van-icon v-if="editing" name="exchange" color="#1989fa" size="18" />
      </div>

      <!-- Empty slots in edit mode -->
      <div v-if="editing">
        <div
          v-for="slot in ALL_SLOTS"
          :key="'edit-' + slot"
          v-show="!editSlots[slot]"
          class="item-row clickable"
          @click="openSelector(slot)"
        >
          <div class="item-slot">{{ SLOT_LABELS[slot] }}</div>
          <div class="item-clothing item-clothing-empty">
            <van-icon name="plus" color="#c8c9cc" size="20" />
            <span class="empty-slot-text">点击添加{{ SLOT_LABELS[slot] }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Actions -->
    <div v-if="!editing" class="actions-section">
      <van-button type="primary" block round @click="startEdit">
        <template #icon><van-icon name="edit" /></template>
        编辑穿搭
      </van-button>
      <van-button type="primary" block round plain @click="cloneOutfit">
        <template #icon><van-icon name="copy-o" /></template>
        复制此搭配
      </van-button>
      <van-button type="danger" block round plain @click="deleteOutfit">
        <template #icon><van-icon name="delete-o" /></template>
        删除搭配
      </van-button>
    </div>

    <!-- Edit mode actions -->
    <div v-if="editing" class="actions-section">
      <van-button type="primary" block round :loading="saving" :disabled="saving" @click="saveEdits">
        <template #icon><van-icon name="success" /></template>
        保存修改
      </van-button>
      <van-button block round plain @click="cancelEdit">
        <template #icon><van-icon name="cross" /></template>
        取消
      </van-button>
    </div>

    <!-- Clothing selector popup -->
    <ClothingSelectorPopup
      :show="showSelector"
      :slot="activeSlot"
      :items="slotItems"
      :selected-id="selectedSlotId"
      @close="showSelector = false"
      @select="selectItem"
    />
  </div>
</template>

<style scoped>
.loading-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
}
.detail-page {
  padding-bottom: 40px;
}
.card-banner {
  background: #fff;
  cursor: pointer;
}
.cover-wrapper {
  position: relative;
  width: 100%;
  aspect-ratio: 3/4;
  overflow: hidden;
}
.cover-img {
  width: 100%;
  height: 100%;
}
.cover-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(transparent, rgba(0,0,0,0.6));
  padding: 40px 16px 12px;
  display: flex;
  align-items: center;
  gap: 6px;
  color: #fff;
  font-size: 13px;
  opacity: 0;
  transition: opacity 0.2s;
}
.cover-wrapper:hover .cover-overlay,
.cover-wrapper:active .cover-overlay {
  opacity: 1;
}
.cover-upload {
  text-align: center;
  padding: 60px 16px;
  color: #999;
  background: #fafafa;
}
.cover-upload p {
  font-size: 15px;
  margin: 12px 0 4px;
  color: #666;
}
.cover-hint {
  font-size: 12px;
  color: #bbb;
}
.name-section {
  padding: 16px;
  background: #fff;
  margin-bottom: 12px;
}
.name-section h2 {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0;
}
.name-section .date {
  font-size: 13px;
  color: #999;
  margin-top: 4px;
}
.items-section {
  background: #fff;
  padding: 16px;
  margin-bottom: 12px;
}
.section-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}
.empty-items {
  padding: 20px 0;
}
.item-row {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f5f5f5;
}
.item-row:last-child {
  border-bottom: none;
}
.item-row.clickable {
  cursor: pointer;
}
.item-row.clickable:active {
  background: #fafafa;
}
.item-slot {
  width: 44px;
  font-size: 13px;
  color: #666;
  flex-shrink: 0;
}
.item-clothing {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
}
.item-clothing-empty {
  gap: 6px;
  color: #c8c9cc;
}
.empty-slot-text {
  font-size: 13px;
  color: #c8c9cc;
}
.item-thumb {
  width: 50px;
  height: 50px;
  border-radius: 6px;
  overflow: hidden;
  background: #f5f5f5;
  flex-shrink: 0;
}
.item-detail {
  flex: 1;
}
.item-name {
  font-size: 14px;
  color: #333;
}
.item-category {
  font-size: 12px;
  color: #999;
  margin-top: 2px;
}
.actions-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 0 16px;
}
</style>
