FROM golang
#Копируем данные в контейнер
RUN mkdir /server
WORKDIR /server
COPY . /server/
#Запускаем сервер
RUN cd /server
RUN go get -u
ENV BASEPATH "/server"
ENV PORT 8080
EXPOSE $PORT
CMD go run main.go