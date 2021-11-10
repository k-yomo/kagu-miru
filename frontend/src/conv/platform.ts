import { ItemSellingPlatform } from '@src/generated/graphql';

export function platFormText(
  platform: ItemSellingPlatform,
  shortName?: boolean
) {
  switch (platform) {
    case ItemSellingPlatform.Rakuten:
      return '楽天';
    case ItemSellingPlatform.YahooShopping:
      return shortName ? 'Yahoo' : 'Yahooショッピング';
  }
}

export function platFormColor(platform: ItemSellingPlatform) {
  switch (platform) {
    case ItemSellingPlatform.Rakuten:
      return 'text-rakuten';
    case ItemSellingPlatform.YahooShopping:
      return 'text-yahoo-shopping';
  }
}
