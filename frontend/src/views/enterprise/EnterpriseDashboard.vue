<template>
  <AppLayout>
    <div class="space-y-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('enterprise.dashboard.title') }}</h1>
      <div v-if="loading" class="flex justify-center py-20"><div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent" /></div>
      <template v-else-if="profile">
        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 xl:grid-cols-4">
          <div class="card flex items-center gap-4">
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-blue-100 dark:bg-blue-900/30"><svg class="h-6 w-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z"/></svg></div>
            <div><p class="text-sm text-gray-500">{{ t('enterprise.dashboard.balance') }}</p><p class="text-xl font-bold">${{ parseFloat(profile.enterprise.balance ?? '0').toFixed(2) }}</p></div>
          </div>
          <div class="card flex items-center gap-4">
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-green-100 dark:bg-green-900/30"><svg class="h-6 w-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z"/></svg></div>
            <div><p class="text-sm text-gray-500">{{ t('enterprise.dashboard.members') }}</p><p class="text-xl font-bold">{{ memberCount }}</p></div>
          </div>
          <div class="card flex items-center gap-4">
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-purple-100 dark:bg-purple-900/30"><svg class="h-6 w-6 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 5.25a3 3 0 013 3m3 0a6 6 0 01-7.029 5.912c-.563-.097-1.159.026-1.563.43L10.5 17.25H8.25v2.25H6v2.25H2.25v-2.818c0-.597.237-1.17.659-1.591l6.499-6.499c.404-.404.527-1 .43-1.563A6 6 0 1121.75 8.25z"/></svg></div>
            <div><p class="text-sm text-gray-500">{{ t('enterprise.dashboard.keys') }}</p><p class="text-xl font-bold">{{ keyCount }}</p></div>
          </div>
          <div class="card flex items-center gap-4">
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-amber-100 dark:bg-amber-900/30"><svg class="h-6 w-6 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z"/></svg></div>
            <div><p class="text-sm text-gray-500">{{ t('enterprise.dashboard.usage') }}</p><p class="text-xl font-bold">${{ monthlyCost.toFixed(2) }}</p></div>
          </div>
        </div>
        <div class="card">
          <h2 class="mb-3 text-lg font-semibold">{{ t('enterprise.dashboard.enterpriseInfo') }}</h2>
          <div class="grid grid-cols-2 gap-4 text-sm sm:grid-cols-4">
            <div><span class="text-gray-500">{{ t('admin.enterprises.form.fullName') }}</span><p class="font-medium">{{ profile.enterprise.name }}</p></div>
            <div><span class="text-gray-500">{{ t('admin.enterprises.columns.status') }}</span><p class="font-medium">{{ profile.enterprise.status }}</p></div>
            <div><span class="text-gray-500">{{ t('enterprise.dashboard.totalRecharged') }}</span><p class="font-medium">${{ parseFloat(profile.enterprise.total_recharged ?? '0').toFixed(2) }}</p></div>
            <div><span class="text-gray-500">{{ t('enterprise.dashboard.myRole') }}</span><p class="font-medium">{{ profile.member?.role === 'enterprise_admin' ? t('admin.enterpriseMembers.role.admin') : t('admin.enterpriseMembers.role.member') }}</p></div>
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
    try {
      const membersRes = await enterpriseAPI.listMembers(1, 1)
      memberCount.value = (membersRes as any).total || 0
    } catch {}
    try {
      const keysRes = await enterpriseAPI.listKeys(1, 1)
      keyCount.value = (keysRes as any).total || 0
    } catch {}
    monthlyCost.value = p.monthly_usage?.total_cost || 0
  } finally { loading.value = false }
})
</script>
