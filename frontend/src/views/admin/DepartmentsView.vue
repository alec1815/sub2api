<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- Page Header -->
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('admin.departments.title') }}</h1>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('admin.departments.description') }}</p>
      </div>

      <!-- Enterprise Selector + Actions -->
      <div class="card">
        <div class="flex flex-wrap items-center gap-3">
          <!-- Enterprise selector -->
          <div class="w-full sm:w-64">
            <label class="mb-1 block text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.departments.selectEnterprise') }}</label>
            <select v-model="selectedEnterpriseId" class="input w-full" @change="handleEnterpriseChange">
              <option :value="0" disabled>{{ t('admin.departments.selectEnterprisePlaceholder') }}</option>
              <option v-for="e in enterpriseOptions" :key="e.id" :value="e.id">{{ e.name }}</option>
            </select>
          </div>
          <!-- Add Root Department -->
          <button v-if="selectedEnterpriseId > 0" class="btn-primary inline-flex items-center gap-2 rounded-lg px-4 py-2.5 text-sm font-semibold mt-5" @click="openCreateModal(null)">
            <Icon name="plus" size="sm" />
            {{ t('admin.departments.addDepartment') }}
          </button>
          <!-- Refresh -->
          <button v-if="selectedEnterpriseId > 0" class="btn-secondary inline-flex items-center gap-2 rounded-lg px-4 py-2.5 text-sm mt-5" @click="loadTree">
            <Icon name="refresh" size="sm" />
            {{ t('admin.departments.refresh') }}
          </button>
        </div>
      </div>

      <!-- Department Tree -->
      <div v-if="selectedEnterpriseId > 0" class="card">
        <!-- Loading -->
        <div v-if="treeLoading" class="flex items-center justify-center py-20">
          <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent" />
        </div>
        <!-- Empty -->
        <div v-else-if="filteredTree.length === 0" class="flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
          <Icon name="inbox" size="xl" class="mb-4 h-16 w-16" />
          <p class="text-sm">{{ t('admin.departments.noDepartments') }}</p>
          <button class="mt-4 text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400" @click="openCreateModal(null)">{{ t('admin.departments.createFirst') }}</button>
        </div>
        <!-- Tree -->
        <div v-else class="p-4">
          <template v-for="node in filteredTree" :key="node.id">
            <DepartmentTreeNode :node="node" :depth="0" :expanded-ids="expandedNodes" @edit="openEditModal" @add-child="openCreateModal" @delete="confirmDelete" @toggle="toggleNode" />
          </template>
        </div>
      </div>

      <!-- No enterprise selected -->
      <div v-else class="card flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
        <Icon name="inbox" size="xl" class="mb-4 h-16 w-16" />
        <p class="text-sm">{{ t('admin.departments.selectEnterpriseHint') }}</p>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showFormModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4" @click.self="closeFormModal">
          <div class="w-full max-w-md rounded-2xl bg-white shadow-2xl dark:bg-dark-800">
            <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dark-700">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ isEditing ? t('admin.departments.editDepartment') : t('admin.departments.addDepartment') }}</h2>
              <button class="rounded-lg p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300" @click="closeFormModal">
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <form class="space-y-4 px-6 py-5" @submit.prevent="submitForm">
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.departments.form.name') }} <span class="text-red-500">*</span></label>
                <input v-model="form.name" type="text" class="input w-full" :placeholder="t('admin.departments.form.namePlaceholder')" :disabled="submitting" />
                <p v-if="formErrors.name" class="mt-1 text-xs text-red-500">{{ formErrors.name }}</p>
              </div>
              <div v-if="!isEditing">
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.departments.form.parent') }}</label>
                <input :value="parentName" type="text" class="input w-full bg-gray-50" disabled />
              </div>
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.departments.form.leader') }}</label>
                <input v-model="form.leader" type="text" class="input w-full" :placeholder="t('admin.departments.form.leaderPlaceholder')" :disabled="submitting" />
              </div>
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.departments.form.phone') }}</label>
                <input v-model="form.phone" type="text" class="input w-full" :placeholder="t('admin.departments.form.phonePlaceholder')" :disabled="submitting" />
              </div>
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.departments.form.email') }}</label>
                <input v-model="form.email" type="email" class="input w-full" :placeholder="t('admin.departments.form.emailPlaceholder')" :disabled="submitting" />
              </div>
              <div class="flex items-center justify-end gap-3 pt-2">
                <button type="button" class="btn-secondary rounded-lg px-4 py-2.5 text-sm font-medium" :disabled="submitting" @click="closeFormModal">{{ t('common.cancel') }}</button>
                <button type="submit" class="btn-primary rounded-lg px-5 py-2.5 text-sm font-semibold" :disabled="submitting">{{ t('common.save') }}</button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Delete Confirm -->
    <ConfirmDialog :show="showDeleteDialog" :title="t('admin.departments.deleteTitle')" :message="t('admin.departments.deleteConfirm', { name: deletingNode?.name })" :danger="true" @confirm="handleDelete" @cancel="showDeleteDialog = false" />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import Icon from '@/components/icons/Icon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { adminAPI } from '@/api/admin'
