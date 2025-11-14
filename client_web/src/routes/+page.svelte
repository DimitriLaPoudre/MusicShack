<script lang="ts">
	import { afterNavigate, goto } from "$app/navigation";

	afterNavigate(async () => {
		goto("/dashboard");
	});
</script>
