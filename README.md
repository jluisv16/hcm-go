# hcm-go

Microservicio base en **Go + Gin** para practicar ingeniería de software y DevOps sobre un caso de negocio real: **mantenimiento de empleados**.  
El proyecto usa **DDD + Arquitectura Hexagonal**, pruebas, contenedores, pipeline de CI, manifiestos de Kubernetes y observabilidad con Prometheus.

## Objetivo del repositorio

Este proyecto está pensado como laboratorio para practicar, de forma progresiva:

- Diseño de servicios mantenibles por capas.
- Separación clara de dominio, casos de uso e infraestructura.
- Calidad de código en pipeline (`fmt`, `vet`, tests y build).
- Empaquetado con Docker.
- Ejecución local multi-servicio con Docker Compose.
- Deploy inicial en Kubernetes con probes de salud.
- Métricas operativas mediante `/metrics` y Prometheus.

## Stack tecnológico

- **Lenguaje:** Go 1.22
- **Framework HTTP:** Gin
- **Arquitectura:** DDD + Hexagonal
- **Contenedores:** Docker (multi-stage)
- **Orquestación local:** Docker Compose
- **CI:** GitHub Actions
- **Orquestación de despliegue:** Kubernetes (Deployment + Service)
- **Observabilidad:** Prometheus

## Arquitectura (DDD + Hexagonal)

El módulo de empleados está organizado por responsabilidades:

1. **Dominio** (`internal/employees/domain`)
   - Entidad `Employee` y reglas de negocio.
   - Puerto `Repository` como contrato de persistencia.
   - Errores de dominio (ej. empleado no encontrado o email duplicado).

2. **Aplicación** (`internal/employees/application`)
   - Casos de uso: listar, obtener, crear, actualizar y eliminar empleados.
   - Coordina reglas de negocio y puertos.

3. **Infraestructura** (`internal/employees/infrastructure`)
   - Adaptador en memoria que implementa `Repository`.
   - Semilla inicial con 10 empleados.

4. **Interfaces de entrada** (`internal/employees/interfaces/http`)
   - Handlers HTTP para exponer el CRUD vía REST.
   - Mapeo entre JSON y casos de uso.

5. **Composición de la aplicación**
   - `cmd/api/main.go` hace el wiring de dependencias.
   - `internal/http/router` monta rutas globales y rutas del módulo.

## Estructura principal del proyecto

- `cmd/api/main.go`: arranque del servidor y apagado graceful.
- `internal/config`: carga de configuración por variables de entorno.
- `internal/http/handlers`: endpoints de salud (`healthz`, `readyz`) y `ping`.
- `internal/http/router`: enrutamiento y middlewares globales.
- `internal/employees/...`: dominio, aplicación, infraestructura y handlers del módulo.
- `.github/workflows/ci.yml`: pipeline de integración continua.
- `Dockerfile`: imagen de producción con multi-stage build.
- `docker-compose.yml`: stack local de app + Prometheus.
- `deploy/k8s`: manifiestos base de Deployment y Service.
- `deploy/prometheus/prometheus.yml`: configuración de scrape para entorno local.

## Modelo de datos: Employee

Cada empleado tiene 8 campos:

1. `id`
2. `first_name`
3. `last_name`
4. `email`
5. `department`
6. `role`
7. `salary`
8. `hire_date`

La aplicación inicia con datos semilla en memoria (`EMP-001` a `EMP-010`).

## Endpoints disponibles

- `GET /healthz`: liveness
- `GET /readyz`: readiness
- `GET /api/v1/ping`: endpoint de verificación rápida
- `GET /api/v1/employees`: listar empleados
- `GET /api/v1/employees/:id`: obtener por id
- `POST /api/v1/employees`: crear empleado
- `PUT /api/v1/employees/:id`: actualizar empleado
- `DELETE /api/v1/employees/:id`: eliminar empleado
- `GET /metrics`: métricas para Prometheus

