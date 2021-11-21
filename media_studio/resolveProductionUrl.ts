import { KAGU_MIRU_URL } from "./config/env"

export default function resolveProductionUrl(document) {
    return `${KAGU_MIRU_URL}/media/posts/${document.slug.current}/preview`
}
