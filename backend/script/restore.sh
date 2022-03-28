# Need to replace file when running
docker exec -i -u postgres gostudy_db_1 pg_restore -C -d postgres < $file