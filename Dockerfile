# build stage
FROM node:lts-alpine as build-stage

WORKDIR /app/ui
COPY ui/package*.json ./
RUN npm install
COPY ui/. .
RUN npm run build

WORKDIR /app/backend
COPY backend/package*.json ./
RUN npm install
COPY backend/. .
RUN npm run build --production

# production stage
FROM nginx as production-stage
RUN apt-get update
RUN apt-get install -y nodejs
COPY --from=build-stage /app/ui/dist/podlox /usr/share/nginx/html
COPY --from=build-stage /app/backend /app
COPY start.sh /app/
EXPOSE 8080
EXPOSE 4200
ENV NODE_ENV=production
ENV PORT=4200
ENV NGINX_PORT=8080

ENTRYPOINT ["/app/start.sh"]