import type { Department } from '@/types/enterprise'
import type { Enterprise } from '@/types/enterprise'

const { t } = useI18n()
const appStore = useAppStore()

// ---- Enterprise selector ----
const selectedEnterpriseId = ref(0)
const enterpriseOptions = ref<Enterprise[]>([])

async function loadEnterprises() {
  try {
    const result = await adminAPI.enterprises.list(1, 100, {})
    enterpriseOptions.value = result.items
  } catch { /* silent */ }
}

function handleEnterpriseChange() {
  if (selectedEnterpriseId.value > 0) loadTree()
}

// ---- Tree ----
const treeLoading = ref(false)
const tree = ref<Department[]>([])
const expandedNodes = reactive<Set<number>>(new Set())
const searchKeyword = ref('')

const filteredTree = computed(() => {
  if (!searchKeyword.value.trim()) return tree.value
  return filterTree(tree.value, searchKeyword.value.toLowerCase())
})

function filterTree(nodes: Department[], keyword: string): Department[] {
  return nodes.reduce<Department[]>((acc, node) => {
    const nameMatch = node.name.toLowerCase().includes(keyword)
    const filteredChildren = filterTree(node.children, keyword)
    if (nameMatch || filteredChildren.length > 0) {
      acc.push({ ...node, children: filteredChildren })
    }
    return acc
  }, [])
}

async function loadTree() {
  if (selectedEnterpriseId.value <= 0) return
  treeLoading.value = true
  try {
    tree.value = await adminAPI.departments.getTree(selectedEnterpriseId.value)
    // Auto-expand first level
    tree.value.forEach(n => expandedNodes.add(n.id))
  } catch {
    appStore.showError(t('admin.departments.failedToLoad'))
  } finally {
    treeLoading.value = false
  }
}

function toggleNode(id: number) {
  if (expandedNodes.has(id)) expandedNodes.delete(id)
  else expandedNodes.add(id)
}

// ---- Form Modal ----
const showFormModal = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)
const parentId = ref<number | null>(null)
const parentName = ref('')
const submitting = ref(false)
const form = ref({ name: '', leader: '', phone: '', email: '' })
const formErrors = ref<Record<string, string>>({})

function getNodeName(id: number): string {
  const find = (nodes: Department[]): string | null => {
    for (const n of nodes) {
      if (n.id === id) return n.name
      const r = find(n.children)
      if (r) return r
    }
    return null
  }
  return find(tree.value) || t('admin.departments.form.rootDepartment')
}

function openCreateModal(parent: Department | null) {
  isEditing.value = false
  editingId.value = null
  parentId.value = parent?.id ?? null
  parentName.value = parent ? parent.name : t('admin.departments.form.rootDepartment')
  form.value = { name: '', leader: '', phone: '', email: '' }
  formErrors.value = {}
  showFormModal.value = true
}

function openEditModal(node: Department) {
  isEditing.value = true
  editingId.value = node.id
  parentId.value = node.parent_id
  parentName.value = getNodeName(node.parent_id)
  form.value = { name: node.name, leader: node.leader || '', phone: node.phone || '', email: node.email || '' }
  formErrors.value = {}
  showFormModal.value = true
}

function closeFormModal() { showFormModal.value = false }

async function submitForm() {
  formErrors.value = {}
  if (!form.value.name.trim()) { formErrors.value.name = t('admin.departments.form.nameRequired'); return }
  submitting.value = true
  try {
    if (isEditing.value) {
      await adminAPI.departments.update(editingId.value!, { name: form.value.name.trim(), leader: form.value.leader.trim() || undefined, phone: form.value.phone.trim() || undefined, email: form.value.email.trim() || undefined })
      appStore.showSuccess(t('admin.departments.updateSuccess'))
    } else {
      await adminAPI.departments.create({ enterprise_id: selectedEnterpriseId.value, parent_id: parentId.value ?? undefined, name: form.value.name.trim(), leader: form.value.leader.trim() || undefined, phone: form.value.phone.trim() || undefined, email: form.value.email.trim() || undefined })
      appStore.showSuccess(t('admin.departments.createSuccess'))
    }
    closeFormModal()
    await loadTree()
  } catch {
    appStore.showError(t('admin.departments.saveFailed'))
  } finally {
    submitting.value = false
  }
}

// ---- Delete ----
const showDeleteDialog = ref(false)
const deletingNode = ref<Department | null>(null)

function confirmDelete(node: Department) { deletingNode.value = node; showDeleteDialog.value = true }

async function handleDelete() {
  if (!deletingNode.value) return
  submitting.value = true
  try {
    await adminAPI.departments.delete(deletingNode.value.id)
    appStore.showSuccess(t('admin.departments.deleteSuccess'))
    showDeleteDialog.value = false; deletingNode.value = null
    await loadTree()
  } catch {
    appStore.showError(t('admin.departments.deleteFailed'))
  } finally {
    submitting.value = false
  }
}

onMounted(() => { loadEnterprises() })
</script>
