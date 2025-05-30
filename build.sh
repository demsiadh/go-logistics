sudo docker build -t go-logistics-app .
sudo docker run -d --name go-logistics-app -p 8080:8080 go-logistics-app
