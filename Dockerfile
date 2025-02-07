FROM golang:1.23.6

# Set the working directory inside the container
WORKDIR /app

# Install Git for cloning repositories
RUN apt-get update && apt-get install -y git vim docker.io && rm -rf /var/lib/apt/lists/*

# Install Docker Compose
ARG DOCKER_COMPOSE_VERSION=2.24.2
RUN curl -L "https://github.com/docker/compose/releases/download/v${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
RUN chmod +x /usr/local/bin/docker-compose

# Create a non-root user and switch to it
RUN useradd -m devuser
USER devuser

# Expose port (optional, for web servers)
EXPOSE 8080
