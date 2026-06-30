<template>
  <AppLayout>
    <TablePageLayout>
      <!-- Filters Row: Search, Status Filter, and Actions -->
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <!-- Left: Search + Filters -->
          <div class="flex-1 sm:max-w-64">
            <input
              v-model="searchQuery"
              type="text"
              :placeholder="t('admin.enterprises.searchPlaceholder')"
              class="input"
              @input="onSearchInput"
            />
          </div>
          <Select
            v-model="statusFilter"
            :options="filterStatusOptions"
            class="w-36"
            @change="loadEnterprises"
          />

          <!-- Right: Action buttons -->
          <div class="flex flex-1 flex-wrap items-center justify-end gap-2">
            <!-- Column Settings Dropdown -->
            <div class="relative" ref="columnDropdownRef">
              <button
                @click="showColumnDropdown = !showColumnDropdown"
                class="btn btn-secondary px-2 md:px-3"
                :title="t('admin.enterprises.columnSettings')"
              >
                <svg class="h-4 w-4 md:mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 4.5v15m6-15v15m-10.875 0h15.75c.621 0 1.125-.504 1.125-1.125V5.625c0-.621-.504-1.125-1.125-1.125H4.125C3.504 4.5 3 5.004 3 5.625v12.75c0 .621.504 1.125 1.125 1.125z" />
                </svg>
                <span class="hidden md:inline">{{ t('admin.enterprises.columnSettings') }}</span>
              </button>
              <!-- Dropdown menu -->
              <div
                v-if="showColumnDropdown"
                class="absolute right-0 top-full z-50 mt-1 max-h-80 w-48 overflow-y-auto rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-dark-600 dark:bg-dark-800"
              >
                <button
                  v-for="col in toggleableColumns"
                  :key="col.key"
                  :disabled="isForcedVisibleColumn(col.key)"
                  @click="toggleColumn(col.key)"
                  :class="[
                    'flex w-full items-center justify-between px-4 py-2 text-left text-sm',
                    isForcedVisibleColumn(col.key)
                      ? 'cursor-not-allowed text-gray-400 dark:text-gray-500'
                      : 'text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700'
                  ]"
                  :title="isForcedVisibleColumn(col.key) ? t('admin.enterprises.columnAlwaysVisible') : ''"
                >
                  <span>{{ col.label }}</span>
                  <Icon
                    v-if="isColumnVisible(col.key)"
                    name="check"
                    size="sm"
                    :class="isForcedVisibleColumn(col.key) ? 'text-gray-400 dark:text-gray-500' : 'text-primary-500'"
                    :stroke-width="2"
                  />
                </button>
              </div>
            </div>

            <button
              @click="loadEnterprises"
              :disabled="loading"
              class="btn btn-secondary"
              :title="t('common.refresh')"
            >
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
            <button class="btn btn-primary" @click="openCreateModal">
              <Icon name="plus" size="md" class="mr-1" />
              {{ t('admin.enterprises.createEnterprise') }}
            </button>
          </div>
        </div>
      </template>

      <!-- Table -->
      <template #table>
        <DataTable :data="enterprises" :columns="columns" :loading="loading" :actions-count="4">
          <!-- Name -->
          <template #cell-name="{ value, row }">
            <div class="flex items-center gap-2">
              <div class="flex h-8 w-8 items-center justify-center rounded-full bg-primary-100 dark:bg-primary-900/30">
                <span class="text-sm font-medium text-primary-700 dark:text-primary-300">{{ value.charAt(0).toUpperCase() }}</span>
              </div>
              <div class="flex flex-col">
                <span class="font-medium text-gray-900 dark:text-white">{{ value }}</span>
                <span v-if="row.short_name" class="text-xs text-gray-400 dark:text-gray-500">{{ row.short_name }}</span>
              </div>
            </div>
          </template>
          <!-- Admin Email -->
          <template #cell-admin_email="{ value }">
            <span class="text-sm text-gray-700 dark:text-gray-300">{{ value || '-' }}</span>
          </template>
          <!-- Contact Name -->
          <template #cell-contact_name="{ value }">
            <span class="text-sm text-gray-700 dark:text-gray-300">{{ value || '-' }}</span>
          </template>
          <!-- Contact Phone -->
          <template #cell-contact_phone="{ value }">
            <span class="text-sm text-gray-700 dark:text-gray-300">{{ value || '-' }}</span>
          </template>
          <!-- Member Count -->
          <template #cell-member_count="{ value }">
            <span class="text-sm text-gray-700 dark:text-gray-300">{{ value ?? 0 }}</span>
          </template>
          <!-- Balance -->
          <template #cell-balance="{ value, row }">
            <div class="flex items-center gap-2">
              <span class="cursor-pointer font-medium text-gray-900 hover:text-primary-600 dark:text-white dark:hover:text-primary-400" @click="handleBalanceHistory(row)" :title="t('admin.enterprises.balanceHistoryTitle')">${{ parseFloat(value ?? '0').toFixed(2) }}</span>
              <button @click.stop="handleDeposit(row)" class="rounded px-1.5 py-0.5 text-xs text-emerald-600 hover:bg-emerald-50 dark:text-emerald-400 dark:hover:bg-emerald-900/20" :title="t('admin.enterprises.deposit')">{{ t('admin.enterprises.deposit') }}</button>
            </div>
          </template>
          <!-- Concurrency -->
          <template #cell-concurrency="{ value }">
            <span class="text-sm text-gray-700 dark:text-gray-300">{{ value ?? 0 }}</span>
          </template>
          <!-- Status -->
          <template #cell-status="{ value }">
            <div class="flex items-center gap-1.5">
              <span :class="['inline-block h-2 w-2 rounded-full', value === 'active' ? 'bg-green-500' : 'bg-red-500']" />
              <span class="text-sm text-gray-700 dark:text-gray-300">{{ value === 'active' ? t('admin.enterprises.status.active') : t('admin.enterprises.status.disabled') }}</span>
            </div>
          </template>
          <!-- Created At -->
          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateTime(value) }}</span>
          </template>
          <!-- Actions -->
          <template #cell-actions="{ row }">
            <div class="flex items-center gap-1">
              <button @click="openEditModal(row)" class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-primary-600 dark:hover:bg-dark-700 dark:hover:text-primary-400">
                <Icon name="edit" size="sm" />
                <span class="text-xs">{{ t('common.edit') }}</span>
              </button>
              <button @click="handleViewMembers(row)" class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-amber-600 dark:hover:bg-dark-700 dark:hover:text-amber-400">
                <Icon name="users" size="sm" />
                <span class="text-xs">{{ t('admin.enterprises.viewMembers') }}</span>
              </button>
              <button @click="row.status === 'active' ? confirmDisable(row) : confirmActivate(row)" :class="['flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors', row.status === 'active' ? 'hover:bg-orange-50 hover:text-orange-600 dark:hover:bg-orange-900/20 dark:hover:text-orange-400' : 'hover:bg-green-50 hover:text-green-600 dark:hover:bg-green-900/20 dark:hover:text-green-400']">
                <Icon v-if="row.status === 'active'" name="ban" size="sm" />
                <Icon v-else name="checkCircle" size="sm" />
                <span class="text-xs">{{ row.status === 'active' ? t('admin.enterprises.disable') : t('admin.enterprises.activate') }}</span>
              </button>
              <button @click="openActionMenu(row, $event)" class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:hover:bg-dark-700 dark:hover:text-white" :class="{ 'bg-gray-100 text-gray-900 dark:bg-dark-700 dark:text-white': activeMenuId === row.id }">
                <Icon name="more" size="sm" />
                <span class="text-xs">{{ t('common.more') }}</span>
              </button>
            </div>
          </template>
          <!-- Empty State -->
          <template #empty>
            <EmptyState :title="t('admin.enterprises.noEnterprisesYet')" :description="t('admin.enterprises.noEnterprisesHint')" :action-text="t('admin.enterprises.createEnterprise')" @action="openCreateModal" />
          </template>
        </DataTable>
      </template>

      <!-- Pagination -->
      <template #pagination>
        <Pagination v-if="pagination.total > 0" :page="pagination.page" :page-size="pagination.page_size" :total="pagination.total" @update:page="handlePageChange" @update:page-size="handlePageSizeChange" />
      </template>
    </TablePageLayout>

    <!-- Action Menu -->
    <Teleport to="body">
      <div v-if="activeMenuId !== null && menuPosition" class="fixed z-50 w-48 rounded-lg border border-gray-200 bg-white shadow-lg dark:border-dark-600 dark:bg-dark-800" :style="{ top: menuPosition.top + 'px', left: menuPosition.left + 'px' }" @click.stop>
        <div class="py-1">
          <template v-for="enterprise in enterprises" :key="'menu-' + enterprise.id">
            <template v-if="enterprise.id === activeMenuId">
              <button @click="handleViewMembers(enterprise); closeActionMenu()" class="flex w-full items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700">
                <Icon name="users" size="sm" class="text-gray-400" :stroke-width="2" />
                {{ t('admin.enterprises.viewMembers') }}
              </button>
              <div class="my-1 border-t border-gray-100 dark:border-dark-700" />
              <button @click="handleDeposit(enterprise); closeActionMenu()" class="flex w-full items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700">
                <Icon name="plus" size="sm" class="text-emerald-500" :stroke-width="2" />
                {{ t('admin.enterprises.deposit') }}
              </button>
              <button @click="handleWithdraw(enterprise); closeActionMenu()" class="flex w-full items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700">
                <Icon name="dollar" size="sm" class="text-red-500" :stroke-width="2" />
                {{ t('admin.enterprises.withdraw') }}
              </button>
              <button @click="handleBalanceHistory(enterprise); closeActionMenu()" class="flex w-full items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700">
                <Icon name="chart" size="sm" class="text-blue-500" :stroke-width="2" />
                {{ t('admin.enterprises.balanceHistoryTitle') }}
              </button>
              <div class="my-1 border-t border-gray-100 dark:border-dark-700" />
              <button @click="handleDelete(enterprise); closeActionMenu()" class="flex w-full items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20">
                <Icon name="trash" size="sm" :stroke-width="2" />
                {{ t('common.delete') }}
              </button>
            </template>
          </template>
        </div>
      </div>
    </Teleport>

    <!-- Create Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-black/60 backdrop-blur-sm p-4">
          <div class="my-8 w-full max-w-lg rounded-2xl bg-white shadow-2xl dark:bg-dark-800">
            <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dark-700">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.enterprises.createEnterprise') }}</h2>
              <button class="rounded-lg p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300" @click="closeCreateModal">
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <form class="space-y-5 px-6 py-5" @submit.prevent="submitCreateForm">
              <!-- 企业基本信息 -->
              <div class="space-y-4">
                <div class="grid gap-4 sm:grid-cols-2">
                  <div class="sm:col-span-2">
                    <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.fullName') }} <span class="text-red-500">*</span></label>
                    <input v-model="createForm.name" type="text" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.fullNamePlaceholder')" :disabled="submitting" />
                    <p v-if="createErrors.name" class="mt-1 text-xs text-red-500">{{ createErrors.name }}</p>
                  </div>
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.shortName') }}</label>
                    <input v-model="createForm.short_name" type="text" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.shortNamePlaceholder')" :disabled="submitting" />
                  </div>
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.creditCode') }}</label>
                    <input v-model="createForm.credit_code" type="text" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.creditCodePlaceholder')" :disabled="submitting" />
                  </div>
                </div>
              </div>

              <hr class="border-gray-100 dark:border-dark-700" />

              <!-- 管理员账号 -->
              <div class="space-y-4">
                <h3 class="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('admin.enterprises.form.adminSection') }}</h3>
                <div class="grid gap-4 sm:grid-cols-2">
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.adminEmail') }} <span class="text-red-500">*</span></label>
                    <input v-model="createForm.admin_email" type="email" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.adminEmailPlaceholder')" :disabled="submitting" />
                    <p v-if="createErrors.adminEmail" class="mt-1 text-xs text-red-500">{{ createErrors.adminEmail }}</p>
                  </div>
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.adminName') }}</label>
                    <input v-model="createForm.admin_name" type="text" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.adminNamePlaceholder')" :disabled="submitting" />
                  </div>
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.password') }} <span class="text-red-500">*</span></label>
                    <div class="relative">
                      <input v-model="createForm.password" :type="showPassword ? 'text' : 'password'" class="w-full rounded-lg border border-gray-200 bg-white py-2.5 pl-3 pr-10 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.passwordPlaceholder')" :disabled="submitting" />
                      <button type="button" class="absolute right-2.5 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600" @click="showPassword = !showPassword">
                        <svg v-if="!showPassword" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178zM15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
                        <svg v-else class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88"/></svg>
                      </button>
                    </div>
                    <p v-if="createErrors.password" class="mt-1 text-xs text-red-500">{{ createErrors.password }}</p>
                  </div>
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.confirmPassword') }} <span class="text-red-500">*</span></label>
                    <div class="relative">
                      <input v-model="createForm.confirmPassword" :type="showConfirmPassword ? 'text' : 'password'" class="w-full rounded-lg border border-gray-200 bg-white py-2.5 pl-3 pr-10 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.confirmPasswordPlaceholder')" :disabled="submitting" />
                      <button type="button" class="absolute right-2.5 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600" @click="showConfirmPassword = !showConfirmPassword">
                        <svg v-if="!showConfirmPassword" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178zM15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
                        <svg v-else class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88"/></svg>
                      </button>
                    </div>
                    <p v-if="createErrors.confirmPassword" class="mt-1 text-xs text-red-500">{{ createErrors.confirmPassword }}</p>
                  </div>
                </div>
              </div>

              <hr class="border-gray-100 dark:border-dark-700" />

              <!-- 联系人信息 -->
              <div class="space-y-4">
                <h3 class="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('admin.enterprises.form.contactSection') }}</h3>
                <div class="grid gap-4 sm:grid-cols-2">
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.contactName') }} <span class="text-red-500">*</span></label>
                    <input v-model="createForm.contact_name" type="text" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.contactNamePlaceholder')" :disabled="submitting" />
                    <p v-if="createErrors.contact_name" class="mt-1 text-xs text-red-500">{{ createErrors.contact_name }}</p>
                  </div>
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.contactPhone') }} <span class="text-red-500">*</span></label>
                    <input v-model="createForm.contact_phone" type="text" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.contactPhonePlaceholder')" :disabled="submitting" />
                    <p v-if="createErrors.contact_phone" class="mt-1 text-xs text-red-500">{{ createErrors.contact_phone }}</p>
                  </div>
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.contactEmail') }} <span class="text-red-500">*</span></label>
                    <input v-model="createForm.contact_email" type="email" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.contactEmailPlaceholder')" :disabled="submitting" />
                    <p v-if="createErrors.contact_email" class="mt-1 text-xs text-red-500">{{ createErrors.contact_email }}</p>
                  </div>
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.address') }}</label>
                    <input v-model="createForm.address" type="text" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.addressPlaceholder')" :disabled="submitting" />
                  </div>
                </div>
              </div>

              <hr class="border-gray-100 dark:border-dark-700" />

              <!-- 企业分类 -->
              <div class="space-y-4">
                <h3 class="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('admin.enterprises.form.categorySection') }}</h3>
                <div class="grid gap-4 sm:grid-cols-2">
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.scale') }}</label>
                    <Select v-model="createForm.scale" :options="scaleOptions" :disabled="submitting" />
                  </div>
                  <div>
                    <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.industry') }}</label>
                    <Select v-model="createForm.industry" :options="industryOptions" :disabled="submitting" />
                  </div>
                </div>
              </div>

              <hr class="border-gray-100 dark:border-dark-700" />

              <!-- 备注 -->
              <div class="space-y-2">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.notes') }}</label>
                <textarea v-model="createForm.notes" rows="2" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500 resize-none" :placeholder="t('admin.enterprises.form.notesPlaceholder')" :disabled="submitting" />
              </div>

              <div class="rounded-lg border border-blue-100 bg-blue-50 px-4 py-3 dark:border-blue-900/40 dark:bg-blue-900/20">
                <p class="text-xs text-blue-700 dark:text-blue-300">{{ t('admin.enterprises.createHint') }}</p>
              </div>

              <div class="flex items-center justify-end gap-3 pt-2">
                <button type="button" class="btn-secondary rounded-lg px-4 py-2.5 text-sm font-medium" :disabled="submitting" @click="closeCreateModal">{{ t('common.cancel') }}</button>
                <button type="submit" class="btn-primary rounded-lg px-5 py-2.5 text-sm font-semibold" :disabled="submitting">
                  <span v-if="submitting" class="mr-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                  {{ t('admin.enterprises.confirmCreate') }}
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
        <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-black/60 backdrop-blur-sm p-4">
          <div class="my-8 w-full max-w-lg rounded-2xl bg-white shadow-2xl dark:bg-dark-800">
            <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dark-700">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.enterprises.editEnterprise') }}</h2>
              <button class="rounded-lg p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300" @click="closeEditModal">
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <form class="space-y-5 px-6 py-5" @submit.prevent="submitEditForm">
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.fullName') }} <span class="text-red-500">*</span></label>
                <input v-model="editForm.name" type="text" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.namePlaceholder')" :disabled="submitting" />
                <p v-if="editErrors.name" class="mt-1 text-xs text-red-500">{{ editErrors.name }}</p>
              </div>
              <div class="grid gap-4 sm:grid-cols-2">
                <div>
                  <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.shortName') }}</label>
                  <input v-model="editForm.short_name" type="text" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.shortNamePlaceholder')" :disabled="submitting" />
                </div>
                <div>
                  <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.creditCode') }}</label>
                  <input v-model="editForm.credit_code" type="text" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.creditCodePlaceholder')" :disabled="submitting" />
                </div>
              </div>
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.address') }}</label>
                <input v-model="editForm.address" type="text" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.addressPlaceholder')" :disabled="submitting" />
              </div>
              <div class="grid gap-4 sm:grid-cols-2">
                <div>
                  <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.scale') }}</label>
                  <Select v-model="editForm.scale" :options="scaleOptions" :disabled="submitting" />
                </div>
                <div>
                  <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.industry') }}</label>
                  <Select v-model="editForm.industry" :options="industryOptions" :disabled="submitting" />
                </div>
              </div>
              <hr class="border-gray-100 dark:border-dark-700" />
              <div class="grid gap-4 sm:grid-cols-2">
                <div>
                  <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.contactName') }}</label>
                  <input v-model="editForm.contact_name" type="text" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.contactNamePlaceholder')" :disabled="submitting" />
                </div>
                <div>
                  <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.contactPhone') }} <span class="text-red-500">*</span></label>
                  <input v-model="editForm.contact_phone" type="text" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.contactPhonePlaceholder')" :disabled="submitting" />
                  <p v-if="editErrors.contact_phone" class="mt-1 text-xs text-red-500">{{ editErrors.contact_phone }}</p>
                </div>
                <div class="sm:col-span-2">
                  <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.contactEmail') }} <span class="text-red-500">*</span></label>
                  <input v-model="editForm.contact_email" type="email" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500" :placeholder="t('admin.enterprises.form.emailPlaceholder')" :disabled="submitting" />
                  <p v-if="editErrors.contact_email" class="mt-1 text-xs text-red-500">{{ editErrors.contact_email }}</p>
                </div>
              </div>
              <div>
                <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.enterprises.form.notes') }}</label>
                <textarea v-model="editForm.notes" rows="3" class="w-full rounded-lg border border-gray-200 bg-white px-3 py-2.5 text-sm text-gray-900 placeholder-gray-400 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500 resize-none" :placeholder="t('admin.enterprises.form.notesPlaceholder')" :disabled="submitting" />
              </div>
              <div class="flex items-center justify-end gap-3 pt-2">
                <button type="button" class="btn-secondary rounded-lg px-4 py-2.5 text-sm font-medium" :disabled="submitting" @click="closeEditModal">{{ t('common.cancel') }}</button>
                <button type="submit" class="btn-primary rounded-lg px-5 py-2.5 text-sm font-semibold" :disabled="submitting">
                  <span v-if="submitting" class="mr-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                  {{ t('common.save') }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Disable/Activate Confirm -->
    <ConfirmDialog :show="showDisableModal" :title="t('admin.enterprises.disableConfirmTitle')" :message="t('admin.enterprises.disableConfirmMessage', { name: targetEnterprise?.name })" :danger="true" @confirm="handleDisable" @cancel="showDisableModal = false" />
    <ConfirmDialog :show="showActivateModal" :title="t('admin.enterprises.activateConfirmTitle')" :message="t('admin.enterprises.activateConfirmMessage', { name: targetEnterprise?.name })" @confirm="handleActivate" @cancel="showActivateModal = false" />

    <!-- Delete Confirm -->
    <ConfirmDialog :show="showDeleteDialog" :title="t('admin.enterprises.deleteEnterprise')" :message="t('admin.enterprises.deleteConfirm', { name: deletingEnterprise?.name })" :danger="true" @confirm="confirmDelete" @cancel="showDeleteDialog = false" />

    <!-- Balance Modal -->
    <EnterpriseBalanceModal :show="showBalanceModal" :enterprise="balanceEnterprise" :operation="balanceOperation" @close="showBalanceModal = false" @success="loadEnterprises" />

    <!-- Balance History Modal -->
    <EnterpriseBalanceHistoryModal :show="showBalanceHistoryModal" :enterprise="balanceHistoryEnterprise" @close="showBalanceHistoryModal = false" @deposit="handleDeposit(balanceHistoryEnterprise!)" @withdraw="handleWithdraw(balanceHistoryEnterprise!)" />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { formatDateTime } from '@/utils/format'
