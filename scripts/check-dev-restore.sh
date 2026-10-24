#!/bin/sh
set -eu
# Run only inside the disposable postgres-test service after the backend tests.
# No live application database, original uploads or production credentials.
test "${ARTEX_DISPOSABLE_TEST_DB:-}" = "1"
source_db=artex_test_db
target_db="artex_restore_check_$(date +%s)"
pg_dump -U artex -d "$source_db" -Fc -f /tmp/artex-restore-check.dump
createdb -U artex "$target_db"
pg_restore -U artex --exit-on-error --no-owner -d "$target_db" /tmp/artex-restore-check.dump
for database in "$source_db" "$target_db"; do
  # Table names come only from this test schema. Compare complete row contents,
  # with deterministic ordering and without emitting the underlying records.
  psql -X -U artex -d "$database" -At -v ON_ERROR_STOP=1 > "/tmp/$database.fingerprint" <<'SQL'
SELECT format('SELECT %L || '':'' || count(*) || '':'' || md5(coalesce(string_agg(row_to_json(t)::text, E''\n'' ORDER BY row_to_json(t)::text), '''')) FROM public.%I t;', tablename, tablename)
FROM pg_tables WHERE schemaname='public' ORDER BY tablename
\gexec
SQL
done
cmp -s "/tmp/$source_db.fingerprint" "/tmp/$target_db.fingerprint"
echo "Local test database dump/restore: schema and row fingerprints match."
