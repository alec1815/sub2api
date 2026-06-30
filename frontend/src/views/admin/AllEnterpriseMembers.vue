<template>
  <AppLayout>
    <div class="flex flex-col items-center justify-center py-40">
      <div class="text-center">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">{{ t('admin.enterpriseMembers.allMembersTitle') }}</h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">{{ t('admin.enterpriseMembers.allMembersHint') }}</p>
        <router-link to="/admin/enterprises" class="mt-6 inline-flex items-center gap-2 rounded-lg bg-primary-500 px-5 py-2.5 text-sm font-medium text-white hover:bg-primary-600">
          {{ t('admin.enterpriseMembers.goToEnterprises') }}
        </router-link>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'

const { t } = useI18n()
</script>
