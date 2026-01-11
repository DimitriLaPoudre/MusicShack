<script lang="ts">
	export let quality: { name: string; color: string } = {
		name: "UNKNOWN",
		color: "grey",
	};
</script>

<span class="video-badge" style="background-color: {quality.color}">
	{quality.name}
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
