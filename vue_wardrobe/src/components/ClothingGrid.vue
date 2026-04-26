<script setup lang="ts">
import type { ClothingItem } from '@/types/clothing'

defineProps<{
  items: ClothingItem[]
}>()

const emit = defineEmits<{
  select: [item: ClothingItem]
}>()
</script>

<template>
  <div class="clothing-grid">
    <div
      v-for="item in items"
      :key="item.id"
      class="grid-item"
      @click="emit('select', item)"
    >
      <div class="image-wrapper">
        <van-image
          :src="item.masked_image || item.original_image"
          fit="cover"
          class="thumb"
          loading-icon="photo-o"
        />
        <div v-if="item.status === 'pending'" class="status-overlay">
          <van-loading color="#fff" size="16" />
        </div>
      </div>
      <div class="item-name">{{ item.name }}</div>
    </div>
  </div>
</template>

<style scoped>
.clothing-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  padding: 12px;
}
.grid-item {
  cursor: pointer;
}
.image-wrapper {
  position: relative;
  width: 100%;
  padding-top: 100%;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
}
.thumb {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}
.status-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.3);
  display: flex;
  align-items: center;
  justify-content: center;
}
.item-name {
  font-size: 12px;
  color: #333;
  text-align: center;
  padding: 4px 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
