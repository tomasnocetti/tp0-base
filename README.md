### Ejercicio N°7:

Modificar el servidor actual para que el mismo permita procesar mensajes y aceptar nuevas conexiones en paralelo. Además, deberá comenzar a persistir la información de los ganadores utilizando la función provista `persist_winners(...)`. Considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de dicha persistencia.

En caso de que el alumno desee implementar un nuevo servidor en Python, deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

### Resolución

Se implementa la libreria de multiprocesing en Python para el Servidor.
