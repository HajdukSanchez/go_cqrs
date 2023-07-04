# IMAGE
FROM postgres:10.3
# COPY INITIAL FILE FOR START DB
COPY up.sql /docker-entrypoint-initdb.d/1.sql
# EXECUTE DB
CMD ["postgres"]