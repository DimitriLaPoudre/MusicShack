<script lang="ts">
	import type { CreateUserRequest } from "$lib/types/request/user";
	import { CreateUser } from "$lib/usecases/user";
	import { toast } from "svelte-sonner";
	import Checkbox from "../bits-ui/Checkbox.svelte";
	import type { UserResponse } from "$lib/types/response/user";

	let req = $state<CreateUserRequest>({
		username: "",
		password: "",
		hi_res: false,
	});

	type Props = {
		afterCreateUser: (user: UserResponse) => void;
	};

	let { afterCreateUser }: Props = $props();
</script>

<form
	onsubmit={async (e) => {
		e.preventDefault();

		try {
			const user = await CreateUser(req);
			afterCreateUser(user);
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
