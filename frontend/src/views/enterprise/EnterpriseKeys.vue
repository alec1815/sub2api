<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-col gap-3">
          <div class="flex flex-wrap items-center gap-3">
            <SearchInput
              v-model="filterSearch"
              :placeholder="t('keys.searchPlaceholder')"
              class="w-full sm:w-64"
              @search="onFilterChange"
            />
            <Select :model-value="filterStatus" class="w-40" :options="statusFilterOptions" @update:model-value="onStatusFilterChange" />
          </div>
        </div>
      </template>

      <template #actions>
        <div class="flex justify-end gap-3">
          <button @click="loadData" :disabled="loading" class="btn btn-secondary" :title="t('common.refresh')">
            <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
          </button>
          <div class="relative" ref="columnDropdownRef">
            <button @click="showColumnDropdown = !showColumnDropdown" class="btn btn-secondary px-2 md:px-3" :title="t('keys.columnSettings')">
              <svg class="h-4 w-4 md:mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 4.5v15m6-15v15m-10.875 0h15.75c.621 0 1.125-.504 1.125-1.125V5.625c0-.621-.504-1.125-1.125-1.125H4.125C3.504 4.5 3 5.004 3 5.625v12.75c0 .621.504 1.125 1.125 1.125z" />
              </svg>
              <span class="hidden md:inline">{{ t('keys.columnSettings') }}</span>
            </button>
            <div v-if="showColumnDropdown" class="absolute right-0 top-full z-50 mt-1 max-h-80 w-48 overflow-y-auto rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-dark-600 dark:bg-dark-800">
              <button v-for="col in toggleableColumns" :key="col.key" @click="toggleColumn(col.key)" class="flex w-full items-center justify-between px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700">
                <span>{{ col.label }}</span>
                <Icon v-if="isColumnVisible(col.key)" name="check" size="sm" class="text-primary-500" :stroke-width="2" />
              </button>
            </div>
          </div>
          <button @click="openCreateModal" class="btn btn-primary">
            <Icon name="plus" size="md" class="mr-2" />
            {{ t('enterprise.keys.createKey') }}
          </button>
        </div>
      </template>

      <template #table>
        <DataTable :columns="visibleColumns" :data="filteredKeys" :loading="loading" :server-side-sort="true" default-sort-key="created_at" default-sort-order="desc" @sort="onSortChange">
          <template #cell-name="{ value }">
            <span class="font-medium text-gray-900 dark:text-white">{{ value }}</span>
          </template>
          <template #cell-key_prefix="{ value }">
            <code class="text-xs">{{ value }}***</code>
          </template>
          <template #cell-status="{ row }">
            <span :class="['badge text-xs', row.status === 'active' ? 'badge-success' : 'badge-danger']">
              {{ row.status === 'active' ? t('common.active') : t('common.disabled') }}
            </span>
          </template>
          <template #cell-assigned_to="{ row }">
            <span class="text-sm text-gray-600 dark:text-gray-400">{{ row.assigned_member_name || row.assigned_member_email || '-' }}</span>
          </template>
          <template #cell-groups="{ row }">
            <div class="flex flex-wrap gap-1">
              <span v-for="g in row.groups" :key="g.id" class="badge badge-gray text-xs">{{ g.name }}</span>
              <span v-if="!row.groups?.length" class="text-sm text-gray-400">-</span>
            </div>
          </template>
          <template #cell-bound_tool="{ row }">
            <span class="text-sm text-gray-600 dark:text-gray-400">{{ row.bound_tool || row.usage_purpose || '-' }}</span>
          </template>
          <template #cell-quota="{ row }">
            <span class="text-sm text-gray-600 dark:text-gray-400">{{ row.quota_used ?? 0 }} / {{ row.quota || '∞' }}</span>
          </template>
          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-500 dark:text-gray-500">{{ formatDateTime(value) }}</span>
          </template>
          <template #cell-actions="{ row }">
            <div class="flex items-center gap-2">
              <button v-if="row.status === 'active'" class="text-sm text-amber-600 hover:text-amber-700" :title="t('enterprise.keys.disable')" @click="confirmToggle(row)"><Icon name="ban" size="sm" /></button>
              <button v-else class="text-sm text-success-600 hover:text-success-700" :title="t('enterprise.keys.enable')" @click="confirmToggle(row)"><Icon name="checkCircle" size="sm" /></button>
              <button class="text-sm text-danger-600 hover:text-danger-700" :title="t('enterprise.keys.delete')" @click="confirmDelete(row)"><Icon name="trash" size="sm" /></button>
            </div>
          </template>
          <template #empty>
            <EmptyState :title="t('enterprise.keys.noKeys')" :description="t('enterprise.keys.createFirstHint')" :action-text="t('enterprise.keys.createKey')" @action="openCreateModal" />
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination v-if="pagination.total > 0" v-model:page="pagination.page" :total="pagination.total" :page-size="pagination.page_size" @update:page="onPageChange" />
      </template>
    </TablePageLayout>

    <!-- Create Key Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showCreateModal" class="modal-overlay" @click.self="closeCreateModal">
          <div class="modal-content max-w-lg">
            <div class="modal-header">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('enterprise.keys.createKey') }}</h3>
              <button class="text-gray-400 hover:text-gray-600" @click="closeCreateModal"><Icon name="x" size="md" /></button>
            </div>
            <div class="modal-body space-y-4">
              <div v-if="createdKey" class="rounded-xl border border-success-200 bg-success-50 p-4 dark:border-success-800 dark:bg-success-900/20">
                <p class="text-sm font-medium text-success-700 dark:text-success-400">{{ t('enterprise.keys.keyCreated') }}</p>
                <code class="mt-2 block break-all rounded bg-white dark:bg-dark-900 p-2 text-xs text-gray-900 dark:text-white">{{ createdKey.key_full }}</code>
                <p class="mt-2 text-xs text-success-600">{{ t('enterprise.keys.keyWarning') }}</p>
                <button class="btn btn-primary btn-sm mt-3 w-full" @click="closeCreateModal">{{ t('common.close') }}</button>
              </div>
              <template v-else>
                <div><label class="input-label">{{ t('enterprise.keys.fields.name') }} *</label><input v-model="createForm.name" type="text" class="input w-full" :placeholder="t('enterprise.keys.fields.namePlaceholder')" /></div>
                <div><label class="input-label">{{ t('enterprise.keys.fields.assignedTo') }}</label><input v-model.number="createForm.assigned_to" type="number" class="input w-full" :placeholder="t('enterprise.keys.fields.assignedToPlaceholder')" /></div>
                <div><label class="input-label">{{ t('enterprise.keys.fields.usagePurpose') }}</label><input v-model="createForm.usage_purpose" type="text" class="input w-full" :placeholder="t('enterprise.keys.fields.usagePurposePlaceholder')" /></div>
                <div>
                  <label class="input-label">{{ t('enterprise.keys.fields.boundTool') }}</label>
                  <select v-model="createForm.bound_tool" class="input w-full">
                    <option value="">{{ t('enterprise.keys.fields.boundToolNone') }}</option>
                    <option value="cursor">Cursor</option>
                    <option value="trae">Trae</option>
                    <option value="claude_code">Claude Code</option>
                    <option value="codex">Codex</option>
                    <option value="opencode">OpenCode</option>
                    <option value="pixso">Pixso</option>
                    <option value="other">{{ t('common.other') }}</option>
                  </select>
                </div>
                <div><label class="input-label">{{ t('enterprise.keys.fields.quota') }}</label><input v-model.number="createForm.quota" type="number" step="0.01" class="input w-full" :placeholder="t('enterprise.keys.fields.quotaPlaceholder')" /></div>
              </template>
            </div>
            <div v-if="!createdKey" class="modal-footer">
              <button class="btn btn-secondary" @click="closeCreateModal">{{ t('common.cancel') }}</button>
              <button class="btn btn-primary" :disabled="submitting" @click="submitCreate">{{ submitting ? t('common.creating') : t('common.create') }}</button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <ConfirmDialog :show="showToggleConfirm" :title="toggleTarget?.status === 'active' ? t('enterprise.keys.disableKey') : t('enterprise.keys.enableKey')" :message="toggleTarget?.status === 'active' ? t('enterprise.keys.disableMessage', { name: toggleTarget?.name ?? '' }) : t('enterprise.keys.enableMessage', { name: toggleTarget?.name ?? '' })" @confirm="executeToggle" @cancel="showToggleConfirm = false" />
    <ConfirmDialog :show="showDeleteConfirm" :title="t('enterprise.keys.deleteKey')" :message="t('enterprise.keys.deleteMessage', { name: deleteTarget?.name ?? '' })" :confirm-text="t('common.delete')" :dangerous="true" @confirm="executeDelete" @cancel="showDeleteConfirm = false" />
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
import SearchInput from '@/components/common/SearchInput.vue'
import Select from '@/components/common/Select.vue'
import Pagination from '@/components/common/Pagination.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import enterpriseAdminAPI from '@/api/enterprise'
import type { EnterpriseKey, CreateEnterpriseKeyResponse } from '@/types/enterprise'
import type { Column } from '@/components/common/types'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(false)
const submitting = ref(false)
const keys = ref<EnterpriseKey[]>([])
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

