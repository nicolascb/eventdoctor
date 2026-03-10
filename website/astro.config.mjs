// @ts-check
import starlight from '@astrojs/starlight';
import { defineConfig } from 'astro/config';

// https://astro.build/config
export default defineConfig({
	site: 'https://nicolascb.github.io',
	base: '/eventdoctor',
	integrations: [
		starlight({
			title: 'EventDoctor',
			description: 'Keep your event-driven architecture documentation alive.',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/nicolascb/eventdoctor' }],
			sidebar: [
				{
					label: 'Getting Started',
					items: [
						{ label: 'Introduction', slug: 'getting-started/introduction' },
						{ label: 'Quick Start', slug: 'getting-started/quick-start' },
						{ label: 'Spec File Reference', slug: 'getting-started/spec-reference' },
					],
				},
				{
					label: 'Core Concepts',
					items: [
						{ label: 'How It Works', slug: 'concepts/how-it-works' },
						{ label: 'Topics & Events', slug: 'concepts/topics-and-events' },
						{ label: 'Producers & Consumers', slug: 'concepts/producers-and-consumers' },
					],
				},
				{
					label: 'Guides',
					items: [
						{ label: 'CLI Usage', slug: 'guides/cli' },
						{ label: 'CI/CD Integration', slug: 'guides/cicd' },
						{ label: 'Web UI', slug: 'guides/web-ui' },
					],
				},
				{
					label: 'Architecture',
					items: [
						{ label: 'Overview', slug: 'architecture/overview' },
						{ label: 'API Reference', slug: 'architecture/api-reference' },
					],
				},
				{
					label: 'Contributing',
					items: [
						{ label: 'How to Contribute', slug: 'contributing/how-to-contribute' },
					],
				},
			],
		}),
	],
});