import Icon from '@/components/icons/Icon.vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Select from '@/components/common/Select.vue'
import EnterpriseBalanceModal from '@/components/admin/enterprise/EnterpriseBalanceModal.vue'
import EnterpriseBalanceHistoryModal from '@/components/admin/enterprise/EnterpriseBalanceHistoryModal.vue'
import { adminAPI } from '@/api/admin'
import type { Enterprise, EnterpriseStatus, EnterpriseScale, EnterpriseIndustry } from '@/types/enterprise'
import type { Column } from '@/components/common/types'

const { t } = useI18n()
const router = useRouter()
const appStore = useAppStore()

// ---- State ----
const loading = ref(false)
const submitting = ref(false)
const enterprises = ref<Enterprise[]>([])
const searchQuery = ref('')
const statusFilter = ref('')
let searchTimer: ReturnType<typeof setTimeout> | null = null
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const showColumnDropdown = ref(false)
const columnDropdownRef = ref<HTMLElement | null>(null)

// ---- Columns ----
const allColumns = computed<Column[]>(() => [
  { key: 'name', label: t('admin.enterprises.columns.name'), sortable: true },
  { key: 'admin_email', label: t('admin.enterprises.columns.adminEmail'), sortable: false },
  { key: 'contact_name', label: t('admin.enterprises.columns.contactName'), sortable: false },
  { key: 'contact_phone', label: t('admin.enterprises.columns.contactPhone'), sortable: false },
  { key: 'member_count', label: t('admin.enterprises.columns.memberCount'), sortable: true },
  { key: 'balance', label: t('admin.enterprises.columns.balance'), sortable: true },
  { key: 'concurrency', label: t('admin.enterprises.columns.concurrency'), sortable: true },
  { key: 'status', label: t('admin.enterprises.columns.status'), sortable: true },
  { key: 'created_at', label: t('admin.enterprises.columns.createdAt'), sortable: true },
  { key: 'actions', label: t('admin.enterprises.columns.actions'), sortable: false },
])

