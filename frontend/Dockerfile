FROM node:18 AS builder

WORKDIR /frontend

# Install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy source code and build the app
COPY . .
RUN npm run build

# Use a minimal image for the final container
FROM nginx:alpine

WORKDIR /usr/share/nginx/html
COPY --from=builder /frontend/dist .

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]