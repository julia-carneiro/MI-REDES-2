
# Definindo as imagens e os containers
SERVER1_NAME = server1
SERVER2_NAME = server2
SERVER3_NAME = server3
NETWORK_NAME = minha-rede

# Definindo as portas
SERVER1_PORT = 8000
SERVER2_PORT = 8001
SERVER3_PORT = 8002

# Imagem dos servidores (substitua pela imagem correta)
SERVER1_IMAGE = server1
SERVER2_IMAGE = server2
SERVER3_IMAGE = server3

.PHONY: all build run clean

# Target para construir os servidores
build:
	docker build -t $(SERVER1_IMAGE) ./server1
	docker build -t $(SERVER2_IMAGE) ./server2
	docker build -t $(SERVER3_IMAGE) ./server3

# Target para criar a rede se não existir
create-network:
	docker network inspect $(NETWORK_NAME) || docker network create $(NETWORK_NAME)

# Target para executar os servidores
run: create-network build
	docker run -d -p $(SERVER1_PORT):8000 --name $(SERVER1_NAME) --network $(NETWORK_NAME) $(SERVER1_IMAGE)
	docker run -d -p $(SERVER2_PORT):8001 --name $(SERVER2_NAME) --network $(NETWORK_NAME) $(SERVER2_IMAGE)
	docker run -d -p $(SERVER3_PORT):8002 --name $(SERVER3_NAME) --network $(NETWORK_NAME) $(SERVER3_IMAGE)

# Target para parar e remover os containers
clean:
	docker stop $(SERVER1_NAME) $(SERVER2_NAME) $(SERVER3_NAME) || true
	docker rm $(SERVER1_NAME) $(SERVER2_NAME) $(SERVER3_NAME) || true
