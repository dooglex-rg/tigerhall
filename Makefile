makemigrations:
	python ./python/db_manager_proj/manage.py makemigrations

migrate:
	python ./python/db_manager_proj/manage.py migrate

start_app:
	cd golang && go build -o ./bin/app ./app/ 
	./golang/bin/app

start_task:
	cd golang && go build -o ./bin/tasks ./tasks/
	./golang/bin/tasks