<script>
	import { CircleArrowDown } from "lucide-svelte";

	export let color = "#00FF00";
</script>

<div class="grid place-content-center">
	<p style="color:{color};">
		<CircleArrowDown size={20} />
	</p>
</div>
