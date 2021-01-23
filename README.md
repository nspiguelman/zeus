# zeus

## DESARROLLO

Para poder usar la aplicación primero debe levantarse la base de datos y luego la API.

```bash
    docker-compose -f ./docker/sql/docker-compose.yml up --build -d
    docker-compose -f ./docker/redis/docker-compose.yml up --build -d 
    docker-compose up --build
```
