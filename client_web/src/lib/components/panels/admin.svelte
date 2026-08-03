<script lang="ts">
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Error } from "$lib/components/ui/error/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import {
		ApiError,
		createUser,
		deleteUser,
		getMe,
		listUsers,
	} from "$lib/api";
	import { Plus, Shield, Trash } from "@lucide/svelte";
	import { onMount } from "svelte";
	import { isAdmin, openPanel, usersList } from "$lib/stores";

	let {
		anchor = undefined,
	}: {
		anchor?: HTMLElement | null;
	} = $props();

	let error = $state<string | null>(null);
	let userForm = $state({ username: "", password: "", hi_res: true });

	function message(e: unknown, fallback: string) {
		return e instanceof ApiError ? e.message : fallback;
	}

	async function loadUsers() {
		try {
			usersList.set(await listUsers());
			error = null;
		} catch (e) {
			error = message(e, "Failed to reload user list");
		}
	}

	async function createUserForm(event: SubmitEvent) {
		event.preventDefault();
		try {
			await createUser({
				username: userForm.username,
				password: userForm.password,
				hi_res: userForm.hi_res,
			});
			userForm = { username: "", password: "", hi_res: true };
			error = null;
			await loadUsers();
		} catch (e) {
			error = message(e, "Failed to create user");
		}
	}

	async function removeUser(id: string) {
		try {
			await deleteUser(id);
			error = null;
			await loadUsers();
		} catch (e) {
			error = message(e, "Failed to delete user");
		}
	}

	onMount(async () => {
		if ($isAdmin) {
			if ($usersList === null) await loadUsers();
			return;
		}
		try {
			const me = await getMe();
			if (me.role !== "admin") return;
			isAdmin.set(true);
			if ($usersList === null) await loadUsers();
		} catch {
			// not authenticated
		}
	});
</script>

{#if $isAdmin}
	<Popover.Root
		open={$openPanel === "admin"}
		onOpenChange={(v) => openPanel.set(v ? "admin" : null)}
	>
		<Popover.Trigger>
			{#snippet child({ props })}
				<Button {...props} variant="hover-full" size="icon-panel">
					<Shield />
				</Button>
			{/snippet}
		</Popover.Trigger>
		<Popover.Content customAnchor={anchor} class="w-full">
			<h1 class="text-xs font-semibold tracking-widest uppercase">Admin</h1>
			<Error text={error} />
			<h2 class="text-xs font-semibold tracking-widest uppercase">Users</h2>
			<form
				class="grid grid-cols-[1fr_auto] items-stretch gap-2"
				onsubmit={createUserForm}
			>
				<div class="grid grid-cols-2 gap-2">
					<Input
						type="text"
						placeholder="Username"
						bind:value={userForm.username}
					/>
					<Input
						type="password"
						placeholder="Password"
						bind:value={userForm.password}
					/>
				</div>
				<Button type="submit" variant="hover-full" size="icon">
					<Plus />
				</Button>
			</form>
			{#if $usersList === null}
				<p class="py-2 text-center text-xs text-muted-foreground">
					Loading...
				</p>
			{:else if $usersList.length === 0}
				<p class="py-2 text-center text-xs text-muted-foreground">
					No users
				</p>
			{:else}
				<div class="flex flex-col gap-1">
					{#each $usersList as user}
						<div class="grid grid-cols-[1fr_auto] items-stretch gap-2">
							<div
								class="flex items-center gap-2 p-4 hover:shadow-[inset_0_0_0_1px_var(--foreground)]"
							>
								<p class="break-words">{user.username}</p>
								{#if user.role === "admin"}
									<p class="text-xs text-muted-foreground">
										(admin)
									</p>
								{/if}
							</div>
							<Button
								variant="hover-full"
								size="icon"
								onclick={() => removeUser(user.id)}
							>
								<Trash />
							</Button>
						</div>
					{/each}
				</div>
			{/if}
		</Popover.Content>
	</Popover.Root>
{/if}
