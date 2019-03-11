# Build Instructions

## How to run the server (FIRST TIME ONLY)

1. Install docker depending on your operating system

		https://docs.docker.com/install/#reporting-security-issues

		MAC:
		https://docs.docker.com/docker-for-mac/install/
		WIN:
		https://docs.docker.com/docker-for-windows/install/

2. Move into your project directory

		cd smart-contract/

3. Install the server using docker

		docker-compose up -d

4. Your server will be up at localhost:9000



## How to run the server (NEXT TIME)

1. just run:

	    docker-compose start

2. to view the logs, use:

	    docker logs smart-contract -f

	quit the logs using ctrl-c

3. to stop the server:

	    docker-compose stop


### Notes:
1. If you want docker-compose to run the server in the background, you can use

        docker-compose up -d
   -d command stands for detached mode
   Now, the docker will run on the background, to attach the log back, you can use `docker attach smart-contract` or `docker-compose logs -f`.
   If you want to stop the server, you need to stop the container using `docker-compose stop`

2. If you accidentally press ctrl-c, and you wish to attach to the docker logs again. you can use:

		docker attach smart-contract
	and you will continue to see the output of backend server. Or, you can use:

		docker-compose logs -f
	to aggregate all container output up until now (including mysql)

3. You may need to run using `sudo` for all of the commands above. 
    Especially if you find errors like: 
        
        "ERROR: Couldn't connect to Docker daemon - you might need to run docker-machine start default. "
    when running  `docker-compose up`
    

4. For WINDOWS user:
	the port forwarding mechanism for docker is not working on windows, so you need to use docker machine's IP. to get the IP, use:

		docker-machine ip

	and then it will output the IP of your docker machine, e.g: 192.168.99.100
	next you can access the django using that IP:

		http://192.168.99.100:8000/
		



## How to stop the server

1. Stop smart-contract project container

	```
	docker-compose stop
	```

2. See whether the container stopped

	```
	sudo docker ps -a
	```

## Other How to:

Run commands on the docker container: 
        
```
sudo docker exec -ti smart-contract (enter commands here)
```

### Manual Build Instructions: 

Note: Please do not use this except if it is really necessary! Use the default docker-compose is already fine.

1. Build the docker container 

        sudo docker build -t smart-contract .
    
2. Run the docker container

        sudo docker run -P --rm -it smart-contract


