<script lang="ts">
	import {
		Button,
		Input,
		Card,
		CardHeader,
		CardTitle,
		CardContent,
		Table,
		TableHeader,
		TableBody,
		TableRow,
		TableHead,
		TableCell,
		Alert,
		AlertTitle,
		AlertDescription,
		Badge,
		Separator
	} from '$lib/components/ui';
</script>

<div class="min-h-screen bg-slate-50 p-8">
	<div class="max-w-6xl mx-auto space-y-8">
		<div>
			<h1 class="text-4xl font-bold text-slate-900 mb-2">shadcn-svelte Component Test</h1>
			<p class="text-slate-600">Testing all installed components</p>
		</div>

		<Separator />

		<!-- Buttons -->
		<Card>
			<CardHeader>
				<CardTitle>Buttons</CardTitle>
			</CardHeader>
			<CardContent>
				<div class="flex flex-wrap gap-2">
					<Button variant="default">Default</Button>
					<Button variant="destructive">Destructive</Button>
					<Button variant="outline">Outline</Button>
					<Button variant="secondary">Secondary</Button>
					<Button variant="ghost">Ghost</Button>
					<Button variant="link">Link</Button>
				</div>
			</CardContent>
		</Card>

		<!-- Inputs -->
		<Card>
			<CardHeader>
				<CardTitle>Input</CardTitle>
			</CardHeader>
			<CardContent>
				<div class="space-y-2">
					<Input type="text" placeholder="Enter your name..." />
					<Input type="email" placeholder="Email address" />
					<Input type="password" placeholder="Password" />
				</div>
			</CardContent>
		</Card>

		<!-- Badges -->
		<Card>
			<CardHeader>
				<CardTitle>Badges</CardTitle>
			</CardHeader>
			<CardContent>
				<div class="flex flex-wrap gap-2">
					<Badge variant="default">Default</Badge>
					<Badge variant="secondary">Secondary</Badge>
					<Badge variant="destructive">Destructive</Badge>
					<Badge variant="outline">Outline</Badge>
				</div>
			</CardContent>
		</Card>

		<!-- Alerts -->
		<Card>
			<CardHeader>
				<CardTitle>Alerts</CardTitle>
			</CardHeader>
			<CardContent>
				<div class="space-y-4">
					<Alert>
						<AlertTitle>Default Alert</AlertTitle>
						<AlertDescription>
							This is a default alert message with some information.
						</AlertDescription>
					</Alert>
					<Alert variant="destructive">
						<AlertTitle>Error Alert</AlertTitle>
						<AlertDescription>
							This is a destructive alert showing an error message.
						</AlertDescription>
					</Alert>
				</div>
			</CardContent>
		</Card>

		<!-- Table -->
		<Card>
			<CardHeader>
				<CardTitle>Table</CardTitle>
			</CardHeader>
			<CardContent>
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>ID</TableHead>
							<TableHead>Name</TableHead>
							<TableHead>Email</TableHead>
							<TableHead>Status</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow>
							<TableCell>1</TableCell>
							<TableCell>John Doe</TableCell>
							<TableCell>john@example.com</TableCell>
							<TableCell><Badge variant="secondary">Active</Badge></TableCell>
						</TableRow>
						<TableRow>
							<TableCell>2</TableCell>
							<TableCell>Jane Smith</TableCell>
							<TableCell>jane@example.com</TableCell>
							<TableCell><Badge variant="secondary">Active</Badge></TableCell>
						</TableRow>
						<TableRow>
							<TableCell>3</TableCell>
							<TableCell>Bob Johnson</TableCell>
							<TableCell>bob@example.com</TableCell>
							<TableCell><Badge variant="outline">Inactive</Badge></TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</CardContent>
		</Card>
	</div>
</div>
