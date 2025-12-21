<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { adminFetch } from "$lib/functions/fetch";
	import { Plus, Trash } from "lucide-svelte";

	let errorPassword = $state<null | string>(null);
	let errorUser = $state<null | string>(null);
	let errorLogout = $state<null | string>(null);

	let inputNewPassword = $state<null | string>(null);
	let inputOldPassword = $state<null | string>(null);
	let inputUserUsername = $state<null | string>(null);
	let inputUserPassword = $state<null | string>(null);

	let users = $state<null | any>(null);

	afterNavigate(() => {
		loadUsers();
	});

	async function changePassword() {
		try {
			const res = await adminFetch("/admin/password", "PUT", {
				oldPassword: inputOldPassword,
				newPassword: inputNewPassword,
			});
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to update password");
			}
			inputNewPassword = null;
			inputOldPassword = null;
			errorPassword = null;
		} catch (e) {
			errorPassword =
				e instanceof Error
					? e.message
					: "Failed to update admin password";
		}
	}

	async function loadUsers() {
		try {
			const res = await adminFetch("/users");
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to fetch users");
			}
			users = body;
		} catch (e) {
			errorUser =
				e instanceof Error ? e.message : "Failed to reload user list";
		}
	}

	async function createUser() {
		try {
			const res = await adminFetch("/users", "POST", {
				username: inputUserUsername,
				password: inputUserPassword,
			});
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to create user");
			}
			inputUserUsername = null;
			inputUserPassword = null;
			await loadUsers();
			errorUser = null;
		} catch (e) {
			errorUser =
				e instanceof Error ? e.message : "Failed to create user";
		}
	}

	async function deleteUser(userId: string) {
		try {
			const res = await adminFetch(`/users/${userId}`, "DELETE");
			const body = await res.json();
			if (!res.ok) {
				throw new Error(body.error || "Failed to delete user");
			}
			await loadUsers();
			errorUser = null;
		} catch (e) {
			errorUser =
				e instanceof Error ? e.message : "Failed to delete user";
		}
	}

	async function logout() {
		try {
			const res = await adminFetch(`/admin/logout`, "POST");
			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || "error while trying to logout");
			}
			goto("/admin/login");
			errorLogout = null;
		} catch (e) {
			errorLogout = e instanceof Error ? e.message : "Failed to logout";
			return;
		}
	}
</script>

<div class="body">
	<h1 class="title">Admin Dashboard</h1>
	<h2>Password</h2>
	<div class="change-password">
		{#if errorPassword}
			<p class="error">
				{errorPassword}
			</p>
		{/if}
		<form class="form" onsubmit={changePassword}>
			<div class="inputs">
				<input
					placeholder="New Password"
					bind:value={inputNewPassword}
				/>
				<input
					placeholder="Actual Password"
					bind:value={inputOldPassword}
				/>
			</div>
			<button><Plus /></button>
		</form>
	</div>
	<h2>Users</h2>
	<div class="users">
		{#if errorUser}
			<p class="error">
				{errorUser}
			</p>
		{/if}
		<form
			class="form"
			onsubmit={async () => {
				await createUser();
			}}
		>
			<div class="inputs">
				<input placeholder="Username" bind:value={inputUserUsername} />
				<input placeholder="Password" bind:value={inputUserPassword} />
			</div>
			<button><Plus /></button>
		</form>
		{#if !users}
			<p class="loading">Loading...</p>
		{:else}
			<div class="items">
				{#each users as user}
					<div class="item">
						<div class="data">
							<p>{user.Username}</p>
						</div>
						<button
							onclick={async () => {
								await deleteUser(user.ID);
								await loadUsers();
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

	.body {
		display: flex;
		flex-direction: column;
		align-items: center;
		width: 70vw;
		height: 100vh;
		margin: 0 auto;
		gap: 16px;

		.title {
			text-align: center;
			margin-top: 20px;
		}

		.form {
			display: grid;
			grid-template-columns: 1fr auto;
			gap: 8px;
			width: 100%;
			align-items: stretch;
			container-type: inline-size;

			.inputs {
				display: grid;
				grid-template-columns: 1fr 1fr;
				gap: 8px;
			}
			button {
				aspect-ratio: 1/1;
			}

			@container (max-width: 420px) {
				.inputs {
					grid-template-columns: 1fr;
				}
			}
		}

		.change-password {
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
					.data:hover {
						outline: 1px solid #ffffff;
						outline-offset: -1px;
					}
					button {
						aspect-ratio: 1/1;
					}
					@container (max-width: 420px) {
						.data {
							grid-template-columns: 1fr;
						}
					}
				}
			}
		}

		.logout {
			width: 100%;
			border-color: var(--err);
		}
		.logout:hover {
			background-color: var(--err);
		}
	}
</style>
