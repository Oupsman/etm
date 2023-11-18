# Start from the latest Golang base image
FROM golang:1.21.4-alpine AS base

# Add Maintainer Info
LABEL maintainer="Oupsman <oupsman@oupsman.fr>"


FROM base AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

FROM builder AS runner

ENV PORT=8080

# Expose port 8080 to the outside
EXPOSE ${PORT}

# Healthcheck
HEALTHCHECK --interval=10s --timeout=3s --start-period=20s \
  CMD wget --no-verbose --tries=1 --no-check-certificate http://localhost:"$PORT"/api/v1/healthcheck || exit 1

# Command to run the executable
CMD ["/app/main"]