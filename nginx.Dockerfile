FROM nginx:latest as base

COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./cert.pem /ssl/cert.pem
COPY ./key.pem /ssl/key.pem