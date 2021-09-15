import type { NextPage } from 'next';
import Head from 'next/head';
import SEOMeta from '@src/components/SEOMeta';

const Home: NextPage = () => {
  return (
    <div>
      <SEOMeta
        title="カグミル - 家具検索サービス"
        excludeSiteTitle
        description="カグミルはオンラインで買える家具を横断で検索出来るサービスです。"
        // img={{ srcPath: TopImg.src }}
      />
    </div>
  );
};

export default Home;