// Column Settings
const hiddenColumns = reactive<Set<string>>(new Set())
const DEFAULT_HIDDEN_COLUMNS = ['contact_phone']
const HIDDEN_COLUMNS_KEY = 'enterprise-hidden-columns'
const FORCED_VISIBLE_COLUMNS = new Set(['name', 'actions'])

const toggleableColumns = computed(() =>
  allColumns.value.filter(col => !FORCED_VISIBLE_COLUMNS.has(col.key))
)

function loadSavedColumns() {
  try {
    const saved = localStorage.getItem(HIDDEN_COLUMNS_KEY)
    if (saved) {
      const parsed = JSON.parse(saved)
      hiddenColumns.clear()
      ;(parsed as string[]).forEach(key => hiddenColumns.add(key))
      return
    }
  } catch { /* ignore */ }
  DEFAULT_HIDDEN_COLUMNS.forEach(key => hiddenColumns.add(key))
}

function saveColumnsToStorage() {
  localStorage.setItem(HIDDEN_COLUMNS_KEY, JSON.stringify([...hiddenColumns]))
}

function toggleColumn(key: string) {
  if (FORCED_VISIBLE_COLUMNS.has(key)) return
  if (hiddenColumns.has(key)) {
    hiddenColumns.delete(key)
  } else {
    hiddenColumns.add(key)
  }
  saveColumnsToStorage()
}

