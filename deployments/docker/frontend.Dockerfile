FROM node:24-alpine AS builder
WORKDIR /app
ARG VITE_OTEL_COLLECTOR_URL=/otlp-http/v1/traces
ENV VITE_OTEL_COLLECTOR_URL=$VITE_OTEL_COLLECTOR_URL
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginxinc/nginx-unprivileged:alpine-slim
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
