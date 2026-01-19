# Step 1: Assembly
FROM golang:alpine3.20 AS builder

WORKDIR /app

# Copy only go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the remaining files
COPY .. .

# Building binaries
RUN mkdir -p /app/build && go build -o ./build/k8s_worker

# Stage 2: The final image
FROM alpine:3.21

WORKDIR /app

# Copy the binaries from the previous step
COPY --from=builder /app/build/. .

# Making binaries executable
RUN chmod +x ./k8s_worker

CMD ["./k8s_worker"]