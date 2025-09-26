# Use official Node.js LTS image
FROM node:20-alpine

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json / yarn.lock
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the frontend source
COPY . .

# Expose Vite dev server port
EXPOSE 5173

# Start Vite dev server and bind to 0.0.0.0
CMD ["npm", "run", "dev", "--", "--host"]
