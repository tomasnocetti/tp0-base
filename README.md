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
