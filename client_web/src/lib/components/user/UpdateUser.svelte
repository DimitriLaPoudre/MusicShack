<script lang="ts">
	import type { UpdateUserRequest } from "$lib/types/request/user";
	import { toast } from "svelte-sonner";
	import Checkbox from "../bits-ui/Checkbox.svelte";
	import type { UserResponse } from "$lib/types/response/user";
	import { UpdateUser } from "$lib/usecases/user";

	let req = $state<UpdateUserRequest>({
		username: undefined,
		password: undefined,
		hi_res: undefined,
	});

	type Props = {
		id: string;
		afterUpdateUser: (user: UserResponse) => void;
	};

	let { id, afterUpdateUser }: Props = $props();
</script>

<form
	onsubmit={async (e) => {
		e.preventDefault();

		try {
			const user = await UpdateUser(id, req);
			afterUpdateUser(user);
		} catch (e) {
			return toast.error(
				e instanceof Error ? e.message : "Unknown error",
			);
		}
	}}
>
	<input
		class=""
		name={"username"}
		bind:value={req.username}
		placeholder="Username"
	/>
	<input
		class=""
		name={"password"}
		bind:value={req.password}
		placeholder="Password"
	/>
	<Checkbox bind:checked={req.hi_res} name={"hi_res"} id={"hi_res"}>
		HiRes
	</Checkbox>
	<button>Create</button>
</form>
