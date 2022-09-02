#!/bin/sh
node --version
ls -lisa /app
nginx -g 'daemon off;' & 
node /app/dist/server.js