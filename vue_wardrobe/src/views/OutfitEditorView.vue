<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useClothingStore } from '@/stores/clothing'
import { useOutfitStore } from '@/stores/outfit'
import { useCategoryStore } from '@/stores/category'
import { useImageUrl } from '@/composables/useImageUrl'
import { SLOT_CATEGORY_MAP, SLOT_LABELS, ALL_SLOTS } from '@/types/category'
import { generateOutfitName } from '@/utils/format'
import type { SlotType } from '@/types/category'
import type { ClothingItem } from '@/types/clothing'
import type { OutfitEditorState } from '@/types/outfit'
import { showToast, showLoadingToast, closeToast } from 'vant'

const router = useRouter()
const clothingStore = useClothingStore()
const outfitStore = useOutfitStore()
const categoryStore = useCategoryStore()
const { resolve } = useImageUrl()

const editor = reactive<OutfitEditorState>({
  name: generateOutfitName(),
  slots: {
    outer: { clothing: null },
    top: { clothing: null },
    bottom: { clothing: null },
    skirt: { clothing: null },
    shoes: { clothing: null },
    accessory: { clothing: null },
  },
})

const showSelector = ref(false)
const activeSlot = ref<SlotType | null>(null)
const saving = ref(false)

const slotItems = computed<ClothingItem[]>(() => {
  if (!activeSlot.value) return []
  const catId = SLOT_CATEGORY_MAP[activeSlot.value]
  return clothingStore.getItemsByCategory(catId)
})

const selectedSlotId = computed(() => {
  if (!activeSlot.value) return null
  return editor.slots[activeSlot.value]?.clothing?.id || null
})

const slotRules = computed(() => {
  const hasTop = !!editor.slots.top.clothing
  const hasBottom = !!editor.slots.bottom.clothing
  const hasSkirt = !!editor.slots.skirt.clothing
  const conflict = hasBottom && hasSkirt
  const missingRequired = !hasTop || (!hasBottom && !hasSkirt)

  let message = ''
  if (conflict) message = '裤子和裙子不能同时选择'
  else if (!hasTop) message = '请选择上衣'
  else if (!hasBottom && !hasSkirt) message = '请选择裤子或裙子'

  return { message, valid: !conflict && !missingRequired }
})

function openSelector(slot: SlotType) {
  activeSlot.value = slot
  showSelector.value = true
}

function selectItem(item: ClothingItem) {
  if (!activeSlot.value) return
  editor.slots[activeSlot.value].clothing = item
  showSelector.value = false
}

async function saveOutfit() {
  if (!slotRules.value.valid) {
    showToast(slotRules.value.message)
    return
  }

  saving.value = true
  showLoadingToast({ message: '保存中...', forbidClick: true })

  try {
    const items: Record<string, string> = {}
    for (const slot of ALL_SLOTS) {
      const clothing = editor.slots[slot].clothing
      if (clothing) {
        items[slot] = clothing.id
      }
    }

    if (!items.top || !items.bottom) {
      closeToast()
      showToast('请确保上衣和裤子已选择')
      saving.value = false
      return
    }

    await outfitStore.createOutfit({
      name: editor.name,
      items,
    })

    closeToast()
    showToast('搭配保存成功')
    router.push('/outfit')
  } catch (e) {
    closeToast()
    showToast('保存失败，请重试')
    console.error(e)
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="editor-page">
    <!-- 6 Slot Grid -->
    <div class="slot-grid">
      <div
        v-for="slot in ALL_SLOTS"
        :key="slot"
        class="slot-wrapper"
        @click="openSelector(slot)"
      >
        <div class="slot-block" :class="{ selected: !!editor.slots[slot].clothing }">
          <div class="slot-image">
            <img
              v-if="editor.slots[slot].clothing"
              :src="resolve(editor.slots[slot].clothing!.masked_image || editor.slots[slot].clothing!.original_image)"
              :alt="editor.slots[slot].clothing!.name"
            />
            <div v-else class="slot-placeholder">
              <van-icon name="plus" size="28" color="#c8c9cc" />
            </div>
          </div>
          <div class="slot-label">{{ SLOT_LABELS[slot] }}</div>
          <div v-if="editor.slots[slot].clothing" class="slot-name">
            {{ editor.slots[slot].clothing!.name }}
          </div>
        </div>
      </div>
    </div>

    <!-- Validation message -->
    <div v-if="!slotRules.valid" class="rule-hint">
      <van-icon name="info-o" color="#ee0a24" />
      <span>{{ slotRules.message }}</span>
    </div>

    <!-- Preview -->
    <OutfitPreview :slots="editor.slots" />

    <!-- Name input -->
    <div class="name-section">
      <van-field
        v-model="editor.name"
        label="搭配名称"
        placeholder="输入搭配名称"
        clearable
      />
    </div>

    <!-- Save button -->
    <div class="save-section">
      <van-button
        type="primary"
        block
        size="large"
        round
        :disabled="!slotRules.valid || saving"
        :loading="saving"
        @click="saveOutfit"
      >
        保存搭配
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
.editor-page {
  padding: 12px 0 40px;
}
.slot-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  padding: 0 12px;
  margin-bottom: 12px;
}
.slot-block {
  background: #fff;
  border-radius: 10px;
  border: 2px solid #f0f0f0;
  padding: 8px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s;
}
.slot-block.selected {
  border-color: #1989fa;
}
.slot-image {
  width: 100%;
  padding-top: 100%;
  position: relative;
  background: #fafafa;
  border-radius: 6px;
  overflow: hidden;
}
.slot-image img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.slot-placeholder {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.slot-label {
  font-size: 13px;
  color: #666;
  margin-top: 6px;
}
.slot-name {
  font-size: 11px;
  color: #999;
  margin-top: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.rule-hint {
  padding: 8px 16px;
  background: #fff2f0;
  color: #ee0a24;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0 12px 12px;
  border-radius: 8px;
}
.name-section {
  margin: 0 12px 16px;
}
.name-section :deep(.van-field) {
  background: #fff;
  border-radius: 8px;
}
.save-section {
  padding: 0 16px 20px;
}
</style>
