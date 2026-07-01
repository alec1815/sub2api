<template>
  <AppLayout>
    <div class="space-y-6">
      <div v-if="loading" class="flex items-center justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent" />
      </div>
      <template v-else-if="profile">
        <!-- Stats Cards -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
          <div class="card">
            <div class="flex items-center gap-3">
              <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-100 dark:bg-blue-900/30">
                <Icon name="dollar" size="md" class="text-blue-600 dark:text-blue-400" />
              </div>
              <div>
                <p class="text-xs text-gray-500">{{ t('enterprise.dashboard.balance') }}</p>
                <p class="text-lg font-bold">${{ Number(profile.enterprise?.balance ?? 0).toFixed(2) }}</p>
              </div>
            </div>
          </div>
          <div class="card">
            <div class="flex items-center gap-3">
              <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-green-100 dark:bg-green-900/30">
                <Icon name="users" size="md" class="text-green-600 dark:text-green-400" />
              </div>
              <div>
                <p class="text-xs text-gray-500">{{ t('enterprise.dashboard.members') }}</p>
                <p class="text-lg font-bold">{{ memberCount }}</p>
              </div>
            </div>
          </div>
          <div class="card">
            <div class="flex items-center gap-3">
              <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-purple-100 dark:bg-purple-900/30">
                <Icon name="key" size="md" class="text-purple-600 dark:text-purple-400" />
              </div>
              <div>
                <p class="text-xs text-gray-500">{{ t('enterprise.dashboard.keys') }}</p>
                <p class="text-lg font-bold">{{ keyCount }}</p>
              </div>
            </div>
          </div>
          <div class="card">
            <div class="flex items-center gap-3">
              <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-amber-100 dark:bg-amber-900/30">
                <Icon name="chart" size="md" class="text-amber-600 dark:text-amber-400" />
              </div>
              <div>
                <p class="text-xs text-gray-500">{{ t('enterprise.dashboard.usage') }}</p>
                <p class="text-lg font-bold">${{ monthlyCost.toFixed(2) }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Enterprise Info -->
        <div class="card">
          <h2 class="mb-3 text-base font-semibold text-gray-700 dark:text-gray-300">{{ t('enterprise.dashboard.enterpriseInfo') }}</h2>
          <div class="grid grid-cols-2 gap-4 text-sm sm:grid-cols-3">
            <div><span class="text-gray-500">{{ t('admin.enterprises.form.fullName') }}</span><p class="mt-0.5 font-medium">{{ profile.enterprise?.name || '-' }}</p></div>
            <div><span class="text-gray-500">{{ t('enterprise.dashboard.totalRecharged') }}</span><p class="mt-0.5 font-medium">${{ Number(profile.enterprise?.total_recharged ?? 0).toFixed(2) }}</p></div>
            <div><span class="text-gray-500">{{ t('enterprise.dashboard.myRole') }}</span><p class="mt-0.5 font-medium">{{ profile.member?.role === 'enterprise_admin' ? t('admin.enterpriseMembers.role.admin') : t('admin.enterpriseMembers.role.member') }}</p></div>
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

onMounted(async () => {
  try {
    const p = await enterpriseAPI.getProfile()
    profile.value = p
    monthlyCost.value = Number(p.monthly_usage?.total_cost ?? 0)
    try {
      const membersRes = await enterpriseAPI.listMembers(1, 1)
      memberCount.value = (membersRes as any).total || 0
    } catch {}
    try {
      const keysRes = await enterpriseAPI.listKeys(1, 1)
      keyCount.value = (keysRes as any).total || 0
    } catch {}
  } finally { loading.value = false }
})
</script>
