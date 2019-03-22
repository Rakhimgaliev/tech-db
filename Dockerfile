FROM ubuntu:18.04

LABEL name="Rakhimgaliyev Temirlan"
LABEL email="rakhimgaliev56@gmail.com"

ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

USER root
RUN apt update
RUN apt install -y golang-1.10
RUN apt install -y git

ENV GOROOT /usr/lib/go-1.10
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR $GOPATH/src/github.com/Rakhimgaliev/tech-db-forum/
ADD . $GOPATH/src/github.com/Rakhimgaliev/tech-db-forum/
RUN go install ./forum/main/
EXPOSE 5000

# ENV POSTGRESQLVERSION 10
# RUN apt install -y postgresql-$POSTGRESQLVERSION

# USER postgres
# RUN /etc/init.d/postgresql start
# RUN psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';"
# RUN createdb -O docker docker
# RUN psql docker -a -f scheme.sql
# RUN /etc/init.d/postgresql stop

# USER root
# RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$POSTGRESQLVERSION/main/pg_hba.conf &&\
#     echo "listen_addresses='*'" >> /etc/postgresql/$POSTGRESQLVERSION/main/postgresql.conf &&\
#     echo "shared_buffers=256MB" >> /etc/postgresql/$POSTGRESQLVERSION/main/postgresql.conf &&\
#     echo "unix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/$POSTGRESQLVERSION/main/postgresql.conf
# EXPOSE 5432

# VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# CMD service postgresql start && go run ./forum/main -port=5000  -db=postgres://docker:docker@localhost/docker

WORKDIR $GOPATH/bin/
CMD main