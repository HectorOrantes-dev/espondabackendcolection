package application

import (
	"sync"

	"coleccionbackend/src/feature/coleccion/domain"
)

// uploadImagesParallel sube todas las imágenes concurrentemente (cada una en
// su propia goroutine) y espera a que terminen. Mantiene el orden original.
// Si alguna falla, elimina las que sí se subieron (rollback) y retorna el error.
func uploadImagesParallel(svc domain.ImageService, images []ImageInput) (urls []string, ids []string, err error) {
	n := len(images)
	if n == 0 {
		return nil, nil, nil
	}

	resultURLs := make([]string, n)
	resultIDs := make([]string, n)
	resultErrs := make([]error, n)

	var wg sync.WaitGroup
	for i, img := range images {
		wg.Add(1)
		go func(i int, img ImageInput) {
			defer wg.Done()
			url, id, uploadErr := svc.Upload(img.Filename, img.Content, img.MimeType)
			resultURLs[i], resultIDs[i], resultErrs[i] = url, id, uploadErr
		}(i, img)
	}
	wg.Wait()

	// Si alguna subida falló, hacer rollback de las exitosas.
	for _, e := range resultErrs {
		if e != nil {
			for j, id := range resultIDs {
				if resultErrs[j] == nil && id != "" {
					_ = svc.Delete(id)
				}
			}
			return nil, nil, e
		}
	}

	return resultURLs, resultIDs, nil
}
