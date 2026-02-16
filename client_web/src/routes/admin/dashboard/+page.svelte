<script lang="ts">
	import { goto } from "$app/navigation";
	import { adminFetch } from "$lib/functions/fetch";
	import type { RequestAdminPassword, RequestUser } from "$lib/types/request";
	import type {
		StatusResponse,
		AdminUsersResponse,
	} from "$lib/types/response";
	import { Plus, Trash } from "lucide-svelte";
	import { onMount } from "svelte";

	let errorPassword = $state<null | string>(null);
	let errorUser = $state<null | string>(null);
	let errorLogout = $state<null | string>(null);

	let inputAdminPassword = $state<RequestAdminPassword>({
		oldPassword: "",
		newPassword: "",
	});

	let inputUser = $state<RequestUser>({
		username: "",
		password: "",
		hiRes: true,
	});
	let users = $state<null | AdminUsersResponse>(null);

	onMount(() => {
		loadUsers();
	});

	async function changePassword() {
		try {
			const data = await adminFetch<StatusResponse>(
				"/admin/password",
				"PUT",
				inputAdminPassword,
			);
			if ("error" in data) {
				throw new Error(data.error || "Failed to update password");
			}
			inputAdminPassword = { oldPassword: "", newPassword: "" };
			errorPassword = null;
			await logout();
		} catch (e) {
			errorPassword =
				e instanceof Error
					? e.message
					: "Failed to update admin password";
		}
	}

	async function loadUsers() {
		try {
			const data = await adminFetch<AdminUsersResponse>("/users");
			if ("error" in data) {
				throw new Error(data.error || "Failed to fetch users");
			}
			users = data;
			errorUser = null;
		} catch (e) {
			errorUser =
				e instanceof Error ? e.message : "Failed to reload user list";
		}
	}

	async function createUser() {
		try {
			const data = await adminFetch<StatusResponse>(
				"/users",
				"POST",
				inputUser,
			);
			if ("error" in data) {
				throw new Error(data.error || "Failed to create user");
			}
			inputUser = { username: "", password: "", hiRes: true };
			errorUser = null;
			await loadUsers();
		} catch (e) {
			errorUser =
				e instanceof Error ? e.message : "Failed to create user";
		}
	}

	async function deleteUser(userId: number) {
		try {
			const data = await adminFetch<StatusResponse>(
				`/users/${userId}`,
				"DELETE",
			);
			if ("error" in data) {
				throw new Error(data.error || "Failed to delete user");
			}
			errorUser = null;
			await loadUsers();
		} catch (e) {
			errorUser =
				e instanceof Error ? e.message : "Failed to delete user";
		}
	}

	async function logout() {
		try {
			const data = await adminFetch<StatusResponse>(
				`/admin/logout`,
				"POST",
			);
			if ("error" in data) {
				throw new Error(data.error || "error while trying to logout");
			}
			goto("/admin/login");
			errorLogout = null;
		} catch (e) {
			errorLogout = e instanceof Error ? e.message : "Failed to logout";
		}
	}
</script>

<svelte:head>
	<title>Admin | Dashboard - MusicShack</title>
</svelte:head>

<div
	class="flex flex-col items-center w-[clamp(320px,70vw+20px,1200px)] h-screen mx-auto gap-4"
>
	<h1 class="text-center font-extrabold mt-5">Admin Dashboard</h1>
	<h2 class="font-extrabold">Password</h2>
	<div class="flex flex-col items-center gap-4 w-full">
		{#if errorPassword}
			<p class="text-center bg-err p-2 w-full">
				{errorPassword}
			</p>
		{/if}
		<form
			class="grid grid-cols-[1fr_auto] gap-3 w-full items-stretch @container"
			onsubmit={(e) => {
				e.preventDefault();
				changePassword();
			}}
		>
			<div
				class="grid grid-cols-2 @[0px]:grid-cols-1 @[520px]:grid-cols-2 gap-3"
			>
				<input
					placeholder="New Password"
					bind:value={inputAdminPassword.newPassword}
				/>
				<input
					placeholder="Actual Password"
					bind:value={inputAdminPassword.oldPassword}
				/>
			</div>
			<button class="hover-full w-14 h-15"><Plus /></button>
		</form>
	</div>
	<h2 class="font-extrabold">Users</h2>
	<div class="flex flex-col items-center gap-4 w-full">
		{#if errorUser}
			<p class="text-center bg-err p-2 w-full">
				{errorUser}
			</p>
		{/if}
		<form
			class="grid grid-cols-[1fr_auto] gap-3 w-full items-stretch @container"
			onsubmit={async (e) => {
				e.preventDefault();
				await createUser();
			}}
		>
			<div
				class="grid grid-cols-2 @[0px]:grid-cols-1 @[520px]:grid-cols-2 gap-3"
			>
				<input placeholder="Username" bind:value={inputUser.username} />
				<input placeholder="Password" bind:value={inputUser.password} />
			</div>
			<button class="hover-full w-14 h-15"><Plus /></button>
		</form>
		{#if !users}
			<p class="text-center">Loading...</p>
		{:else}
			<div class="flex flex-col gap-1 w-full">
				{#each users as user}
					<div
						class="grid grid-cols-[1fr_auto] gap-2 items-stretch @container"
					>
						<div
							class="hover-soft grid grid-cols-1 gap-2 items-center p-4"
						>
							<p>{user.username}</p>
						</div>
						<button
							class="hover-full w-14 h-15"
							onclick={async () => {
								await deleteUser(user.id);
							}}
						>
							<Trash />
						</button>
					</div>
				{/each}
			</div>
		{/if}
	</div>
	{#if errorLogout}
		<p class="text-center bg-err p-2 w-full">
			{errorLogout}
		</p>
	{/if}
	<button
		class="w-full shadow-[inset_0_0_0_1px_var(--err)] hover:bg-err"
		onclick={logout}
	>
		Logout
	</button>
</div>
