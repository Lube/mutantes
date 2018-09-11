## Instalar

Asegurarnos de tener el binario de go en el PATH y una instancia de Redis, de no tenerla instalada podemos hacerlo con docker facilmente:

    docker run --name redis_db -p 6379:6379 -d redis

Podemos cambiar el puerto y la DSN de redis en config/app.yml

Si tenemos make instalado podemos compilar con make build o ejecutar los tests con make test, de lo contrario podemos hacer go build en la raiz del proyecto.

# Dependencias

- Routing framework: [ozzo-routing](https://github.com/go-ozzo/ozzo-routing)
- Database: [redis-go](https://github.com/redis/redis-go)
- Data validation: [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
- Logging: [logrus](https://github.com/Sirupsen/logrus)
- Configuration: [viper](https://github.com/spf13/viper)
- Dependency management: [dep](https://github.com/golang/dep)
- Testing: [testify](https://github.com/stretchr/testify)

## Endpoints

- `GET /stats`: servicio de estadísticas
- `POST /mutant`: analiza el adn de un posible mutante

## Estructura del proyecto

- `models`: contiene las estructuras de datos para la comunicación entre capas.
- `services`: contiene la lógica de negocios principal de la aplicación.
- `components`: contiene lógica pura del dominio.
- `daos`: contiene la capa dao que interactua con la capa de persistencia.
- `apis`: contiene la capa de APIs que conecta rutas HTTP con servicios de la aplicación.

- `app`: contiene los middlewares de routing y configuración
- `errors`: contiene representacion de errores y funciones para su manejo
- `testUtils`: contiene funciones comunes para testear

El punto de entrada de la aplicación es el archivo `server.go`. Que hace lo siguiente:

- carga configuracion externa
- establece conexion con redis
- instancia servicios e inyecta dependencias
- inicia el servidor http

Documentación adicional en components/IsMutant.md