const isColumnVisible = (key: string) => !hiddenColumns.has(key)
const isForcedVisibleColumn = (key: string) => FORCED_VISIBLE_COLUMNS.has(key)

const columns = computed<Column[]>(() =>
  allColumns.value.filter(col =>
    FORCED_VISIBLE_COLUMNS.has(col.key) || !hiddenColumns.has(col.key)
  )
)

// Load saved column visibility
loadSavedColumns()

// ---- Options ----
const filterStatusOptions = computed(() => [
  { value: '', label: t('admin.enterprises.allStatus') },
  { value: 'active', label: t('admin.enterprises.status.active') },
  { value: 'disabled', label: t('admin.enterprises.status.disabled') },
])

const scaleOptions = computed(() => [
  { value: '', label: t('admin.enterprises.form.scalePlaceholder') },
  { value: 'micro', label: t('admin.enterprises.scales.micro') },
  { value: 'small', label: t('admin.enterprises.scales.small') },
  { value: 'medium', label: t('admin.enterprises.scales.medium') },
  { value: 'large', label: t('admin.enterprises.scales.large') },
])

const industryOptions = computed(() => [
  { value: '', label: t('admin.enterprises.form.industryPlaceholder') },
  { value: 'internet', label: t('admin.enterprises.industries.internet') },
  { value: 'finance', label: t('admin.enterprises.industries.finance') },
  { value: 'education', label: t('admin.enterprises.industries.education') },
  { value: 'healthcare', label: t('admin.enterprises.industries.healthcare') },
  { value: 'manufacturing', label: t('admin.enterprises.industries.manufacturing') },
  { value: 'other', label: t('admin.enterprises.industries.other') },
])

