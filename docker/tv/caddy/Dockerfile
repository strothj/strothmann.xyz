FROM alpine:3.1

RUN mkdir /etc/caddy && mkdir -p /opt/caddy && mkdir /var/www
ADD Caddyfile /etc/caddy
ADD caddy /opt/caddy
ADD www /var/www
RUN chmod 0755 /etc/caddy && \
  chmod 644 /etc/caddy/Caddyfile && \
  chmod 755 /opt/caddy/caddy
CMD ["/opt/caddy/caddy", "-conf=/etc/caddy/Caddyfile"]
EXPOSE 80
