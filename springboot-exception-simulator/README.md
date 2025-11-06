# Spring Boot Exception Simulator

一个用于模拟各种异常的 Spring Boot Demo：HTTP错误、延迟、OOM、CPU占用等。

## 运行

```bash
mvn spring-boot:run
```

或打包运行：

```bash
mvn package -DskipTests
java -jar target/springboot-exception-simulator-1.0.0.jar
```

## 要求

- Java 21（项目已配置为使用 Java 21 编译）。确保本地或 CI 环境使用 JDK 21，否则构建/运行会失败。

## 运行

```bash
mvn spring-boot:run
```

或打包运行：

```bash
mvn package -DskipTests
java -jar target/springboot-exception-simulator-1.0.0.jar
```

> 注意：如果你的系统默认 Java 不是 21，可以在运行前切换 JAVA_HOME。例如（示例路径，请根据你的系统调整）：
>
> ```bash
> export JAVA_HOME=/usr/local/sdkman/candidates/java/21.0.7-ms
> export PATH=$JAVA_HOME/bin:$PATH
> ```

## Swagger / OpenAPI

本项目使用 springdoc OpenAPI（springdoc-openapi-starter-webmvc-ui）自动生成 API 文档并提供 Swagger UI。

启动应用后可访问：

- Swagger UI: http://localhost:8080/swagger-ui/index.html
- OpenAPI JSON: http://localhost:8080/v3/api-docs

示例快速使用：

```bash
# 使用 JDK21（示例路径，请根据你的环境调整）
export JAVA_HOME=/usr/local/sdkman/candidates/java/21.0.7-ms
export PATH=$JAVA_HOME/bin:$PATH

# 本地运行
mvn spring-boot:run

# 或已打包的可执行 jar
# java -jar target/springboot-exception-simulator-1.0.0.jar
```

如果启动后无法访问 Swagger UI：

- 确认应用已成功启动且监听 8080（或在 `application.yml` 中配置的端口）。
- 打开 `http://localhost:8080/v3/api-docs` 看是否返回 JSON。如果没有，检查应用日志中的启动异常。
- 确认运行时使用 JDK21（`java -version` / `mvn -v` 输出应显示 Java 21）。

## 接口示例

| URL | 说明 |
|------|-----|
| /simulate/runtime | 运行时异常 |
| /simulate/null | 空指针异常 |
| /simulate/delay | 延迟响应 |
| /simulate/oom | 模拟OOM |
| /simulate/cpu | 模拟CPU高占用 |
| /simulate/random | 随机异常 |
| /simulate/status/503 | 自定义HTTP状态码 |
| /actuator/prometheus | Prometheus指标 |

| URL | 说明 |
|------|-----|
| /simulate/runtime | 运行时异常 |
| /simulate/null | 空指针异常 |
| /simulate/delay | 延迟响应 |
| /simulate/oom | 模拟OOM |
| /simulate/cpu | 模拟CPU高占用 |
| /simulate/random | 随机异常 |
| /simulate/status/503 | 自定义HTTP状态码 |
| /actuator/prometheus | Prometheus指标 |
