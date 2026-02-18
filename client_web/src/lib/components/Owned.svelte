<script>
	import { DownloadIcon } from "lucide-svelte";

	export let color = "#00FF00";
</script>

<div class="grid place-content-center">
	<p style="color:{color};">
		<DownloadIcon size={20} />
	</p>
</div>
