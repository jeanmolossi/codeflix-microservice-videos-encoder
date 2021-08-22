# Codeflix

## Microserviço de encoder de vídeos

---

Este serviço tem como objetivo converter vídeos enviados em .mp4 para o formato mpeg dash. O formato mpeg dash por sua
vez é o melhor indicado para streaming de vídeos sob demanda

# Tabela de conteúdo

- 1 - [Primeiros passos](#primeiros-passos)

# Primeiros passos

Certifique-se de ter instalado [Docker](https://docs.docker.com/engine/install/ubuntu/)
e [Docker Compose](https://docs.docker.com/compose/install/) em sua máquina.

Faça o clone deste repositório:

$ `git clone https://github.com/jeanmolossi/codeflix-microservice-videos-encoder.git`

Acesse o diretório do repositório:

$ `cd codeflix-microservice-videos-encoder`

$ `cp .env.example .env`

Execute o comando do docker para subir os containers:

$ `docker-compose up -d`

Após todo o ambiente, você pode visualizar o processamento da fila:

Para isso configure suas variáveis de ambiente com as credenciais da AWS e um bucket com um arquivo.

**No cURL abaixo você deve substituir o payload com o `file_path` correspondente ao nome de seu arquivo no bucket**

```shell
curl -X POST 'http://localhost:15672/api/exchanges/%2F/amq.default/publish' \
-H 'Content-Type: application/json' \
-H 'Authorization: Basic cmFiYml0bXE6cmFiYml0bXE=' \
-H 'Cookie: m=2258:cmFiYml0bXE6cmFiYml0bXE%253D;' \
-d '{
    "vhost":"/",
    "name": "amq.default",
    "properties": {
        "delivery_mode": 1,
        "headers": {}
    },
    "routing_key": "videos",
    "delivery_mode": "1",
    "payload": "{\r\n\"resource_id\": \"a466c9bc-5e3c-435c-88fd-ee48ad33a9a2\",\r\n\"file_path\": \"BigBuckBunny.mp4\"\r\n}",
    "headers": {},
    "payload_encoding":"string"
}'

```