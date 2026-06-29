<template>
  <AppLayout>
    <div class="mx-auto max-w-3xl space-y-6">
      <!-- Page Header -->
      <div>
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
          {{ t('enterprise.settings.title') }}
        </h1>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ t('enterprise.settings.description') }}
        </p>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-600 border-t-transparent"></div>
      </div>

      <template v-else-if="settings">
        <!-- Enterprise Info Card -->
        <section class="card border border-gray-100 bg-white p-6 dark:border-dark-700 dark:bg-dark-900/50">
          <div class="mb-5 flex items-start justify-between gap-4">
            <div>
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('enterprise.settings.basicInfo') }}
              </h3>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                {{ t('enterprise.settings.basicInfoDescription') }}
              </p>
            </div>
            <button
              v-if="!editing"
              class="btn btn-secondary btn-sm"
              @click="startEdit"
            >
              <Icon name="edit" size="sm" class="mr-1" />
              {{ t('common.edit') }}
            </button>
          </div>

          <!-- Display Mode -->
          <div v-if="!editing" class="grid gap-4 sm:grid-cols-2">
            <InfoField :label="t('enterprise.settings.fields.name')" :value="settings.name" />
            <InfoField :label="t('enterprise.settings.fields.shortName')" :value="settings.short_name" />
            <InfoField :label="t('enterprise.settings.fields.creditCode')" :value="settings.credit_code" />
            <InfoField :label="t('enterprise.settings.fields.address')" :value="settings.address" />
            <InfoField :label="t('enterprise.settings.fields.contactName')" :value="settings.contact_name" />
            <InfoField :label="t('enterprise.settings.fields.contactPhone')" :value="settings.contact_phone" />
            <InfoField :label="t('enterprise.settings.fields.contactEmail')" :value="settings.contact_email" />
            <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
              <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                {{ t('enterprise.settings.fields.status') }}
              </p>
              <span :class="['badge mt-1', settings.status === 'active' ? 'badge-success' : 'badge-danger']">
                {{ settings.status === 'active' ? t('common.active') : t('common.disabled') }}
              </span>
            </div>
            <InfoField :label="t('enterprise.settings.fields.scale')" :value="settings.scale || '-'" />
            <InfoField :label="t('enterprise.settings.fields.industry')" :value="settings.industry || '-''" />
            <InfoField :label="t('enterprise.settings.fields.notes')" :value="settings.notes" />
          </div>

          <!-- Edit Mode -->
          <div v-else class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="input-label">{{ t('enterprise.settings.fields.name') }}</label>
              <input v-model="editForm.name" type="text" class="input w-full" />
            </div>
            <div>
              <label class="input-label">{{ t('enterprise.settings.fields.shortName') }}</label>
              <input v-model="editForm.short_name" type="text" class="input w-full" />
            </div>
            <div>
              <label class="input-label">{{ t('enterprise.settings.fields.address') }}</label>
              <input v-model="editForm.address" type="text" class="input w-full" />
            </div>
            <div>
              <label class="input-label">{{ t('enterprise.settings.fields.contactName') }}</label>
              <input v-model="editForm.contact_name" type="text" class="input w-full" />
            </div>
            <div>
              <label class="input-label">{{ t('enterprise.settings.fields.contactPhone') }}</label>
              <input v-model="editForm.contact_phone" type="text" class="input w-full" />
            </div>
            <div>
              <label class="input-label">{{ t('enterprise.settings.fields.contactEmail') }}</label>
              <input v-model="editForm.contact_email" type="email" class="input w-full" />
            </div>
            <div class="sm:col-span-2">
              <label class="input-label">{{ t('enterprise.settings.fields.notes') }}</label>
              <textarea v-model="editForm.notes" rows="3" class="input w-full"></textarea>
            </div>
            <div class="sm:col-span-2 flex justify-end gap-3 pt-2">
              <button class="btn btn-secondary" @click="cancelEdit">
                {{ t('common.cancel') }}
              </button>
              <button class="btn btn-primary" :disabled="saving" @click="saveSettings">
                {{ saving ? t('common.saving') : t('common.save') }}
              </button>
            </div>
          </div>
        </section>

        <!-- Stats Card -->
        <section class="card border border-gray-100 bg-white p-6 dark:border-dark-700 dark:bg-dark-900/50">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t('enterprise.settings.overview') }}
          </h3>
          <div class="mt-4 grid gap-4 sm:grid-cols-3">
            <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
              <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                {{ t('enterprise.settings.memberCount') }}
              </p>
              <p class="mt-1 text-xl font-semibold text-gray-900 dark:text-white">
                {{ settings.member_count }}
              </p>
            </div>
            <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
              <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                {{ t('enterprise.settings.balance') }}
              </p>
              <p class="mt-1 text-xl font-semibold text-gray-900 dark:text-white">
                ${{ Number(settings.balance).toFixed(2) }}
              </p>
            </div>
            <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
              <p class="text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
                {{ t('enterprise.settings.createdAt') }}
              </p>
              <p class="mt-1 text-xl font-semibold text-gray-900 dark:text-white">
                {{ formatDateTime(settings.created_at) }}
              </p>
            </div>
          </div>
        </section>
      </template>

      <EmptyState v-else :message="t('enterprise.settings.noData')" />
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { formatDateTime } from '@/utils/format'
import Icon from '@/components/icons/Icon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import enterpriseAdminAPI from '@/api/enterprise'
import type { EnterpriseSettings as EnterpriseSettingsType } from '@/types/enterprise'

