FROM nginx:1.19-alpine

COPY docker/nginx.conf /etc/nginx/templates/default.conf.template

EXPOSE 80