# BOOKIFY_CD
Proyecto Cómputo Distribuido

- [Especificación del clúster](docs/cluster-specs.md)
- [Design document](docs/design-document.md)

## Frontend mínimo para el gateway

El directorio `frontend/` contiene un sitio estático que consulta `GET /overview/book` y muestra la información consolidada de library, playlist, reviews y shelves.

### Cómo ejecutar el front

1. Asegúrate de tener el gateway corriendo de forma local (`localhost:8090`) o ajusta la URL desde la propia interfaz.
2. Sirve los archivos estáticos:

   ```bash
   cd frontend
   python3 -m http.server 4173
   ```

3. Abre `http://localhost:4173` en tu navegador. Ingresa `bookId` (y opcionalmente `userId`) y el front hará la llamada `fetch` al gateway.

> Si sirves los archivos desde otro host/puerto, actualiza el campo “URL del gateway” dentro de la interfaz.

### Endpoints expuestos por el gateway

- `GET /overview/book?bookId=<id>&userId=<id>`: agrega información de todos los microservicios.
- `POST /books`: crea un libro nuevo en el servicio library. Ejemplo de cuerpo:

  ```json
  {
    "id": "book-99",
    "title": "Título",
    "author": "Autora",
    "pages": 320,
    "edition": "1a"
  }
  ```
