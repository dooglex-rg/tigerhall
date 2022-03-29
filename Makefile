setup:
	python ./python/db_manager_proj/manage.py makemigrations
	python ./python/db_manager_proj/manage.py migrate
	cd golang && go mod tidy

app:
	cd golang/app && go build -o appx . 
	./golang/app/appx

consumer:
	cd golang/tasks && go build -o tasksx .
	./golang/tasks/tasksx