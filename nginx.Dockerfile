FROM nginx:1.25.3-alpine as base

COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./cert.pem /ssl/cert.pem
COPY ./key.pem /ssl/key.pem