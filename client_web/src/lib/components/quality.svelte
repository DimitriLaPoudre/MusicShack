<script lang="ts">
	import type { AudioQuality } from "$lib/types";

	let {
		quality = { name: "UNKNOWN", color: "grey" },
	}: {
		quality: AudioQuality;
	} = $props();
</script>

<span
	class="inline-block w-[10ch] text-center px-1.5 py-0.5 text-xs font-extrabold rounded text-white"
	style="background-color: {quality.color}"
>
	{quality.name}
</span>
