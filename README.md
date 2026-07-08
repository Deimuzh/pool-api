# Piscina Comunitaria Los Ceibos API

API REST construida en Go para administrar una piscina comunitaria: clientes, pagos, reservas, seguridad, mantenimiento y usuarios administrativos.

El proyecto esta organizado por capas para separar responsabilidades y facilitar pruebas:

```txt
HTTP -> Handler -> Service -> Repository -> GORM -> SQLite/PostgreSQL
```

## Stack

- Go 1.25
- Chi router
- GORM
- SQLite para desarrollo local
- PostgreSQL con Docker Compose
- JWT para autenticacion
- bcrypt para contrasenas
- tests con `testing`, `httptest` y `testify`
- GitHub Actions para CI

## Modulos

- **Seguridad**: guardavidas, accesos de clientes e incidentes.
- **Clientes**: clientes, reservas y pagos.
- **Mantenimiento**: equipos, registros de mantenimiento y productos quimicos.
- **Auth**: login, usuarios administrativos, JWT y roles.

## Responsables por modulo

- **Seguridad**: Anthony Joel Mendoza Arcentales.
- **Clientes**: Garcia Cedeno Geovanny Alexander.
- **Mantenimiento**: Lucas Holguin Nathaly Jasmin.
- **Auth, Docker, CI e integracion**: responsabilidad grupal.

## Arquitectura

Diagrama detallado: [`docs/arquitectura.md`](docs/arquitectura.md).

- `cmd/piscina-api/main.go`: ensambla configuracion, base de datos, services, handlers, router y servidor HTTP.
- `internal/models`: structs del dominio con tags JSON/GORM.
- `internal/handlers`: recibe requests HTTP, decodifica JSON y responde status codes.
- `internal/service`: contiene reglas de negocio y validaciones.
- `internal/storage`: interfaces de repository e implementacion GORM.
- `internal/middleware`: autenticacion JWT, roles y CORS.
- `internal/config`: variables de entorno y valores por defecto.
- `web`: interfaz HTML/CSS simple para probar la API.

## Reglas principales

- Solo usuarios autenticados pueden usar las rutas `/api/v1/*`, excepto login.
- El JWT contiene el ID y el rol del usuario.
- Los endpoints de usuarios requieren rol `admin`.
- Un cliente con membresia no necesita pagar entrada diaria.
- Un cliente sin membresia necesita pago de entrada para registrar acceso.
- Un incidente requiere guardavida y cliente validos.
- Un incidente solo se registra si el cliente tiene membresia o acceso autorizado.
- Las reservas solo aceptan duraciones permitidas.
- Los pagos validan concepto y monto mayor a cero.

## Correr local con SQLite

```bash
go mod download
go run ./cmd/piscina-api
```

Servidor:

```txt
http://localhost:8080
```

Login inicial sembrado automaticamente si no existen usuarios:

```txt
email: admin@piscina.com
password: admin123
rol: admin
```

## Correr con Docker Compose

```bash
docker compose up --build
```

Esto levanta:

- API Go en `http://localhost:8080`
- PostgreSQL 16 en el servicio `postgres`
- volumen persistente `postgres-data`

Para detener:

```bash
docker compose down
```

Para borrar tambien los datos:

```bash
docker compose down -v
```

## Variables de entorno

Ver `.env.example`.

Variables principales:

- `PUERTO`: puerto HTTP, por ejemplo `:8080`.
- `DB_DRIVER`: `sqlite` o `postgres`.
- `RUTA_DB`: archivo SQLite local.
- `DB_DSN`: cadena de conexion PostgreSQL.
- `JWT_SECRETO`: secreto para firmar tokens.
- `JWT_DURACION`: duracion del token, por ejemplo `24h`.
- `HTTP_READ_TIMEOUT`: timeout de lectura.
- `HTTP_WRITE_TIMEOUT`: timeout de escritura.

## Endpoints principales

### Auth

| Metodo | Ruta | Descripcion |
| --- | --- | --- |
| POST | `/api/v1/login` | Iniciar sesion y obtener JWT |

Body de login:

```json
{
  "email": "admin@piscina.com",
  "password": "admin123"
}
```

Usar el token en rutas protegidas:

