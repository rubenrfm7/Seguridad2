# Práctica 2 Seguridad en Redes

En primer lugar para preparar el entorno he creado un makefile para preparar todo lo necesario. Tan solo habria que ejecutar:
```bash
sudo make
```
Se puede ejecutar make clean para borrar los archivos que se han descargado:
```bash
make clean
```
He intentado descifrar el archivo de dos diferentes maneras:

- La primera se trata de un ataque por diccionario, la cual no me dio resultado,
pero basicamente este programa coge las contraseñas de un fichero y las diferentes gourutines prueban a descifrar el archivo con dichas contraseñas.
Para ejecutarlos deberemos compilar el programa:
```bash
go build Diccionary.go
```
Y ejecutar el programa:
```bash
./Diccionary
```
Al ejecutarlo nos mostrará las difentes contraseñas que esta probando y al finalizar nos dirá si ha encontrado la contraseña o no y nos mostrará el tiempo que ha tardado en ejecutarse.


- La segunda manera por la cual he conseguido descifrar la contraseña ha sido mediante un descifrado por fuerza bruta, este
  programa basicámente prueba todas las combinaciones hasta una longitud dado por parámetro.
  Para ejecutarlos deberemos compilar el programa:
```bash
go build Bruteforce.go
````
Y ejecutar el programa:
```bash
./Bruteforce
```
Al ejecutarlo nos mostrará las difentes contraseñas que esta probando y al finalizar nos dirá si ha encontrado la contraseña o no y nos mostrará el tiempo que ha tardado en ejecutarse.

Con este método conseguí descifrar la contraseña que es "jlsl" y me he dado cuenta de lo importante que es la paralelización, ya que en la primera prueba que hice usé unicamente dos gourutines y tardo mas de 5 horas en encontrar la contraseña.
Pero al utilizar un número mayor(en mi caso use 32) tardo poco mas de 1 hora y media.
