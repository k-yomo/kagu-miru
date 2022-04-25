import '../styles/globals.css';
import type { AppProps } from 'next/app';
import { ThemeProvider } from 'next-themes';
import NextNProgress from 'nextjs-progressbar';
import { ApolloProvider } from '@apollo/client';
import { usePageView } from '@src/lib/googleAnalytics';
import apolloClient from '@src/lib/apolloClient';
import Header from '@src/components/Header';
import Footer from '@src/components/Footer';
import { ToastProvider } from '@src/contexts/toast';

interface WebVitalsMetric {
  id: string;
  name: string;
  startTime: number;
  value: number;
  label: 'web-vital' | 'custom';
}

export function reportWebVitals({ id, name, label, value }: WebVitalsMetric) {
  // 以下のURLの例のように初期化した場合は、`window.gtag` を使用します:
  // https://github.com/vercel/next.js/blob/canary/examples/with-google-analytics/pages/_app.js
  window.gtag('event', name, {
    event_category:
      label === 'web-vital' ? 'Web Vitals' : 'Next.js custom metric',
    value: Math.round(name === 'CLS' ? value * 1000 : value), // 値は整数にする必要があります
    event_label: id, // 現在のページのロードした一意なID
    non_interaction: true, // バウンス率への影響を回避
  })
}

function MyApp({ Component, pageProps }: AppProps) {
  usePageView();

  return (
    // @ts-ignore
    <ThemeProvider attribute="class" forcedTheme={Component.theme}>
      <ApolloProvider client={apolloClient}>
        <div className="flex flex-col min-h-screen">
          <Header />
          <main className="z-0 grow relative bg-white dark:bg-black">
            <NextNProgress
              color="#06b6d4"
              height={3}
              showOnShallow={false}
              options={{ parent: 'main', showSpinner: false }}
            />
            <ToastProvider>
              <Component {...pageProps} />
            </ToastProvider>
          </main>
          <Footer />
        </div>
      </ApolloProvider>
    </ThemeProvider>
  );
}

export default MyApp;
