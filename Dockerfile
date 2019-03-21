FROM ubuntu:18.04

LABEL name="Rakhimgaliyev Temirlan"
LABEL email="rakhimgaliev56@gmail.com"

ENV TZ=Europe/Moscow

USER root
RUN apt update
RUN apt install -y golang-1.10

ENV GOPATH /user/lib/go-1.10
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR $GOPATH/src/github.com/Rakhimgaliev/tech-db-forum/
ADD . $GOPATH/src/github.com/Rakhimgaliev/tech-db-forum/
# RUN go install .
EXPOSE 5000

ENV POSTGRESQLVERSION 10
RUN apt install -y postgresql-$POSTGRESQLVERSION

USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql docker -a -f scheme.sql &&\
    /etc/init.d/postgresql stop

USER root
RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf &&\
    echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf &&\
    echo "shared_buffers=256MB" >> /etc/postgresql/$PGVER/main/postgresql.conf &&\
    echo "unix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/$PGVER/main/postgresql.conf
EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

CMD service postgresql start && forum-server -port=5000  -db=postgres://docker:docker@localhost/docker
