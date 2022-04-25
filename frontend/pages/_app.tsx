import type { AppProps } from 'next/app'

import globalStyles from '../styles/globals.scss'

function MyApp({ Component, pageProps }: AppProps) {
  return <Component className={globalStyles.body && globalStyles.html} {...pageProps} />
}

export default MyApp
