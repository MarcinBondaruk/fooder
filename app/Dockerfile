FROM php:8.3-fpm-bullseye AS base

RUN set -eux; \
    mkdir -p /var/www/fooder; \
    chown www-data:www-data /var/www/fooder; \
    chmod 1777 /var/www/fooder;

COPY --chown=www-data:www-data ./public/ /var/www/fooder/public
COPY --chown=www-data:www-data ./src/ /var/www/fooder/src
# todo: build this one in separate step
COPY --chown=www-data:www-data ./vendor/ /var/www/fooder/vendor

WORKDIR /var/www/fooder
