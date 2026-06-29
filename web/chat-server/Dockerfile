FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json yarn.lock ./
RUN npm install

COPY . .

ARG VUE_APP_API_BASE_URL
ARG VUE_APP_WS_BASE_URL

ENV VUE_APP_API_BASE_URL=${VUE_APP_API_BASE_URL}
ENV VUE_APP_WS_BASE_URL=${VUE_APP_WS_BASE_URL}

RUN npm run build

FROM nginx:1.27-alpine

COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
