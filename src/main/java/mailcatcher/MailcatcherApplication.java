package mailcatcher;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.subethamail.smtp.helper.SimpleMessageListenerAdapter;
import org.subethamail.smtp.server.SMTPServer;

@SpringBootApplication
public class MailcatcherApplication {

    @Value("${mailcatcher.port}")
    private int port;

    public static void main(String[] args) {
        SpringApplication.run(MailcatcherApplication.class, args);
    }

    @Bean
    public SMTPServer createServer(CatchingMailListener service) {
        SMTPServer server = new SMTPServer(new SimpleMessageListenerAdapter(service));
        server.setPort(port);
        return server;
    }
}
