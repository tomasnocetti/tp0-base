### Ejercicio N°6:

Modificar los clientes para que levanten los datos de los participantes desde los datasets provistos en los archivos de prueba en lugar de las variables de entorno. Cada cliente deberá consultar por todas las personas de un mismo set (los cuales representan a los jugadores de cada agencia) en forma de batch, de manera de poder hacer varias consultas en un solo request. El servidor por otro lado deberá responder con los datos de todos los ganadores del batch, y cada cliente al finalizar las consultas deberá loguear el porcentaje final de jugadores que hayan ganado en su agencia.

### Resolución

Para este ejercicio se hizo una adaptacion tanto del cliente como del servidor. Para probar la ejecucción basta con correr:

Expande el dataset y configura docker-compose para el uso de multiples clientes.

```
sh ./scripts/docker-compose-gen.sh
```

Corre el sistema:

```
make docker-compose-up && make docker-compose-logs
```

### Protocolo

Se detalla en terminos generales la interacion definida entre cliente-servidor.

Para el mensaje de chequeo de ganador se define la siguiente estructura:

```
| Opcode | Length in bytes | Contestants info (ID;FirstName;LastName;Birth)	|
| 1 byte | 	4 bytes 	   |		Dynamic						|
```

Constara de:

- OpCode: codigo para la operación de consulta ( en este caso 1)
- Length: largo del payload en bytes.
- Info: informacion de todos los participante separada por ';' manteniendo el orden estipulado; participantes separados por '|'.

Dicho mensaje recibira una respuesta por parte del servidor que tendra la siguiente estructura:

```
| Length in Bytes |     ID info (|)     |
|    4 byte       |        Dynamic      |
```

El largo del payload con los ids de los ganadores separados por el caracter '|'

### Problemas Interesantes

* En este punto de los ejercicios se encontro un problema de alocacion de memoria en go, se estaba guardando basura en el buffer que recibia el socket por lo que generaba datos incorrectos. El problema no era deterministico debido a que en cada iteracion cambiaba el output de los datos.
