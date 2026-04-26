<script setup lang="ts">
defineProps<{
  categories: { id: number; name: string }[]
  activeId?: number
}>()

const emit = defineEmits<{
  change: [id: number | undefined]
}>()
</script>

<template>
  <div class="category-tabs">
    <div
      class="tab-item"
      :class="{ active: !activeId }"
      @click="emit('change', undefined)"
    >
      全部
    </div>
    <div
      v-for="cat in categories"
      :key="cat.id"
      class="tab-item"
      :class="{ active: activeId === cat.id }"
      @click="emit('change', cat.id)"
    >
      {{ cat.name }}
    </div>
  </div>
</template>

<style scoped>
.category-tabs {
  display: flex;
  overflow-x: auto;
  gap: 0;
  background: #fff;
  padding: 0 12px;
  white-space: nowrap;
  -webkit-overflow-scrolling: touch;
  border-bottom: 1px solid #f0f0f0;
}
.tab-item {
  padding: 10px 14px;
  font-size: 14px;
  color: #666;
  position: relative;
  cursor: pointer;
  flex-shrink: 0;
}
.tab-item.active {
  color: #1989fa;
  font-weight: 600;
}
.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 20px;
  height: 3px;
  background: #1989fa;
  border-radius: 2px;
}
</style>
