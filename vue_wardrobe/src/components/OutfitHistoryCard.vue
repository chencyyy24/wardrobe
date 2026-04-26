<script setup lang="ts">
import type { OutfitSummary } from '@/types/outfit'
import { useImageUrl } from '@/composables/useImageUrl'
import { formatDate } from '@/utils/format'

const props = defineProps<{
  outfit: OutfitSummary
}>()

const emit = defineEmits<{
  click: []
  clone: []
  delete: []
}>()

const { resolve } = useImageUrl()
</script>

<template>
  <div class="history-card">
    <div class="card-image" @click="emit('click')">
      <van-image
        :src="resolve(outfit.card_image)"
        fit="cover"
        class="thumb"
        loading-icon="photo-o"
      >
        <template v-if="!outfit.card_image" #error>
          <div class="card-placeholder">
            <van-icon name="gem-o" size="32" color="#c8c9cc" />
          </div>
        </template>
      </van-image>
    </div>
    <div class="card-info">
      <div class="card-name">{{ outfit.name || '未命名搭配' }}</div>
      <div class="card-meta">
        <span>{{ formatDate(outfit.created_at) }}</span>
        <span class="item-count">{{ outfit.item_count }} 件</span>
      </div>
      <div class="card-actions">
        <van-button size="mini" plain type="primary" icon="copy-o" @click.stop="emit('clone')">
          复制
        </van-button>
        <van-button size="mini" plain type="danger" icon="delete-o" @click.stop="emit('delete')">
          删除
        </van-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.history-card {
  background: #fff;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 1px 4px rgba(0,0,0,0.06);
}
.card-image {
  width: 100%;
  aspect-ratio: 3/4;
  overflow: hidden;
  cursor: pointer;
}
.thumb {
  width: 100%;
  height: 100%;
}
.card-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
}
.card-info {
  padding: 8px 10px;
}
.card-name {
  font-size: 13px;
  font-weight: 500;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.card-meta {
  font-size: 11px;
  color: #999;
  margin: 4px 0 6px;
  display: flex;
  justify-content: space-between;
}
.card-actions {
  display: flex;
  gap: 6px;
}
</style>
