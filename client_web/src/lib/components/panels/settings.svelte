<script lang="ts">
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Error } from "$lib/components/ui/error/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import {
		ApiError,
		createInstance,
		deleteInstance,
		getMe,
		listInstances,
		logout,
		updateMe,
	} from "$lib/api";
	import { Pencil, Plus, Settings, Trash } from "@lucide/svelte";
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import type { UpdateUser } from "$lib/types";
	import { instanceList, openPanel, userData } from "$lib/stores";

	let {
		anchor = undefined,
	}: {
		anchor?: HTMLElement | null;
	} = $props();

	let userError = $state<string | null>(null);
	let userForm = $state({ username: "", password: "", hi_res: false });

	let instanceError = $state<string | null>(null);
	let instanceUrl = $state("");

	function message(e: unknown, fallback: string) {
		return e instanceof ApiError ? e.message : fallback;
	}

	async function loadUser() {
		try {
			const user = await getMe();
			userData.set(user);
			userForm = { username: "", password: "", hi_res: user.hi_res };
			userError = null;
		} catch (e) {
			userError = message(e, "Failed to get user info");
		}
	}

	async function changeUser(event: SubmitEvent) {
		event.preventDefault();

		const patch: UpdateUser = { hi_res: userForm.hi_res };
		if (userForm.username) patch.username = userForm.username;
		if (userForm.password) patch.password = userForm.password;

		try {
			const updated = await updateMe(patch);
			userData.set(updated);
			userForm = { username: "", password: "", hi_res: updated.hi_res };
			userError = null;
			if (patch.username || patch.password) {
				await handleLogout();
			}
		} catch (e) {
			userError = message(e, "Failed to update user info");
		}
	}

	async function loadInstances() {
		try {
			instanceList.set(await listInstances());
			instanceError = null;
		} catch (e) {
			instanceError = message(e, "Failed to reload instances");
		}
	}

	async function addInstance(event: SubmitEvent) {
		event.preventDefault();

		const url = instanceUrl.trim().replace(/\/+$/, "");
		if (!url) {
			instanceError = "fill url with valid value";
			return;
		}

		try {
			await createInstance({ url });
			instanceUrl = "";
			await loadInstances();
		} catch (e) {
			instanceError = message(e, "Failed to add instance");
		}
	}

	async function removeInstance(id: string) {
		try {
			await deleteInstance(id);
			await loadInstances();
		} catch (e) {
			instanceError = message(e, "Failed to delete instance");
		}
	}

	async function handleLogout() {
		try {
			await logout();
			await goto("/login");
		} catch (e) {
			userError = message(e, "Failed to logout");
		}
	}

	onMount(() => {
		if ($userData === null) loadUser();
		else userForm.hi_res = $userData.hi_res;
		if ($instanceList === null) loadInstances();
	});
</script>

<Popover.Root
	open={$openPanel === "settings"}
	onOpenChange={(v) => openPanel.set(v ? "settings" : null)}
>
	<Popover.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="hover-full" size="icon-panel">
				<Settings />
			</Button>
		{/snippet}
	</Popover.Trigger>
	<Popover.Content customAnchor={anchor} class="w-full">
		<h1 class="text-xs font-semibold tracking-widest uppercase">Settings</h1>

		<div class="flex flex-col gap-2">
			<h2 class="text-xs font-semibold tracking-widest uppercase">User</h2>
			<Error text={userError} />
			<form
				class="grid grid-cols-[1fr_auto] items-stretch gap-2"
				onsubmit={changeUser}
			>
				<div class="flex flex-col gap-2">
					<div class="grid grid-cols-2 gap-2">
						<Input
							type="text"
							placeholder={$userData?.username || "username"}
							bind:value={userForm.username}
						/>
						<Input
							type="password"
							placeholder="password"
							bind:value={userForm.password}
						/>
					</div>
					<div class="grid grid-cols-2 gap-2">
						<Button
							type="button"
							variant="hover-soft"
							class={!userForm.hi_res ? "underline" : ""}
							onclick={() => (userForm.hi_res = false)}
						>
							Lossless (recommended)
						</Button>
						<Button
							type="button"
							variant="hover-soft"
							class={userForm.hi_res ? "underline" : ""}
							onclick={() => (userForm.hi_res = true)}
						>
							Hires (advanced)
						</Button>
					</div>
				</div>
				<Button type="submit" variant="hover-full" size="icon">
					<Pencil />
				</Button>
			</form>
		</div>

		<div class="flex flex-col gap-2">
			<h2 class="text-xs font-semibold tracking-widest uppercase">Instances</h2>
			<Error text={instanceError} />
			<form
				class="grid grid-cols-[1fr_auto] items-stretch gap-2"
				onsubmit={addInstance}
			>
				<Input type="text" placeholder="URL" bind:value={instanceUrl} />
				<Button type="submit" variant="hover-full" size="icon">
					<Plus />
				</Button>
			</form>
			{#if $instanceList === null}
				<p class="py-2 text-center text-xs text-muted-foreground">
					Loading...
				</p>
			{:else if $instanceList.length === 0}
				<p class="py-2 text-center text-xs text-muted-foreground">
					No instances
				</p>
			{:else}
				<div class="flex flex-col gap-2">
					{#each $instanceList as instance}
						<div class="grid grid-cols-[1fr_auto] items-stretch gap-2">
							<div
								class="grid grid-cols-[1fr_auto_6ch] items-center gap-3 p-4 hover:shadow-[inset_0_0_0_1px_var(--foreground)]"
							>
								<p class="break-words">{instance.url}</p>
								<p class="break-words">
									{instance.provider}|{instance.plugin}
								</p>
								<p class="justify-self-end">
									{#if instance.ping === null}
										failed
									{:else}
										{instance.ping}ms
									{/if}
								</p>
							</div>
							<Button
								variant="hover-full"
								size="icon"
								onclick={() => removeInstance(instance.id)}
							>
								<Trash />
							</Button>
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<button
			type="button"
			class="w-full py-2 text-xs font-semibold tracking-widest uppercase shadow-[inset_0_0_0_1px_var(--destructive)] text-destructive transition-all hover:bg-destructive hover:text-background hover:shadow-none"
			onclick={handleLogout}
		>
			Logout
		</button>
	</Popover.Content>
</Popover.Root>
