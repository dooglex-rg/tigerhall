setup:
	pip3 install python-dotenv --quiet
	python3 ./python/db_manager_proj/manage.py makemigrations
	python3 ./python/db_manager_proj/manage.py migrate
	cd golang && go mod tidy

mock_db:
	python3 ./python/db_manager_proj/manage.py migrate --database=mock_db

app:
	cd golang/app && go build -o appx . 
	./golang/app/appx

consumer:
	cd golang/tasks && go build -o tasksx .
	./golang/tasks/tasksx

test:
	cd golang/app && go test