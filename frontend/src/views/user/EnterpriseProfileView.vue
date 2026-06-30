<template>
  <AppLayout>
    <div data-testid="enterprise-profile-shell" class="mx-auto max-w-[950px] space-y-6">
      <!-- Enterprise Overview Hero -->
      <section
        data-testid="enterprise-profile-hero"
        class="card overflow-hidden border border-primary-100/80 bg-gradient-to-br from-primary-50 via-white to-amber-50/70 dark:border-primary-900/40 dark:from-primary-950/40 dark:via-dark-900 dark:to-dark-950"
      >
        <div class="px-6 py-6 md:px-8">
          <div class="flex flex-col gap-6 lg:flex-row lg:items-start">
            <!-- Enterprise Logo/Avatar -->
            <div
              class="flex h-20 w-20 shrink-0 items-center justify-center overflow-hidden rounded-[1.75rem] bg-gradient-to-br from-blue-500 to-indigo-600 text-2xl font-bold text-white shadow-lg shadow-blue-500/20"
            >
              <span>{{ enterpriseInitial }}</span>
            </div>

            <div class="min-w-0 flex-1 space-y-5">
              <div class="space-y-3">
                <div class="flex flex-wrap items-center gap-2">
                  <h2 class="truncate text-2xl font-semibold text-gray-900 dark:text-white">
                    {{ profile?.enterprise?.name || t('enterpriseProfile.unnamedEnterprise') }}
                  </h2>
                  <span :class="['badge', myRole === 'enterprise_admin' ? 'badge-primary' : 'badge-gray']">
                    {{ myRole === 'enterprise_admin' ? t('enterpriseProfile.roleAdmin') : t('enterpriseProfile.roleMember') }}
                  </span>
                  <span
                    v-if="profile?.enterprise?.industry"
                    class="badge badge-gray"
                  >
                    {{ profile.enterprise.industry }}
                  </span>
                </div>

                <div class="space-y-1">
                  <p v-if="profile?.enterprise?.contact_email" class="truncate text-sm text-gray-600 dark:text-gray-300">
                    <Icon name="mail" size="sm" class="mr-1 inline align-middle" />
                    {{ profile.enterprise.contact_email }}
                  </p>
                  <p v-if="profile?.enterprise?.contact_phone" class="truncate text-sm text-gray-500 dark:text-gray-400">
                    <Icon name="chatBubble" size="sm" class="mr-1 inline align-middle" />
                    {{ profile.enterprise.contact_phone }}
                  </p>
                </div>
              </div>

              <!-- Metrics -->
              <div class="grid gap-3 sm:grid-cols-4">
                <div
                  data-testid="enterprise-metric-role"
                  class="rounded-2xl bg-white/85 px-4 py-3 shadow-sm ring-1 ring-white/70 dark:bg-dark-900/60 dark:ring-dark-700"
                >
                  <p class="text-xs font-medium uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">
                    {{ t('enterpriseProfile.myRole') }}
                  </p>
                  <p class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
                    {{ myRole === 'enterprise_admin' ? t('enterpriseProfile.roleAdmin') : t('enterpriseProfile.roleMember') }}
                  </p>
                </div>
                <div
                  data-testid="enterprise-metric-department"
                  class="rounded-2xl bg-white/85 px-4 py-3 shadow-sm ring-1 ring-white/70 dark:bg-dark-900/60 dark:ring-dark-700"
                >
                  <p class="text-xs font-medium uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">
                    {{ t('enterpriseProfile.department') }}
                  </p>
                  <p class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
                    {{ profile?.my_department || '-' }}
                  </p>
                </div>
                <div
                  data-testid="enterprise-metric-monthly-calls"
                  class="rounded-2xl bg-white/85 px-4 py-3 shadow-sm ring-1 ring-white/70 dark:bg-dark-900/60 dark:ring-dark-700"
                >
                  <p class="text-xs font-medium uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">
                    {{ t('enterpriseProfile.monthlyCalls') }}
                  </p>
                  <p class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
                    {{ (profile?.monthly_usage?.total_calls ?? 0).toLocaleString() }}
                  </p>
                </div>
                <div
                  data-testid="enterprise-metric-joined"
                  class="rounded-2xl bg-white/85 px-4 py-3 shadow-sm ring-1 ring-white/70 dark:bg-dark-900/60 dark:ring-dark-700"
                >
                  <p class="text-xs font-medium uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">
                    {{ t('enterpriseProfile.joinedSince') }}
                  </p>
                  <p class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
                    {{ joinedSinceLabel }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Details section -->
      <div class="space-y-6">
        <div data-testid="enterprise-main-column" class="space-y-6">
          <!-- Basic Info Panel -->
          <section
            data-testid="enterprise-basics-panel"
            class="card border border-gray-100 bg-white/90 p-6 dark:border-dark-700 dark:bg-dark-900/50"
          >
            <div class="mb-5">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('enterpriseProfile.basicsTitle') }}
              </h3>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                {{ t('enterpriseProfile.basicsDescription') }}
              </p>
            </div>

            <div class="grid gap-4 sm:grid-cols-2">
              <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
                <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                  {{ t('enterpriseProfile.fields.name') }}
                </p>
                <p class="mt-1 font-medium text-gray-900 dark:text-white">
                  {{ profile?.enterprise?.name || '-' }}
                </p>
              </div>
              <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
                <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                  {{ t('enterpriseProfile.fields.shortName') }}
                </p>
                <p class="mt-1 font-medium text-gray-900 dark:text-white">
                  {{ profile?.enterprise?.short_name || '-' }}
                </p>
              </div>
              <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
                <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                  {{ t('enterpriseProfile.fields.industry') }}
                </p>
                <p class="mt-1 font-medium text-gray-900 dark:text-white">
                  {{ profile?.enterprise?.industry || '-' }}
                </p>
              </div>
              <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
                <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                  {{ t('enterpriseProfile.fields.size') }}
                </p>
                <p class="mt-1 font-medium text-gray-900 dark:text-white">
                  {{ sizeLabel || '-' }}
                </p>
              </div>
            </div>
          </section>

          <!-- Contact Info Panel -->
          <section
            data-testid="enterprise-contact-panel"
            class="card border border-gray-100 bg-white/90 p-6 dark:border-dark-700 dark:bg-dark-900/50"
          >
            <div class="mb-5">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('enterpriseProfile.contactTitle') }}
              </h3>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                {{ t('enterpriseProfile.contactDescription') }}
              </p>
            </div>

            <div class="grid gap-4 sm:grid-cols-2">
              <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
                <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                  {{ t('enterpriseProfile.fields.contactName') }}
                </p>
                <p class="mt-1 font-medium text-gray-900 dark:text-white">
                  {{ profile?.enterprise?.contact_name || '-' }}
                </p>
              </div>
              <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
                <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                  {{ t('enterpriseProfile.fields.contactEmail') }}
                </p>
                <p class="mt-1 font-medium text-gray-900 dark:text-white">
                  {{ profile?.enterprise?.contact_email || '-' }}
                </p>
              </div>
              <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
                <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                  {{ t('enterpriseProfile.fields.contactPhone') }}
                </p>
                <p class="mt-1 font-medium text-gray-900 dark:text-white">
                  {{ profile?.enterprise?.contact_phone || '-' }}
                </p>
              </div>
              <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
                <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                  {{ t('enterpriseProfile.fields.address') }}
                </p>
                <p class="mt-1 font-medium text-gray-900 dark:text-white truncate">
                  {{ profile?.enterprise?.address || '-' }}
                </p>
              </div>
            </div>
          </section>

          <!-- Usage Summary Panel -->
          <section
            v-if="profile?.monthly_usage"
            data-testid="enterprise-usage-panel"
            class="card border border-gray-100 bg-white/90 p-6 dark:border-dark-700 dark:bg-dark-900/50"
          >
            <div class="mb-5">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('enterpriseProfile.usageTitle') }}
              </h3>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                {{ t('enterpriseProfile.usageDescription') }}
              </p>
            </div>

            <div class="grid gap-4 sm:grid-cols-3">
              <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
                <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                  {{ t('enterpriseProfile.usageMetrics.totalRequests') }}
                </p>
                <p class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
                  {{ (profile.monthly_usage.total_calls ?? 0).toLocaleString() }}
                </p>
              </div>
              <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
                <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                  {{ t('enterpriseProfile.usageMetrics.totalTokens') }}
                </p>
                <p class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
                  -
                </p>
              </div>
              <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
                <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                  {{ t('enterpriseProfile.usageMetrics.totalCost') }}
                </p>
                <p class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
                  ${{ Number(profile.monthly_usage.total_cost ?? 0).toFixed(2) }}
                </p>
              </div>
            </div>
          </section>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-600 border-t-transparent"></div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import Icon from '@/components/icons/Icon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import enterpriseAdminAPI from '@/api/enterprise'
