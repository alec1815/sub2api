<template>
  <AppLayout>
    <div class="mx-auto max-w-5xl space-y-6">
      <!-- Page Header -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
            {{ t('enterprise.finance.title') }}
          </h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {{ t('enterprise.finance.description') }}
          </p>
        </div>
        <button class="btn btn-secondary btn-sm" @click="loadData">
          <Icon name="refresh" size="sm" class="mr-1" />
          {{ t('common.refresh') }}
        </button>
      </div>

      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-600 border-t-transparent"></div>
      </div>

      <template v-else>
        <!-- Balance Cards -->
        <div class="grid gap-4 sm:grid-cols-3">
          <div class="card border border-primary-100 bg-gradient-to-br from-primary-50 to-white p-6 dark:border-primary-900/40 dark:from-primary-950/30 dark:to-dark-900">
            <p class="text-xs font-medium uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">
              {{ t('enterprise.finance.currentBalance') }}
            </p>
            <p class="mt-2 text-3xl font-bold text-primary-700 dark:text-primary-400">
              ${{ formatBalance(finance?.balance) }}
            </p>
          </div>
          <div class="card border border-gray-100 bg-white p-6 dark:border-dark-700 dark:bg-dark-900/50">
            <p class="text-xs font-medium uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">
              {{ t('enterprise.finance.totalRecharged') }}
            </p>
            <p class="mt-2 text-3xl font-bold text-gray-900 dark:text-white">
              ${{ formatBalance(finance?.total_recharged) }}
            </p>
          </div>
          <div class="card border border-gray-100 bg-white p-6 dark:border-dark-700 dark:bg-dark-900/50">
            <p class="text-xs font-medium uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">
              {{ t('enterprise.finance.totalCost') }}
            </p>
            <p class="mt-2 text-3xl font-bold text-gray-900 dark:text-white">
              ${{ formatBalance(finance?.usage_summary?.total_cost) }}
            </p>
          </div>
        </div>

        <!-- Subscriptions -->
        <section v-if="finance?.subscriptions?.length" class="card border border-gray-100 bg-white p-6 dark:border-dark-700 dark:bg-dark-900/50">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t('enterprise.finance.subscriptionsTitle') }}
          </h3>
          <div class="mt-4 overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b border-gray-100 dark:border-dark-700">
                  <th class="whitespace-nowrap px-3 py-2 text-left font-medium text-gray-500">{{ t('enterprise.finance.planName') }}</th>
                  <th class="whitespace-nowrap px-3 py-2 text-left font-medium text-gray-500">{{ t('enterprise.finance.groupName') }}</th>
                  <th class="whitespace-nowrap px-3 py-2 text-left font-medium text-gray-500">{{ t('enterprise.finance.status') }}</th>
                  <th class="whitespace-nowrap px-3 py-2 text-left font-medium text-gray-500">{{ t('enterprise.finance.dailyUsage') }}</th>
                  <th class="whitespace-nowrap px-3 py-2 text-left font-medium text-gray-500">{{ t('enterprise.finance.monthlyUsage') }}</th>
                  <th class="whitespace-nowrap px-3 py-2 text-left font-medium text-gray-500">{{ t('enterprise.finance.period') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="sub in finance.subscriptions"
                  :key="sub.id"
                  class="border-b border-gray-50 dark:border-dark-800"
                >
                  <td class="px-3 py-3 font-medium text-gray-900 dark:text-white">{{ sub.plan_name }}</td>
                  <td class="px-3 py-3 text-gray-600 dark:text-gray-400">{{ sub.group_name }}</td>
                  <td class="px-3 py-3">
                    <span :class="['badge text-xs', sub.status === 'active' ? 'badge-success' : sub.status === 'expired' ? 'badge-danger' : 'badge-warning']">
                      {{ t(`enterprise.finance.status${sub.status.charAt(0).toUpperCase() + sub.status.slice(1)}`) }}
                    </span>
                  </td>
                  <td class="px-3 py-3 text-gray-600 dark:text-gray-400">${{ sub.daily_usage_usd ?? '-' }}</td>
                  <td class="px-3 py-3 text-gray-600 dark:text-gray-400">${{ sub.monthly_usage_usd ?? '-' }}</td>
                  <td class="px-3 py-3 text-xs text-gray-500">
                    {{ formatDateTime(sub.starts_at) }} ~ {{ formatDateTime(sub.expires_at) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <!-- Usage Breakdown -->
        <section v-if="finance?.usage_summary" class="grid gap-6 lg:grid-cols-2">
          <!-- By Member -->
          <div v-if="finance.usage_summary.by_member?.length" class="card border border-gray-100 bg-white p-6 dark:border-dark-700 dark:bg-dark-900/50">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('enterprise.finance.byMember') }}
            </h3>
            <div class="mt-4 space-y-3">
              <div
                v-for="item in finance.usage_summary.by_member.slice(0, 10)"
                :key="item.member_id"
                class="flex items-center justify-between rounded-lg bg-gray-50 px-4 py-2.5 dark:bg-dark-800"
              >
                <span class="text-sm text-gray-700 dark:text-gray-300">{{ item.member_name || `#${item.member_id}` }}</span>
                <div class="flex items-center gap-3">
                  <span class="text-xs text-gray-400">{{ item.calls }} {{ t('enterprise.finance.calls') }}</span>
                  <span class="text-sm font-medium text-gray-900 dark:text-white">${{ Number(item.cost).toFixed(4) }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- By Model -->
          <div v-if="finance.usage_summary.by_model?.length" class="card border border-gray-100 bg-white p-6 dark:border-dark-700 dark:bg-dark-900/50">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('enterprise.finance.byModel') }}
            </h3>
            <div class="mt-4 space-y-3">
              <div
                v-for="item in finance.usage_summary.by_model.slice(0, 10)"
                :key="item.model"
                class="flex items-center justify-between rounded-lg bg-gray-50 px-4 py-2.5 dark:bg-dark-800"
              >
                <span class="text-sm font-mono text-gray-700 dark:text-gray-300">{{ item.model }}</span>
                <div class="flex items-center gap-3">
                  <span class="text-xs text-gray-400">{{ item.calls }} {{ t('enterprise.finance.calls') }}</span>
                  <span class="text-sm font-medium text-gray-900 dark:text-white">${{ Number(item.cost).toFixed(4) }}</span>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- Empty State -->
        <EmptyState
          v-if="!finance"
          :message="t('enterprise.finance.noData')"
        />
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { formatDateTime } from '@/utils/format'
import Icon from '@/components/icons/Icon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import enterpriseAdminAPI from '@/api/enterprise'
import type { EnterpriseFinance } from '@/types/enterprise'

const { t } = useI18n()
const appStore = useAppStore()

// ---- State ----
const loading = ref(false)
const finance = ref<EnterpriseFinance | null>(null)

// ---- Methods ----
function formatBalance(val?: string): string {
  if (!val) return '0.00'
  return Number(val).toFixed(2)
}

async function loadData() {
  loading.value = true
  try {
    finance.value = await enterpriseAdminAPI.getFinance()
  } catch (err: any) {
    appStore.showToast?.(err?.message ?? t('common.loadError'), 'error')
  } finally {
    loading.value = false
  }
}

// ---- Init ----
onMounted(() => {
  loadData()
})
</script>
