<script setup lang="ts">
import type { SlotType } from '@/types/category'
import type { ClothingItem } from '@/types/clothing'
import { useImageUrl } from '@/composables/useImageUrl'
import { SLOT_LABELS } from '@/types/category'

const props = defineProps<{
  slot: SlotType
  selected: ClothingItem | null
}>()

const emit = defineEmits<{
  select: [slot: SlotType]
}>()

const { resolve } = useImageUrl()
</script>

<template>
  <div class="slot-block" :class="{ selected: !!selected }" @click="emit('select', slot)">
    <div class="slot-image">
      <img v-if="selected" :src="resolve(selected.masked_image || selected.original_image)" :alt="selected.name" />
      <div v-else class="slot-placeholder">
        <van-icon name="plus" size="24" color="#c8c9cc" />
      </div>
    </div>
    <div class="slot-label">{{ SLOT_LABELS[slot] }}</div>
    <div v-if="selected" class="slot-clothing-name">{{ selected.name }}</div>
  </div>
</template>

<style scoped>
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
.slot-clothing-name {
  font-size: 11px;
  color: #999;
  margin-top: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
