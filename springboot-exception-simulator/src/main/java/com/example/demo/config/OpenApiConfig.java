package com.example.demo.config;

import io.swagger.v3.oas.models.ExternalDocumentation;
import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.info.License;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class OpenApiConfig {

    @Bean
    public OpenAPI springShopOpenAPI() {
        return new OpenAPI()
                .info(new Info().title("Spring Boot Exception Simulator API")
                        .description("Demo API to simulate exceptions and failures")
                        .version("v1.0.0")
                        .contact(new Contact().name("Example API Team").email("api@example.com")))
                .externalDocs(new ExternalDocumentation()
                        .description("Project README")
                        .url("https://github.com/sguo2025/springboot-exception-simulator"));
    }
}
