FROM eclipse-temurin:21-jdk
WORKDIR /app
COPY target/*.jar app.jar
EXPOSE 8888

# 可通过环境变量灵活设置 JVM 参数
ENV JAVA_OPTS="-XX:+PreserveFramePointer"

ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -jar /app/app.jar"]

