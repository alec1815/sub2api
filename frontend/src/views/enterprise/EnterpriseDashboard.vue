<template>
  <AppLayout>
    <div class="space-y-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('enterprise.dashboard.title') }}</h1>
      <div v-if="loading" class="flex justify-center py-20"><div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent" /></div>
      <template v-else-if="profile">
        <!-- Stats Cards -->
        <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
          <div class="card p-4"><div class="flex items-center gap-3"><div class="rounded-lg bg-emerald-100 p-2 dark:bg-emerald-900/30"><Icon name="dollar" size="md" class="text-emerald-600" /></div><div><p class="text-xs text-gray-500">{{ t('enterprise.dashboard.balance') }}</p><p class="text-xl font-bold text-emerald-600">${{ Number(profile.enterprise?.balance ?? 0).toFixed(2) }}</p></div></div></div>
          <div class="card p-4"><div class="flex items-center gap-3"><div class="rounded-lg bg-blue-100 p-2 dark:bg-blue-900/30"><Icon name="key" size="md" class="text-blue-600" /></div><div><p class="text-xs text-gray-500">{{ t('enterprise.dashboard.keys') }}</p><p class="text-xl font-bold text-gray-900 dark:text-white">{{ keyCount }}</p></div></div></div>
          <div class="card p-4"><div class="flex items-center gap-3"><div class="rounded-lg bg-green-100 p-2 dark:bg-green-900/30"><Icon name="users" size="md" class="text-green-600" /></div><div><p class="text-xs text-gray-500">{{ t('enterprise.dashboard.members') }}</p><p class="text-xl font-bold text-gray-900 dark:text-white">{{ memberCount }}</p></div></div></div>
          <div class="card p-4"><div class="flex items-center gap-3"><div class="rounded-lg bg-amber-100 p-2 dark:bg-amber-900/30"><Icon name="chart" size="md" class="text-amber-600" /></div><div><p class="text-xs text-gray-500">{{ t('enterprise.dashboard.usage') }}</p><p class="text-xl font-bold text-gray-900 dark:text-white">${{ Number(monthlyCost).toFixed(2) }}</p></div></div></div>
        </div>

        <!-- Charts (from enterprise usage APIs) -->
        <template v-if="!chartsLoading && (modelStats.length || trend.length)">
          <div class="card p-4">
            <div class="flex flex-wrap items-center gap-4 mb-4">
              <span class="text-sm font-medium">{{ t('admin.dashboard.timeRange') }}:</span>
              <div class="flex items-center gap-2"><span class="text-xs text-gray-500">{{ startDate }}</span><span class="text-xs">~</span><span class="text-xs text-gray-500">{{ endDate }}</span></div>
              <div class="ml-auto flex items-center gap-2">
                <span class="text-sm font-medium">{{ t('admin.dashboard.granularity') }}:</span>
                <Select v-model="granularity" :options="chartGranularityOptions" @change="loadCharts" />
              </div>
            </div>
          </div>
          <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
            <div class="card p-4"><h3 class="mb-3 text-sm font-semibold">{{ t('admin.dashboard.modelDistribution') }}</h3><div class="h-64"><model-distribution-chart :model-stats="modelStats" :loading="false" /></div></div>
            <div class="card p-4"><h3 class="mb-3 text-sm font-semibold">{{ t('admin.dashboard.tokenUsageTrend') }}</h3><div class="h-64"><token-usage-trend :trend-data="trend" :loading="false" /></div></div>
          </div>
        </template>

        <!-- Quick Actions -->
        <div class="card p-6">
          <h2 class="mb-3 text-base font-semibold text-gray-700 dark:text-gray-300">{{ profile.enterprise?.name || '-' }}</h2>
          <div class="flex flex-wrap gap-2">
            <router-link to="/enterprise/keys" class="btn btn-primary btn-sm">{{ t('enterprise.keys.createKey') }}</router-link>
            <router-link to="/enterprise/members" class="btn btn-secondary btn-sm">{{ t('enterprise.members.title') }}</router-link>
            <router-link to="/enterprise/finance" class="btn btn-secondary btn-sm">{{ t('nav.enterpriseFinance') }}</router-link>
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
const memberCount = ref(0)
const keyCount = ref(0)
const monthlyCost = ref(0)
const modelStats = ref<any[]>([])
const trend = ref<any[]>([])

const granularity = ref('hour')
const chartGranularityOptions = computed(() => [
  { value: 'day', label: t('admin.dashboard.day') },
  { value: 'hour', label: t('admin.dashboard.hour') },
])

const formatLD = (d: Date) => `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`
const startDate = ref(formatLD(new Date(Date.now() - 24*3600000)))
const endDate = ref(formatLD(new Date()))

async function loadCharts() {
  chartsLoading.value = true
  try {
    const [m, t2] = await Promise.all([
      enterpriseAPI.getEnterpriseModelStats(startDate.value, endDate.value),
      enterpriseAPI.getEnterpriseUsageTrend(startDate.value, endDate.value, granularity.value, 12),
    ])
    modelStats.value = (m as any).data || []
    trend.value = (t2 as any).data?.trend || []
  } catch { /* non-critical */ }
  chartsLoading.value = false
}

onMounted(async () => {
  try {
    const p = await enterpriseAPI.getProfile()
    profile.value = p
    monthlyCost.value = Number(p.monthly_usage?.total_cost ?? 0)
  } catch {}
  try { const r = await enterpriseAPI.listMembers(1, 1); memberCount.value = (r as any).total ?? 0 } catch {}
  try { const r = await enterpriseAPI.listKeys(1, 1); keyCount.value = (r as any).total ?? 0 } catch {}
  loading.value = false
  loadCharts()
})
</script>
