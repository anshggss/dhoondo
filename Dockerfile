FROM rabbitmq:4.0.9-alpine
RUN rabbitmq-plugins enable --offline rabbitmq_stomp
