FROM eclipse-temurin:21-jdk
WORKDIR /app

# 拷贝应用 jar
COPY target/springboot-exception-simulator-1.0.0.jar app.jar

# 拷贝 OpenTelemetry Java Agent
COPY target/opentelemetry-javaagent.jar /app/opentelemetry-javaagent.jar

EXPOSE 8888

# 可通过环境变量灵活设置 JVM 参数
ENV JAVA_OPTS="-XX:+PreserveFramePointer -javaagent:/app/opentelemetry-javaagent.jar"

# 可通过环境变量灵活设置 OpenTelemetry 配置
ENV OTEL_SERVICE_NAME=my-java-service
ENV OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://http://13.49.159.64:8080/v1/traces
ENV OTEL_METRICS_EXPORTER=none
ENV OTEL_LOGS_EXPORTER=none

ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -jar /app/app.jar"]
