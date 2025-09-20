run:
	docker-compose up -d postgreSQL
postgres:
	docker exec -it newpost psql -U user -d postgresA
createdb:
	docker exec -it newpost createdb --username=user --owner=user go-sample 

