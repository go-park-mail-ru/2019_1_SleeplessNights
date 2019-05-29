FROM golang
#Копируем данные в контейнер
RUN mkdir /server
WORKDIR /server
COPY . /server/
#Запускаем сервер
RUN cd /server
RUN go mod tidy
RUN go get -u
RUN cd chat_microservice && go generate -n && cd ..
RUN cd game_microservice && go generate -n && cd ..
RUN cd main_microservice && go generate -n && cd ..
RUN cd user_microservice && go generate -n && cd ..
RUN cd shared/config && go generate -n && cd ../..
RUN cd shared/console && go generate -n && cd ../..
RUN cd shared/errors && go generate -n && cd ../..
RUN cd shared/logger && go generate -n && cd ../..
RUN cd shared/middleware && go generate -n && cd ../..
RUN cd shared/services && go generate -n && cd ../..
ENV BASEPATH "/server"