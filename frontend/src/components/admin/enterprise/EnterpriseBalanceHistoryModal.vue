<template>
  <BaseDialog :show="show" :title="t('admin.enterprises.balanceHistoryTitle')" width="wide" @close="$emit('close')">
    <div v-if="enterprise" class="space-y-4">
      <div class="flex items-center gap-3 rounded-xl bg-gray-50 p-4 dark:bg-dark-700">
        <div class="flex h-10 w-10 items-center justify-center rounded-full bg-primary-100"><span class="text-lg font-medium text-primary-700">{{ enterprise.name.charAt(0).toUpperCase() }}</span></div>
        <div class="flex-1"><p class="font-medium text-gray-900">{{ enterprise.name }}</p><p class="text-sm text-gray-500">{{ t('admin.enterprises.currentBalance') }}: ${{ parseFloat(enterprise.balance ?? '0').toFixed(2) }}</p></div>
      </div>

      <div v-if="loading" class="flex justify-center py-8"><svg class="h-8 w-8 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/></svg></div>

      <div v-else-if="history.length === 0" class="py-8 text-center"><p class="text-sm text-gray-500">{{ t('admin.enterprises.noBalanceHistory') }}</p></div>

      <div v-else class="max-h-[24rem] space-y-2 overflow-y-auto">
        <div v-for="item in history" :key="item.id" class="flex items-center justify-between rounded-lg border p-3 dark:border-dark-600">
          <div class="flex items-center gap-3">
            <div :class="['flex h-8 w-8 items-center justify-center rounded-lg', item.operation === 'add' ? 'bg-emerald-100 dark:bg-emerald-900/30' : 'bg-red-100 dark:bg-red-900/30']">
              <Icon :name="item.operation === 'add' ? 'plus' : 'dollar'" size="sm" :class="item.operation === 'add' ? 'text-emerald-600' : 'text-red-600'" />
            </div>
            <div>
              <p class="text-sm font-medium">{{ item.operation === 'add' ? t('admin.enterprises.deposit') : t('admin.enterprises.withdraw') }}</p>
              <p v-if="item.notes" class="max-w-xs truncate text-xs text-gray-400" :title="item.notes">{{ item.notes }}</p>
              <p class="text-xs text-gray-400">{{ formatDateTime(item.created_at) }}</p>
            </div>
          </div>
          <p :class="['text-sm font-semibold', item.operation === 'add' ? 'text-emerald-600' : 'text-red-600']">{{ item.operation === 'add' ? '+' : '-' }}${{ item.amount.toFixed(2) }}</p>
        </div>
      </div>

      <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 pt-2">
        <button :disabled="page <= 1" class="btn btn-secondary px-3 py-1 text-sm" @click="loadHistory(page - 1)">{{ t('pagination.previous') }}</button>
        <span class="text-sm text-gray-500">{{ page }} / {{ totalPages }}</span>
        <button :disabled="page >= totalPages" class="btn btn-secondary px-3 py-1 text-sm" @click="loadHistory(page + 1)">{{ t('pagination.next') }}</button>
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { formatDateTime } from '@/utils/format'
import { adminAPI } from '@/api/admin'
import type { Enterprise } from '@/types/enterprise'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'

const props = defineProps<{ show: boolean; enterprise: Enterprise | null }>()
defineEmits(['close'])
const { t } = useI18n()

interface HistoryItem { id: number; amount: number; operation: string; notes: string; created_at: string }
const history = ref<HistoryItem[]>([])
const loading = ref(false)
const page = ref(1)
const total = ref(0)
const pageSize = 15

const totalPages = computed(() => Math.ceil(total.value / pageSize) || 1)

watch(() => props.show, (v) => {
  if (v && props.enterprise) { loadHistory(1) }
})

const loadHistory = async (p: number) => {
  if (!props.enterprise) return
  loading.value = true; page.value = p
  try {
    const res = await adminAPI.enterprises.getBalanceHistory(props.enterprise.id, p, pageSize)
    history.value = res.items || []
    total.value = res.total || 0
  } finally { loading.value = false }
}
</script>
