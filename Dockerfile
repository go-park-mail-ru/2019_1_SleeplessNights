FROM golang
#Копируем данные в контейнер
RUN mkdir /server
WORKDIR /server
COPY . /server/
#Запускаем сервер
RUN cd /server
RUN go mod tidy
RUN go get -u
#RUN go generate ./...
ENV BASEPATH "/server"