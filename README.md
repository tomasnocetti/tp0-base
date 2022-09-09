### Ejercicio N°8:

Agregar en los clientes una consulta por el número total de ganadores de todas las agencias, por lo cual se deberá modificar el servidor para poder llevar el trackeo de dicha información.

En caso de que alguna agencia consulte a la central antes de que esta haya completado el procesamiento de las demás, deberá recibir una respuesta parcial con el número de agencias que aún no hayan finalizado su carga de datos y volver a consultar tras N segundos.

### Protocolo

Para esta parte del ejercicio se crea un nuevo OpCode dentro del protocolo que indica la intencion de recibir estadisticas.

```
| Opcode |
| 1 byte |
```

Constara de:

- OpCode: codigo para la operación de consulta ( en este caso 1)

Dicho mensaje recibira una respuesta por parte del servidor que tendra la siguiente estructura:

```
|    Partial      |     Winners     |
|    1 byte       |      4 bytes    |
```

### Mecanismos de sincronizacion

Para la parte de sincronizacion de escritura del Archivo se uso un IPC Lock provisto en la libreria de Multiprocessing. Se puede ver su implementacion en server/common/persistance.

Asi mismo para saber cual es la totalidad de clientes activos en la plataforma, la misma entidad Persistance mantiene un contador activo de los clientes que estan `Chequeando Participantes`, se puede ver esto en la linea 37 del archivo `client_connection.py`. Si ese contador no es 0 y un cliente decide solicitar estadisticas, el resultado sera enviado como parcial y el cliente tendra que reintentar luego de 8 segundos.
