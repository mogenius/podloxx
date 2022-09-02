# Common build stage
FROM node:16-alpine3.14 as common-build-stage

COPY . ./app

WORKDIR /app

RUN npm install

RUN npm run build

EXPOSE 3000

FROM common-build-stage as production-build-stage

ENV NODE_ENV production
ENV PORT 3000

RUN chown -R 1000:1000 /app/dist
USER 1000

CMD node dist/server.js
