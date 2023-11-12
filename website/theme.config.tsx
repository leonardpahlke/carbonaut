import type { DocsThemeConfig} from 'nextra-theme-docs';
import { useConfig } from 'nextra-theme-docs'
import { useRouter } from 'next/router'

const logo = (<>
  <span style={{ marginLeft: '.4em', fontWeight: 800 }}>
    Carbonaut
  </span>
</>
)

const config: DocsThemeConfig = {
  project: {
    link: 'https://github.com/carbonaut-cloud/carbonaut'
  },
  docsRepositoryBase: 'https://github.com/carbonaut-cloud/carbonaut-docs/blob/main/',
  useNextSeoProps() {
    const { asPath } = useRouter()
    if (asPath !== '/') {
      return {
        titleTemplate: '%s | Carbonaut'
      }
    }
  },
  logo,
  head: function useHead() {
    const { title } = useConfig()
    const socialCard = 'https://carbonaut.cloud/carbonaut-banner.png'
    return (
      <>
        <meta name="msapplication-TileColor" content="#fff" />
        <meta name="theme-color" content="#fff" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <meta httpEquiv="Content-Language" content="en" />
        <meta
          name="description"
          content="Carbonaut your Cloud Native Energy and Carbon Control Center"
        />
        <meta
          name="og:description"
          content="Carbonaut your Cloud Native Energy and Carbon Control Center"
        />
        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:image" content={socialCard} />
        <meta name="twitter:site:domain" content="carbonaut.cloud" />
        <meta name="twitter:url" content="https://carbonaut.cloud" />
        <meta
          name="og:title"
          content={title ? title + ' - Carbonaut' : 'Carbonaut'}
        />
        <meta name="og:image" content={socialCard} />
        <meta name="apple-mobile-web-app-title" content="Carbonaut" />
        <link rel="icon" href="/logo.svg" type="image/svg+xml" />
        <link
          rel="icon"
          href="/logo-light-no-font.svg"
          type="image/svg+xml"
          media="(prefers-color-scheme: dark)"
        />
        <link
          rel="icon"
          href="/logo-dark-no-font.svg"
          type="image/svg+xml"
          media="(prefers-color-scheme: light)"
        />
      </>
    )
  },
  // banner: {
  //   key: 'v0.X-release',
  //   text: (
  //     <a href="https://carbonaut.cloud/docs" target="_blank" rel="noreferrer">
  //       🎉 Carbonaut v0.X is released. Read more →
  //     </a>
  //   )
  // },
  primaryHue: 104,
  editLink: {
    text: 'Edit this page on GitHub →'
  },
  feedback: {
    content: 'Question? Give us feedback →',
    labels: 'feedback'
  },
  sidebar: {
    titleComponent({ title, type }) {
      if (type === 'separator') {
        return <span className="cursor-default">{title}</span>
      }
      return <>{title}</>
    },
    defaultMenuCollapseLevel: 1,
    toggleButton: true,
  },
  footer: {
    text: (
      <div className="flex w-full flex-col items-center sm:items-start">
        <p className="mt-6 text-xs">
          © {new Date().getFullYear()} The Carbonaut Project.
        </p>
      </div>
    )
  }
}

export default config
