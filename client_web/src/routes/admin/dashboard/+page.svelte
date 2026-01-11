<script lang="ts">
	import { goto } from "$app/navigation";
	import { adminFetch } from "$lib/functions/fetch";
	import type { RequestAdminPassword, RequestUser } from "$lib/types/request";
	import type { StatusResponse, UsersResponse } from "$lib/types/response";
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
	let users = $state<null | UsersResponse>(null);

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
			const data = await adminFetch<UsersResponse>("/users");
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

<div class="body">
	<h1 class="title">Admin Dashboard</h1>
	<h2 class="sub-title">Password</h2>
	<div class="change-password">
		{#if errorPassword}
			<p class="error">
				{errorPassword}
			</p>
		{/if}
		<form
			class="form"
			onsubmit={(e) => {
				e.preventDefault();
				changePassword();
			}}
		>
			<div class="inputs">
				<input
					placeholder="New Password"
					bind:value={inputAdminPassword.newPassword}
				/>
				<input
					placeholder="Actual Password"
					bind:value={inputAdminPassword.oldPassword}
				/>
			</div>
			<button class="hover-full"><Plus /></button>
		</form>
	</div>
	<h2 class="sub-title">Users</h2>
	<div class="users">
		{#if errorUser}
			<p class="error">
				{errorUser}
			</p>
		{/if}
		<form
			class="form"
			onsubmit={async (e) => {
				e.preventDefault();
				await createUser();
			}}
		>
			<div class="inputs">
				<input placeholder="Username" bind:value={inputUser.username} />
				<input placeholder="Password" bind:value={inputUser.password} />
			</div>
			<button class="hover-full"><Plus /></button>
		</form>
		{#if !users}
			<p class="loading">Loading...</p>
		{:else}
			<div class="items">
				{#each users as user}
					<div class="item">
						<div class="data hover-soft">
							<p>{user.username}</p>
						</div>
						<button
							class="hover-full"
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
		<p class="error">
			{errorLogout}
		</p>
	{/if}
	<button class="logout" onclick={logout}> Logout </button>
</div>

<style>
	.loading {
		text-align: center;
	}

	.error {
		text-align: center;
		background-color: var(--err);
		padding: 0.5rem;
		width: 100%;
	}

	.sub-title {
		font-weight: bolder;
	}

	.body {
		display: flex;
		flex-direction: column;
		align-items: center;
		width: clamp(320px, 70vw + 20px, 1200px);
		height: 100vh;
		margin: 0 auto;
		gap: 16px;

		.title {
			text-align: center;
			font-weight: bolder;
			margin-top: 20px;
		}
		.form {
			display: grid;
			grid-template-columns: 1fr auto;
			gap: 0.75rem;
			width: 100%;
			align-items: stretch;
			container-type: inline-size;

			.inputs {
				display: grid;
				grid-template-columns: 1fr 1fr;
				gap: 0.75rem;
			}

			@container (max-width: 520px) {
				.inputs {
					grid-template-columns: 1fr;
				}
			}
		}
		.change-password {
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 16px;
			width: 100%;
		}
		.users {
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 16px;
			width: 100%;

			.items {
				display: flex;
				flex-direction: column;
				gap: 4px;
				width: 100%;

				.item {
					display: grid;
					grid-template-columns: 1fr auto;
					gap: 8px;
					align-items: stretch;
					container-type: inline-size;

					.data {
						display: grid;
						grid-template-columns: 1fr;
						gap: 8px;
						align-items: center;
						padding: 1rem;
					}
					button {
						aspect-ratio: 1/1;
					}
					@container (max-width: 520px) {
						.data {
							grid-template-columns: 1fr;
						}
					}
				}
			}
		}
		.logout {
			width: 100%;
			box-shadow: inset 0 0 0 1px var(--err);
		}
		.logout:hover {
			background-color: var(--err);
		}
	}
</style>
