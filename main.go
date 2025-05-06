package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	authorName = "Mykhailo Kapustianyk"
)

// Struktura dla danych pogodowych
type WeatherResponse struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

// Struktura wynikowa dla odpowiedzi API
type WeatherResult struct {
	City        string  `json:"city"`
	Country     string  `json:"country"`
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feels_like"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
	Icon        string  `json:"icon"`
}

// Mapa krajów i miast
var countries = map[string][]string{
	"Poland":  {"Warsaw", "Krakow", "Gdansk", "Lublin", "Wroclaw"},
	"Germany": {"Berlin", "Munich", "Hamburg", "Frankfurt"},
	"UK":      {"London", "Manchester", "Liverpool", "Edinburgh"},
	"USA":     {"New York", "Los Angeles", "Chicago", "Miami"},
}

func main() {
	// Dodanie flagi healthcheck
	healthCheck := flag.Bool("health-check", false, "Wykonaj test zdrowia aplikacji")
	flag.Parse()

	// Jeśli flaga health-check jest ustawiona, wykonanie testu zdrowia
	if *healthCheck {
		// Prosty test zdrowia - sprawdza czy aplikacja może się uruchomić
		fmt.Println("Health check passed!")
		os.Exit(0)
	}

	// Pobranie portu z zmiennej środowiskowej lub użycie domyślnego
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Zapisanie informacji o uruchomieniu
	startTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] Aplikacja uruchomiona przez: %s\n", startTime, authorName)
	fmt.Printf("[%s] Nasłuchiwanie na porcie TCP: %s\n", startTime, port)

	// Dodanie endpointu healthcheck
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Główna strona
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html, err := os.ReadFile("index.html")
		if err != nil {
			http.Error(w, "Nie można wczytać strony", http.StatusInternalServerError)
			return
		}

		// Przekształcenie mapy krajów na format JSON
		countriesJSON, err := json.Marshal(countries)
		if err != nil {
			http.Error(w, "Błąd przetwarzania danych", http.StatusInternalServerError)
			return
		}

		// Zamiana placeholderu w HTML na dane JSON
		htmlStr := string(html)
		htmlStr = strings.Replace(htmlStr, "{{ countries|tojson|safe }}", string(countriesJSON), 1)

		// Ustaw odpowiednie nagłówki
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlStr))
	})

	// Endpoint API pogodowego
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Tylko metoda POST jest dozwolona", http.StatusMethodNotAllowed)
			return
		}

		// Parsowanie danych formularza
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Nie można przetworzyć formularza", http.StatusBadRequest)
			return
		}

		country := r.FormValue("country")
		city := r.FormValue("city")

		// Klucz API z zmiennej środowiskowej
		apiKey := os.Getenv("WEATHER_API_KEY")
		if apiKey == "" {
			apiKey = "your_api_key_here" // domyślna wartość
		}

		// Zakodowanie nazwy miasta w URL
		cityQuery := url.QueryEscape(city)

		// Wywołanie API OpenWeatherMap
		apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", cityQuery, apiKey)

		resp, err := http.Get(apiUrl)
		if err != nil {
			http.Error(w, "Błąd podczas pobierania danych pogodowych", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			http.Error(w, "Nie udało się pobrać danych pogodowych", http.StatusBadRequest)
			return
		}

		// Odczytanie odpowiedzi z API
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Błąd odczytu odpowiedzi", http.StatusInternalServerError)
			return
		}

		// Odkodowanie danych JSON
		var weatherData WeatherResponse
		err = json.Unmarshal(body, &weatherData)
		if err != nil {
			http.Error(w, "Błąd dekodowania JSON", http.StatusInternalServerError)
			return
		}

		// Utworzenie rezultatu
		result := WeatherResult{
			City:        city,
			Country:     country,
			Temperature: weatherData.Main.Temp,
			FeelsLike:   weatherData.Main.FeelsLike,
			Humidity:    weatherData.Main.Humidity,
			WindSpeed:   weatherData.Wind.Speed,
		}

		// Dane o ikonie i opisie, jeśli są dostępne
		if len(weatherData.Weather) > 0 {
			result.Description = weatherData.Weather[0].Description
			result.Icon = weatherData.Weather[0].Icon
		}

		// Odpowiedz z danymi JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// Uruchom serwer
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
