import { ItemSellingPlatform } from '@src/generated/graphql';

export const allPlatforms = [
  ItemSellingPlatform.Rakuten,
  ItemSellingPlatform.YahooShopping,
  ItemSellingPlatform.PaypayMall,
];

export function platFormText(
  platform: ItemSellingPlatform,
  shortName?: boolean
) {
  switch (platform) {
    case ItemSellingPlatform.Rakuten:
      return '楽天';
    case ItemSellingPlatform.YahooShopping:
      return shortName ? 'Yahoo' : 'Yahooショッピング';
    case ItemSellingPlatform.PaypayMall:
      return shortName ? 'PayPay' : 'PayPayモール';
  }
}

export function platFormColor(platform: ItemSellingPlatform) {
  switch (platform) {
    case ItemSellingPlatform.Rakuten:
      return 'text-rakuten';
    case ItemSellingPlatform.YahooShopping:
      return 'text-yahoo-shopping';
    case ItemSellingPlatform.PaypayMall:
      return 'text-paypay-mall';
  }
}
