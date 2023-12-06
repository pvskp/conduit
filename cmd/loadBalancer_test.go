package cmd

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRRLoadBalancer_ServeHTTP(t *testing.T) {
	// Criar hosts fictícios para o balanceador de carga
	hosts := []url.URL{
		*urlMustParse("http://localhost:3000"),
		*urlMustParse("http://localhost:3001"),
	}

	rr := NewRRLoadBalancer(hosts)

	// Criar uma requisição HTTP de teste
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	for i := range hosts {
		w := httptest.NewRecorder()
		rr.ServeHTTP(w, req)

		// Ajustar o índice para verificar o host anterior
		prevIndex := (rr.Index + len(rr.Hosts) - 1) % len(rr.Hosts)

		// Verificar se a requisição foi encaminhada para o host correto
		if rr.Hosts[prevIndex].String() != hosts[i].String() {
			t.Errorf("Requisição %d encaminhada para o host errado: got %v, want %v", i, rr.Hosts[prevIndex], hosts[i])
		}
	}
}

// urlMustParse é um helper para criar URLs sem ter que lidar com erros no teste
func urlMustParse(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return u
}
