package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProxyHandler_ServeHTTP(t *testing.T) {
	// Configurar um servidor HTTP de teste
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}))
	defer testServer.Close()

	// Criar uma requisição para o proxy
	req, err := http.NewRequest("GET", testServer.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Gravar a resposta do proxy
	rr := httptest.NewRecorder()
	handler := proxyHandler{}

	// Simular a chamada ao proxy
	handler.ServeHTTP(rr, req)

	// Ler o corpo da resposta
	result, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Verificar se a resposta é o que esperamos
	if !strings.Contains(string(result), "Hello, world!") {
		t.Errorf("Unexpected response body: %s", result)
	}
}
