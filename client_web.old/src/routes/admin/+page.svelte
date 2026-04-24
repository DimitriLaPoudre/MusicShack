<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { adminFetch } from "$lib/functions/fetch";
	import type { StatusResponse } from "$lib/types/response";

	afterNavigate(async () => {
		const data = await adminFetch<StatusResponse>("/admin");
		if ("error" in data) {
			return;
		}
		goto("/admin/dashboard");
	});
</script>
