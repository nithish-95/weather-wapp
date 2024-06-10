## WEATHER - APP
-- https://weather.app.nithish.net
### Running this application in the your Local Meachine

### 1.1 Clone the Repo

-- git clone git@github.com:nithish-95/weather-wapp.git

### 1.2 Install the  Chi package 

-- go get -u github.com/go-chi/chi/v5

### 1.3 Run the Application using Make (it will perfrom : clean --> build --> run )

-- make

### Your application will be available 
--  http://localhost:3000

## Using Docker 

### Building and running your application

-- docker build --progress plain --no-cache -t weatherapp1 .


### Expose the application to the given port

docker run -p 3000:3000 weatherapp1 

### Your application will be available 
--  http://localhost:3000