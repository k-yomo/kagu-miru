import AffiliateLinkPreview from "../components/AffiliateLinkPreview"

export default {
  name: 'affiliateLink',
  title: 'Item',
  type: 'object',
  fields: [
    {
      name: 'html',
      title: 'HTML',
      type: 'text',
    },
  ],
  preview: {
    select: {
      html: 'html'
    },
    component: AffiliateLinkPreview
  }
}