const { t } = useI18n()
const appStore = useAppStore()

// ---- InfoField component (inline) ----
const InfoField = {
  props: { label: String, value: [String, Number] },
  setup(props: { label: string; value?: string | number }) {
    return () => h('div', {
      class: 'rounded-xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30'
    }, [
      h('p', {
        class: 'text-xs font-medium uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500'
      }, props.label),
      h('p', {
        class: 'mt-1 font-medium text-gray-900 dark:text-white truncate'
      }, props.value || '-'),
    ])
  }
}

// ---- State ----
const loading = ref(false)
const saving = ref(false)
const editing = ref(false)
const settings = ref<EnterpriseSettingsType | null>(null)

const editForm = reactive({
  name: '',
  short_name: '',
  address: '',
  contact_name: '',
  contact_phone: '',
  contact_email: '',
  notes: '',
})

// ---- Methods ----
async function loadData() {
  loading.value = true
  try {
    settings.value = await enterpriseAdminAPI.getSettings()
  } catch (err: any) {
    appStore.showToast?.(err?.message ?? t('common.loadError'), 'error')
  } finally {
    loading.value = false
  }
}

function startEdit() {
  if (!settings.value) return
  editForm.name = settings.value.name ?? ''
  editForm.short_name = settings.value.short_name ?? ''
  editForm.address = settings.value.address ?? ''
  editForm.contact_name = settings.value.contact_name ?? ''
  editForm.contact_phone = settings.value.contact_phone ?? ''
  editForm.contact_email = settings.value.contact_email ?? ''
  editForm.notes = settings.value.notes ?? ''
  editing.value = true
}

function cancelEdit() {
  editing.value = false
}

async function saveSettings() {
  saving.value = true
  try {
    const updated = await enterpriseAdminAPI.updateSettings({
      name: editForm.name,
      short_name: editForm.short_name,
      address: editForm.address,
      contact_name: editForm.contact_name,
      contact_phone: editForm.contact_phone,
      contact_email: editForm.contact_email,
      notes: editForm.notes,
    })
    settings.value = updated
    editing.value = false
    appStore.showToast?.(t('enterprise.settings.saved'), 'success')
  } catch (err: any) {
    appStore.showToast?.(err?.message ?? t('common.saveError'), 'error')
  } finally {
    saving.value = false
  }
}

// ---- Init ----
onMounted(() => {
  loadData()
})
</script>
