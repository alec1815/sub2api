<template>
  <BaseDialog :show="show" :title="t('admin.enterprises.apiKeys')" width="wide" @close="$emit('close')">
    <div v-if="enterprise" class="space-y-4">
      <div class="flex items-center gap-3 rounded-xl bg-gray-50 p-4 dark:bg-dark-700">
        <div class="flex h-10 w-10 items-center justify-center rounded-full bg-primary-100"><span class="text-lg font-medium text-primary-700">{{ enterprise.name.charAt(0).toUpperCase() }}</span></div>
        <div class="flex-1"><p class="font-medium text-gray-900">{{ enterprise.name }}</p></div>
      </div>

      <div v-if="loading" class="flex justify-center py-8"><svg class="h-8 w-8 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/></svg></div>

      <div v-else-if="items.length === 0" class="py-8 text-center"><p class="text-sm text-gray-500">{{ t('admin.enterprises.noApiKeys') }}</p></div>

      <div v-else class="max-h-[24rem] space-y-2 overflow-y-auto">
        <div v-for="item in items" :key="item.id" class="rounded-lg border p-3 dark:border-dark-600">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="text-sm font-medium text-gray-900 dark:text-white">{{ item.name }}</span>
              <span :class="['inline-block h-1.5 w-1.5 rounded-full', item.status === 'active' ? 'bg-green-500' : 'bg-gray-400']" />
            </div>
            <span class="text-xs text-gray-400">{{ formatDateTime(item.created_at) }}</span>
          </div>
          <p class="mt-1 font-mono text-xs text-gray-500 dark:text-gray-400">{{ item.key.slice(0, 12) }}...{{ item.key.slice(-6) }}</p>
        </div>
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { formatDateTime } from '@/utils/format'
import { adminAPI } from '@/api/admin'
import type { Enterprise } from '@/types/enterprise'
import BaseDialog from '@/components/common/BaseDialog.vue'

const props = defineProps<{ show: boolean; enterprise: Enterprise | null }>()
defineEmits(['close'])
const { t } = useI18n()

interface KeyItem { id: number; name: string; key: string; status: string; created_at: string }
const items = ref<KeyItem[]>([])
const loading = ref(false)

watch(() => props.show, (v) => {
  if (v && props.enterprise) { loadKeys() }
})

const loadKeys = async () => {
  if (!props.enterprise) return
  loading.value = true
  try {
    const res = await adminAPI.enterprises.getEnterpriseKeys(props.enterprise.id)
    items.value = (res as any).items || []
  } finally { loading.value = false }
}
</script>
