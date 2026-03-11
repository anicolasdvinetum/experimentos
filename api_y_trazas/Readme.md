# ejercicio

Vamos a intentar montar un servicio que llame a una api, haciendo un GET. 
Debe usar k6 para hacer una prueba de carga y jaeger para visualizar las trazas
Todo ello en contenedores docker. 

Puertos relevantes: 
- 8080 aplicación web. Solo tiene dos rutas, /hello y /health
- 16686 jaeger, para visualizar trazas

En principio, te posicionas en la carpeta api_y_trazas, haces 

    docker compose up

y deberían levantarse los tres contenedores sin problemas

## Reflexiones pre y durante el desarrollo

Supongo que habrá 3 contenedores, uno con k6, uno con jaeger y uno con el servicio que reciba las peticiones de k6 (y del navegador, es una aplicación web al final)

Viendo la página de jaeger, dice

    Your applications must be instrumented before they can send tracing data to Jaeger. We recommend using the OpenTelemetry instrumentation and SDKs. 

Asumo que se refiere a que la aplicación tiene que generar los "logs" (trazas) de algún modo, en concreto usando OpenTelemetry

Para el docker-compose, jaeger se levanta tal cual desde su imagen, se mapean puertos (interfaz visual, recibir rcp y recibir http, probablemente solo usemos el de http) y después el de k6 simplemente ejecuta un script

Parece que el principal reto será el usar opentelemetry. Estoy mirando la página (https://opentelemetry.io/docs/languages/go/getting-started/) para ver cómo exactamente se maneja
