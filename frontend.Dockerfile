FROM node:14.17.0-slim AS build
ENV NODE_ENV=production

WORKDIR /build

COPY frontend/package.json frontend/package-lock.json ./
RUN npm install --production=false

COPY ./frontend .
RUN npm run build

FROM gcr.io/distroless/nodejs:14
ENV NODE_ENV=production
WORKDIR /app

COPY frontend/next.config.js ./
COPY --from=build /build/public ./public
COPY --from=build /build/.next ./.next
COPY --from=build /build/node_modules ./node_modules

CMD ["node_modules/.bin/next", "start"]