```txt
Authorization: Bearer <token>
```

### Usuarios

Requiere JWT con rol `admin`.

| Metodo | Ruta | Descripcion |
| --- | --- | --- |
| GET | `/api/v1/usuarios/` | Listar usuarios |
| POST | `/api/v1/usuarios/` | Crear usuario |
| GET | `/api/v1/usuarios/{id}` | Obtener usuario |
| PUT | `/api/v1/usuarios/{id}` | Actualizar usuario |
| DELETE | `/api/v1/usuarios/{id}` | Eliminar usuario |

### Seguridad

| Metodo | Ruta | Descripcion |
| --- | --- | --- |
| GET | `/api/v1/guardavidas/` | Listar guardavidas |
| POST | `/api/v1/guardavidas/` | Crear guardavida |
| GET | `/api/v1/guardavidas/{id}` | Obtener guardavida |
| PUT | `/api/v1/guardavidas/{id}` | Actualizar guardavida |
| DELETE | `/api/v1/guardavidas/{id}` | Eliminar guardavida |
| GET | `/api/v1/accesos/` | Listar accesos |
| POST | `/api/v1/accesos/` | Registrar acceso por cliente |
| DELETE | `/api/v1/accesos/{id}` | Eliminar acceso |
| GET | `/api/v1/incidentes/` | Listar incidentes |
| POST | `/api/v1/incidentes/` | Crear incidente |
| GET | `/api/v1/incidentes/{id}` | Obtener incidente |
| PUT | `/api/v1/incidentes/{id}` | Actualizar incidente |
| DELETE | `/api/v1/incidentes/{id}` | Eliminar incidente |

### Clientes

| Metodo | Ruta | Descripcion |
| --- | --- | --- |
| GET | `/api/v1/clientes/` | Listar clientes |
| POST | `/api/v1/clientes/` | Crear cliente |
| GET | `/api/v1/clientes/{id}` | Obtener cliente |
| PUT | `/api/v1/clientes/{id}` | Actualizar cliente |
| DELETE | `/api/v1/clientes/{id}` | Eliminar cliente |
| GET | `/api/v1/reservas/` | Listar reservas |
| POST | `/api/v1/reservas/` | Crear reserva |
| GET | `/api/v1/pagos/` | Listar pagos |
| POST | `/api/v1/pagos/` | Crear pago |

### Mantenimiento

| Metodo | Ruta | Descripcion |
| --- | --- | --- |
| GET | `/api/v1/equipos/` | Listar equipos |
| POST | `/api/v1/equipos/` | Crear equipo |
| GET | `/api/v1/mantenimientos/` | Listar mantenimientos |
| POST | `/api/v1/mantenimientos/` | Crear mantenimiento |
| GET | `/api/v1/quimicos/` | Listar quimicos |
| POST | `/api/v1/quimicos/` | Crear producto quimico |

## Tests

Ejecutar todo:

```bash
go test ./...
```

Con cobertura:

```bash
go test ./... -cover
```

El modulo Seguridad incluye pruebas de:

- service con mocks;
- handler con `httptest`;
- repository GORM con SQLite en memoria.

## Postman

La coleccion exportada para probar la API esta en:

```txt
docs/postman/piscina-api.postman_collection.json
```

Flujo recomendado para demo:

1. Ejecutar `docker compose up --build`.
2. Importar la coleccion en Postman.
3. Ejecutar `Auth / Login admin` para guardar el JWT en la variable `token`.
4. Probar endpoints protegidos de Seguridad, Clientes y Mantenimiento.

## CI

La pipeline de GitHub Actions esta en `.github/workflows/ci.yml` y ejecuta:

```txt
go vet ./...
go test ./... -cover
go build ./cmd/piscina-api
```

## Documentacion del Hito 3

- Diagrama de arquitectura: [`docs/arquitectura.md`](docs/arquitectura.md).
- Coleccion Postman: [`docs/postman/piscina-api.postman_collection.json`](docs/postman/piscina-api.postman_collection.json).
- Documento de cierre: [`docs/cierre.md`](docs/cierre.md).

## Estado pendiente

- Aumentar cobertura total de tests.
- Completar nombres de responsables por modulo si el docente lo solicita explicitamente.
