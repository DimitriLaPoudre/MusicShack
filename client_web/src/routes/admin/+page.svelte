<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";
	import { adminFetch } from "$lib/functions/fetch";

	afterNavigate(async () => {
		const res = await adminFetch("/admin");
		if (!res.ok) {
			return;
		}
		goto("/admin/dashboard");
	});
</script>
