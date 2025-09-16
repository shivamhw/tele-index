FROM oven/bun:latest AS js_builder

WORKDIR /build
COPY package.json bun.lockb tsconfig.json ./
COPY src/ ./src
COPY public/ ./public
COPY .env .env
RUN bun install
RUN bun run build
RUN ls

FROM golang:1.24-alpine AS go_builder

WORKDIR /build
COPY backend/ .
RUN go build -o tele-index
RUN ls


FROM alpine:latest
WORKDIR /app
COPY --from=js_builder /build/build ./build
COPY --from=go_builder /build/tele-index ./tele-index

ENTRYPOINT [ "./tele-index", "-static", "build" ]
# CMD ["./tele-index", "-static", "/app/build"]