import type { EnterpriseProfile } from '@/types/enterprise'

const { t } = useI18n()
const appStore = useAppStore()

// ---- State ----
const loading = ref(false)
const profile = ref<EnterpriseProfile | null>(null)

// ---- Computed ----
const enterpriseInitial = computed(() => {
  return (profile.value?.enterprise?.name || 'E').charAt(0).toUpperCase()
})

const myRole = computed(() => {
  return profile.value?.my_role ?? 'enterprise_member'
})

const sizeLabel = computed(() => {
  const scale = profile.value?.enterprise?.scale
  if (!scale) return '-'
  const labels: Record<string, string> = {
    micro: t('enterpriseProfile.sizes.micro'),
    small: t('enterpriseProfile.sizes.small'),
    medium: t('enterpriseProfile.sizes.medium'),
    large: t('enterpriseProfile.sizes.large'),
  }
  return labels[scale] || scale
})

const joinedSinceLabel = computed(() => {
  const raw = profile.value?.my_joined_at?.trim()
  if (!raw) return '-'
  const date = new Date(raw)
  if (Number.isNaN(date.getTime())) return '-'
  return new Intl.DateTimeFormat(undefined, {
    year: 'numeric',
    month: 'short',
  }).format(date)
})

// ---- Methods ----
async function loadData() {
  loading.value = true
  try {
    profile.value = await enterpriseAdminAPI.getProfile()
  } catch (err: any) {
    appStore.showToast('error', err?.message ?? t('common.loadError'))
  } finally {
    loading.value = false
  }
}

// ---- Init ----
onMounted(() => {
  loadData()
})
</script>



