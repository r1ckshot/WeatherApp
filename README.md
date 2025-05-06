# Zadanie 1: Aplikacja Pogodowa w kontenerze Docker

## Autor
Mykhailo Kapustianyk

![image alt](https://github.com/r1ckshot/WeatherApp/blob/3ae0c038be9983b3bd4b19c3b24619f2b179976b/Poprawno%C5%9B%C4%87%20dzia%C5%82ania.png)

## Opis projektu
Projekt zawiera aplikację pogodową napisaną w Go, która umożliwia sprawdzanie aktualnej pogody dla różnych miast w wybranych krajach. Aplikacja korzysta z API OpenWeatherMap do pobierania danych pogodowych.

## Wymagania
- Klucz API OpenWeatherMap (można uzyskać za darmo na stronie [OpenWeatherMap](https://openweathermap.org/api))
- Docker

## Technologie
- Go (backend)
- HTML/CSS/JavaScript (frontend)
- Docker (konteneryzacja)

## Struktura projektu
- `main.go` - kod źródłowy aplikacji
- `index.html` - plik HTML z interfejsem użytkownika
- `Dockerfile` - konfiguracja budowania obrazu Docker

## Optymalizacja obrazu Docker
Obraz Docker został zoptymalizowany pod względem rozmiaru i wydajności:
1. Wykorzystano podejście multi-stage build
2. Skompilowano aplikację Go statycznie z optymalizacjami
3. Użyto narzędzia UPX do kompresji pliku wykonywalnego
4. Zastosowano bazowy obraz `scratch` (najmniejszy możliwy)
5. Dodano HEALTHCHECK dla monitorowania stanu aplikacji

## Wyniki optymalizacji
- Liczba warstw obrazu: 3
- Rozmiar obrazu: 4.44MB

## Instrukcje uruchomienia

### Budowanie obrazu
```bash
docker build -t weather-app:go .
```

### Uruchamianie kontenera
```bash
docker run -d -p 5000:5000 -e WEATHER_API_KEY=your_api_key_here --name weather-container weather-app:go
```

### Sprawdzenie logów
```bash
docker logs weather-container
```

### Sprawdzenie liczby warstw i rozmiaru obrazu
```bash
docker image inspect weather-app:go --format='{{.RootFS.Layers}}' | wc -w
docker image ls weather-app:go
```

### Sprawdzenie statusu healthcheck
```bash
docker inspect weather-container | grep -A 10 "Health"
```

## Dostęp do aplikacji
Po uruchomieniu kontenera aplikacja będzie dostępna pod adresem: http://localhost:5000

## Uwagi
- Należy zastąpić `your_api_key_here` własnym kluczem API OpenWeatherMap
- Aplikacja obsługuje miasta w Polsce, Niemczech, Wielkiej Brytanii i USA
