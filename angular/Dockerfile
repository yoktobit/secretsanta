# 1. Build
FROM node:lts-alpine AS build
WORKDIR /app

COPY package.json /app/package.json
RUN yarn install
COPY . /app
ENV PRODENV=/app/src/environments/environment.prod.ts
RUN yarn run envsub $PRODENV $PRODENV
RUN yarn build --configuration=production

# 2. Deploy
FROM nginx:latest
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist/secretsantang/. /usr/share/nginx/html
EXPOSE 80 443

HEALTHCHECK --interval=60s --timeout=30s --start-period=1s --retries=3 CMD curl --fail http://localhost/health || exit 1
CMD ["nginx", "-g", "daemon off;"]
