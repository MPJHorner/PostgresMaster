<script lang="ts">
	import {
		Button,
		Card,
		CardHeader,
		CardTitle,
		CardContent,
		Badge,
		Separator
	} from '$lib/components/ui';

	// GitHub releases base URL - this will be updated once we have actual releases
	const RELEASES_BASE_URL = 'https://github.com/MPJHorner/PostgresMaster/releases/latest/download';

	const platforms = [
		{
			name: 'Windows (x64)',
			icon: '🪟',
			filename: 'postgres-proxy-windows-amd64.exe',
			command: '.\\postgres-proxy-windows-amd64.exe'
		},
		{
			name: 'macOS (Intel)',
			icon: '🍎',
			filename: 'postgres-proxy-darwin-amd64',
			command: './postgres-proxy-darwin-amd64'
		},
		{
			name: 'macOS (Apple Silicon)',
			icon: '🍎',
			filename: 'postgres-proxy-darwin-arm64',
			command: './postgres-proxy-darwin-arm64'
		},
		{
			name: 'Linux (x64)',
			icon: '🐧',
			filename: 'postgres-proxy-linux-amd64',
			command: './postgres-proxy-linux-amd64'
		}
	];

	function downloadProxy(filename: string) {
		const url = `${RELEASES_BASE_URL}/${filename}`;
		window.open(url, '_blank');
	}
</script>

