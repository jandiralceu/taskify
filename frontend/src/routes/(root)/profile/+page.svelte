<script lang="ts">
  import { goto } from '$app/navigation'
  import { resolve } from '$app/paths'
  import {
    LoaderCircle,
    Camera,
    ShieldCheck,
    User,
    KeyRound,
    Trash2,
    TriangleAlert,
    Eye,
    EyeOff,
  } from '@lucide/svelte'
  import { Dialog, Portal } from '@skeletonlabs/skeleton-svelte'
  import {
    createProfileQuery,
    uploadAvatarMutation,
    updateProfileMutation,
    changePasswordMutation,
    deleteProfileMutation,
  } from '$lib/state/user.svelte'
  import { toaster } from '$lib/state/toast.svelte'
  import Input from '$lib/components/Input.svelte'

  const profileQuery = createProfileQuery()
  const uploadAvatar = uploadAvatarMutation()
  const updateProfile = updateProfileMutation()
  const changePassword = changePasswordMutation()
  const deleteProfile = deleteProfileMutation()

  // Avatar upload
  let fileInput = $state<HTMLInputElement>(null!)
  let avatarPreview = $state<string | null>(null)
  const MAX_FILE_SIZE = 5 * 1024 * 1024 // 5 MB

  function onAvatarClick() {
    fileInput.click()
  }

  async function onFileSelected(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return

    if (!file.type.startsWith('image/')) {
      toaster.error({
        title: 'Invalid File',
        description: 'Please select an image file (JPEG, PNG or WebP).',
      })
      return
    }

    if (file.size > MAX_FILE_SIZE) {
      toaster.error({
        title: 'File Too Large',
        description: 'Avatar must be smaller than 5 MB.',
      })
      return
    }

    // Show local preview immediately
    avatarPreview = URL.createObjectURL(file)

    try {
      await uploadAvatar.mutateAsync(file)
      toaster.success({
        title: 'Avatar Updated',
        description: 'Your profile picture has been saved.',
      })
    } catch {
      avatarPreview = null
      toaster.error({
        title: 'Upload Failed',
        description: 'Could not upload your avatar. Please try again.',
      })
    } finally {
      // Reset input so the same file can be re-selected if needed
      fileInput.value = ''
    }
  }

  // Edit info state
  let firstName = $state('')
  let lastName = $state('')

  // Change password state
  let oldPassword = $state('')
  let newPassword = $state('')
  let confirmPassword = $state('')
  let showOldPassword = $state(false)
  let showNewPassword = $state(false)

  // Delete confirmation
  let isDeleteOpen = $state(false)
  let deleteConfirmText = $state('')
  const DELETE_PHRASE = 'delete my account'

  $effect(() => {
    if (profileQuery.data) {
      firstName = profileQuery.data.firstName
      lastName = profileQuery.data.lastName
    }
  })

  function formatDate(dateStr: string) {
    return new Intl.DateTimeFormat('en-GB', {
      day: '2-digit',
      month: 'long',
      year: 'numeric',
    }).format(new Date(dateStr))
  }

  async function handleUpdateInfo(e: SubmitEvent) {
    e.preventDefault()
    try {
      await updateProfile.mutateAsync({ firstName, lastName })
      toaster.success({
        title: 'Profile Updated',
        description: 'Your personal information has been saved.',
      })
    } catch {
      toaster.error({
        title: 'Update Failed',
        description: 'Could not update your profile. Please try again.',
      })
    }
  }

  async function handleChangePassword(e: SubmitEvent) {
    e.preventDefault()
    if (newPassword !== confirmPassword) {
      toaster.error({
        title: 'Password Mismatch',
        description: 'New passwords do not match.',
      })
      return
    }
    try {
      await changePassword.mutateAsync({ oldPassword, newPassword })
      oldPassword = ''
      newPassword = ''
      confirmPassword = ''
      toaster.success({
        title: 'Password Changed',
        description: 'Your password has been updated successfully.',
      })
    } catch {
      toaster.error({
        title: 'Change Failed',
        description:
          'Could not change your password. Check your current password and try again.',
      })
    }
  }

  async function handleDeleteAccount() {
    try {
      await deleteProfile.mutateAsync()
      goto(resolve('/signin'))
    } catch {
      toaster.error({
        title: 'Delete Failed',
        description: 'Could not delete your account. Please try again.',
      })
    }
  }

  const canConfirmDelete = $derived(
    deleteConfirmText.toLowerCase() === DELETE_PHRASE
  )
