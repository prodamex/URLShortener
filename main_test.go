package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleForm(t *testing.T) {
	// Cas de test pour une requête GET
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	http.HandlerFunc(handleForm).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("HandleForm a retourné un code de statut incorrect: got %v, attendu %v",
			status, http.StatusOK)
	}

	expected := "URL Shortener"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("HandleForm n'a pas retourné le corps attendu. Attendu : %s, Obtenu : %s", expected, rr.Body.String())
	}

	// Cas de test pour une requête POST
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()

	http.HandlerFunc(handleForm).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("HandleForm a retourné un code de statut incorrect: got %v, attendu %v",
			status, http.StatusSeeOther)
	}
}


func TestHandleRedirect(t *testing.T) {
	req, err := http.NewRequest("GET", "/short/ArTeql", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Crée un ResponseRecorder pour enregistrer la réponse.
	rr := httptest.NewRecorder()

	// Appelle la fonction de gestion pour tester.
	handleRedirect(rr, req)

	// Vérifie le code de statut de la réponse.
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("HandleRedirect a retourné un code de statut incorrect: got %v, attendu %v",
			status, http.StatusMovedPermanently)
	}

	// Vérifie que la redirection est effectuée vers l'URL d'origine
	expectedURL := "https://www.youtube.com/"
	if location := rr.Header().Get("Location"); location != expectedURL {
		t.Errorf("HandleRedirect n'a pas effectué la redirection correctement. Attendu : %s, Obtenu : %s", expectedURL, location)
	}
}
