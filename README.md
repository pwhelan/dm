# DM: docker-machine (slim)

  dm es un programa simple para manejar los archivos para conectarse
remotamente a docker.

## Como usar dm

Para listas las maquinas configuradas con dm:

  $ dm ls

Para exportar la configuracion de una maquina:

  $ dm export MACHINE > MACHINE.pem

Para importar una configuracion de una maquina (exportada con dm):

  $ dm import MACHINE < MACHINE.pem

Para cargar la configuracion de una maquina en el shell actual:

  $ eval $(dm env MACHINE) # bash
  $ eval (dm env MACHINE) # fish
