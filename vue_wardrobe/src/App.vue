<script setup lang="ts">
import AppLayout from '@/components/AppLayout.vue'
</script>

<template>
  <AppLayout />
</template>

<style>
:root {
  --van-primary-color: #1989fa;
  --van-success-color: #07c160;
  --van-danger-color: #ee0a24;
  --van-warning-color: #ff976a;
  --van-text-color: #333;
  --van-text-color-2: #666;
  --van-text-color-3: #999;
  --van-border-color: #f0f0f0;
  --van-background-2: #f7f8fa;
}

body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB',
    'Microsoft YaHei', 'Helvetica Neue', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  background: #f7f8fa;
  color: #333;
}

* {
  box-sizing: border-box;
}

img {
  -webkit-user-drag: none;
  user-select: none;
}
</style>
