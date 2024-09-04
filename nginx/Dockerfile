FROM nginx:stable-alpine AS base

COPY ./nginx.conf /etc/nginx/nginx.conf

RUN set -eux; \
    rm /etc/nginx/conf.d/default.conf; \
    mkdir -p /var/www/fooder/public; \
    chown -R nginx:nginx /var/www/fooder; \
    touch /var/www/fooder/public/index.php;