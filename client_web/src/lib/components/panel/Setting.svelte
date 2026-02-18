<script lang="ts">
	import { goto } from "$app/navigation";
	import { apiFetch } from "$lib/functions/fetch";
	import { onMount } from "svelte";
	import { Pencil, Plus, Trash } from "lucide-svelte";
	import type { RequestInstance, RequestUser } from "$lib/types/request";
	import type {
		InstancesResponse,
		StatusResponse,
		UserResponse,
	} from "$lib/types/response";
	import { instanceList, userData } from "$lib/stores/panel/setting";

	let errorUser = $state<null | string>(null);
	let inputUser = $state<RequestUser>({
		username: "",
		password: "",
		hiRes: $userData?.hiRes || false,
	});

	let errorInstances = $state<null | string>(null);
	let inputInstance = $state<RequestInstance>({ url: "" });

	onMount(() => {
		loadInstance();
		getUser();
	});

	async function getUser() {
		try {
			const data = await apiFetch<UserResponse>("/me");
			$userData = data;
			inputUser.hiRes = data.hiRes;
			errorUser = null;
		} catch (e) {
			errorUser =
				e instanceof Error ? e.message : "Failed to get user info";
		}
	}

	async function changeUser(event: SubmitEvent) {
		event.preventDefault();
		try {
			const data = await apiFetch<UserResponse>("/me", "PUT", inputUser);
			$userData = data;
			errorUser = null;
			if (inputUser.username !== "" || inputUser.password !== "") {
				inputUser = {
					username: "",
					password: "",
					hiRes: data.hiRes,
				};
				await logout();
			} else {
				inputUser = {
					username: "",
					password: "",
					hiRes: data.hiRes,
				};
			}
		} catch (e) {
			errorUser =
				e instanceof Error ? e.message : "Failed to update user info";
		}
	}

	async function loadInstance() {
		try {
			const data = await apiFetch<InstancesResponse>(`/instances`);
			$instanceList = data;
			errorInstances = null;
		} catch (e) {
			errorInstances =
				e instanceof Error
					? e.message
					: "Failed to reload instances queue";
		}
	}

	async function addInstance(event: SubmitEvent) {
		event.preventDefault();
		try {
			inputInstance.url = inputInstance.url.trim();
			if (inputInstance.url.endsWith("/")) {
				inputInstance.url = inputInstance.url.substring(
					0,
					inputInstance.url.length - 1,
				);
			}

			if (!inputInstance.url) {
				errorInstances = "fill url with valid value";
				return;
			}

			await apiFetch<StatusResponse>(`/instances`, "POST", inputInstance);
			inputInstance = { url: "" };
			loadInstance();
		} catch (e) {
			errorInstances =
				e instanceof Error ? e.message : "Failed to add instance";
			return;
		}
	}

	async function deleteInstance(id: number) {
		try {
			await apiFetch<StatusResponse>(`/instances/${id}`, "DELETE");
			loadInstance();
		} catch (e) {
			errorInstances =
				e instanceof Error ? e.message : "Failed to delete instance";
		}
	}

	async function logout() {
		try {
			await apiFetch<StatusResponse>(`/logout`, "POST");
			goto("/login");
		} catch (e) {
			errorUser = e instanceof Error ? e.message : "Failed to logout";
		}
	}
</script>

<div class="flex flex-col gap-3">
	<h1 class="font-extrabold">Settings</h1>
	<div class="flex flex-col gap-2">
		<h2 class="font-extrabold">User</h2>
		<div class="flex flex-col p-3 gap-3">
			{#if errorUser}
				<p class="text-center bg-err p-2">
					{errorUser}
				</p>
			{/if}
			<form
				class="grid grid-cols-[1fr_auto] gap-2 items-stretch @container"
				onsubmit={changeUser}
			>
				<div class="flex flex-col gap-2">
					<div
						class="grid grid-cols-2 @max-[520px]:grid-cols-1 gap-2"
					>
						<input
							placeholder={$userData?.username || "username"}
							bind:value={inputUser.username}
						/>
						<input
							placeholder="password"
							bind:value={inputUser.password}
						/>
					</div>
					<div
						class="grid grid-cols-2 @max-[520px]:grid-cols-1 gap-2"
					>
						<button
							type="button"
							class="py-4 hover:shadow-[inset_0_0_0_1px_var(--fg)] focus:outline-none active:bg-fg active:text-bg"
							class:underline={inputUser.hiRes !== true}
							onclick={() => {
								inputUser.hiRes = false;
							}}
						>
							LOSSLESS (recommended)
						</button>
						<button
							type="button"
							class="py-4 hover:shadow-[inset_0_0_0_1px_var(--fg)] focus:outline-none active:bg-fg active:text-bg"
							class:underline={inputUser.hiRes === true}
							onclick={() => {
								inputUser.hiRes = true;
							}}
						>
							HIRES (advanced)
						</button>
					</div>
				</div>
				<button class="hover-full">
					<Pencil />
				</button>
			</form>
		</div>
	</div>
	<div class="flex flex-col gap-2">
		<h2 class="font-extrabold">Instances</h2>
		<div class="flex flex-col p-3 gap-5">
			{#if errorInstances}
				<p class="text-center bg-err p-2">
					{errorInstances}
				</p>
			{/if}
			<form
				class="grid grid-cols-[1fr_auto] gap-2 items-stretch @container"
				onsubmit={addInstance}
			>
				<input
					class="w-full"
					placeholder="URL"
					bind:value={inputInstance.url}
				/>
				<button class="hover-full"><Plus /></button>
			</form>
			{#if !$instanceList}
				<p class="text-center">Loading...</p>
			{:else}
				<div class="flex flex-col gap-2">
					{#each $instanceList as instance}
						<div
							class="grid grid-cols-[1fr_auto] gap-2 items-stretch @container"
						>
							<div
								class="hover-soft grid grid-cols-[1fr_auto_6ch] @max-[520px]:grid-cols-1 gap-3 items-center p-4"
							>
								<p class="warp-break-words">{instance.url}</p>
								<p class="warp-break-words">
									{instance.provider}|{instance.api}
								</p>
								<p
									class="justify-self-end @max-[520px]:justify-self-start"
								>
									{#if instance.ping === 0}
										failed
									{:else}
										{instance.ping}ms
									{/if}
								</p>
							</div>
							<button
								class="hover-full"
								onclick={() => deleteInstance(instance.id)}
							>
								<Trash />
							</button>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
	<button
		class="w-full py-3 shadow-[inset_0_0_0_1px_var(--err)] hover:bg-err"
		onclick={logout}
	>
		Logout
	</button>
</div>
