# mailcatcher

#### *mailcatcher* is a small self hosted SMTP server which catches all incoming mails and sends them to a defined mail address. Can be used with dyndns to create own addresses for trash mails. Works fine on Raspberry PI 3!


## Installation with docker
* copy docker-compose.yml
* change variables (see bellow)
* run with "docker-compose up -d"
* configure port forwarding for internet port 25 in your router (for example map internet port 25 to your Raspberry PI's port 1025)

#### Variables:
| Name | Description |
| ---- |------       |
| MC_PORT | SMTP listening port. Must match mapped port of container. |
| MC_HOST | email with this host name will be accepted (typically your dyndns host name) |
| MC_REDIRECT_TO | destination address (all mails will be redirected to this address |
| MC_SENDER_MAIL=yyy@googlemail.com | This address will be used for mail sending |
| MC_SMTP_HOST | use this SMTP server |
| MC_SMTP_PORT | with SMTP Port |
| MC_SMTP_USER | Authentication for SMTP server |
| MC_SMTP_PASSWORD| Authentication for SMTP server |