<div class="min-h-screen bg-gradient-to-b from-slate-50 to-slate-100">
	<!-- Hero Section -->
	<div class="container mx-auto px-4 py-16 sm:py-24">
		<div class="text-center max-w-4xl mx-auto space-y-6">
			<Badge variant="secondary" class="mb-4">100% Browser-Based • Zero Backend</Badge>
			<h1 class="text-5xl sm:text-6xl font-bold text-slate-900 tracking-tight">
				PostgreSQL Client<br />That Lives in Your Browser
			</h1>
			<p class="text-xl text-slate-600 max-w-2xl mx-auto">
				A lightning-fast, secure SQL client that requires zero backend infrastructure. Connect to any
				Postgres database through a tiny local proxy.
			</p>
			<div class="flex flex-col sm:flex-row gap-4 justify-center mt-8">
				<Button
					size="lg"
					class="text-lg px-8"
					onclick={() => {
						document.getElementById('download')?.scrollIntoView({ behavior: 'smooth' });
					}}
				>
					Get Started
				</Button>
				<Button
					variant="outline"
					size="lg"
					class="text-lg px-8"
					onclick={() => window.open('https://github.com/MPJHorner/PostgresMaster', '_blank')}
				>
					View on GitHub
				</Button>
			</div>
		</div>
	</div>

	<Separator class="my-16" />

	<!-- Features Section -->
	<div class="container mx-auto px-4 py-16">
		<h2 class="text-3xl font-bold text-center text-slate-900 mb-12">Why PostgresMaster?</h2>
		<div class="grid md:grid-cols-3 gap-8 max-w-6xl mx-auto">
			<!-- Feature 1: Any Postgres -->
			<Card class="border-2 hover:border-slate-300 transition-colors">
				<CardHeader>
					<div class="text-4xl mb-4">🎯</div>
					<CardTitle>Any Postgres Database</CardTitle>
				</CardHeader>
				<CardContent>
					<p class="text-slate-600">
						Connect to any remote Postgres server - AWS RDS, Google Cloud SQL, DigitalOcean, or your
						own server. Works with all versions of PostgreSQL.
					</p>
				</CardContent>
			</Card>

			<!-- Feature 2: Lightning Fast -->
			<Card class="border-2 hover:border-slate-300 transition-colors">
				<CardHeader>
					<div class="text-4xl mb-4">⚡</div>
					<CardTitle>Lightning Fast</CardTitle>
				</CardHeader>
				<CardContent>
					<p class="text-slate-600">
						Sub-2-second load times. Intelligent SQL autocomplete with schema-aware suggestions.
						Execute queries instantly with real-time results.
					</p>
				</CardContent>
			</Card>

			<!-- Feature 3: Secure -->
			<Card class="border-2 hover:border-slate-300 transition-colors">
				<CardHeader>
					<div class="text-4xl mb-4">🔒</div>
					<CardTitle>100% Secure</CardTitle>
				</CardHeader>
				<CardContent>
					<p class="text-slate-600">
						Your credentials never leave your machine. The local proxy handles all database
						connections. No third-party servers involved.
					</p>
				</CardContent>
			</Card>
		</div>
	</div>

	<Separator class="my-16" />

	<!-- Get Started Section -->
	<div class="container mx-auto px-4 py-16">
		<h2 class="text-3xl font-bold text-center text-slate-900 mb-12">Get Started in 3 Steps</h2>
		<div class="max-w-4xl mx-auto space-y-8">
			<!-- Step 1 -->
			<Card>
				<CardHeader>
					<div class="flex items-center gap-4">
						<Badge class="text-lg px-4 py-2">1</Badge>
						<CardTitle>Download the Proxy</CardTitle>
					</div>
				</CardHeader>
				<CardContent>
					<p class="text-slate-600 mb-4">
						Download the tiny (~5MB) proxy binary for your operating system. No installation
						required.
					</p>
					<div class="bg-slate-900 text-slate-100 p-4 rounded-lg font-mono text-sm">
						# Choose your platform below
					</div>
				</CardContent>
			</Card>

			<!-- Step 2 -->
			<Card>
				<CardHeader>
					<div class="flex items-center gap-4">
						<Badge class="text-lg px-4 py-2">2</Badge>
						<CardTitle>Run the Proxy</CardTitle>
					</div>
				</CardHeader>
				<CardContent>
					<p class="text-slate-600 mb-4">
						Run the proxy with your Postgres connection string or use interactive mode.
					</p>
					<div class="bg-slate-900 text-slate-100 p-4 rounded-lg font-mono text-sm space-y-2">
						<div class="text-slate-400"># Option A: Connection String</div>
						<div>./postgres-proxy "postgres://user:pass@host:5432/dbname"</div>
						<div class="mt-4 text-slate-400"># Option B: Interactive Mode</div>
						<div>./postgres-proxy</div>
						<div class="text-slate-400 text-xs mt-2">
							# Follow the prompts to enter your connection details
						</div>
					</div>
				</CardContent>
			</Card>

			<!-- Step 3 -->
			<Card>
				<CardHeader>
					<div class="flex items-center gap-4">
						<Badge class="text-lg px-4 py-2">3</Badge>
						<CardTitle>Open in Browser</CardTitle>
					</div>
				</CardHeader>
				<CardContent>
					<p class="text-slate-600 mb-4">
						The proxy will print a URL with a secure secret. Open it in your browser and start
						querying!
					</p>
					<div class="bg-slate-900 text-slate-100 p-4 rounded-lg font-mono text-sm space-y-1">
						<div class="text-green-400">✓ Connected to database!</div>
						<div class="text-slate-400">→ Open in browser:</div>
						<div class="text-blue-400 break-all">
							http://localhost:5173?secret=a1b2c3d4e5f6...
						</div>
					</div>
				</CardContent>
			</Card>
		</div>
	</div>

	<Separator class="my-16" />

	<!-- Download Section -->
	<div id="download" class="container mx-auto px-4 py-16">
		<h2 class="text-3xl font-bold text-center text-slate-900 mb-4">Download PostgresMaster</h2>
		<p class="text-center text-slate-600 mb-12 max-w-2xl mx-auto">
			Choose your platform and download the proxy binary. All binaries are under 10MB.
		</p>
		<div class="grid md:grid-cols-2 lg:grid-cols-4 gap-4 max-w-6xl mx-auto">
			{#each platforms as platform}
				<Card class="hover:shadow-lg transition-shadow">
					<CardContent class="pt-6">
						<div class="text-center space-y-4">
							<div class="text-5xl">{platform.icon}</div>
							<h3 class="font-semibold text-slate-900">{platform.name}</h3>
							<Button
								variant="outline"
								class="w-full"
								onclick={() => downloadProxy(platform.filename)}
							>
								Download
							</Button>
							<div class="text-xs text-slate-500 font-mono break-all">
								{platform.filename}
							</div>
						</div>
					</CardContent>
				</Card>
			{/each}
		</div>
		<div class="mt-8 text-center">
			<p class="text-sm text-slate-500">
				All downloads are served from GitHub Releases •
				<a
					href="https://github.com/MPJHorner/PostgresMaster/releases"
					target="_blank"
					class="text-blue-600 hover:underline"
				>
					View all releases
				</a>
			</p>
		</div>
	</div>

	<Separator class="my-16" />

	<!-- Usage Example Section -->
	<div class="container mx-auto px-4 py-16">
		<h2 class="text-3xl font-bold text-center text-slate-900 mb-12">Usage Example</h2>
		<div class="max-w-4xl mx-auto">
			<Card>
				<CardContent class="pt-6">
					<div class="bg-slate-900 text-slate-100 p-6 rounded-lg font-mono text-sm space-y-3">
						<div class="text-slate-400"># Start the proxy</div>
						<div class="text-blue-400">$ ./postgres-proxy</div>
						<div class="text-slate-300">Host: db.example.com</div>
						<div class="text-slate-300">Port [5432]: 5432</div>
						<div class="text-slate-300">Database: myapp_production</div>
						<div class="text-slate-300">Username: admin</div>
						<div class="text-slate-300">Password: ****</div>
						<div class="text-slate-300">SSL Mode [prefer]: require</div>
						<div class="mt-4 text-slate-400">Connecting to database...</div>
						<div class="text-green-400">✓ Connected successfully!</div>
						<div class="text-slate-300">
							→ Open in browser: http://localhost:5173?secret=a1b2c3d4e5f6g7h8
						</div>
						<div class="mt-4 text-slate-400">
							# Now open the URL in your browser and start querying!
						</div>
					</div>
				</CardContent>
			</Card>
		</div>
	</div>

	<Separator class="my-16" />

	<!-- Footer -->
	<footer class="container mx-auto px-4 py-12 border-t">
		<div class="max-w-6xl mx-auto">
			<div class="grid md:grid-cols-3 gap-8 mb-8">
				<!-- Tech Stack -->
				<div>
					<h3 class="font-semibold text-slate-900 mb-4">Built With</h3>
					<div class="space-y-2 text-sm text-slate-600">
						<div>• SvelteKit + TypeScript</div>
						<div>• Monaco Editor (VS Code)</div>
						<div>• shadcn-svelte + Tailwind CSS</div>
						<div>• Go + pgx (Postgres driver)</div>
					</div>
				</div>

				<!-- Features -->
				<div>
					<h3 class="font-semibold text-slate-900 mb-4">Features</h3>
					<div class="space-y-2 text-sm text-slate-600">
						<div>• SQL syntax highlighting</div>
						<div>• Intelligent autocomplete</div>
						<div>• Schema introspection</div>
						<div>• Query execution & results</div>
					</div>
				</div>

				<!-- Links -->
				<div>
					<h3 class="font-semibold text-slate-900 mb-4">Resources</h3>
					<div class="space-y-2 text-sm text-slate-600">
						<a
							href="https://github.com/MPJHorner/PostgresMaster"
							target="_blank"
							class="block hover:text-blue-600"
						>
							GitHub Repository
						</a>
						<a
							href="https://github.com/MPJHorner/PostgresMaster/issues"
							target="_blank"
							class="block hover:text-blue-600"
						>
							Report Issues
						</a>
						<a
							href="https://github.com/MPJHorner/PostgresMaster/blob/main/prd.md"
							target="_blank"
							class="block hover:text-blue-600"
						>
							Documentation
						</a>
					</div>
				</div>
			</div>

			<!-- Copyright -->
			<div class="text-center text-sm text-slate-500 pt-8 border-t">
				<p>
					PostgresMaster • Licensed under
					<a
						href="https://github.com/MPJHorner/PostgresMaster/blob/main/LICENSE"
						target="_blank"
						class="text-blue-600 hover:underline"
					>
						AGPL-3.0
					</a>
				</p>
				<p class="mt-2">Built with ❤️ for the PostgreSQL community</p>
			</div>
		</div>
	</footer>
</div>

<style>
	:global(html) {
		scroll-behavior: smooth;
	}
</style>
