package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	urlsCollection *mongo.Collection
	stats           Statistics
)

type Statistics struct {
	TotalShortenedLinks int
	ClicksPerShortLink  map[string]int
}

type URLMapping struct {
	ShortKey    string `bson:"shortKey"`
	OriginalURL string `bson:"originalURL"`
	Clicks      int    `bson:"clicks"`
}

func getTotalShortenedLinks() (int64, error) {
    count, err := urlsCollection.CountDocuments(context.Background(), bson.M{})
    if err != nil {
        return 0, err
    }
    return count, nil
}
func incrementTotalShortenedLinks() error {
    count, err := getTotalShortenedLinks()
    if err != nil {
        return err
    }
    stats.TotalShortenedLinks = int(count)
    return nil
}

// Fonction pour incrémenter le nombre de clics pour une clé de raccourcissement donnée
func incrementClicks(shortKey string) {
	stats.ClicksPerShortLink[shortKey]++
	// Met à jour le nombre de clics dans la base de données
	err := updateClicksInDatabase(shortKey, stats.ClicksPerShortLink[shortKey])
	if err != nil {
		log.Println("Error updating clicks in database:", err)
	}
}


// Fonction pour mettre à jour le nombre de clics dans la base de données
func updateClicksInDatabase(shortKey string, clickCount int) error {
	// Définit le filtre pour identifier le document correspondant à la clé de raccourcissement
	filter := bson.M{"shortKey": shortKey}

	// Définit la mise à jour pour incrémenter le nombre de clics
	update := bson.M{"$set": bson.M{"clicks": clickCount}}

	// Effectue la mise à jour dans la base de données
	_, err := urlsCollection.UpdateOne(context.Background(), filter, update)

	return err
}

func getStatistics() Statistics {
	return stats
}

func main() {
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/short/", handleRedirect)

	fmt.Println("URL Shortener is running on :3030")
	http.ListenAndServe(":3030", nil)
}

func init() {
	stats = Statistics{
		ClicksPerShortLink: make(map[string]int),
	}

    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal("Error connecting to MongoDB:", err)
    }

    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal("Error pinging MongoDB:", err)
    }

    urlsCollection = client.Database("urlshortener").Collection("urls")
    if urlsCollection == nil {
        log.Fatal("Error initializing MongoDB collection")
    }
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	  // Incrémente le nombre total de liens raccourcis
	  err := incrementTotalShortenedLinks()
	  if err != nil {
		  log.Println("Error incrementing total shortened links:", err)
	  }
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)
		return
	}
	currentStats := getStatistics()
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortener</title>
		</head>
		<body>
		<h2>URL Shortener</h2>
		<p>Number of Shortened Links: `, currentStats.TotalShortenedLinks, `</p>
		<form method="post" action="/shorten">
			<input type="url" name="url" placeholder="Enter a URL" required>
			<input type="submit" value="Shorten">
		</form>
		</body>
		</html>
	`)
}


// Fonction pour récupérer la clé de raccourcissement et le nombre de clics à partir de l'URL d'origine
func findURLMappingByOriginalURL(originalURL string) (URLMapping, bool) {
	var result URLMapping
	err := urlsCollection.FindOne(context.Background(), bson.M{"originalURL": originalURL}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		// Aucun document trouvé
		return URLMapping{}, false
	} else if err != nil {
		log.Println("Error finding URL mapping:", err)
		return URLMapping{}, false
	}
	return result, true
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	var shortenedURL string
	var clickCount int
	// Vérifie si l'URL d'origine a déjà été raccourcie
	urlMapping, found := findURLMappingByOriginalURL(originalURL)

	if found {
		shortKey := urlMapping.ShortKey
		incrementClicks(shortKey)
		clickCount = urlMapping.Clicks
		shortenedURL = fmt.Sprintf("http://localhost:3030/short/%s", shortKey)
	} else {
		shortKey := generateShortKey()
		err := storeURLMapping(shortKey, originalURL)
		if err != nil {
			http.Error(w, "Error storing URL mapping", http.StatusInternalServerError)
			return
		}
		incrementTotalShortenedLinks()
		clickCount = urlMapping.Clicks
		shortenedURL = fmt.Sprintf("http://localhost:3030/short/%s", shortKey)
	}

	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortener</title>
		</head>
		<body>
			<h2>URL Shortener</h2>
			<p>Original URL: %s</p>
			<p>Shortened URL: <a href="%s">%s</a></p>
			<p>Clicks: %d</p>
		</body>
		</html>
	`, originalURL, shortenedURL, shortenedURL, clickCount)
}

func storeURLMapping(shortKey, originalURL string) error {
	urlMapping := URLMapping{
		ShortKey:    shortKey,
		OriginalURL: originalURL,
		Clicks:      0, 
	}

	_, err := urlsCollection.InsertOne(context.Background(), urlMapping)
	return err
}

func retrieveOriginalURL(shortKey string) (string, error) {
	// Récupère l'URL d'origine de la collection MongoDB
	var result struct {
		OriginalURL string `bson:"originalURL"`
	}
	err := urlsCollection.FindOne(context.Background(), bson.M{"shortKey": shortKey}).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.OriginalURL, nil
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := strings.TrimPrefix(r.URL.Path, "/short/")
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}
	incrementClicks(shortKey)

	// Récupère l'URL d'origine de la base de données MongoDB
	originalURL, err := retrieveOriginalURL(shortKey)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Shortened key not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving original URL", http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

// Fonction pour générer une clé de raccourcissement unique
func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

