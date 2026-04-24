<script lang="ts">
	export let quality: { name: string; color: string } = {
		name: "UNKNOWN",
		color: "grey",
	};
</script>

<span
	class="inline-block w-[10ch] text-center px-1.5 py-0.5 text-xs font-extrabold rounded text-white"
	style="background-color: {quality.color}"
>
	{quality.name}
</span>
