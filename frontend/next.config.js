/** @type {import('next').NextConfig} */
module.exports = {
  reactStrictMode: true,
  images: {
    domains: [
        'thumbnail.image.rakuten.co.jp',
        'via.placeholder.com'
    ],
  },
async headers() {
    return [
        {
            source: '/(.*)',
            headers: securityHeaders,
        },
    ]
},
}
