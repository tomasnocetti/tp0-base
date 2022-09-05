### Ejercicio N°5:

Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

Por el lado de los clientes (quienes emularán _agencias de quiniela_) deberán recibir como variables de entorno los siguientes datos de una persona: nombre, apellido, documento y fecha de nacimiento. Dichos datos deberán ser enviados al servidor para saber si corresponden a los de un ganador, información que deberá loguearse. Por el lado del servidor (que emulará la _central de Lotería Nacional_), deberán recibirse los datos enviados desde los clientes y analizar si corresponden a los de un ganador utilizando la función provista `is_winner(...)`, para luego responderles.

Deberá implementarse un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:

- Serialización de los datos.
- Definición de un protocolo para el envío de los mensajes.
- Correcto encapsulamiento entre el modelo de dominio y la capa de transmisión.
- Empleo correcto de sockets, incluyendo manejo de errores y evitando el fenómeno conocido como _short-read_.

### Resolución

Para este ejercicio se hizo una adaptacion tanto del cliente como del servidor. Para probar la ejecucción basta con correr:

```
make docker-compose-up && make docker-compose-logs
```

De esta manera se correra una iteracion del cliente que consultara los datos de un participante. Los datos estan definidos en el archivo de configuracion `config.yaml` del cliente.

### Protocolo

Se detalla en terminos generales la interacion definida entre cliente-servidor.

Para el mensaje de chequeo de ganador se define la siguiente estructura:

```
| Opcode | Length in bytes | Client1 info (ID;FirstName;LastName;Birth)	|
| 1 byte | 	4 bytes 	   |		Dynamic						|
```

Constara de:

- OpCode: codigo para la operación de consulta ( en este caso 1)
- Length: largo del payload en bytes.
- Info: informacion del participante serpara por ';' manteniendo el orden estipulado.

Dicho mensaje recibira una respuesta por parte del servidor que tendra la siguiente estructura:

```
| Winners |
| 2 byte  |
```

Determina la cantidad de ganadores, para este punto del ejercicio esto puede valer 1 o 0.
