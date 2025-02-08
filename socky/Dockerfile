FROM node:20 AS build-env
COPY . /app
WORKDIR /app

RUN npm ci
RUN npm run build

FROM gcr.io/distroless/nodejs20-debian12
COPY --from=build-env /app /app
WORKDIR /app

ENV NODE_ENV=production
CMD ["dist/index.js"]
