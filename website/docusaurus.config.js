module.exports = {
  title: 'ARGO web-api Documentation',
  tagline: 'Learn how the argo-web-api works',
  url: 'https://argoeu.github.io',
  baseUrl: '/argo-web-api/',
  onBrokenLinks: 'throw',
  favicon: 'img/favicon.ico',
  organizationName: 'ARGOeu', // Usually your GitHub org/user name.
  projectName: 'argo-web-api', // Usually your repo name.
  themeConfig: {
    navbar: {
      title: 'ARGO WEB API',
      logo: {
        alt: 'argo-web-api logo',
        src: 'img/argo-web-api.png',
      },
      items: [
        {
          to: 'docs/',
          activeBasePath: 'docs',
          label: 'Docs',
          position: 'left',
        },
        {
          href: 'https://github.com/ARGOeu/argo-web-api',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [{
            to: 'docs/',
            activeBasePath: 'docs',
            label: 'Explore Documentation',
            position: 'left',
          },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'Github',
              href: 'https://github.com/ARGOeu/argo-web-api',
            }
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/ARGOeu/argo-web-api',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} GRNET`,
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};
