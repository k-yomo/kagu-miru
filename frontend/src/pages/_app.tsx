import '../styles/globals.css';
import type { AppProps } from 'next/app';
import { ThemeProvider } from 'next-themes';
import { ApolloProvider } from '@apollo/client';
import { usePageView } from '@src/lib/googleAnalytics';
import apolloClient from '@src/lib/apolloClient';
import Header from '@src/components/Header';
import Footer from '@src/components/Footer';

function MyApp({ Component, pageProps }: AppProps) {
  usePageView();

  return (
    <ThemeProvider attribute="class">
      <ApolloProvider client={apolloClient}>
        <div className="flex flex-col min-h-screen">
          <Header />
          <main className="z-0 flex-grow relative bg-white dark:bg-black">
            <Component {...pageProps} />
          </main>
          <Footer />
        </div>
      </ApolloProvider>
    </ThemeProvider>
  );
}

export default MyApp;
