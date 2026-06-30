<template>
  <AppLayout>
    <TablePageLayout :title="t('enterprise.keys.title')" :description="t('enterprise.keys.description')">
      <!-- Filters -->
      <template #filters>
        <div class="flex items-center gap-3">
          <button class="btn btn-primary btn-sm" @click="openCreateModal">
            <Icon name="plus" size="sm" class="mr-1" />
            {{ t('enterprise.keys.createKey') }}
          </button>
          <button class="btn btn-secondary btn-sm" @click="loadData">
            <Icon name="refresh" size="sm" class="mr-1" />
            {{ t('common.refresh') }}
          </button>
        </div>
      </template>

      <!-- Table -->
      <DataTable
        :columns="columns"
        :data="keys"
        :loading="loading"
        :empty-message="t('enterprise.keys.noKeys')"
      >
        <template #cell-name="{ row }">
          <span class="text-sm font-medium text-gray-900 dark:text-white">{{ row.name }}</span>
        </template>
        <template #cell-key_prefix="{ row }">
          <code class="text-xs bg-gray-100 dark:bg-dark-800 px-1.5 py-0.5 rounded">{{ row.key_prefix }}***</code>
        </template>
        <template #cell-status="{ row }">
          <span :class="['badge text-xs', row.status === 'active' ? 'badge-success' : 'badge-danger']">
            {{ row.status === 'active' ? t('common.active') : t('common.disabled') }}
          </span>
        </template>
        <template #cell-assigned_to="{ row }">
          <span class="text-sm text-gray-600 dark:text-gray-400">
            {{ row.assigned_member_name || row.assigned_member_email || '-' }}
          </span>
        </template>
        <template #cell-groups="{ row }">
          <div class="flex flex-wrap gap-1">
            <span
              v-for="g in row.groups"
              :key="g.id"
              class="badge badge-gray text-xs"
            >
              {{ g.name }}
            </span>
            <span v-if="!row.groups?.length" class="text-sm text-gray-400">-</span>
          </div>
        </template>
        <template #cell-bound_tool="{ row }">
          <span class="text-sm text-gray-600 dark:text-gray-400">{{ row.bound_tool || row.usage_purpose || '-' }}</span>
        </template>
        <template #cell-quota="{ row }">
          <span class="text-sm text-gray-600 dark:text-gray-400">
            {{ row.quota_used ?? 0 }} / {{ row.quota || '∞' }}
          </span>
        </template>
        <template #cell-created_at="{ row }">
          <span class="text-sm text-gray-500 dark:text-gray-500">{{ formatDateTime(row.created_at) }}</span>
        </template>
        <template #cell-actions="{ row }">
          <div class="flex items-center gap-2">
            <button
              v-if="row.status === 'active'"
              class="text-sm text-amber-600 hover:text-amber-700 dark:text-amber-400"
              :title="t('enterprise.keys.disable')"
              @click="confirmToggle(row)"
            >
              <Icon name="ban" size="sm" />
            </button>
            <button
              v-else
              class="text-sm text-success-600 hover:text-success-700 dark:text-success-400"
              :title="t('enterprise.keys.enable')"
              @click="confirmToggle(row)"
            >
              <Icon name="checkCircle" size="sm" />
            </button>
            <button
              class="text-sm text-danger-600 hover:text-danger-700 dark:text-danger-400"
              :title="t('enterprise.keys.delete')"
              @click="confirmDelete(row)"
            >
              <Icon name="trash" size="sm" />
            </button>
          </div>
        </template>
      </DataTable>

      <template #pagination>
        <Pagination
          v-if="pagination.total > 0"
          v-model:current="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          @change="onPageChange"
        />
      </template>
    </TablePageLayout>

    <!-- Create Key Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showCreateModal" class="modal-overlay" @click.self="closeCreateModal">
          <div class="modal-content max-w-lg">
            <div class="modal-header">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('enterprise.keys.createKey') }}
              </h3>
              <button class="text-gray-400 hover:text-gray-600" @click="closeCreateModal">
                <Icon name="x" size="md" />
              </button>
            </div>
            <div class="modal-body space-y-4">
              <!-- Created key display -->
              <div v-if="createdKey" class="rounded-xl border border-success-200 bg-success-50 p-4 dark:border-success-800 dark:bg-success-900/20">
                <p class="text-sm font-medium text-success-700 dark:text-success-400">
                  {{ t('enterprise.keys.keyCreated') }}
                </p>
                <code class="mt-2 block break-all rounded bg-white dark:bg-dark-900 p-2 text-xs text-gray-900 dark:text-white">
                  {{ createdKey.key_full }}
                </code>
                <p class="mt-2 text-xs text-success-600 dark:text-success-500">
                  {{ t('enterprise.keys.keyWarning') }}
                </p>
                <button class="btn btn-primary btn-sm mt-3 w-full" @click="closeCreateModal">
                  {{ t('common.close') }}
                </button>
              </div>

              <!-- Create form -->
              <template v-else>
                <div>
                  <label class="input-label">{{ t('enterprise.keys.fields.name') }} *</label>
                  <input
                    v-model="createForm.name"
                    type="text"
                    class="input w-full"
                    :placeholder="t('enterprise.keys.fields.namePlaceholder')"
                  />
                </div>
                <div>
                  <label class="input-label">{{ t('enterprise.keys.fields.groupIds') }}</label>
                  <input
                    v-model="createForm.groupIdsInput"
                    type="text"
                    class="input w-full"
                    :placeholder="t('enterprise.keys.fields.groupIdsPlaceholder')"
                  />
                  <p class="mt-1 text-xs text-gray-400">{{ t('enterprise.keys.fields.groupIdsHint') }}</p>
                </div>
                <div>
                  <label class="input-label">{{ t('enterprise.keys.fields.assignedTo') }}</label>
                  <input
                    v-model.number="createForm.assigned_to"
                    type="number"
                    class="input w-full"
                    :placeholder="t('enterprise.keys.fields.assignedToPlaceholder')"
                  />
                </div>
                <div>
                  <label class="input-label">{{ t('enterprise.keys.fields.usagePurpose') }}</label>
                  <input
                    v-model="createForm.usage_purpose"
                    type="text"
                    class="input w-full"
                    :placeholder="t('enterprise.keys.fields.usagePurposePlaceholder')"
                  />
                </div>
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
                <div>
                  <label class="input-label">{{ t('enterprise.keys.fields.quota') }}</label>
                  <input
                    v-model.number="createForm.quota"
                    type="number"
                    step="0.01"
                    class="input w-full"
                    :placeholder="t('enterprise.keys.fields.quotaPlaceholder')"
                  />
                </div>
              </template>
            </div>
            <div v-if="!createdKey" class="modal-footer">
              <button class="btn btn-secondary" @click="closeCreateModal">{{ t('common.cancel') }}</button>
              <button class="btn btn-primary" :disabled="submitting" @click="submitCreate">
                {{ submitting ? t('common.creating') : t('common.create') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Confirm Toggle Dialog -->
    <ConfirmDialog
      :show="showToggleConfirm"
      :title="toggleTarget?.status === 'active' ? t('enterprise.keys.disableKey') : t('enterprise.keys.enableKey')"
      :message="toggleTarget?.status === 'active' ? t('enterprise.keys.disableMessage', { name: toggleTarget?.name ?? '' }) : t('enterprise.keys.enableMessage', { name: toggleTarget?.name ?? '' })"
      :confirm-text="t('common.confirm')"
      :cancel-text="t('common.cancel')"
      @confirm="executeToggle"
      @cancel="showToggleConfirm = false"
    />

    <!-- Confirm Delete Dialog -->
    <ConfirmDialog
      :show="showDeleteConfirm"
      :title="t('enterprise.keys.deleteKey')"
      :message="t('enterprise.keys.deleteMessage', { name: deleteTarget?.name ?? '' })"
      :confirm-text="t('common.delete')"
      :cancel-text="t('common.cancel')"
      :dangerous="true"
      @confirm="executeDelete"
      @cancel="showDeleteConfirm = false"
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
import enterpriseAdminAPI from '@/api/enterprise'
import type { EnterpriseKey, CreateEnterpriseKeyResponse } from '@/types/enterprise'
import type { Column } from '@/components/common/types'

const { t } = useI18n()
const appStore = useAppStore()

// ---- State ----
const loading = ref(false)
const submitting = ref(false)
const keys = ref<EnterpriseKey[]>([])
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

// ---- Columns ----
const columns = computed<Column[]>(() => [
  { key: 'name', label: t('enterprise.keys.columns.name'), sortable: false },
  { key: 'key_prefix', label: t('enterprise.keys.columns.keyPrefix'), sortable: false },
  { key: 'status', label: t('enterprise.keys.columns.status'), sortable: false },
  { key: 'assigned_to', label: t('enterprise.keys.columns.assignedTo'), sortable: false },
  { key: 'groups', label: t('enterprise.keys.columns.groups'), sortable: false },
  { key: 'bound_tool', label: t('enterprise.keys.columns.boundTool'), sortable: false },
  { key: 'quota', label: t('enterprise.keys.columns.quota'), sortable: false },
  { key: 'created_at', label: t('enterprise.keys.columns.createdAt'), sortable: true },
  { key: 'actions', label: t('enterprise.keys.columns.actions'), sortable: false },
])

// ---- Create Modal ----
const showCreateModal = ref(false)
const createdKey = ref<CreateEnterpriseKeyResponse | null>(null)
const createForm = reactive({
  name: '',
  groupIdsInput: '',
  assigned_to: undefined as number | undefined,
  usage_purpose: '',
  bound_tool: '' as string,
  quota: undefined as number | undefined,
})

// ---- Toggle Confirm ----
const showToggleConfirm = ref(false)
const toggleTarget = ref<EnterpriseKey | null>(null)

// ---- Delete Confirm ----
const showDeleteConfirm = ref(false)
const deleteTarget = ref<EnterpriseKey | null>(null)

// ---- Methods ----
async function loadData() {
  loading.value = true
  try {
    const res = await enterpriseAdminAPI.listKeys(pagination.page, pagination.page_size)
    keys.value = res.data ?? []
    pagination.total = res.total ?? 0
  } catch (err: any) {
    appStore.showToast(err?.message ?? t('common.loadError'), 'error')
  } finally {
    loading.value = false
  }
}

function onPageChange(page: number) {
  pagination.page = page
  loadData()
}

// ---- Create ----
function openCreateModal() {
  createdKey.value = null
  createForm.name = ''
  createForm.groupIdsInput = ''
  createForm.assigned_to = undefined
  createForm.usage_purpose = ''
  createForm.bound_tool = ''
  createForm.quota = undefined
  showCreateModal.value = true
}

function closeCreateModal() {
  showCreateModal.value = false
  if (createdKey.value) {
    createdKey.value = null
    loadData()
  }
}

async function submitCreate() {
  submitting.value = true
  try {
    if (!createForm.name.trim()) {
      appStore.showToast(t('enterprise.keys.nameRequired'), 'error')
      submitting.value = false
      return
    }
    const groupIds = createForm.groupIdsInput
      ? createForm.groupIdsInput.split(',').map(s => parseInt(s.trim(), 10)).filter(n => !isNaN(n))
      : []

    const res = await enterpriseAdminAPI.createKey({
      name: createForm.name.trim(),
      group_ids: groupIds,
      assigned_to: createForm.assigned_to,
      usage_purpose: createForm.usage_purpose.trim() || undefined,
      bound_tool: (createForm.bound_tool as any) || undefined,
      quota: createForm.quota,
    })
    createdKey.value = res
    appStore.showToast(t('enterprise.keys.created'), 'success')
  } catch (err: any) {
    appStore.showToast(err?.message ?? t('common.saveError'), 'error')
  } finally {
    submitting.value = false
  }
}

// ---- Toggle ----
function confirmToggle(key: EnterpriseKey) {
  toggleTarget.value = key
  showToggleConfirm.value = true
}

async function executeToggle() {
  if (!toggleTarget.value) return
  try {
    await enterpriseAdminAPI.toggleKey(toggleTarget.value.id)
    appStore.showToast(t('enterprise.keys.toggled'), 'success')
    showToggleConfirm.value = false
    await loadData()
  } catch (err: any) {
    appStore.showToast(err?.message ?? t('common.error'), 'error')
  }
}

// ---- Delete ----
function confirmDelete(key: EnterpriseKey) {
  deleteTarget.value = key
  showDeleteConfirm.value = true
}

async function executeDelete() {
  if (!deleteTarget.value) return
  try {
    await enterpriseAdminAPI.deleteKey(deleteTarget.value.id)
    appStore.showToast(t('enterprise.keys.deleted'), 'success')
    showDeleteConfirm.value = false
    await loadData()
  } catch (err: any) {
    appStore.showToast(err?.message ?? t('common.error'), 'error')
  }
}

// ---- Init ----
onMounted(() => {
  loadData()
})
</script>