// Filters
const filterSearch = ref('')
const filterStatus = ref('')
const statusFilterOptions = computed(() => [
  { value: '', label: t('admin.users.allStatus') },
  { value: 'active', label: t('common.active') },
  { value: 'disabled', label: t('common.disabled') },
])

// Columns
const allColumns = computed<Column[]>(() => [
  { key: 'name', label: t('enterprise.keys.columns.name'), sortable: true },
  { key: 'key_prefix', label: t('enterprise.keys.columns.keyPrefix'), sortable: false },
  { key: 'status', label: t('enterprise.keys.columns.status'), sortable: true },
  { key: 'assigned_to', label: t('enterprise.keys.columns.assignedTo'), sortable: false },
  { key: 'groups', label: t('enterprise.keys.columns.groups'), sortable: false },
  { key: 'bound_tool', label: t('enterprise.keys.columns.boundTool'), sortable: false },
  { key: 'quota', label: t('enterprise.keys.columns.quota'), sortable: false },
  { key: 'created_at', label: t('enterprise.keys.columns.createdAt'), sortable: true },
  { key: 'actions', label: t('enterprise.keys.columns.actions'), sortable: false },
])

// Column settings
const hiddenColumns = reactive<Set<string>>(new Set())
const COLUMN_KEY = 'enterprise-keys-hidden-columns'
const FORCED_VISIBLE_COLUMNS = new Set(['name', 'actions'])
const toggleableColumns = computed(() => allColumns.value.filter(c => !FORCED_VISIBLE_COLUMNS.has(c.key)))
const visibleColumns = computed<Column[]>(() => allColumns.value.filter(c => FORCED_VISIBLE_COLUMNS.has(c.key) || !hiddenColumns.has(c.key)))
const showColumnDropdown = ref(false)
const columnDropdownRef = ref<HTMLElement | null>(null)

