package org.bal.starter;

import org.bal.config.ConsumerConfiguration;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication(scanBasePackageClasses = {ConsumerConfiguration.class}, scanBasePackages = "org.org.bal")
public class App {

	public static void main(String[] args) {
		SpringApplication.run(App.class, args);
	}
}