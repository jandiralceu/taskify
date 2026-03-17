<script lang="ts">
  import { goto } from '$app/navigation'
  import { resolve } from '$app/paths'
  import { resolveAvatarUrl } from '$lib/utils/avatar'
  import {
    LoaderCircle,
    Camera,
    ShieldCheck,
    User,
    KeyRound,
    Trash2,
    TriangleAlert,
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
  import Button from '$lib/components/Button.svelte'

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
                src={avatarPreview ?? resolveAvatarUrl(user.avatarUrl) ?? ''}
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
              {uploadAvatar.isPending
                ? 'text-primary-400'
                : 'text-surface-600'}"
            >
              {#if uploadAvatar.isPending}
                <LoaderCircle size={13} class="animate-spin" />
              {:else}
                <Camera size={13} />
              {/if}
            </button>
          </div>

          <div class="min-w-0 flex-1">
            <p class="truncate text-lg font-normal text-surface-900">
              {user.firstName}
              {user.lastName}
            </p>
            <p class="truncate text-sm text-surface-600">{user.email}</p>
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
              <span class="text-xs text-surface-600"
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
          <div class="mb-6 flex items-center gap-2">
            <User size={24} class="text-surface-700" />
            <h3 class="text-sm font-medium text-surface-800">
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
            <Input id="ro-email" label="Email" value={user.email} disabled />
            <Input
              id="ro-role"
              label="Role"
              value={user.role.charAt(0).toUpperCase() + user.role.slice(1)}
              disabled
            />
          </div>

          <div class="flex justify-end pt-1">
            <Button
              type="submit"
              loading={updateProfile.isPending}
              loadingText="Saving..."
            >
              Save Changes
            </Button>
          </div>
        </form>

        <!-- Change password -->
        <form
          onsubmit={handleChangePassword}
          class="space-y-5 rounded-2xl border border-surface-200 bg-white p-6"
        >
          <div class="mb-6 flex items-center gap-2">
            <KeyRound size={24} class="text-surface-700" />
            <h3 class="text-sm font-medium text-surface-800">
              Change Password
            </h3>
          </div>

          <Input
            id="oldPassword"
            type="password"
            label="Current Password"
            placeholder="Current password"
            bind:value={oldPassword}
            required
          />

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <Input
              id="newPassword"
              type="password"
              label="New Password"
              placeholder="New password"
              bind:value={newPassword}
              required
              minlength={8}
            />
            <Input
              id="confirmPassword"
              type="password"
              label="Confirm New Password"
              placeholder="Confirm new password"
              bind:value={confirmPassword}
              required
              error={confirmPassword && confirmPassword !== newPassword
                ? 'Passwords do not match'
                : undefined}
            />
          </div>

          <div class="flex justify-end pt-1">
            <Button
              type="submit"
              disabled={!!confirmPassword && confirmPassword !== newPassword}
              loading={changePassword.isPending}
              loadingText="Updating..."
            >
              Update Password
            </Button>
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
            <Button
              variant="danger"
              onclick={handleDeleteAccount}
              disabled={!canConfirmDelete}
              loading={deleteProfile.isPending}
              loadingText="Deleting..."
            >
              <Trash2 size={15} />
              Delete My Account
            </Button>
          </div>
      </Dialog.Content>
    </Dialog.Positioner>
  </Portal>
</Dialog>
