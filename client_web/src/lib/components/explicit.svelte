<script lang="ts">
	let {
		color = "#000",
		bg = "#fff",
	}: {
		color?: string;
		bg?: string;
	} = $props();
</script>

<div
	class="grid place-content-center min-w-4.5 w-4.5 min-h-4.5 h-4.5 rounded-sm"
	style="background-color:{bg};"
>
	<p class="text-[13px]" style="color:{color};">E</p>
</div>
