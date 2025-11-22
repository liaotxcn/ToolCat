<template>
  <button 
    @click="handleClick"
    :class="['plugin-item', 'clickable-item', { active: isActive }]"
    :title="description || name"
  >
    <span class="item-title">{{ name }}</span>
    <span v-if="description" class="item-desc">{{ description }}</span>
    <span v-if="badgeCount !== undefined" class="item-badge">{{ badgeCount }}</span>
  </button>
</template>

<script>
export default {
  name: 'PluginItem',
  props: {
    name: {
      type: String,
      required: true
    },
    description: String,
    badgeCount: Number,
    isActive: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    handleClick() {
      this.$emit('select', this.name);
    }
  }
}
</script>

<style scoped>
.plugin-item {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-background);
  color: var(--color-text-primary);
  cursor: pointer;
  text-align: left;
  width: 100%;
  will-change: transform, background-color, border-color;
}

.plugin-item:hover {
  background: var(--color-background-hover);
  border-color: var(--color-primary-light);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.plugin-item.active {
  background: var(--color-primary-light);
  color: var(--color-text-primary);
  border-color: var(--color-primary);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
  font-weight: 600;
}

.item-title {
  font-weight: 500;
  font-size: var(--font-size-sm);
  margin-bottom: var(--space-1);
}

.item-desc {
  font-size: var(--font-size-xs);
  opacity: 0.8;
  line-height: 1.4;
}

.item-badge {
  align-self: flex-end;
  background: rgba(99, 102, 241, 0.2);
  color: var(--color-primary);
  font-size: var(--font-size-xs);
  padding: 2px 6px;
  border-radius: var(--radius-sm);
  margin-top: var(--space-1);
}

.plugin-item.active .item-badge {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}
</style>