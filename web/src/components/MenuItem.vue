<template>
  <button 
    @click="handleClick"
    :class="['menu-item', 'clickable-item', 'sidebar-indicator', { active: isActive }]"
  >
    <span class="menu-icon">{{ icon }}</span>
    <span class="menu-label">{{ label }}</span>
  </button>
</template>

<script>
export default {
  name: 'MenuItem',
  props: {
    icon: {
      type: String,
      required: true
    },
    label: {
      type: String,
      required: true
    },
    value: {
      type: String
    },
    isActive: {
      type: Boolean,
      default: false
    },
    index: {
      type: Number,
      default: 0
    }
  },
  computed: {
    animationDelay() {
      return 0.1 + (this.index * 0.05);
    }
  },
  methods: {
    handleClick() {
      // 优先使用value属性，如果没有则使用label的小写形式
      this.$emit('select', this.value || this.label.toLowerCase());
    }
  },
  mounted() {
    this.$el.style.animationDelay = `${this.animationDelay}s`;
  }
}
</script>

<style scoped>
.menu-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  width: 100%;
  padding: var(--space-3) var(--space-4);
  border: none;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  text-align: left;
  border-radius: var(--radius-md);
  animation: slideIn 0.3s ease forwards;
  opacity: 0;
  transform: translateY(10px);
}

.menu-item:hover {
  background: linear-gradient(90deg, 
    rgba(99, 102, 241, 0.1) 0%, 
    rgba(99, 102, 241, 0.05) 100%
  );
  color: var(--primary);
  transform: translateX(2px);
}

.menu-item.active {
  background: linear-gradient(90deg, 
    rgba(99, 102, 241, 0.15) 0%, 
    rgba(99, 102, 241, 0.08) 100%
  );
  color: var(--primary);
  font-weight: var(--font-weight-semibold);
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.15);
}

.menu-icon {
  font-size: 1.3em;
  width: 28px;
  text-align: center;
  transition: transform 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 28px;
  background: rgba(99, 102, 241, 0.05);
  border-radius: var(--radius-md);
}

.menu-item:hover .menu-icon {
  transform: scale(1.1);
  background: rgba(99, 102, 241, 0.15);
}

.menu-item.active .menu-icon {
  background: rgba(99, 102, 241, 0.2);
}

.menu-label {
  flex: 1;
  position: relative;
  transition: color 0.3s ease;
}
</style>