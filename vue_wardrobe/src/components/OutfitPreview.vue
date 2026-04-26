<script setup lang="ts">
import type { ClothingItem } from '@/types/clothing'
import type { SlotType } from '@/types/category'
import { SLOT_LABELS } from '@/types/category'
import { useImageUrl } from '@/composables/useImageUrl'

const props = defineProps<{
  slots: Record<SlotType, { clothing: ClothingItem | null }>
}>()

const { resolve } = useImageUrl()

function hasItem(slot: SlotType): boolean {
  return !!props.slots[slot]?.clothing
}
</script>

<template>
  <div class="preview-container">
    <div class="preview-title">穿搭预览</div>
    <div class="preview-area">
      <div class="preview-inner">
        <!-- 外套层（最上层） -->
        <div v-if="hasItem('outer')" class="preview-layer layer-outer">
          <img :src="resolve(slots.outer.clothing!.masked_image || slots.outer.clothing!.original_image)" alt="外套" />
        </div>
        <!-- 配饰层 -->
        <div v-if="hasItem('accessory')" class="preview-layer layer-accessory">
          <img :src="resolve(slots.accessory.clothing!.masked_image || slots.accessory.clothing!.original_image)" alt="配饰" />
        </div>
        <!-- 上衣层 -->
        <div v-if="hasItem('top')" class="preview-layer layer-top">
          <img :src="resolve(slots.top.clothing!.masked_image || slots.top.clothing!.original_image)" alt="上衣" />
        </div>
        <!-- 下装层（裤子/裙子） -->
        <div v-if="hasItem('bottom')" class="preview-layer layer-bottom">
          <img :src="resolve(slots.bottom.clothing!.masked_image || slots.bottom.clothing!.original_image)" alt="裤子" />
        </div>
        <div v-else-if="hasItem('skirt')" class="preview-layer layer-bottom">
          <img :src="resolve(slots.skirt.clothing!.masked_image || slots.skirt.clothing!.original_image)" alt="裙子" />
        </div>
        <!-- 鞋子层 -->
        <div v-if="hasItem('shoes')" class="preview-layer layer-shoes">
          <img :src="resolve(slots.shoes.clothing!.masked_image || slots.shoes.clothing!.original_image)" alt="鞋子" />
        </div>
        <!-- 无选择提示 -->
        <div
          v-if="!hasItem('top') && !hasItem('bottom') && !hasItem('skirt')"
          class="preview-empty"
        >
          <van-icon name="photo-o" size="48" color="#ddd" />
          <p>请选择衣物开始搭配</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.preview-container {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  margin: 0 12px 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}
.preview-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
  text-align: center;
}
.preview-area {
  background: #f9f9f9;
  border-radius: 10px;
  padding: 20px;
  min-height: 200px;
}
.preview-inner {
  position: relative;
  width: 140px;
  height: 200px;
  margin: 0 auto;
}
.preview-layer {
  position: absolute;
  width: 100%;
  text-align: center;
}
.preview-layer img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}
.layer-outer {
  top: 0;
  height: 100%;
  z-index: 5;
}
.layer-accessory {
  top: -8px;
  right: -10px;
  width: 50px;
  height: 50px;
  z-index: 6;
}
.layer-top {
  top: 10%;
  height: 40%;
  z-index: 4;
}
.layer-bottom {
  top: 45%;
  height: 45%;
  z-index: 3;
}
.layer-shoes {
  bottom: 2%;
  height: 20%;
  z-index: 2;
}
.preview-empty {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  width: 100%;
}
.preview-empty p {
  color: #bbb;
  font-size: 13px;
  margin-top: 8px;
}
</style>
