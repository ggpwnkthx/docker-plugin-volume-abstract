FROM debian:stable-slim
COPY docker-plugin-volume-abstract /bin/docker-plugin-volume-abstract
RUN chmod +x /bin/docker-plugin-volume-abstract