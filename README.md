# mailcatcher

#### *mailcatcher* is a small java based SMTP server which catches all incoming mails and sends them to a defined mail address. Can be used with dyndns to create own addresses for trash mails. Works fine on Raspberry PI 2 or 3!


## Installation and configuration
* Download jar file
* put your application.configuration in the same directory
* run with Java8: java -jar -Xms32m -Xmx64m mailcatcher.jar
* configure port forwarding for internet port 25 in your router (for example map internet port 25 to your Raspberry PI's port 1025)

#### Application.properties:
```
# SMPT listener port, should be mapped to external port 25.
mailcatcher.port = 1025
# email with this host name will be accepted (for example dyndns url)
mailcatcher.host = mytrash.dyndns.org
# all accepted mail will be redirected to this address
mailcatcher.redirectTo = yourMail@Address

#spring mail configuration (Example for Gmail)
spring.mail.host = smtp.gmail.com
spring.mail.username = xxx@googlemail.com
spring.mail.password = xxx
spring.mail.properties.mail.smtp.auth = true
spring.mail.properties.mail.smtp.socketFactory.port = 465
spring.mail.properties.mail.smtp.socketFactory.class = javax.net.ssl.SSLSocketFactory
spring.mail.properties.mail.smtp.socketFactory.fallback = false
```