// ---- Action Menu ----
const activeMenuId = ref<number | null>(null)
const menuPosition = ref<{ top: number; left: number } | null>(null)

function openActionMenu(enterprise: Enterprise, event: MouseEvent) {
  activeMenuId.value = enterprise.id
  const btn = event.currentTarget as HTMLElement
  const rect = btn.getBoundingClientRect()
  const menuWidth = 192
  let left = rect.right - menuWidth
  if (left < 8) left = rect.left
  menuPosition.value = { top: rect.bottom + 4, left }
}

function closeActionMenu() {
  activeMenuId.value = null
  menuPosition.value = null
}

function handleClickOutside(e: MouseEvent) {
  if (columnDropdownRef.value && !columnDropdownRef.value.contains(e.target as Node)) {
    showColumnDropdown.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  loadSavedColumns()
  loadEnterprises()
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

// ---- Data loading ----
async function loadEnterprises() {
  loading.value = true
  try {
    const filters: { search?: string; status?: EnterpriseStatus } = {}
    if (searchQuery.value) filters.search = searchQuery.value
    if (statusFilter.value) filters.status = statusFilter.value as EnterpriseStatus
    const result = await adminAPI.enterprises.list(pagination.page, pagination.page_size, filters)
    enterprises.value = result.items
    pagination.total = result.total
  } catch (error: any) {
    console.error('Failed to load enterprises:', error)
    appStore.showError(error.response?.data?.detail || t('admin.enterprises.failedToLoad'))
  } finally {
    loading.value = false
  }
}

function onSearchInput() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => { pagination.page = 1; loadEnterprises() }, 350)
}

