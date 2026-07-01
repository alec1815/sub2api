<template>
  <AppLayout>
    <div class="space-y-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('enterprise.dashboard.title') }}</h1>

      <div v-if="loading" class="flex justify-center py-20">
        <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent" />
      </div>

      <template v-else>
        <!-- Row 1: Core Stats (4 cards) -->
        <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
          <div class="card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-emerald-100 p-2 dark:bg-emerald-900/30">
                <Icon name="dollar" size="md" class="text-emerald-600" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500">{{ t('admin.dashboard.balance') }}</p>
                <p class="text-xl font-bold text-emerald-600">${{ formatNumber(profile?.enterprise?.balance) }}</p>
                <p class="text-xs text-gray-500">{{ t('admin.dashboard.totalRecharged') }}: ${{ formatNumber(profile?.enterprise?.total_recharged) }}</p>
              </div>
            </div>
          </div>

          <div class="card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-blue-100 p-2 dark:bg-blue-900/30">
                <Icon name="key" size="md" class="text-blue-600" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500">{{ t('admin.dashboard.apiKeys') }}</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">{{ keyCount }}</p>
                <p class="text-xs text-gray-500">{{ t('admin.dashboard.accountCount') }}: {{ memberCount }}</p>
              </div>
            </div>
          </div>

          <div class="card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-green-100 p-2 dark:bg-green-900/30">
                <Icon name="chart" size="md" class="text-green-600" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500">{{ t('admin.dashboard.todayRequests') }}</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">{{ stats?.today_requests || 0 }}</p>
                <p class="text-xs text-gray-500">{{ t('common.total') }}: {{ formatNumber(stats?.total_requests) }}</p>
              </div>
            </div>
          </div>

          <div class="card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-amber-100 p-2 dark:bg-amber-900/30">
                <Icon name="cube" size="md" class="text-amber-600" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500">{{ t('admin.dashboard.todayTokens') }}</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">{{ formatTokens(stats?.today_tokens) }}</p>
                <p class="text-xs">
                  <span class="text-emerald-600">${{ stats?.today_actual_cost?.toFixed(4) || '0.0000' }}</span>
                  <span class="text-gray-400"> / </span>
                  <span class="text-gray-400">${{ stats?.today_cost?.toFixed(4) || '0.0000' }}</span>
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Row 2: Time range filter -->
        <div class="card p-4">
          <div class="flex flex-wrap items-center gap-4">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.dashboard.timeRange') }}:</span>
            <div class="flex items-center gap-2">
              <Select :model-value="dateRange" :options="dateRangeOptions" @update:model-value="onDateRangeChange" />
            </div>
            <button @click="loadSnapshot" :disabled="chartsLoading" class="btn btn-secondary">
              {{ t('common.refresh') }}
            </button>
            <div class="ml-auto flex items-center gap-2">
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.dashboard.granularity') }}:</span>
              <div class="w-28">
                <Select v-model="granularity" :options="granularityOptions" @update:model-value="loadSnapshot" />
              </div>
            </div>
          </div>
        </div>

        <!-- Row 3: Model distribution + Token trend -->
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <div class="card p-4">
            <h3 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.dashboard.modelDistribution') }}</h3>
            <div class="h-64">
              <ModelDistributionChart v-if="models.length || !chartsLoading" :model-stats="models" :loading="chartsLoading" />
              <div v-else class="flex h-full items-center justify-center text-sm text-gray-500 dark:text-gray-400">
                {{ t('admin.dashboard.noData') }}
              </div>
            </div>
          </div>
          <div class="card p-4">
            <h3 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.dashboard.tokenUsageTrend') }}</h3>
            <div class="h-64">
              <TokenUsageTrend v-if="trend.length || !chartsLoading" :trend-data="trend" :loading="chartsLoading" />
              <div v-else class="flex h-full items-center justify-center text-sm text-gray-500 dark:text-gray-400">
                {{ t('admin.dashboard.noData') }}
              </div>
            </div>
          </div>
        </div>

        <!-- Row 4: Recent usage (Top 12) -->
        <div class="card p-4">
          <h3 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.dashboard.recentUsage') }} (Top 12)</h3>
          <div class="overflow-y-auto" style="max-height: 16rem;">
            <table v-if="usersTrend.length > 0" class="w-full text-sm">
              <thead class="bg-gray-50 text-xs uppercase text-gray-500 dark:bg-dark-700">
                <tr>
                  <th class="px-3 py-2 text-left">User</th>
                  <th class="px-3 py-2 text-right">Requests</th>
                  <th class="px-3 py-2 text-right">Tokens</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="u in usersTrend" :key="u.user_id" class="border-t dark:border-dark-700">
                  <td class="px-3 py-2 text-gray-700 dark:text-gray-300">{{ u.email || u.username || ('User#' + u.user_id) }}</td>
                  <td class="px-3 py-2 text-right">{{ u.requests }}</td>
                  <td class="px-3 py-2 text-right">{{ formatTokens(u.tokens) }}</td>
                </tr>
              </tbody>
            </table>
            <div v-else class="flex h-32 items-center justify-center text-sm text-gray-500 dark:text-gray-400">
              {{ t('admin.dashboard.noData') }}
            </div>
          </div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import Select from '@/components/common/Select.vue'
