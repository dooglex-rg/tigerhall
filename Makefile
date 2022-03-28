setup:
	python ./python/db_manager_proj/manage.py makemigrations
	python ./python/db_manager_proj/manage.py migrate
	cd golang && go mod tidy

app:
	cd golang && go build -o ./bin/app.exe ./app/ 
	./golang/bin/app.exe

consumer:
	cd golang && go build -o ./bin/tasks.exe ./tasks/
	./golang/bin/tasks.exe