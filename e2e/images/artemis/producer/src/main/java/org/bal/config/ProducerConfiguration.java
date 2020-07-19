package org.bal.config;

import org.bal.producer.Producer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.jms.annotation.EnableJms;

@Configuration
@EnableJms
public class ProducerConfiguration {

  @Bean
  public Producer sender() {
    return new Producer();
  }
}