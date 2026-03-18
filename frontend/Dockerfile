# Stage 1: Build
FROM node:20-alpine AS builder

WORKDIR /app

# Accept build arguments for public environment variables
# SvelteKit embeds PUBLIC_* variables at build time
ARG PUBLIC_API_URL
ENV PUBLIC_API_URL=$PUBLIC_API_URL

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source files
COPY . .

# Build the application
RUN npm run build

# Stage 2: Run
FROM nginx:stable-alpine

# Copy the custom nginx configuration
# We'll create this file next to handle SPA routing
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Remove default nginx static assets
RUN rm -rf /usr/share/nginx/html/*

# Copy the built files from the builder stage
COPY --from=builder /app/build /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
