<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- Page header -->
      <div class="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('admin.enterpriseMembers.title') }}</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {{ t('admin.enterpriseMembers.description') }}
            <span v-if="enterpriseName" class="font-medium text-gray-700 dark:text-gray-300"> — {{ enterpriseName }}</span>
          </p>
        </div>
        <div class="flex items-center gap-2">
          <button class="btn-secondary inline-flex items-center gap-2 rounded-lg px-4 py-2.5 text-sm font-semibold" @click="goBack">
            <Icon name="arrowLeft" size="sm" />
            {{ t('common.back') }}
          </button>
          <button class="btn-primary inline-flex items-center gap-2 rounded-lg px-4 py-2.5 text-sm font-semibold" @click="openCreateModal">
            <Icon name="plus" size="sm" />
            {{ t('admin.enterpriseMembers.registerMember') }}
          </button>
        </div>
      </div>

      <!-- Search & Filter -->
      <div class="card">
        <div class="flex flex-wrap items-center gap-3">
          <div class="flex-1 sm:max-w-64">
            <input
              v-model="searchQuery"
              type="text"
              :placeholder="t('admin.enterpriseMembers.searchPlaceholder')"
              class="input"
              @input="onSearchInput"
            />
          </div>
          <Select
            v-model="statusFilter"
            :options="filterStatusOptions"
            class="w-36"
            @change="loadMembers"
          />
          <div class="flex flex-1 items-center justify-end gap-2">
            <button
              @click="loadMembers"
              :disabled="loading"
              class="btn btn-secondary"
              :title="t('common.refresh')"
            >
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
          </div>
        </div>
      </div>

      <!-- Member table -->
      <div class="card overflow-hidden">
        <div v-if="loading" class="flex items-center justify-center py-20">
          <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent" />
        </div>
        <table v-else-if="members.length > 0" class="w-full min-w-max divide-y divide-gray-200 dark:divide-dark-700">
          <thead class="bg-gray-50 dark:bg-dark-800">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.enterpriseMembers.columns.member') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.enterpriseMembers.columns.role') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.enterpriseMembers.columns.department') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.enterpriseMembers.columns.status') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.enterpriseMembers.columns.joinedAt') }}</th>
              <th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dark-400">{{ t('admin.enterpriseMembers.columns.actions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 bg-white dark:divide-dark-700 dark:bg-dark-900">
            <tr v-for="member in members" :key="member.id" class="hover:bg-gray-50 dark:hover:bg-dark-800/50">
              <td class="px-6 py-4">
                <div class="flex flex-col">
                  <span class="text-sm font-medium text-gray-900 dark:text-white">{{ member.name }}</span>
                  <span class="text-xs text-gray-400 dark:text-gray-500">{{ member.email }}</span>
                </div>
              </td>
              <td class="px-6 py-4">
                <span :class="['inline-flex rounded-full px-2 py-0.5 text-xs font-medium', member.role === 'enterprise_admin' ? 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300' : 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300']">
                  {{ member.role === 'enterprise_admin' ? t('admin.enterpriseMembers.role.admin') : t('admin.enterpriseMembers.role.member') }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-600 dark:text-gray-300">{{ member.department_name || '-' }}</td>
              <td class="px-6 py-4">
                <span :class="['inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium', member.status === 'active' ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300' : member.status === 'pending' ? 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-300' : 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400']">
                  <span :class="['inline-block h-1.5 w-1.5 rounded-full', member.status === 'active' ? 'bg-green-500' : member.status === 'pending' ? 'bg-yellow-500' : 'bg-gray-400']" />
                  {{ member.status === 'active' ? t('admin.enterpriseMembers.status.active') : member.status === 'pending' ? t('admin.enterpriseMembers.status.pending') : t('admin.enterpriseMembers.status.unbound') }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500 dark:text-dark-400">{{ formatDateTime(member.joined_at) }}</td>
              <td class="px-6 py-4 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <button @click="openEditModal(member)" class="rounded-lg p-1.5 text-gray-400 hover:text-primary-600 hover:bg-gray-100 dark:hover:bg-dark-700 dark:hover:text-primary-400">
                    <Icon name="edit" size="sm" />
                  </button>
                  <button v-if="member.status !== 'unbound'" @click="confirmUnbind(member)" class="rounded-lg p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 dark:hover:text-red-400">
                    <Icon name="trash" size="sm" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-else class="flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
          <Icon name="users" size="xl" class="mb-4 h-16 w-16" />
          <p class="text-sm">{{ t('admin.enterpriseMembers.noMembers') }}</p>
          <button class="mt-4 text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400" @click="openCreateModal">{{ t('admin.enterpriseMembers.registerFirst') }}</button>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="pagination.total > 0" class="flex items-center justify-between">
        <span class="text-sm text-gray-500 dark:text-gray-400">{{ t('common.total', { total: pagination.total }) }}</span>
        <Pagination :page="pagination.page" :page-size="pagination.page_size" :total="pagination.total" @update:page="handlePageChange" @update:page-size="handlePageSizeChange" />
      </div>
    </div>

    <!-- Create Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4" @click.self="closeCreateModal">
          <div class="w-full max-w-md rounded-2xl bg-white shadow-2xl dark:bg-dark-800">
            <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dark-700">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.enterpriseMembers.registerMember') }}</h2>
              <button class="rounded-lg p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300" @click="closeCreateModal">
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <form class="space-y-4 px-6 py-5" @submit.prevent="submitCreateForm">
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterpriseMembers.form.email') }} <span class="text-red-500">*</span></label>
                <input v-model="createForm.email" type="email" class="input w-full" :placeholder="t('admin.enterpriseMembers.form.emailPlaceholder')" :disabled="submitting" />
                <p v-if="createErrors.email" class="mt-1 text-xs text-red-500">{{ createErrors.email }}</p>
              </div>
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterpriseMembers.form.username') }}</label>
                <input v-model="createForm.username" type="text" class="input w-full" :placeholder="t('admin.enterpriseMembers.form.usernamePlaceholder')" :disabled="submitting" />
              </div>
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterpriseMembers.form.password') }}</label>
                <input v-model="createForm.password" type="password" class="input w-full" :placeholder="t('admin.enterpriseMembers.form.passwordPlaceholder')" :disabled="submitting" />
              </div>
              <div class="flex items-center justify-end gap-3 pt-2">
                <button type="button" class="btn-secondary rounded-lg px-4 py-2.5 text-sm font-medium" :disabled="submitting" @click="closeCreateModal">{{ t('common.cancel') }}</button>
                <button type="submit" class="btn-primary rounded-lg px-5 py-2.5 text-sm font-semibold" :disabled="submitting">
                  <span v-if="submitting" class="mr-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                  {{ t('admin.enterpriseMembers.confirmCreate') }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Edit Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4" @click.self="closeEditModal">
          <div class="w-full max-w-md rounded-2xl bg-white shadow-2xl dark:bg-dark-800">
            <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dark-700">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.enterpriseMembers.editMember') }}</h2>
              <button class="rounded-lg p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300" @click="closeEditModal">
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <form class="space-y-4 px-6 py-5" @submit.prevent="submitEditForm">
              <div v-if="editingMember" class="rounded-lg bg-gray-50 px-4 py-3 dark:bg-dark-700">
                <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.enterpriseMembers.form.name') }}</p>
                <p class="text-sm font-medium text-gray-900 dark:text-white">{{ editingMember.name }}</p>
                <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.enterpriseMembers.form.email') }}</p>
                <p class="text-sm text-gray-700 dark:text-gray-300">{{ editingMember.email }}</p>
              </div>
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterpriseMembers.form.username') }}</label>
                <input v-model="editForm.username" type="text" class="input w-full" :placeholder="t('admin.enterpriseMembers.form.usernamePlaceholder')" :disabled="submitting" />
              </div>
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterpriseMembers.form.password') }} <span class="text-xs text-gray-400">({{ t('admin.enterpriseMembers.form.passwordHint') }})</span></label>
                <input v-model="editForm.password" type="password" class="input w-full" :placeholder="t('admin.enterpriseMembers.form.passwordPlaceholder')" :disabled="submitting" />
              </div>
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterpriseMembers.form.notes') }}</label>
                <textarea v-model="editForm.notes" rows="2" class="input w-full resize-none" :placeholder="t('admin.enterpriseMembers.form.notesPlaceholder')" :disabled="submitting" />
              </div>
              <div class="flex items-center justify-end gap-3 pt-2">
                <button type="button" class="btn-secondary rounded-lg px-4 py-2.5 text-sm font-medium" :disabled="submitting" @click="closeEditModal">{{ t('common.cancel') }}</button>
                <button type="submit" class="btn-primary rounded-lg px-5 py-2.5 text-sm font-semibold" :disabled="submitting">{{ t('common.save') }}</button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Unbind Confirm -->
    <ConfirmDialog :show="showUnbindDialog" :title="t('admin.enterpriseMembers.unbindTitle')" :message="t('admin.enterpriseMembers.unbindConfirm', { name: unbindingMember?.name })" :danger="true" @confirm="handleUnbind" @cancel="showUnbindDialog = false" />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { formatDateTime } from '@/utils/format'
import Icon from '@/components/icons/Icon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import Pagination from '@/components/common/Pagination.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import Select from '@/components/common/Select.vue'
import { adminAPI } from '@/api/admin'
import type { EnterpriseMember } from '@/types/enterprise'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()

const enterpriseId = Number(route.params.enterpriseId)
const enterpriseName = (route.query.enterpriseName as string) || ''

const loading = ref(false)
const submitting = ref(false)
const members = ref<EnterpriseMember[]>([])
const searchQuery = ref('')
const statusFilter = ref('')
let searchTimer: ReturnType<typeof setTimeout> | null = null
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

// Options
const filterStatusOptions = computed(() => [
  { value: '', label: t('admin.enterpriseMembers.filterAllStatus') },
  { value: 'active', label: t('admin.enterpriseMembers.status.active') },
  { value: 'pending', label: t('admin.enterpriseMembers.status.pending') },
  { value: 'unbound', label: t('admin.enterpriseMembers.status.unbound') },
])

onMounted(() => {
  if (!enterpriseId || enterpriseId <= 0) {
    appStore.showError(t('admin.enterpriseMembers.invalidEnterprise'))
    router.replace({ name: 'AdminEnterprises' })
    return
  }
  loadMembers()
})

async function loadMembers() {
  loading.value = true
  try {
    const filters: Record<string, unknown> = {}
    if (searchQuery.value) filters.search = searchQuery.value
    if (statusFilter.value) filters.status = statusFilter.value
    const result = await adminAPI.enterpriseMembers.list(enterpriseId, pagination.page, pagination.page_size, filters)
    members.value = result.items
    pagination.total = result.total
  } catch {
    appStore.showError(t('admin.enterpriseMembers.failedToLoad'))
  } finally {
    loading.value = false
  }
}

function onSearchInput() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => { pagination.page = 1; loadMembers() }, 350)
}

