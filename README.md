### Ejercicio N°3:

Crear un script que permita testear el correcto funcionamiento del servidor utilizando el comando `netcat`. Dado que el servidor es un EchoServer, se debe enviar un mensaje el servidor y esperar recibir el mismo mensaje enviado. Netcat no debe ser instalado en la máquina host y no se puede exponer puertos del servidor para realizar la comunicación (hint: `docker network`).

### Resolucion

Primero que nada se le agrego un nombre a la red generada por docker-compose-up. El nombre de la red es `testing_net`.

Segundo se creo dentro de la carpeta `/netcat` una imagen sencilla que instancia un linux Alpine e instala netcat.

Tercero esta el script que corre el chequeo. Dicho script recibe opcionalmente el nombre del servicio y el puerto en el que se encuentra el servidor. En caso de no proveerlos se usa el default `server:12345`

```
sh ./test_server.sh [service {server}] [port {12345}]
```
