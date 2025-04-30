FROM codeberg.org/forgejo/forgejo:11
COPY forgejo.db app.ini /data/

RUN chmod -R 777 /data
RUN chown 1000:1000 /data