function handlePageChange(page: number) { pagination.page = page; loadMembers() }
function handlePageSizeChange(size: number) { pagination.page_size = size; pagination.page = 1; loadMembers() }
function goBack() { router.push({ name: 'AdminEnterprises' }) }

// ---- Create form ----
const showCreateModal = ref(false)
const createForm = ref({ email: '', username: '', password: '' })
const createErrors = ref<Record<string, string>>({})

function openCreateModal() {
  createForm.value = { email: '', username: '', password: '' }
  createErrors.value = {}
  showCreateModal.value = true
}

function closeCreateModal() { showCreateModal.value = false }

async function submitCreateForm() {
  createErrors.value = {}
  if (!createForm.value.email.trim()) { createErrors.value.email = t('admin.enterpriseMembers.form.emailRequired'); return }
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(createForm.value.email.trim())) { createErrors.value.email = t('admin.enterpriseMembers.form.emailInvalid'); return }
  submitting.value = true
  try {
    await adminAPI.enterpriseMembers.create(enterpriseId, {
      email: createForm.value.email.trim(),
      username: createForm.value.username.trim() || undefined,
      password: createForm.value.password || undefined,
    })
    appStore.showSuccess(t('admin.enterpriseMembers.createSuccess'))
    closeCreateModal()
    await loadMembers()
  } catch {
    appStore.showError(t('admin.enterpriseMembers.createFailed'))
  } finally {
    submitting.value = false
  }
}

