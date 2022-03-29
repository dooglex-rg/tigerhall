setup:
	python ./python/db_manager_proj/manage.py makemigrations
	python ./python/db_manager_proj/manage.py migrate
	cd golang && go mod tidy

app:
	cd golang && go build -o ./bin/app ./app/ 
	./golang/bin/app

consumer:
	cd golang && go build -o ./bin/tasks ./tasks/
	./golang/bin/tasks