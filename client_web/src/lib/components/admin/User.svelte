<script lang="ts">
	import type { UserResponse } from "$lib/types/response/user";
	import { Trash } from "@lucide/svelte";
	import { DeleteUser } from "$lib/usecases/user";
	import Button from "../bits-ui/Button.svelte";
	import { toast } from "svelte-sonner";

	export let user: UserResponse;
	export let reloadUsers: () => Promise<void>;
</script>

<div class="flex gap-1 border p-1">
	<div class="flex gap-1 border p-1">
		<p>Username</p>
		<p>{user.username}</p>
	</div>
	<div class="flex gap-1 border p-1">
		<p>HiRes</p>
		<p>{user.hi_res}</p>
	</div>
	<div class="border p-1">
		<Button
			onclick={async () => {
				try {
					await DeleteUser(user.id);
					await reloadUsers();
				} catch (e) {
					return toast.error(
						e instanceof Error ? e.message : "Unknown error",
					);
				}
			}}
		>
			<Trash />
		</Button>
	</div>
</div>
