# Etap 1: Budowanie aplikacji 
FROM golang:1.24-alpine AS builder

# Informacje o autorze (standard OCI)
LABEL org.opencontainers.image.authors="Mykhailo Kapustianyk"

# Instalacja UPX do kompresji binarnej
RUN apk add --no-cache upx

# Konfiguracja katalogu roboczego
WORKDIR /build

# Kopiowanie kodu źródłowego
COPY main.go .

# Kompilacja statyczna z optymalizacjami
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o app main.go

# Kompresja pliku wykonywalnego (redukcja ~60-70%)
RUN upx --best --lzma app

# Etap 2: Przygotowanie certyfikatów SSL
FROM alpine:latest AS certs
RUN apk --no-cache add ca-certificates && update-ca-certificates

# Etap 3: Obraz końcowy
FROM scratch

# Informacje o autorze (standard OCI)
LABEL org.opencontainers.image.authors="Mykhailo Kapustianyk"
LABEL org.opencontainers.image.description="Aplikacja Pogodowa z wykorzystaniem API OpenWeatherMap"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.created="2025-05-06"

# Kopiowanie certyfikatów SSL
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Zmienne środowiskowe
ENV PORT=5000

# Kopiowanie skompresowanej binarki i zminifikowanego HTML
COPY --from=builder /build/app /app
COPY index.html /index.html

# Ekspozycja portu
EXPOSE ${PORT}

# Dodanie healthcheck - zmienione z CMD-SHELL na CMD
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app", "-health-check"]

# Uruchomienie aplikacji
ENTRYPOINT ["/app"]