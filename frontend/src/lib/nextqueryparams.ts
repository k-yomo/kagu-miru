import { useMemo } from 'react';
import { useRouter } from 'next/router';

export const useNextQueryParams = (): { [key: string]: string } => {
  const router = useRouter();
  const value = useMemo(() => {
    const queryParamsStr = router.asPath.split('?').slice(1).join('');
    const urlSearchParams = new URLSearchParams(queryParamsStr);
    return Object.fromEntries(urlSearchParams.entries());
  }, [router.asPath]);

  return value;
};
