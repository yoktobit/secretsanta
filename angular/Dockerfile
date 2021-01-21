# 1. Build
FROM node:lts AS build
WORKDIR /app

COPY package.json /app/package.json
RUN npm install
COPY . /app
RUN npm run build --prod

# 2. Deploy
FROM nginx:latest
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist/secretsantang/. /usr/share/nginx/html
EXPOSE 80 443

HEALTHCHECK --interval=60s --timeout=30s --start-period=1s --retries=3 CMD curl --fail http://localhost/health || exit 1
CMD ["nginx", "-g", "daemon off;"]
