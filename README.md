### Ejercicio N°2:

Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo debe ser inyectada al ejemplo y persistida afuera del mismo (hint: `docker volumes`).

### Resolucion

Se mueven los archivos de configuracion a una carpeta config en el root del proyecto. Tanto el servidor como el cliente tienen un archivo de configuracion designado que no se copia a la imagen, sino que se monta como volumen al inicializar los containers con `make docker-compose-up`.

Para probarlo simplemente se debe inicializar todo. Realizar un stop de los containers con `docker stop $(docker ps -q)`, modificar el archivo e inicializar los containers nuevamente con `docker compose -f docker-compose-dev.yaml up`.
