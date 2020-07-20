package org.bal.config;

import org.apache.activemq.artemis.jms.client.ActiveMQConnectionFactory;
import org.bal.consumer.Consumer;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.jms.annotation.EnableJms;
import org.springframework.jms.config.DefaultJmsListenerContainerFactory;

@Configuration
@EnableJms
public class ConsumerConfiguration {


    @Value("${artemis.host}")
    private String brokerHost;

    @Value("${artemis.port}")
    private int brokerPort;

    @Value("${artemis.user}")
    private String user;

    @Value("${artemis.password}")
    private String password;

    @Bean
    public ActiveMQConnectionFactory receiverActiveMQConnectionFactory() {
        ActiveMQConnectionFactory factory = new ActiveMQConnectionFactory("tcp://" + brokerHost+":" + brokerPort);
        factory.setConsumerWindowSize(0);
        factory.setUser(user);
        factory.setPassword(password);
        return factory;
    }

    @Bean
    public DefaultJmsListenerContainerFactory jmsListenerContainerFactory() {
        DefaultJmsListenerContainerFactory factory =
                new DefaultJmsListenerContainerFactory();
        factory.setConnectionFactory(receiverActiveMQConnectionFactory());
        //factory.setConcurrency("3-10");

        return factory;
    }
    @Bean
    public Consumer receiver() {
        return new Consumer();
    }
}