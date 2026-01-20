# Step 1: Assembly
FROM golang:alpine3.22 AS builder

WORKDIR /app

# Copy the remaining files
COPY .. .

# Building binaries
RUN mkdir -p /app/build && go build -o ./build/k8s_worker  ./main.go

# Stage 2: The final image
FROM alpine:3.22

WORKDIR /app

# Copy the binaries from the previous step
COPY --from=builder /app/build/. .

# Making binaries executable
RUN chmod +x ./k8s_worker

CMD ["./k8s_worker"]