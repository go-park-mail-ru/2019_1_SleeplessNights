FROM alpine
RUN apk add --no-cache libc6-compat
COPY ./chat_microservice /chat_ms.bin
CMD [ "/chat_ms.bin" ]