const isColumnVisible = (key: string) => !hiddenColumns.has(key)
const toggleColumn = (key: string) => {
  if (hiddenColumns.has(key)) hiddenColumns.delete(key)
  else hiddenColumns.add(key)
  localStorage.setItem(COLUMN_KEY, JSON.stringify([...hiddenColumns]))
}

function loadSavedColumns() {
  try {
    const saved = localStorage.getItem(COLUMN_KEY)
    if (saved) JSON.parse(saved).forEach((k: string) => hiddenColumns.add(k))
  } catch { /* ignore */ }
}

// Modal state
const showCreateModal = ref(false)
const showToggleConfirm = ref(false)
const showDeleteConfirm = ref(false)
const createdKey = ref<CreateEnterpriseKeyResponse | null>(null)
const toggleTarget = ref<EnterpriseKey | null>(null)
const deleteTarget = ref<EnterpriseKey | null>(null)

const createForm = reactive({
  name: '',
  assigned_to: undefined as number | undefined,
  usage_purpose: '',
  bound_tool: '',
  quota: undefined as number | undefined,
})

// Filtered data (client-side for search; server-side for status/pagination)
const filteredKeys = computed(() => {
  let result = keys.value
  if (filterStatus.value) result = result.filter(k => k.status === filterStatus.value)
  if (filterSearch.value) {
    const q = filterSearch.value.toLowerCase()
    result = result.filter(k => k.name.toLowerCase().includes(q) || (k.key ?? '').toLowerCase().includes(q))
  }
  return result
})

async function loadData() {
  loading.value = true
  try {
    const res = await enterpriseAdminAPI.listKeys(pagination.page, pagination.page_size)
    keys.value = res.items ?? []
    pagination.total = res.total ?? 0
  } catch (err: any) {
    appStore.showError(err?.message ?? t('common.loadError'))
  } finally { loading.value = false }
}

function onFilterChange() { pagination.page = 1; loadData() }
function onStatusFilterChange(v: string | number | boolean | null) { filterStatus.value = String(v); pagination.page = 1; loadData() }
function onPageChange(p: number) { pagination.page = p; loadData() }
function onSortChange() { loadData() }

function openCreateModal() {
  createdKey.value = null
  createForm.name = ''
  createForm.assigned_to = undefined
  createForm.usage_purpose = ''
  createForm.bound_tool = ''
  createForm.quota = undefined
  showCreateModal.value = true
}

function closeCreateModal() {
  showCreateModal.value = false
  if (createdKey.value) { createdKey.value = null; loadData() }
}

async function submitCreate() {
  if (!createForm.name.trim()) { appStore.showError(t('enterprise.keys.nameRequired')); return }
  submitting.value = true
  try {
    const res = await enterpriseAdminAPI.createKey({
      name: createForm.name.trim(),
      group_ids: [],
      assigned_to: createForm.assigned_to,
      usage_purpose: createForm.usage_purpose.trim() || undefined,
      bound_tool: (createForm.bound_tool as any) || undefined,
      quota: createForm.quota,
    })
    createdKey.value = res
  } catch (err: any) {
    appStore.showError(err?.message ?? t('common.saveError'))
  } finally { submitting.value = false }
}

function confirmToggle(k: EnterpriseKey) { toggleTarget.value = k; showToggleConfirm.value = true }
async function executeToggle() {
  if (!toggleTarget.value) return
  try {
    await enterpriseAdminAPI.toggleKey(toggleTarget.value.id)
    showToggleConfirm.value = false
    await loadData()
  } catch (err: any) { appStore.showError(err?.message ?? t('common.error')) }
}

function confirmDelete(k: EnterpriseKey) { deleteTarget.value = k; showDeleteConfirm.value = true }
async function executeDelete() {
  if (!deleteTarget.value) return
  try {
    await enterpriseAdminAPI.deleteKey(deleteTarget.value.id)
    showDeleteConfirm.value = false
    await loadData()
  } catch (err: any) { appStore.showError(err?.message ?? t('common.error')) }
}

onMounted(() => {
  loadSavedColumns()
  loadData()
})
</script>
