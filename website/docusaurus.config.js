// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

const copyrightMessage = `
<p xmlns:cc="http://creativecommons.org/ns#">
   This work by
   <a rel="cc:attributionURL dct:creator" property="cc:attributionName" href="https://www.grnet.gr">
      National Infrastructures for Research and Technology - GRNET S.A.
   </a> and
   <a rel="cc:attributionURL dct:creator" property="cc:attributionName" href="https://www.srce.hr">
      University of Zagreb University Computing Centre (SRCE)
   </a> is licensed under
   <a href="http://creativecommons.org/licenses/by/4.0/?ref=chooser-v1" target="_blank"
      rel="license noopener noreferrer" style="display:inline-block;">
      CC BY 4.0
      <img style="height:22px!important;margin-left:3px;vertical-align:text-bottom;"
         src="https://mirrors.creativecommons.org/presskit/icons/cc.svg?ref=chooser-v1">
      <img style="height:22px!important;margin-left:3px;vertical-align:text-bottom;"
         src="https://mirrors.creativecommons.org/presskit/icons/by.svg?ref=chooser-v1">
   </a>
</p>
`

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'ARGO web-api Documentation',
  tagline: 'Learn how the argo-web-api works',
  url: 'https://argoeu.github.io',
  baseUrl: '/argo-web-api/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',



  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'ARGOeu', // Usually your GitHub org/user name.
  projectName: 'argo-web-api', // Usually your repo name.

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        blog: false,
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          // editUrl:
          //   'https://github.com/facebook/docusaurus/tree/main/packages/create-docusaurus/templates/shared/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        title: 'ARGO WEB API',
        logo: {
          alt: 'argo-web-api logo',
          src: 'img/logo.svg',
        },
        items: [
          {
            to: 'docs/',
            activeBasePath: 'docs',
            label: 'Docs',
            position: 'left',
          },
          {
            href: 'pathname:///openapi/explore',
            label: 'Explore the API',
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
        style: 'light',
        links: [
          {
            title: 'Docs',
            items: [{
              to: 'docs/',
              activeBasePath: 'docs',
              label: 'Explore Documentation',
              position: 'left',
            },
            {
              to: 'pathname:///openapi/explore',
              label: 'Explore the API (openapi v3)',
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
        copyright: copyrightMessage,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },

    }),
  themes: [
    [
      "@easyops-cn/docusaurus-search-local",
      {
        hashed: true,
        language: ["en", "zh"],
        highlightSearchTermsOnTargetPage: true,
        explicitSearchResultPath: true,
        indexBlog: false,
      },
    ],
  ],
};

module.exports = config;
