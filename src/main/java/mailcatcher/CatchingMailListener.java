package mailcatcher;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.stereotype.Service;
import org.subethamail.smtp.helper.SimpleMessageListener;

import javax.mail.Message;
import javax.mail.Session;
import javax.mail.internet.InternetAddress;
import javax.mail.internet.MimeMessage;
import java.io.InputStream;
import java.util.Properties;

@Service
public class CatchingMailListener implements SimpleMessageListener {

    private final Logger log = LoggerFactory.getLogger(getClass());

    @Autowired
    private JavaMailSender mailSender;

    @Value("${mailcatcher.redirectTo}")
    private String redirectTo;

    @Value("${mailcatcher.host}")
    private String host;

    @Override
    public boolean accept(String from, String recipient) {
        return recipient.endsWith(host);
    }

    @Override
    public void deliver(String from, String recipient, InputStream data) {
        log.info("received new mail from {}", from);
        try {
            MimeMessage m = new MimeMessage(Session.getDefaultInstance(new Properties()), data);

            m.setSubject("[MAILCATCHER]" + m.getSubject());
            m.setFrom(from);
            m.setRecipient(Message.RecipientType.TO, new InternetAddress(redirectTo, recipient));

            mailSender.send(m);
        } catch (Exception e) {
            log.error("fatal error", e);
        }

    }
}
