<template>
  <AppLayout>
    <div class="space-y-6 px-4 py-6 sm:px-6 lg:px-8">
      <h1 class="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">{{ t('enterprise.dashboard.title') }}</h1>

      <div v-if="loading" class="flex justify-center py-20">
        <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent" />
      </div>

      <template v-else-if="profile">
        <!-- Row 1: Stats Cards -->
        <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
          <div class="card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-emerald-100 p-2 dark:bg-emerald-900/30">
                <Icon name="dollar" size="md" class="text-emerald-600 dark:text-emerald-400" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('enterprise.dashboard.balance') }}</p>
                <p class="text-xl font-bold text-emerald-600 dark:text-emerald-400">${{ formatBalance(Number(profile.enterprise?.balance ?? 0)) }}</p>
              </div>
            </div>
          </div>

          <div class="card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-blue-100 p-2 dark:bg-blue-900/30">
                <Icon name="key" size="md" class="text-blue-600 dark:text-blue-400" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('enterprise.dashboard.keys') }}</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">{{ keyCount }}</p>
              </div>
            </div>
          </div>

          <div class="card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-green-100 p-2 dark:bg-green-900/30">
                <Icon name="users" size="md" class="text-green-600 dark:text-green-400" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('enterprise.dashboard.members') }}</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">{{ memberCount }}</p>
              </div>
            </div>
          </div>

          <div class="card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-amber-100 p-2 dark:bg-amber-900/30">
                <Icon name="chart" size="md" class="text-amber-600 dark:text-amber-400" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('enterprise.dashboard.usage') }}</p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">${{ Number(monthlyCost).toFixed(2) }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Row 2: Enterprise Info + Quick Actions -->
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <div class="card p-6">
            <h2 class="mb-4 text-base font-semibold text-gray-700 dark:text-gray-300">{{ profile.enterprise?.name || '-' }}</h2>
            <div class="space-y-3 text-sm">
              <div class="flex justify-between"><span class="text-gray-500">{{ t('enterprise.dashboard.totalRecharged') }}</span><span class="font-medium">${{ Number(profile.enterprise?.total_recharged ?? 0).toFixed(2) }}</span></div>
              <div class="flex justify-between"><span class="text-gray-500">{{ t('enterprise.dashboard.myRole') }}</span><span class="font-medium">{{ profile.member?.role === 'enterprise_admin' ? t('admin.enterpriseMembers.role.admin') : t('admin.enterpriseMembers.role.member') }}</span></div>
              <div class="flex justify-between"><span class="text-gray-500">{{ t('admin.enterprises.columns.status') }}</span><span class="font-medium">{{ profile.enterprise?.status || '-' }}</span></div>
            </div>
          </div>

          <div class="card p-6">
            <h2 class="mb-4 text-base font-semibold text-gray-700 dark:text-gray-300">{{ t('nav.enterpriseKeyManagement') }}</h2>
            <div class="flex flex-wrap gap-2">
              <router-link to="/enterprise/keys" class="btn btn-primary btn-sm">{{ t('enterprise.keys.createKey') }}</router-link>
              <router-link to="/enterprise/members" class="btn btn-secondary btn-sm">{{ t('enterprise.members.title') }}</router-link>
              <router-link to="/enterprise/departments" class="btn btn-secondary btn-sm">{{ t('nav.enterpriseDepartments') }}</router-link>
              <router-link to="/enterprise/finance" class="btn btn-secondary btn-sm">{{ t('nav.enterpriseFinance') }}</router-link>
              <router-link to="/enterprise/settings" class="btn btn-secondary btn-sm">{{ t('nav.enterpriseSettings') }}</router-link>
            </div>
          </div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import enterpriseAPI from '@/api/enterprise'

const { t } = useI18n()
const loading = ref(true)
const profile = ref<any>(null)
const memberCount = ref(0)
const keyCount = ref(0)
const monthlyCost = ref(0)

const formatBalance = (v: number) => {
  if (v === 0) return '0.00'
  const s = v.toFixed(8).replace(/\.?0+$/, '')
  const parts = s.split('.')
  if (parts.length === 1) return s + '.00'
  if (parts[1].length === 1) return s + '0'
  return s
}

onMounted(async () => {
  try {
    const p = await enterpriseAPI.getProfile()
    profile.value = p
    monthlyCost.value = Number(p.monthly_usage?.total_cost ?? 0)
  } catch { /* non-critical */ }
  try {
    const r = await enterpriseAPI.listMembers(1, 1)
    memberCount.value = (r as any).total ?? 0
  } catch {}
  try {
    const r = await enterpriseAPI.listKeys(1, 1)
    keyCount.value = (r as any).total ?? 0
  } catch {}
  loading.value = false
})
</script>
