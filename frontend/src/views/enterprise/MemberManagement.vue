<template>
  <AppLayout>
    <TablePageLayout :title="t('enterprise.members.title')" :description="t('enterprise.members.description')">
      <!-- Filters -->
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <!-- Left: Search + Filters -->
          <div class="flex-1 sm:max-w-64">
            <input
              v-model="searchQuery"
              type="text"
              :placeholder="t('enterprise.members.searchPlaceholder')"
              class="input"
              @input="onSearchInput"
            />
          </div>
          <Select
            v-model="statusFilter"
            :options="filterStatusOptions"
            class="w-32"
            @change="loadData"
          />
          <Select
            v-model="roleFilter"
            :options="filterRoleOptions"
            class="w-32"
            @change="loadData"
          />

          <!-- Right: Action buttons -->
          <div class="flex flex-1 flex-wrap items-center justify-end gap-2">
            <button
              @click="loadData"
              :disabled="loading"
              class="btn btn-secondary"
              :title="t('common.refresh')"
            >
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
            <button class="btn btn-primary" @click="openCreateModal">
              <Icon name="plus" size="md" class="mr-1" />
              {{ t('enterprise.members.addMember') }}
            </button>
          </div>
        </div>
      </template>

      <!-- Table -->
      <DataTable
        :columns="columns"
        :data="members"
        :loading="loading"
        :empty-message="t('enterprise.members.noMembers')"
      >
        <template #cell-name="{ row }">
          <span class="text-sm font-medium text-gray-900 dark:text-white">{{ row.name || '-' }}</span>
        </template>
        <template #cell-email="{ row }">
          <span class="text-sm text-gray-600 dark:text-gray-400">{{ row.email }}</span>
        </template>
        <template #cell-role="{ row }">
          <span
            :class="[
              'badge text-xs',
              row.role === 'enterprise_admin' ? 'badge-primary' : 'badge-gray'
            ]"
          >
            {{ row.role === 'enterprise_admin' ? t('enterprise.members.roleAdmin') : t('enterprise.members.roleMember') }}
          </span>
        </template>
        <template #cell-status="{ row }">
          <span
            :class="[
              'badge text-xs',
              row.status === 'active' ? 'badge-success' : row.status === 'pending' ? 'badge-warning' : 'badge-danger'
            ]"
          >
            {{ t(`enterprise.members.status${row.status.charAt(0).toUpperCase() + row.status.slice(1)}`) }}
          </span>
        </template>
        <template #cell-department="{ row }">
          <span class="text-sm text-gray-600 dark:text-gray-400">{{ row.department_name || '-' }}</span>
        </template>
        <template #cell-concurrency="{ row }">
          <span class="text-sm text-gray-600 dark:text-gray-400">{{ row.concurrency }}</span>
        </template>
        <template #cell-joined_at="{ row }">
          <span class="text-sm text-gray-500 dark:text-gray-500">{{ formatDateTime(row.joined_at) }}</span>
        </template>
        <template #cell-actions="{ row }">
          <div class="flex items-center gap-2">
            <button
              class="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400"
              @click="openEditModal(row)"
            >
              <Icon name="edit" size="sm" />
            </button>
            <button
              v-if="row.status !== 'unbound'"
              class="text-sm text-danger-600 hover:text-danger-700 dark:text-danger-400"
              @click="confirmUnbind(row)"
            >
              <Icon name="ban" size="sm" />
            </button>
          </div>
        </template>
      </DataTable>

      <template #pagination>
        <Pagination
          v-if="pagination.total > 0"
          v-model:page="pagination.page"
          :total="pagination.total"
          :pageSize="pagination.page_size"
          @update:page="onPageChange"
        />
      </template>
    </TablePageLayout>

    <!-- Create / Edit Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showFormModal" class="modal-overlay" @click.self="closeFormModal">
          <div class="modal-content max-w-lg">
            <div class="modal-header">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ editingMember ? t('enterprise.members.editMember') : t('enterprise.members.addMember') }}
              </h3>
              <button class="text-gray-400 hover:text-gray-600" @click="closeFormModal">
                <Icon name="x" size="md" />
              </button>
            </div>
            <div class="modal-body space-y-4">
              <div v-if="!editingMember">
                <label class="input-label">{{ t('enterprise.members.fields.email') }} *</label>
                <input
                  v-model="form.email"
                  type="email"
                  class="input w-full"
                  :placeholder="t('enterprise.members.fields.emailPlaceholder')"
                />
              </div>
              <div v-if="!editingMember">
                <label class="input-label">{{ t('enterprise.members.fields.username') }}</label>
                <input
                  v-model="form.username"
                  type="text"
                  class="input w-full"
                  :placeholder="t('enterprise.members.fields.usernamePlaceholder')"
                />
              </div>
              <div v-if="!editingMember">
                <label class="input-label">{{ t('enterprise.members.fields.password') }}</label>
                <input
                  v-model="form.password"
                  type="text"
                  class="input w-full"
                  :placeholder="t('enterprise.members.fields.passwordPlaceholder')"
                />
              </div>
              <div>
                <label class="input-label">{{ t('enterprise.members.fields.department') }}</label>
                <input
                  v-model.number="form.department_id"
                  type="number"
                  class="input w-full"
                  :placeholder="t('enterprise.members.fields.departmentPlaceholder')"
                />
              </div>
              <div>
                <label class="input-label">{{ t('enterprise.members.fields.concurrency') }}</label>
                <input
                  v-model.number="form.concurrency"
                  type="number"
                  min="1"
                  class="input w-full"
                />
              </div>
              <div>
                <label class="input-label">{{ t('enterprise.members.fields.rpmLimit') }}</label>
                <input
                  v-model.number="form.rpm_limit"
                  type="number"
                  min="0"
                  class="input w-full"
                />
              </div>
              <div v-if="editingMember">
                <label class="input-label">{{ t('enterprise.members.fields.notes') }}</label>
                <textarea
                  v-model="form.notes"
                  rows="2"
                  class="input w-full"
                  :placeholder="t('enterprise.members.fields.notesPlaceholder')"
                ></textarea>
              </div>
            </div>
            <div class="modal-footer">
              <button class="btn btn-secondary" @click="closeFormModal">{{ t('common.cancel') }}</button>
              <button class="btn btn-primary" :disabled="submitting" @click="submitForm">
                {{ submitting ? t('common.saving') : t('common.save') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Confirm Unbind Dialog -->
    <ConfirmDialog
      :show="showUnbindConfirm"
      :title="t('enterprise.members.unbindTitle')"
      :message="t('enterprise.members.unbindMessage', { email: unbindTarget?.email ?? '' })"
      :confirm-text="t('common.confirm')"
      :cancel-text="t('common.cancel')"
      :dangerous="true"
      @confirm="executeUnbind"
      @cancel="showUnbindConfirm = false"
    />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { formatDateTime } from '@/utils/format'
import Icon from '@/components/icons/Icon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import Select from '@/components/common/Select.vue'
import enterpriseAdminAPI from '@/api/enterprise'
import type { EnterpriseMember } from '@/types/enterprise'
import type { Column } from '@/components/common/types'

const { t } = useI18n()
const appStore = useAppStore()

// ---- State ----
const loading = ref(false)
const submitting = ref(false)
const members = ref<EnterpriseMember[]>([])
const searchQuery = ref('')
const statusFilter = ref('')
const roleFilter = ref('')
let searchTimer: ReturnType<typeof setTimeout> | null = null
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

// ---- Options ----
const filterStatusOptions = computed(() => [
  { value: '', label: t('enterprise.members.allStatus') },
  { value: 'active', label: t('enterprise.members.statusActive') },
  { value: 'pending', label: t('enterprise.members.statusPending') },
  { value: 'unbound', label: t('enterprise.members.statusUnbound') },
])

const filterRoleOptions = computed(() => [
  { value: '', label: t('enterprise.members.allRoles') },
  { value: 'enterprise_admin', label: t('enterprise.members.roleAdmin') },
  { value: 'enterprise_member', label: t('enterprise.members.roleMember') },
])

// ---- Columns ----
const columns = computed<Column[]>(() => [
  { key: 'name', label: t('enterprise.members.columns.name'), sortable: false },
  { key: 'email', label: t('enterprise.members.columns.email'), sortable: false },
  { key: 'role', label: t('enterprise.members.columns.role'), sortable: false },
  { key: 'status', label: t('enterprise.members.columns.status'), sortable: false },
  { key: 'department', label: t('enterprise.members.columns.department'), sortable: false },
  { key: 'concurrency', label: t('enterprise.members.columns.concurrency'), sortable: false },
  { key: 'joined_at', label: t('enterprise.members.columns.joinedAt'), sortable: true },
  { key: 'actions', label: t('enterprise.members.columns.actions'), sortable: false },
])

// ---- Form Modal ----
const showFormModal = ref(false)
const editingMember = ref<EnterpriseMember | null>(null)
const form = reactive({
  email: '',
  username: '',
  password: '',
  department_id: undefined as number | undefined,
  concurrency: 1,
  rpm_limit: 0,
  notes: '',
})

// ---- Unbind Confirm ----
const showUnbindConfirm = ref(false)
const unbindTarget = ref<EnterpriseMember | null>(null)

// ---- Methods ----
async function loadData() {
  loading.value = true
  try {
    const filters: Record<string, string> = {}
    if (searchQuery.value.trim()) filters.search = searchQuery.value.trim()
    if (statusFilter.value) filters.status = statusFilter.value
    if (roleFilter.value) filters.role = roleFilter.value

    const res = await enterpriseAdminAPI.listMembers(pagination.page, pagination.page_size, filters)
    members.value = res.items ?? []
    pagination.total = res.total ?? 0
  } catch (err: any) {
    appStore.showToast('error', err?.message ?? t('common.loadError'))
  } finally {
    loading.value = false
  }
}

function onSearchInput() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    pagination.page = 1
    loadData()
  }, 300)
}

