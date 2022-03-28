if [[ ! -e postgres-backup ]]; then
    mkdir postgres-backup
fi
d=$(date)
docker exec -u postgres gostudy_db_1 pg_dump -Fc gostudy > postgres-backup/db"${d// /_}".dump