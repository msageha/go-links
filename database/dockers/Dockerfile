FROM --platform=linux/x86_64 mysql:8.0

ENV MYSQL_ROOT_PASSWORD rootpass
ENV MYSQL_DATABASE go_links
ENV MYSQL_USER mysql
ENV MYSQL_PASSWORD password

COPY ./config/my.conf /etc/mysql/conf.d/my.cnf
COPY ./schema/db.sql /docker-entrypoint-initdb.d/0_db.sql
COPY ./schema/table.sql /docker-entrypoint-initdb.d/1_table.sql

#RUN mysql -u $MYSQL_USER -p$MYSQL_PASSWORD < ./schema/db.sql
#RUN mysql -u $MYSQL_USER -p$MYSQL_PASSWORD < ./schema/table.sql
