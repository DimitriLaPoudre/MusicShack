<script lang="ts">
	import { LoginAdmin } from "$lib/usecases/admin";

	let password = $state("");
</script>

<form
	class=""
	onsubmit={(e) => {
		e.preventDefault();
		LoginAdmin(password);
	}}
>
	<div class="">
		<input
			class=""
			id="password"
			type="password"
			bind:value={password}
			placeholder="password"
		/>
	</div>
	<div class="">
		<button type="submit">Login</button>
	</div>
</form>
