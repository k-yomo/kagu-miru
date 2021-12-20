import React from 'react';
import { useRouter } from 'next/router';
import SEOMeta from '@src/components/SEOMeta';

export default function ContactPage() {
  const router = useRouter();
  return (
    <>
      <SEOMeta
        title="お問い合わせ"
        description="カグミルのお問い合わせページです。"
        path={router.asPath}
      />
      <div className="mx-auto max-w-[1000px] rounded-lg">
        <div className="my-8">
          <h1 className="my-6 sm:my-8 text-3xl text-center font-bold">
            お問い合わせフォーム
          </h1>
          <div className="mx-auto bg-white">
            <iframe
              src="https://docs.google.com/forms/d/e/1FAIpQLSfmOF55J2f10N3F9w8KcT4RDlb8utquCqDUZH35Wr5vfWB9FA/viewform?embedded=true"
              width="100%"
              height="900"
              frameBorder="0"
              marginHeight={0}
              marginWidth={0}
            >
              読み込んでいます…
            </iframe>
          </div>
        </div>
      </div>
    </>
  );
}
