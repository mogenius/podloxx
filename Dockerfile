# build stage
FROM node:lts-alpine as build-stage-ui
WORKDIR /app/ui
COPY ui/package*.json ./
RUN cd ui; npm install
COPY ui/. .
RUN cd ui; npm run build

# production stage
FROM nginxinc/nginx-unprivileged:stable-alpine as production-stage
COPY --from=build-stage-ui /app/dist/angular-app /usr/share/nginx/html
EXPOSE 8080
USER 101
CMD ["nginx", "-g", "daemon off;"]