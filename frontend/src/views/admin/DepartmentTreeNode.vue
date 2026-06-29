<template>
  <div>
    <!-- Node Row -->
    <div
      :style="{ paddingLeft: depth * 24 + 'px' }"
      class="flex items-center gap-2 rounded-lg py-2.5 pr-3 transition-colors hover:bg-gray-50 dark:hover:bg-dark-800/50"
    >
      <!-- Expand/Collapse -->
      <button
        v-if="hasChildren"
        @click="$emit('toggle', node.id)"
        class="flex h-6 w-6 items-center justify-center rounded text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
      >
        <Icon :name="isExpanded ? 'chevronDown' : 'chevronRight'" size="sm" />
      </button>
      <span v-else class="w-6" />

      <!-- Node Info -->
      <div class="flex flex-1 items-center gap-3">
        <Icon name="inbox" size="md" class="text-gray-400 dark:text-gray-500" />
        <div class="flex flex-col">
          <span class="text-sm font-medium text-gray-900 dark:text-white">{{ node.name }}</span>
          <span class="text-xs text-gray-400 dark:text-gray-500">
            {{ node.member_count }} {{ t('admin.departments.members') }}
            <template v-if="node.leader"> · {{ node.leader }}</template>
          </span>
        </div>
        <span
          :class="[
            'ml-auto inline-flex rounded-full px-2 py-0.5 text-xs font-medium',
            node.status === 'active'
              ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300'
              : 'bg-gray-100 text-gray-500 dark:bg-gray-800 dark:text-gray-400'
          ]"
        >
          {{ node.status === 'active' ? t('admin.departments.status.active') : t('admin.departments.status.disabled') }}
        </span>
      </div>

      <!-- Actions -->
      <div class="flex items-center gap-0.5 ml-2">
        <button
          @click="$emit('add-child', node)"
          class="rounded-lg p-1.5 text-gray-400 hover:text-primary-600 hover:bg-gray-100 dark:hover:bg-dark-700 dark:hover:text-primary-400"
          :title="t('admin.departments.addChild')"
        >
          <Icon name="plus" size="sm" />
        </button>
        <button
          @click="$emit('edit', node)"
          class="rounded-lg p-1.5 text-gray-400 hover:text-primary-600 hover:bg-gray-100 dark:hover:bg-dark-700 dark:hover:text-primary-400"
          :title="t('common.edit')"
        >
          <Icon name="edit" size="sm" />
        </button>
        <button
          @click="$emit('delete', node)"
          class="rounded-lg p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 dark:hover:text-red-400"
          :title="t('common.delete')"
        >
          <Icon name="trash" size="sm" />
        </button>
      </div>
    </div>

    <!-- Children (recursive) -->
    <template v-if="hasChildren && isExpanded">
      <DepartmentTreeNode
        v-for="child in node.children"
        :key="child.id"
        :node="child"
        :depth="depth + 1"
        :expanded-ids="expandedIds"
        @edit="(n) => $emit('edit', n)"
        @add-child="(n) => $emit('add-child', n)"
        @delete="(n) => $emit('delete', n)"
        @toggle="(id) => $emit('toggle', id)"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, type Ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import type { Department } from '@/types/enterprise'

const props = defineProps<{
  node: Department
  depth: number
  expandedIds?: Set<number>
}>()

defineEmits<{
  edit: [node: Department]
  addChild: [node: Department]
  delete: [node: Department]
  toggle: [id: number]
}>()

const { t } = useI18n()

const hasChildren = computed(() => props.node.children && props.node.children.length > 0)
const isExpanded = computed(() => props.expandedIds?.has(props.node.id) ?? false)
</script>