### Ejemplo de payload para crear/actualizar

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@hcm.local",
  "department": "IT",
  "role": "Platform Engineer",
  "salary": 3100,
  "hire_date": "2024-01-15"
}
```

## Configuración por variables de entorno

Usa `.env.example` como base:

- `APP_NAME` (default: `hcm-go`)
- `APP_ENV` (default: `development`)
- `HTTP_PORT` (default: `8080`)
- `SHUTDOWN_TIMEOUT` (default: `10s`)

## Cómo ejecutar el proyecto

### 1) Local (sin Docker)

```bash
cp .env.example .env
make run
```

### 2) Validaciones de calidad

```bash
make fmt
make vet
make test
make build
```

### 3) Con Docker Compose (app + Prometheus)

```bash
make compose-up
```

Servicios disponibles:

- App: `http://localhost:8080`
- Prometheus: `http://localhost:9090`

Para detener:

```bash
make compose-down
```

## CI actual (GitHub Actions)

El workflow de `ci` ejecuta:

1. `go mod download`
2. Revisión de formato con `gofmt`
3. `go vet ./...`
4. `go test ./...`
5. `go build ./cmd/api`

## Guía práctica DevOps (paso a paso)

Esta sección propone una ruta de práctica realista para entrenar habilidades DevOps con este repositorio.

### Fase 1: Desarrollo local reproducible

1. Clona el repositorio.
2. Copia variables de entorno:
   - `cp .env.example .env`
3. Levanta la app en local:
   - `make run`
4. Prueba endpoints de salud y negocio:
   - `GET /healthz`
   - `GET /readyz`
   - `GET /api/v1/employees`

**Qué practicas aquí:** entorno reproducible, smoke test básico y validación manual.

### Fase 2: Calidad automática antes de integrar

1. Ejecuta:
   - `make fmt`
   - `make vet`
   - `make test`
   - `make build`
2. Corrige cualquier fallo antes de abrir PR.

**Qué practicas aquí:** shift-left quality, disciplina de pre-merge y feedback rápido.

### Fase 3: Contenerización

1. Construye la imagen:
   - `make docker-build`
2. Corre el contenedor y valida que responde en `:8080`.
3. Verifica que la imagen es mínima y no ejecuta como root.

**Qué practicas aquí:** build reproducible, runtime seguro, empaquetado para despliegue.

### Fase 4: Observabilidad base con Prometheus

1. Levanta stack:
   - `make compose-up`
2. Verifica métricas:
   - `http://localhost:8080/metrics`
3. En Prometheus (`http://localhost:9090`), consulta:
   - `up`
   - `go_goroutines`
   - `process_resident_memory_bytes`
4. Genera tráfico al API y observa cambios en series temporales.

**Qué practicas aquí:** telemetría mínima, lectura de señales runtime y troubleshooting inicial.

### Fase 5: Integración continua

1. Haz cambios pequeños en una rama.
2. Abre PR para disparar el workflow de CI.
3. Revisa logs de cada etapa y corrige si falla algo.

**Qué practicas aquí:** gobernanza de calidad, trazabilidad de builds y cultura de PR.

### Fase 6: Despliegue en Kubernetes (entorno de práctica)

1. Prepara un clúster local (`kind` o `minikube`).
2. Construye/carga la imagen local en el clúster.
3. Aplica manifiestos:
   - `kubectl apply -f deploy/k8s/deployment.yaml`
   - `kubectl apply -f deploy/k8s/service.yaml`
4. Verifica:
   - Pods en `Running`
   - Readiness/Liveness probes estables
   - Servicio accesible dentro del clúster

**Qué practicas aquí:** despliegue declarativo, probes, replicas y operación básica en K8s.

## Próximos pasos recomendados (nivel intermedio)

Si quieres llevar la práctica DevOps a un nivel más real:

1. Exigir cobertura mínima en CI.
2. Versionar y publicar imagen Docker por tag/commit SHA.
3. Agregar estrategia de release (por ejemplo, tags semánticos).
4. Añadir Grafana y dashboards.
5. Incorporar trazas con OpenTelemetry.
6. Definir alertas operativas y SLOs básicos.

## Licencia

Uso educativo/práctica interna (ajusta esta sección según tu licencia final del proyecto).
