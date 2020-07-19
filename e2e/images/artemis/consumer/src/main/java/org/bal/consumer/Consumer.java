package org.bal.consumer;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.jms.annotation.JmsListener;

import static java.lang.Thread.sleep;

public class Consumer {

    private static final Logger LOGGER =
            LoggerFactory.getLogger(Consumer.class);

    @JmsListener(destination = "test", concurrency = "1")
    public void receive(String message) throws IllegalAccessException {
      try {
        sleep(50);
        LOGGER.info("received message='{}' ......", message);
      } catch (InterruptedException e) {
        e.printStackTrace();
      }

    }
}