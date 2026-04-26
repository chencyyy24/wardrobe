<script setup lang="ts">
import type { ClothingItem } from '@/types/clothing'
import type { SlotType } from '@/types/category'
import { SLOT_LABELS } from '@/types/category'
import { useImageUrl } from '@/composables/useImageUrl'
import { ref } from 'vue'

const props = defineProps<{
  show: boolean
  slot: SlotType | null
  items: ClothingItem[]
  selectedId: string | null
}>()

const emit = defineEmits<{
  close: []
  select: [item: ClothingItem]
}>()

const { resolve } = useImageUrl()
const title = ref('选择衣物')
</script>

<template>
  <van-popup
    :show="show"
    position="bottom"
    round
    closeable
    @close="emit('close')"
    style="max-height: 60vh; min-height: 40vh;"
  >
    <div class="selector-header">
      <span class="selector-title">{{ slot ? `选择${SLOT_LABELS[slot]}` : '选择衣物' }}</span>
    </div>
    <div class="selector-content">
      <div v-if="items.length === 0" class="empty-tip">
        <van-empty description="暂无可用衣物，请先上传" />
      </div>
      <div v-else class="clothing-list">
        <div
          v-for="item in items"
          :key="item.id"
          class="clothing-option"
          :class="{ active: item.id === selectedId }"
          @click="emit('select', item)"
        >
          <van-image
            :src="resolve(item.masked_image || item.original_image)"
            fit="cover"
            class="option-image"
          />
          <div class="option-name">{{ item.name }}</div>
        </div>
      </div>
    </div>
  </van-popup>
</template>

<style scoped>
.selector-header {
  padding: 16px 16px 8px;
  font-size: 16px;
  font-weight: 600;
  text-align: center;
}
.selector-content {
  padding: 8px 16px 24px;
}
.clothing-list {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}
.clothing-option {
  border: 2px solid #f0f0f0;
  border-radius: 10px;
  padding: 8px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s;
}
.clothing-option.active {
  border-color: #1989fa;
  background: #f0f8ff;
}
.option-image {
  width: 100%;
  aspect-ratio: 1;
  border-radius: 6px;
  background: #fafafa;
}
.option-name {
  font-size: 12px;
  color: #333;
  margin-top: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.empty-tip {
  padding: 20px 0;
}
</style>
