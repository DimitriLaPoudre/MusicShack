<script lang="ts">
	import type { UserResponse } from "$lib/types/response/user";
	import { onMount } from "svelte";
	import User from "./User.svelte";
	import { GetUsers } from "$lib/usecases/user";
	import { toast } from "svelte-sonner";
	import Dialog from "../bits-ui/Dialog.svelte";
	import Button from "../bits-ui/Button.svelte";
	import CreateUser from "./CreateUser.svelte";

	let users = $state<UserResponse[]>([]);
	let createUserOpen = $state<boolean>(false);

	onMount(async () => {
		try {
			users = await GetUsers();
		} catch (e) {
			return toast.error(
				e instanceof Error ? e.message : "Unknown error",
			);
		}
	});

	function afterCreateUser(user: UserResponse) {
		users.push(user);
		createUserOpen = false;
	}

</script>

<div class="flex flex-col border">
	<h2>Users</h2>

	<Dialog title={"test"} bind:open={createUserOpen}>
		{#snippet trigger()}
			<Button>Add</Button>
		{/snippet}

		<CreateUser {afterCreateUser} />
	</Dialog>


	<!-- <div> -->
	<!-- 	{#if !createUserOpen} -->
	<!-- 		<button -->
	<!-- 			onclick={(e) => { -->
	<!-- 				e.preventDefault(); -->
	<!-- 				createUserOpen = true; -->
	<!-- 			}} -->
	<!-- 		> -->
	<!-- 			Create User -->
	<!-- 		</button> -->
	<!-- 	{:else} -->
	<!-- 		<CreateUser afterCreateUser /> -->
	<!-- 	{/if} -->
	<!-- </div> -->

	<div class="flex flex-col gap-2 border p-1">
		{#each users as user}
			<User
				{user}
				reloadUsers={async () => {
					console.log("ntm");
					users = await GetUsers();
				}}
			/>
		{/each}
	</div>
</div>
