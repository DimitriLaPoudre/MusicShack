<script lang="ts">
	import * as Card from "$lib/components/ui/card/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import {
		FieldGroup,
		Field,
		FieldLabel,
	} from "$lib/components/ui/field/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { login } from "$lib/api";
	import { goto } from "$app/navigation";
	import type { LoginForm } from "$lib/types";
	import Checkbox from "./ui/checkbox/checkbox.svelte";
	import Error from "./ui/error/error.svelte";

	let loginForm = $state<LoginForm>({
		username: "",
		password: "",
		remember: true,
	});
	let error = $state<string | null>(null);
	let loading = $state(false);

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();

		loading = true;
		error = null;

		try {
			await login(loginForm);
			await goto("/dashboard");
		} catch {
			error = "Invalid Credentials";
		} finally {
			loading = false;
		}
	}
</script>

<Card.Root class="mx-auto w-full max-w-sm">
	<Card.Header>
		<Card.Title class="text-2xl">Login</Card.Title>
		<Card.Description>
			Enter your username below to login to your account
		</Card.Description>
	</Card.Header>
	<Card.Content>
		<form onsubmit={handleSubmit}>
			<FieldGroup>
				<Field>
					<FieldLabel>Username</FieldLabel>
					<Input
						type="text"
						placeholder="totally_normal_user"
						bind:value={loginForm.username}
						required
					/>
				</Field>
				<Field>
					<FieldLabel>Password</FieldLabel>
					<Input
						type="password"
						placeholder="***************"
						bind:value={loginForm.password}
						required
					/>
				</Field>
				<Field>
					<div class="flex items-center gap-2">
						<Checkbox bind:checked={loginForm.remember} />
						<FieldLabel>Remeber me</FieldLabel>
					</div>
				</Field>

				<Error text={error} />

				<Field>
					<Button
						type="submit"
						variant="hover-full"
						disabled={loading}
					>
						{loading ? "Login..." : "Login"}
					</Button>
				</Field>
			</FieldGroup>
		</form>
	</Card.Content>
</Card.Root>
