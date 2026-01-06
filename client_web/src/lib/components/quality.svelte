<script lang="ts">
	export let quality: string = "LOW";
	const colors: Record<string, string> = {
		LOW: "#ff4b4b",
		HIGH: "#ff993d",
		LOSSLESS: "#eee400",
		HIRES: "#47ee37",
	};
</script>

<span class="video-badge" style="background-color: {colors[quality] || 'gray'}">
	{quality}
</span>

<style>
	.video-badge {
		display: inline-block;
		width: 10ch;
		text-align: center;
		padding: 2px 6px;
		font-size: 0.75rem;
		font-weight: bolder;
		border-radius: 4px;
		color: white;
	}
</style>
