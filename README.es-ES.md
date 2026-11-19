![Build](https://github.com/0xERR0R/mailcatcher/workflows/Build/badge.svg)

# mailcatcher

#### *mailcatcher* es un pequeño servidor SMTP auto-alojado que captura todos los correos entrantes y los envía a una dirección de correo definida. Puede usarse con dyndns para crear direcciones propias para correos basura. ¡Funciona perfectamente en Raspberry PI 3!


## Instalación con docker
* copiar docker-compose.yml
* cambiar variables (ver más abajo)
* ejecutar con "docker-compose up -d"
* configurar el reenvío de puertos para el puerto de internet 25 en tu router (por ejemplo, mapear el puerto de internet 25 al puerto 1025 de tu Raspberry PI)

#### Variables:
| Nombre | Descripción |
| ---- |------       |
| MC_PORT | Puerto de escucha SMTP. Debe coincidir con el puerto mapeado del contenedor. |
| MC_HOST | se aceptarán correos con este nombre de host (normalmente tu nombre de host de dyndns) |
| MC_REDIRECT_TO | dirección de destino (todos los correos serán redirigidos a esta dirección) |
| MC_SENDER_MAIL | Esta dirección se utilizará para el envío de correos |
| MC_SMTP_HOST | usar este servidor SMTP |
| MC_SMTP_PORT | con puerto SMTP |
| MC_SMTP_USER | Autenticación para el servidor SMTP |
| MC_SMTP_PASSWORD| Autenticación para el servidor SMTP |


## Sugerencias

Si estás utilizando el SMTP de GMail, por favor activa el uso de aplicaciones menos seguras en https://support.google.com/accounts/answer/6010255?hl=en .
