<script lang="ts">
	import { LoginAdmin } from "$lib/usecases/admin";
	import { toast } from "svelte-sonner";

	let password = $state<string>("");
</script>

<form
	class=""
	onsubmit={async (e) => {
		e.preventDefault();

		try {
			await LoginAdmin(password);
		} catch (e) {
			return toast.error(
				e instanceof Error ? e.message : "Unknown error",
			);
		}

		password = "";
	}}
>
	<div class="">
		<input
			class=""
			type="password"
			bind:value={password}
			placeholder="password"
		/>
	</div>
	<div class="">
		<button type="submit">Login</button>
	</div>
</form>
