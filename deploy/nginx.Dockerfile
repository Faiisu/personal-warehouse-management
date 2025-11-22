# Build the React frontend
FROM node:20-alpine AS builder
WORKDIR /app

COPY frontend/package*.json ./
RUN npm ci

COPY frontend ./

ARG BACKEND_IP=localhost:3000
ENV VITE_BACKEND_IP=${BACKEND_IP}
ENV NODE_ENV=production

RUN npm run build

# Serve with Nginx
FROM nginx:1.27-alpine

COPY deploy/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
