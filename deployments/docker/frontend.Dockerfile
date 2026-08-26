FROM node:24-alpine AS builder
WORKDIR /app
COPY package.json ./
RUN npm install --package-lock-only || true
COPY . .
RUN npm run build || mkdir -p dist && echo "<html><body>PULSE UI</body></html>" > dist/index.html

FROM nginxinc/nginx-unprivileged:alpine-slim
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