function handlePageChange(page: number) {
  pagination.page = page
  loadEnterprises()
}

function handlePageSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  loadEnterprises()
}

// ---- Create form ----
const showCreateModal = ref(false)
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const createForm = ref({
  name: '',
  short_name: '',
  credit_code: '',
  admin_email: '',
  admin_name: '',
  password: '',
  confirmPassword: '',
  contact_name: '',
  contact_phone: '',
  contact_email: '',
  address: '',
  scale: '' as EnterpriseScale | '',
  industry: '' as EnterpriseIndustry | '',
  notes: '',
})
const createErrors = ref<Record<string, string>>({})

function openCreateModal() {
  createForm.value = { name: '', short_name: '', credit_code: '', admin_email: '', admin_name: '', password: '', confirmPassword: '', contact_name: '', contact_phone: '', contact_email: '', address: '', scale: '', industry: '', notes: '' }
  createErrors.value = {}
  showPassword.value = false
  showConfirmPassword.value = false
  showCreateModal.value = true
}

function closeCreateModal() { showCreateModal.value = false }

async function submitCreateForm() {
  createErrors.value = {}
  if (!createForm.value.name.trim()) { createErrors.value.name = t('admin.enterprises.form.fullNameRequired'); return }
  if (!createForm.value.admin_email.trim()) { createErrors.value.adminEmail = t('admin.enterprises.form.adminEmailRequired'); return }
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(createForm.value.admin_email.trim())) { createErrors.value.adminEmail = t('admin.enterprises.form.adminEmailInvalid'); return }
  if (!createForm.value.password) { createErrors.value.password = t('admin.enterprises.form.passwordRequired'); return }
  if (createForm.value.password.length < 8 || createForm.value.password.length > 20) { createErrors.value.password = t('admin.enterprises.form.passwordInvalid'); return }
  if (createForm.value.password !== createForm.value.confirmPassword) { createErrors.value.confirmPassword = t('admin.enterprises.form.confirmPasswordMismatch'); return }
  if (!createForm.value.contact_name.trim()) { createErrors.value.contact_name = t('admin.enterprises.form.contactNameRequired'); return }
  if (!createForm.value.contact_phone.trim()) { createErrors.value.contact_phone = t('admin.enterprises.form.contactPhoneRequired'); return }
  if (!createForm.value.contact_email.trim()) { createErrors.value.contact_email = t('admin.enterprises.form.contactEmailRequired'); return }
  if (!emailRegex.test(createForm.value.contact_email.trim())) { createErrors.value.contact_email = t('admin.enterprises.form.contactEmailInvalid'); return }

  submitting.value = true
  try {
    await adminAPI.enterprises.create({
      name: createForm.value.name.trim(),
      short_name: createForm.value.short_name.trim() || undefined,
      credit_code: createForm.value.credit_code.trim() || undefined,
      address: createForm.value.address.trim() || undefined,
      scale: (createForm.value.scale as EnterpriseScale) || undefined,
      industry: (createForm.value.industry as EnterpriseIndustry) || undefined,
      contact_name: createForm.value.contact_name.trim(),
      contact_phone: createForm.value.contact_phone.trim(),
      contact_email: createForm.value.contact_email.trim(),
      admin_email: createForm.value.admin_email.trim(),
      admin_name: createForm.value.admin_name.trim() || undefined,
      admin_password: createForm.value.password,
      admin_password_confirm: createForm.value.confirmPassword,
      notes: createForm.value.notes.trim() || undefined,
    })
    appStore.showSuccess(t('admin.enterprises.createSuccess'))
    closeCreateModal()
    await loadEnterprises()
  } catch (error: any) {
    console.error('Failed to create enterprise:', error)
    appStore.showError(error.response?.data?.detail || t('admin.enterprises.errors.createFailed'))
  } finally {
    submitting.value = false
  }
}

