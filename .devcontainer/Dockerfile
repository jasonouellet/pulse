FROM mcr.microsoft.com/devcontainers/go:1.25-bookworm

# Install Node.js & npm for React/Vite frontend development
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# Install sqlc for Go SQL code generation
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install golang-migrate CLI
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/migrate

# Install Air for Go live reloading
RUN go install github.com/air-verse/air@latest
