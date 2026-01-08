<script>
	export let color = "#000";
	export let bg = "#fff";
</script>

<div style="background-color:{bg};">
	<p style="color:{color};">E</p>
</div>

<style>
	div {
		display: grid;
		place-content: center;
		min-width: 18px;
		width: 18px;
		min-height: 18px;
		height: 18px;
		border-radius: 10%;

		p {
			font-size: 13px;
		}
	}
</style>