// ---- Edit form ----
const showEditModal = ref(false)
const editingEnterprise = ref<Enterprise | null>(null)
const editForm = ref({ name: '', short_name: '', credit_code: '', address: '', scale: '' as EnterpriseScale | '', industry: '' as EnterpriseIndustry | '', contact_name: '', contact_phone: '', contact_email: '', notes: '' })
const editErrors = ref<Record<string, string>>({})

function openEditModal(enterprise: Enterprise) {
  editingEnterprise.value = enterprise
  editForm.value = { name: enterprise.name, short_name: enterprise.short_name || '', credit_code: enterprise.credit_code || '', address: enterprise.address || '', scale: enterprise.scale || '', industry: enterprise.industry || '', contact_name: enterprise.contact_name || '', contact_phone: enterprise.contact_phone || '', contact_email: enterprise.contact_email || '', notes: enterprise.notes || '' }
  editErrors.value = {}
  showEditModal.value = true
}

function closeEditModal() { showEditModal.value = false; editingEnterprise.value = null }

async function submitEditForm() {
  editErrors.value = {}
  if (!editForm.value.name.trim()) { editErrors.value.name = t('admin.enterprises.form.nameRequired'); return }
  if (!editForm.value.contact_phone.trim()) { editErrors.value.contact_phone = t('admin.enterprises.form.phoneRequired'); return }
  if (!editForm.value.contact_email.trim()) { editErrors.value.contact_email = t('admin.enterprises.form.emailRequired'); return }
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(editForm.value.contact_email.trim())) { editErrors.value.contact_email = t('admin.enterprises.form.adminEmailInvalid'); return }

  submitting.value = true
  try {
    await adminAPI.enterprises.update(editingEnterprise.value!.id, {
      name: editForm.value.name.trim(),
      short_name: editForm.value.short_name.trim() || undefined,
      credit_code: editForm.value.credit_code.trim() || undefined,
      address: editForm.value.address.trim() || undefined,
      scale: (editForm.value.scale as EnterpriseScale) || undefined,
      industry: (editForm.value.industry as EnterpriseIndustry) || undefined,
      contact_name: editForm.value.contact_name.trim() || undefined,
      contact_phone: editForm.value.contact_phone.trim(),
      contact_email: editForm.value.contact_email.trim(),
      notes: editForm.value.notes.trim() || undefined,
    })
    appStore.showSuccess(t('admin.enterprises.updateSuccess'))
    closeEditModal()
    await loadEnterprises()
  } catch (error: any) {
    console.error('Failed to update enterprise:', error)
    appStore.showError(error.response?.data?.detail || t('admin.enterprises.errors.createFailed'))
  } finally {
    submitting.value = false
  }
}

