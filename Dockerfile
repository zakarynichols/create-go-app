FROM golang:1.23.6

# Set the working directory inside the container
WORKDIR /app

# Install Git for cloning repositories
RUN apt-get update && apt-get install -y git vim && rm -rf /var/lib/apt/lists/*

# Create a non-root user and switch to it
RUN useradd -m devuser
USER devuser

# Expose port (optional, for web servers)
EXPOSE 8080