</script>

<svelte:head>
  <title>My Profile - Taskify</title>
</svelte:head>

<div class="flex h-full flex-col pt-8">
  <header class="px-8 pb-8">
    <h2 class="text-4xl leading-tight tracking-tight text-surface-900">
      <span class="font-light">My</span>
      <span class="font-normal"> Profile</span>
    </h2>
  </header>

  <div class="flex-1 overflow-y-auto px-8 pb-8">
    {#if profileQuery.isPending}
      <div
        class="flex flex-col items-center justify-center py-24 text-surface-400"
      >
        <LoaderCircle size={28} class="mb-3 animate-spin" />
        <span class="text-sm font-medium">Loading profile...</span>
      </div>
    {:else if profileQuery.isError}
      <div
        class="flex flex-col items-center justify-center py-24 text-rose-500"
      >
        <p class="text-sm font-medium">Failed to load profile.</p>
      </div>
    {:else if profileQuery.data}
      {@const user = profileQuery.data}
      <div class="max-w-2xl space-y-6">
        <!-- Identity card -->
        <div
          class="flex items-center gap-5 rounded-2xl border border-surface-200 bg-white p-6"
        >
          <!-- Avatar -->
          <div class="relative shrink-0">
            {#if avatarPreview || user.avatarUrl}
              <img
                src={avatarPreview ?? user.avatarUrl}
                alt="{user.firstName} {user.lastName}"
                class="size-20 rounded-2xl border border-surface-200 object-cover"
              />
            {:else}
              <div
                class="flex size-20 items-center justify-center rounded-2xl border border-primary-100 bg-primary-50"
              >
                <span class="text-2xl font-bold text-primary-600">
                  {user.firstName[0]}{user.lastName[0]}
                </span>
              </div>
            {/if}

            <!-- Hidden file input -->
            <input
              bind:this={fileInput}
              type="file"
              accept="image/jpeg,image/png,image/webp"
              class="hidden"
              onchange={onFileSelected}
            />

            <button
              type="button"
              onclick={onAvatarClick}
              disabled={uploadAvatar.isPending}
              title="Change profile picture"
              class="absolute -right-2 -bottom-2 flex size-7 items-center justify-center rounded-lg border border-surface-200 bg-white shadow-sm transition-colors hover:border-primary-300 hover:text-primary-500 disabled:cursor-not-allowed disabled:opacity-60
              {uploadAvatar.isPending ? 'text-primary-400' : 'text-surface-400'}"
            >
              {#if uploadAvatar.isPending}
                <LoaderCircle size={13} class="animate-spin" />
              {:else}
                <Camera size={13} />
              {/if}
            </button>
          </div>

          <div class="min-w-0 flex-1">
            <p class="truncate text-lg font-semibold text-surface-900">
              {user.firstName}
              {user.lastName}
            </p>
            <p class="truncate text-sm text-surface-500">{user.email}</p>
            <div class="mt-2 flex items-center gap-2">
              {#if user.role === 'admin'}
                <span
                  class="inline-flex items-center gap-1 rounded-full bg-violet-50 px-2.5 py-1 text-[11px] font-semibold text-violet-700"
                >
                  <ShieldCheck size={11} /> Admin
                </span>
              {:else}
                <span
                  class="inline-flex items-center gap-1 rounded-full bg-sky-50 px-2.5 py-1 text-[11px] font-semibold text-sky-700"
                >
                  <User size={11} /> Employee
                </span>
              {/if}
              <span class="text-xs text-surface-400"
                >Member since {formatDate(user.createdAt)}</span
              >
            </div>
          </div>
        </div>

        <!-- Edit personal info -->
        <form
          onsubmit={handleUpdateInfo}
          class="space-y-5 rounded-2xl border border-surface-200 bg-white p-6"
        >
          <div class="mb-1 flex items-center gap-2">
            <User size={16} class="text-surface-400" />
            <h3 class="text-sm font-semibold text-surface-900">
              Personal Information
            </h3>
          </div>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <Input
              id="firstName"
              label="First Name"
              placeholder="First name"
              bind:value={firstName}
              required
            />
            <Input
              id="lastName"
              label="Last Name"
              placeholder="Last name"
              bind:value={lastName}
              required
            />
          </div>

          <!-- Read-only fields -->
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div class="space-y-1">
              <label
                for="ro-email"
                class="block text-sm font-medium text-surface-700">Email</label
              >
              <input
                id="ro-email"
                type="text"
                value={user.email}
                disabled
                class="block h-12 w-full cursor-not-allowed rounded-xl border border-surface-200 bg-surface-100 px-4 text-sm text-surface-400"
              />
            </div>
            <div class="space-y-1">
              <label
                for="ro-role"
                class="block text-sm font-medium text-surface-700">Role</label
              >
              <input
                id="ro-role"
                type="text"
                value={user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                disabled
                class="block h-12 w-full cursor-not-allowed rounded-xl border border-surface-200 bg-surface-100 px-4 text-sm text-surface-400"
              />
            </div>
          </div>

          <div class="flex justify-end pt-1">
            <button
              type="submit"
              disabled={updateProfile.isPending}
              class="flex items-center gap-2 rounded-xl bg-primary-500 px-8 py-2.5 text-sm font-bold text-white shadow-lg shadow-primary-500/20 transition-all hover:bg-primary-700 active:scale-95 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {#if updateProfile.isPending}
                <LoaderCircle size={15} class="animate-spin" />
                Saving...
              {:else}
                Save Changes
              {/if}
            </button>
          </div>
        </form>

        <!-- Change password -->
        <form
          onsubmit={handleChangePassword}
          class="space-y-5 rounded-2xl border border-surface-200 bg-white p-6"
        >
          <div class="mb-1 flex items-center gap-2">
            <KeyRound size={16} class="text-surface-400" />
            <h3 class="text-sm font-semibold text-surface-900">
              Change Password
            </h3>
          </div>

          <div class="space-y-1">
            <label
              for="oldPassword"
              class="block text-sm font-medium text-surface-700"
              >Current Password</label
            >
            <div class="relative">
              <input
                id="oldPassword"
                type={showOldPassword ? 'text' : 'password'}
                placeholder="Current password"
                bind:value={oldPassword}
                required
                class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 pr-10 text-sm text-surface-900 placeholder-surface-500 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 focus:outline-none"
              />
              <button
                type="button"
                onclick={() => (showOldPassword = !showOldPassword)}
                class="absolute inset-y-0 right-0 flex items-center pr-3 text-surface-400 transition-colors hover:text-primary-500"
                tabindex="-1"
              >
                {#if showOldPassword}<Eye size={18} />{:else}<EyeOff
                    size={18}
                  />{/if}
              </button>
            </div>
          </div>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div class="space-y-1">
              <label
                for="newPassword"
                class="block text-sm font-medium text-surface-700"
                >New Password</label
              >
              <div class="relative">
                <input
                  id="newPassword"
                  type={showNewPassword ? 'text' : 'password'}
                  placeholder="New password"
                  bind:value={newPassword}
                  required
                  minlength={8}
                  class="block h-12 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 pr-10 text-sm text-surface-900 placeholder-surface-500 transition-all focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 focus:outline-none"
                />
                <button
                  type="button"
                  onclick={() => (showNewPassword = !showNewPassword)}
                  class="absolute inset-y-0 right-0 flex items-center pr-3 text-surface-400 transition-colors hover:text-primary-500"
                  tabindex="-1"
                >
                  {#if showNewPassword}<Eye size={18} />{:else}<EyeOff
                      size={18}
                    />{/if}
                </button>
              </div>
            </div>

            <div class="space-y-1">
              <label
                for="confirmPassword"
                class="block text-sm font-medium text-surface-700"
                >Confirm New Password</label
              >
              <input
                id="confirmPassword"
                type="password"
                placeholder="Confirm new password"
                bind:value={confirmPassword}
                required
                class="block h-12 w-full rounded-xl border transition-all
								{confirmPassword && confirmPassword !== newPassword
                  ? 'border-rose-400 focus:border-rose-400 focus:ring-rose-400/20'
                  : 'border-surface-300 focus:border-primary-500 focus:ring-primary-500/20'}
								bg-surface-50 px-4 text-sm text-surface-900 placeholder-surface-500 focus:ring-2 focus:outline-none"
              />
              {#if confirmPassword && confirmPassword !== newPassword}
                <p class="mt-1 text-xs font-medium text-rose-500">
                  Passwords do not match
                </p>
              {/if}
            </div>
          </div>

          <div class="flex justify-end pt-1">
            <button
              type="submit"
              disabled={changePassword.isPending ||
                (!!confirmPassword && confirmPassword !== newPassword)}
              class="flex items-center gap-2 rounded-xl bg-primary-500 px-8 py-2.5 text-sm font-bold text-white shadow-lg shadow-primary-500/20 transition-all hover:bg-primary-700 active:scale-95 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {#if changePassword.isPending}
                <LoaderCircle size={15} class="animate-spin" />
                Updating...
              {:else}
                Update Password
              {/if}
            </button>
          </div>
        </form>

        <!-- Danger zone -->
        <div class="space-y-4 rounded-2xl border border-rose-200 bg-white p-6">
          <div class="flex items-center gap-2">
            <TriangleAlert size={16} class="text-rose-500" />
            <h3 class="text-sm font-semibold text-rose-600">Danger Zone</h3>
          </div>

          <div
            class="flex items-center justify-between rounded-xl border border-rose-100 bg-rose-50 p-4"
          >
            <div>
              <p class="text-sm font-medium text-surface-800">Delete Account</p>
              <p class="mt-0.5 text-xs text-surface-500">
                Permanently remove your account and all associated data. This
                cannot be undone.
              </p>
            </div>
            <button
              type="button"
              onclick={() => {
                deleteConfirmText = ''
                isDeleteOpen = true
              }}
              class="ml-4 shrink-0 rounded-xl border border-rose-300 px-4 py-2 text-sm font-bold text-rose-600 transition-all hover:border-rose-500 hover:bg-rose-500 hover:text-white"
            >
              Delete Account
            </button>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>

<!-- Delete Account Confirmation Dialog -->
<Dialog
  role="alertdialog"
  open={isDeleteOpen}
  onOpenChange={e => {
    if (!e.open) isDeleteOpen = false
  }}
>
  <Portal>
    <Dialog.Backdrop
      class="fixed inset-0 z-50 bg-surface-950/40 backdrop-blur-sm"
    />
    <Dialog.Positioner
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
    >
      <Dialog.Content
        class="w-full max-w-md space-y-4 rounded-2xl border border-surface-100 bg-white p-6 shadow-2xl"
      >
        <div class="flex items-start gap-4">
          <div
            class="flex size-10 shrink-0 items-center justify-center rounded-xl bg-rose-50 text-rose-500"
          >
            <Trash2 size={20} />
          </div>
          <div>
            <Dialog.Title class="text-base font-bold text-surface-900">
              Delete Your Account
            </Dialog.Title>
            <Dialog.Description class="mt-1 text-sm text-surface-500">
              This will permanently delete your account. To confirm, type
              <span class="font-semibold text-surface-700 select-all"
                >"{DELETE_PHRASE}"</span
              > below.
            </Dialog.Description>
          </div>
        </div>

        <input
          type="text"
          placeholder={DELETE_PHRASE}
          bind:value={deleteConfirmText}
          class="block h-11 w-full rounded-xl border border-surface-300 bg-surface-50 px-4 text-sm text-surface-900 placeholder-surface-400 transition-all focus:border-rose-400 focus:ring-2 focus:ring-rose-400/20 focus:outline-none"
        />

        <div class="flex items-center justify-end gap-3 pt-1">
          <Dialog.CloseTrigger
            type="button"
            class="rounded-xl px-5 py-2.5 text-sm font-medium text-surface-500 transition-all hover:bg-surface-50 hover:text-surface-900"
          >
            Cancel
          </Dialog.CloseTrigger>
          <button
            type="button"
            onclick={handleDeleteAccount}
            disabled={!canConfirmDelete || deleteProfile.isPending}
            class="flex items-center gap-2 rounded-xl bg-rose-500 px-5 py-2.5 text-sm font-bold text-white transition-all hover:bg-rose-600 active:scale-95 disabled:cursor-not-allowed disabled:opacity-40"
          >
            {#if deleteProfile.isPending}
              <LoaderCircle size={15} class="animate-spin" />
              Deleting...
            {:else}
              <Trash2 size={15} />
              Delete My Account
            {/if}
          </button>
        </div>
      </Dialog.Content>
    </Dialog.Positioner>
  </Portal>
</Dialog>
