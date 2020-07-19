package org.bal.starter;

import org.bal.config.ProducerConfiguration;
import org.bal.producer.Producer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import static java.lang.Thread.sleep;

@SpringBootApplication(scanBasePackageClasses = {ProducerConfiguration.class}, scanBasePackages = "org.bal")
public class App implements CommandLineRunner {

	@Autowired
	private Producer producer;

	public static void main(String[] args) {
		SpringApplication.run(App.class, args);
	}


	@Override
	public void run(String... args) throws Exception {
		for (int i = 0; i < 1000; i++){
			producer.send("Message is: " + System.currentTimeMillis());
			sleep(10);
		}
	}
}