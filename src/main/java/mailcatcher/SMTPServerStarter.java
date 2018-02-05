package mailcatcher;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.subethamail.smtp.server.SMTPServer;

import javax.annotation.PostConstruct;
import javax.annotation.PreDestroy;

@Component
public class SMTPServerStarter {
    private final SMTPServer server;

    @Autowired
    public SMTPServerStarter(SMTPServer server) {
        this.server = server;
    }

    @PostConstruct
    public void start() {
        server.start();
    }

    @PreDestroy
    public void stop() {
        server.stop();
    }
}
