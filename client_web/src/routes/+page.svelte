<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { apiFetch } from "$lib/functions/fetch";
	import type { StatusResponse } from "$lib/types/response";

	afterNavigate(async () => {
		const data = await apiFetch<StatusResponse>("/me");
		if ("error" in data) {
			return;
		}
		goto("/dashboard");
	});
</script>
