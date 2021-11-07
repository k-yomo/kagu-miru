
export default function resolveProductionUrl(document) {
    return `${process.env.SANITY_STUDIO_KAGU_MIRU_URL}/media/posts/${document.slug.current}/preview`
}