import ModelDistributionChart from '@/components/charts/ModelDistributionChart.vue'
import TokenUsageTrend from '@/components/charts/TokenUsageTrend.vue'
import enterpriseAPI from '@/api/enterprise'

const { t } = useI18n()
const loading = ref(true)
const chartsLoading = ref(false)
const profile = ref<any>(null)
const stats = ref<any>(null)
const trend = ref<any[]>([])
const models = ref<any[]>([])
const usersTrend = ref<any[]>([])
const keyCount = ref(0)
const memberCount = ref(0)

const dateRange = ref('1')
const granularity = ref('day')
const dateRangeOptions = computed(() => [
  { value: '1', label: t('admin.dashboard.last24h') },
  { value: '7', label: t('admin.dashboard.last7d') },
  { value: '30', label: t('admin.dashboard.last30d') },
])
const granularityOptions = computed(() => [
  { value: 'day', label: t('admin.dashboard.day') },
  { value: 'hour', label: t('admin.dashboard.hour') },
])

const formatLD = (d: Date) => `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`
const startDate = computed(() => formatLD(new Date(Date.now() - Number(dateRange.value) * 86400000)))
const endDate = computed(() => formatLD(new Date()))

const formatNumber = (v: number) => (v == null ? '0' : Number(v).toLocaleString('en-US'))
const formatTokens = (v: number) => {
  if (v == null) return '0'
  const n = Number(v)
  if (n < 1000) return n.toString()
  if (n < 1e6) return (n / 1000).toFixed(1) + 'K'
  if (n < 1e9) return (n / 1e6).toFixed(1) + 'M'
  return (n / 1e9).toFixed(1) + 'B'
}

async function loadSnapshot() {
  chartsLoading.value = true
  try {
    const res = await enterpriseAPI.getEnterpriseDashboardSnapshot({
      start_date: startDate.value,
      end_date: endDate.value,
      granularity: granularity.value,
    }) as any
    const d = res.data ?? res
    stats.value = d.stats
    trend.value = d.trend || []
    models.value = d.models || []
    usersTrend.value = d.users_trend || []
  } finally { chartsLoading.value = false }
}

function onDateRangeChange(v: string | number | boolean | null) { dateRange.value = String(v); loadSnapshot() }

onMounted(async () => {
  try {
    profile.value = await enterpriseAPI.getProfile()
  } catch {}
  try { const r = await enterpriseAPI.listMembers(1, 1); memberCount.value = (r as any).total ?? 0 } catch {}
  try { const r = await enterpriseAPI.listKeys(1, 1); keyCount.value = (r as any).total ?? 0 } catch {}
  loading.value = false
  loadSnapshot()
})
</script>