function onPageChange(page: number) {
  pagination.page = page
  loadData()
}

// ---- Form Modal ----
function openCreateModal() {
  editingMember.value = null
  form.email = ''
  form.username = ''
  form.password = ''
  form.department_id = undefined
  form.concurrency = 1
  form.rpm_limit = 0
  form.notes = ''
  showFormModal.value = true
}

function openEditModal(member: EnterpriseMember) {
  editingMember.value = member
  form.department_id = member.department_id ?? undefined
  form.concurrency = member.concurrency
  form.rpm_limit = member.rpm_limit
  form.notes = member.notes ?? ''
  showFormModal.value = true
}

function closeFormModal() {
  showFormModal.value = false
  editingMember.value = null
}

async function submitForm() {
  submitting.value = true
  try {
    if (editingMember.value) {
      await enterpriseAdminAPI.updateMember(editingMember.value.id, {
        department_id: form.department_id,
        concurrency: form.concurrency,
        rpm_limit: form.rpm_limit,
        notes: form.notes || undefined,
      })
      appStore.showToast('success', t('enterprise.members.updated'))
    } else {
      if (!form.email.trim()) {
        appStore.showToast('error', t('enterprise.members.emailRequired'))
        submitting.value = false
        return
      }
      await enterpriseAdminAPI.createMember({
        email: form.email.trim(),
        username: form.username.trim() || undefined,
        password: form.password || undefined,
        department_id: form.department_id,
        concurrency: form.concurrency,
        rpm_limit: form.rpm_limit,
      })
      appStore.showToast('success', t('enterprise.members.created'))
    }
    closeFormModal()
    await loadData()
  } catch (err: any) {
    appStore.showToast('error', err?.message ?? t('common.saveError'))
  } finally {
    submitting.value = false
  }
}

// ---- Unbind ----
function confirmUnbind(member: EnterpriseMember) {
  unbindTarget.value = member
  showUnbindConfirm.value = true
}

async function executeUnbind() {
  if (!unbindTarget.value) return
  try {
    await enterpriseAdminAPI.unbindMember(unbindTarget.value.id)
    appStore.showToast('success', t('enterprise.members.unbound'))
    showUnbindConfirm.value = false
    await loadData()
  } catch (err: any) {
    appStore.showToast('error', err?.message ?? t('common.error'))
  }
}

// ---- Init ----
onMounted(() => {
  loadData()
})
</script>





