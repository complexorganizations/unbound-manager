FROM debian:latest
LABEL maintainer="Prajwal Koirala <prajwalkoirala23@protonmail.com>"
EXPOSE 53/tcp
RUN apt-get update && \
    apt-get install curl -y && \
    curl https://raw.githubusercontent.com/complexorganizations/unbound-manager/main/unbound-manager.sh --create-dirs -o /usr/local/bin/unbound-manager.sh && \
    chmod +x /usr/local/bin/unbound-manager.sh
