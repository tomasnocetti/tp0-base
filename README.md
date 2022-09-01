## Ejercicio N°1.1:

Definir un script (en el lenguaje deseado) que permita crear una definición de DockerCompose con una cantidad configurable de clientes.

### Resolución

Se debe ejecutar:

```
sh ./scripts/docker-compose-gen.sh <num-of-clients> [output-file-name]
```

El parametro output-file-name es opcional, en caso de no escribirse el archivo creado tendra nombre 'docker-compose-dev.yaml'

### Consideraciones

- Todos los clientes tendran la misma configuración
