package database

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // registra el decodificador PNG
	"io"
	"os"
	"strings"

	"golang.org/x/image/draw"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	// maxImageWidth es el ancho máximo en píxeles; las imágenes más anchas se
	// redimensionan manteniendo la proporción para reducir el peso. 1024px es
	// suficiente para verse bien y hace que la subida a Drive sea más rápida.
	maxImageWidth = 1024
	// jpegQuality es la calidad de compresión JPEG (1-100). 70 baja bastante el
	// peso del archivo sin pérdida visible para fotos de vehículos.
	jpegQuality = 70
)

type GoogleDriveService struct {
	service  *drive.Service
	folderID string
}

// NewGoogleDriveService crea el cliente de Drive autenticado con OAuth2
// usando la cuenta personal del usuario (no Service Account), para que los
// archivos usen la cuota gratuita de 15 GB de la cuenta de Google.
func NewGoogleDriveService() (*GoogleDriveService, error) {
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	refreshToken := os.Getenv("GOOGLE_OAUTH_REFRESH_TOKEN")
	folderID := os.Getenv("GOOGLE_DRIVE_FOLDER_ID")

	if clientID == "" || clientSecret == "" || refreshToken == "" {
		return nil, fmt.Errorf("faltan variables OAuth: GOOGLE_OAUTH_CLIENT_ID, GOOGLE_OAUTH_CLIENT_SECRET y/o GOOGLE_OAUTH_REFRESH_TOKEN")
	}
	if folderID == "" {
		return nil, fmt.Errorf("GOOGLE_DRIVE_FOLDER_ID no está configurado")
	}

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveFileScope},
	}

	ctx := context.Background()
	token := &oauth2.Token{RefreshToken: refreshToken}
	client := conf.Client(ctx, token)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error creando servicio de Google Drive: %w", err)
	}

	return &GoogleDriveService{service: srv, folderID: folderID}, nil
}

// Upload comprime la imagen, la sube a Google Drive y retorna la URL pública
// y el fileID.
func (g *GoogleDriveService) Upload(filename string, content []byte, mimeType string) (url string, fileID string, err error) {
	// Comprimir y redimensionar antes de subir. Si falla (ej: no es imagen
	// válida), se sube el contenido original.
	if compressed, jpgName, cErr := compressImage(content, filename); cErr == nil {
		content = compressed
		filename = jpgName
		mimeType = "image/jpeg"
	}

	f := &drive.File{
		Name:    filename,
		Parents: []string{g.folderID},
	}

	created, err := g.service.Files.Create(f).
		Media(bytes.NewReader(content)).
		Fields("id").
		Do()
	if err != nil {
		return "", "", fmt.Errorf("error subiendo imagen a Drive: %w", err)
	}

	// Hacer el archivo público (lectura para cualquiera)
	permission := &drive.Permission{Type: "anyone", Role: "reader"}
	if _, err = g.service.Permissions.Create(created.Id, permission).Do(); err != nil {
		_ = g.Delete(created.Id)
		return "", "", fmt.Errorf("error estableciendo permisos en Drive: %w", err)
	}

	// URL que sí permite incrustar la imagen en <img> (el formato uc?export=view
	// ya no funciona para hotlinking; lh3.googleusercontent.com sí).
	publicURL := fmt.Sprintf("https://lh3.googleusercontent.com/d/%s", created.Id)
	return publicURL, created.Id, nil
}

// Download descarga el contenido de un archivo de Google Drive por su ID.
func (g *GoogleDriveService) Download(fileID string) ([]byte, error) {
	resp, err := g.service.Files.Get(fileID).Download()
	if err != nil {
		return nil, fmt.Errorf("error descargando imagen de Drive: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo imagen de Drive: %w", err)
	}
	return data, nil
}

// Delete elimina un archivo de Google Drive por su ID.
func (g *GoogleDriveService) Delete(fileID string) error {
	if err := g.service.Files.Delete(fileID).Do(); err != nil {
		return fmt.Errorf("error eliminando imagen de Drive: %w", err)
	}
	return nil
}

// compressImage decodifica la imagen, la redimensiona si excede maxImageWidth
// y la recodifica como JPEG. Retorna el contenido comprimido y el nuevo nombre
// con extensión .jpg.
func compressImage(content []byte, filename string) ([]byte, string, error) {
	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		return nil, "", err
	}

	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	if w > maxImageWidth {
		newW := maxImageWidth
		newH := h * maxImageWidth / w
		dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
		draw.CatmullRom.Scale(dst, dst.Bounds(), img, b, draw.Over, nil)
		img = dst
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return nil, "", err
	}

	return buf.Bytes(), jpgFilename(filename), nil
}

// jpgFilename reemplaza la extensión del archivo por .jpg.
func jpgFilename(filename string) string {
	if i := strings.LastIndex(filename, "."); i >= 0 {
		filename = filename[:i]
	}
	if filename == "" {
		filename = "imagen"
	}
	return filename + ".jpg"
}