// ---- Edit form ----
const showEditModal = ref(false)
const editingMember = ref<EnterpriseMember | null>(null)
const editForm = ref({ username: '', password: '', notes: '' })

function openEditModal(member: EnterpriseMember) {
  editingMember.value = member
  editForm.value = { username: member.name || '', password: '', notes: member.notes || '' }
  showEditModal.value = true
}

function closeEditModal() { showEditModal.value = false; editingMember.value = null }

async function submitEditForm() {
  if (!editingMember.value) return
  submitting.value = true
  try {
    await adminAPI.enterpriseMembers.update(enterpriseId, editingMember.value.id, {
      username: editForm.value.username.trim() || undefined,
      password: editForm.value.password || undefined,
      notes: editForm.value.notes.trim() || undefined,
    })
    appStore.showSuccess(t('admin.enterpriseMembers.updateSuccess'))
    closeEditModal()
    await loadMembers()
  } catch {
    appStore.showError(t('admin.enterpriseMembers.updateFailed'))
  } finally {
    submitting.value = false
  }
}

// ---- Unbind ----
const showUnbindDialog = ref(false)
const unbindingMember = ref<EnterpriseMember | null>(null)

function confirmUnbind(member: EnterpriseMember) { unbindingMember.value = member; showUnbindDialog.value = true }

async function handleUnbind() {
  if (!unbindingMember.value) return
  submitting.value = true
  try {
    await adminAPI.enterpriseMembers.unbind(enterpriseId, unbindingMember.value.id)
    appStore.showSuccess(t('admin.enterpriseMembers.unbindSuccess'))
    showUnbindDialog.value = false; unbindingMember.value = null
    await loadMembers()
  } catch {
    appStore.showError(t('admin.enterpriseMembers.unbindFailed'))
  } finally {
    submitting.value = false
  }
}
</script>
