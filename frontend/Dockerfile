# Build Stage
FROM node:18-alpine AS builder
WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm install

# Copy application files and build
COPY . .
RUN npm run build

# Serve Stage
FROM nginx:alpine
# Copy built static files to Nginx public folder
COPY --from=builder /app/dist /usr/share/nginx/html
# Copy custom Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
