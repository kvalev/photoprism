# syntax=docker/dockerfile:1.4
FROM base

# Populate the instance with the acceptance test data
RUN cd ${PHOTOPRISM_STORAGE_PATH} && \
    wget -c https://dl.photoprism.app/qa/acceptance.tar.gz -O - | tar -xz

# Start app.
CMD ["/opt/photoprism/bin/photoprism", "start"]
