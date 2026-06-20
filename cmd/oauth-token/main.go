// Programa de un solo uso para obtener el GOOGLE_OAUTH_REFRESH_TOKEN.
//
// Uso:
//  1. Configura GOOGLE_OAUTH_CLIENT_ID y GOOGLE_OAUTH_CLIENT_SECRET en el .env
//  2. Ejecuta: go run ./cmd/oauth-token
//  3. Se abrirá tu navegador, inicia sesión con la cuenta de Google donde
//     quieres guardar las imágenes y acepta los permisos.
//  4. Copia el refresh token que aparece en consola al .env como
//     GOOGLE_OAUTH_REFRESH_TOKEN.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const redirectURL = "http://localhost:8090/callback"

func main() {
	_ = godotenv.Load()

	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatal("Configura GOOGLE_OAUTH_CLIENT_ID y GOOGLE_OAUTH_CLIENT_SECRET en el .env primero")
	}

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURL,
		Scopes:       []string{drive.DriveFileScope},
	}

	// AccessTypeOffline + prompt=consent garantiza que Google devuelva refresh token.
	authURL := conf.AuthCodeURL("state-token",
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)

	codeCh := make(chan string)
	srv := &http.Server{Addr: ":8090"}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "no se recibió el código", http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "¡Listo! Ya puedes cerrar esta pestaña y volver a la terminal.")
		codeCh <- code
	})

	go func() {
		_ = srv.ListenAndServe()
	}()

	fmt.Println("Abriendo el navegador para autorizar...")
	fmt.Println("Si no se abre, visita manualmente esta URL:")
	fmt.Println(authURL)
	openBrowser(authURL)

	code := <-codeCh
	_ = srv.Shutdown(context.Background())

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		log.Fatal("Error intercambiando el código:", err)
	}

	fmt.Println("\n========================================")
	fmt.Println("Copia esto a tu .env:")
	fmt.Printf("\nGOOGLE_OAUTH_REFRESH_TOKEN=%s\n", token.RefreshToken)
	fmt.Println("\n========================================")
}

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	_ = exec.Command(cmd, args...).Start()
}
