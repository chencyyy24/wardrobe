<script setup lang="ts">
import type { ClothingItem } from '@/types/clothing'
import { useImageUrl } from '@/composables/useImageUrl'
import { computed } from 'vue'

const props = defineProps<{
  item: ClothingItem
}>()

const emit = defineEmits<{
  delete: [id: string]
}>()

const { resolve } = useImageUrl()

const displayImage = computed(() => {
  return resolve(props.item.original_image)
})
</script>

<template>
  <div class="clothing-card" @click="$emit('delete', item.id)">
    <div class="image-wrapper">
      <van-image
        :src="displayImage"
        fit="cover"
        class="thumb"
        loading-icon="photo-o"
      />
    </div>
    <div class="card-name">{{ item.name }}</div>
  </div>
</template>

<style scoped>
.clothing-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 4px rgba(0,0,0,0.06);
}
.image-wrapper {
  position: relative;
  width: 100%;
  padding-top: 100%;
}
.thumb {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}
.card-name {
  padding: 6px 8px;
  font-size: 13px;
  color: #333;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
