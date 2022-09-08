### Ejercicio N°4:

Modificar el cliente y el servidor para que el programa termine de forma gracefully al recibir la signal SIGTERM. Terminar la aplicación de forma gracefully implica que todos los sockets y threads/procesos de la aplicación deben cerrarse/joinearse antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker-compose down`).

### Resolucion

Para probarlo se debe correr

```
make docker-compose-up
```

y luego se puede detener todos los servicios o individualmente.

```
make docker-compose-down
```

```
docker stop client1/server
```

### Nota

Si se viene del ejercicio 3 puede ser necesario correr

```
docker network rm testing_net
```
