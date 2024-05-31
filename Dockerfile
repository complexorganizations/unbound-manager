FROM debian:buster
LABEL maintainer="Prajwal Koirala <prajwalkoirala23@protonmail.com>"

# Install necessary packages
RUN apt-get update && \
    apt-get install -y curl && \
    rm -rf /var/lib/apt/lists/*

# Download and set up unbound-manager
RUN curl -sSL https://raw.githubusercontent.com/complexorganizations/unbound-manager/main/unbound-manager.sh -o /usr/local/bin/unbound-manager.sh && \
    chmod +x /usr/local/bin/unbound-manager.sh

# Expose port 53
EXPOSE 53/tcp