// ---- Status Toggle ----
const showDisableModal = ref(false)
const showActivateModal = ref(false)
const targetEnterprise = ref<Enterprise | null>(null)

function confirmDisable(enterprise: Enterprise) { targetEnterprise.value = enterprise; showDisableModal.value = true }

function confirmActivate(enterprise: Enterprise) {
  if (enterprise.status !== 'disabled') { appStore.showWarning(t('admin.enterprises.errors.alreadyActive')); return }
  targetEnterprise.value = enterprise; showActivateModal.value = true
}

async function handleDisable() {
  if (!targetEnterprise.value || targetEnterprise.value.status !== 'active') { appStore.showWarning(t('admin.enterprises.errors.alreadyDisabled')); showDisableModal.value = false; return }
  submitting.value = true
  try {
    await adminAPI.enterprises.toggleStatus(targetEnterprise.value.id)
    appStore.showSuccess(t('admin.enterprises.disableSuccess'))
    showDisableModal.value = false; targetEnterprise.value = null
    await loadEnterprises()
  } catch (error: any) { 
    console.error('Failed to disable enterprise:', error)
    appStore.showError(error.response?.data?.detail || t('admin.enterprises.errors.operationFailed'))
  }
  finally { submitting.value = false }
}

async function handleActivate() {
  if (!targetEnterprise.value) return
  submitting.value = true
  try {
    await adminAPI.enterprises.toggleStatus(targetEnterprise.value.id)
    appStore.showSuccess(t('admin.enterprises.activateSuccess'))
    showActivateModal.value = false; targetEnterprise.value = null
    await loadEnterprises()
  } catch (error: any) { 
    console.error('Failed to activate enterprise:', error)
    appStore.showError(error.response?.data?.detail || t('admin.enterprises.errors.operationFailed'))
  }
  finally { submitting.value = false }
}

// ---- Delete ----
const showDeleteDialog = ref(false)
const deletingEnterprise = ref<Enterprise | null>(null)

function handleDelete(enterprise: Enterprise) { deletingEnterprise.value = enterprise; showDeleteDialog.value = true }

async function confirmDelete() {
  if (!deletingEnterprise.value) return
  submitting.value = true
  try {
    await adminAPI.enterprises.delete(deletingEnterprise.value.id)
    appStore.showSuccess(t('common.success'))
    showDeleteDialog.value = false; deletingEnterprise.value = null
    await loadEnterprises()
  } catch (error: any) { 
    console.error('Failed to delete enterprise:', error)
    appStore.showError(error.response?.data?.detail || t('admin.enterprises.errors.operationFailed'))
  }
  finally { submitting.value = false }
}

// ---- View Members ----
function handleViewMembers(enterprise: Enterprise) {
  router.push({ name: 'AdminEnterpriseMembers', params: { enterpriseId: enterprise.id }, query: { enterpriseName: enterprise.name } })
}

// ---- Balance Modal ----
const showBalanceModal = ref(false)
const balanceEnterprise = ref<Enterprise | null>(null)
const balanceOperation = ref<'add' | 'subtract'>('add')
const showBalanceHistoryModal = ref(false)
const balanceHistoryEnterprise = ref<Enterprise | null>(null)

function handleDeposit(enterprise: Enterprise) {
  balanceEnterprise.value = enterprise
  balanceOperation.value = 'add'
  showBalanceModal.value = true
}

function handleWithdraw(enterprise: Enterprise) {
  balanceEnterprise.value = enterprise
  balanceOperation.value = 'subtract'
  showBalanceModal.value = true
}

function handleBalanceHistory(enterprise: Enterprise) {
  balanceHistoryEnterprise.value = enterprise
  showBalanceHistoryModal.value = true
}
</script>
