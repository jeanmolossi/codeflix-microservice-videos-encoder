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

$ `./run.sh`

Após todo o ambiente, você pode visualizar o processamento da fila:

Para isso configure suas variáveis de ambiente com as credenciais da AWS e um bucket com um arquivo.

**Para publicar mensagens, é imprescindível que seja no seguinte formato:**

```json5
{
  // UUID V4
  "resource_id": "a466c9bc-5e3c-435c-88fd-ee48ad33a9a2",
  // Nome [e caminho]? do arquivo no bucket
  "file_path": "BigBuckBunny.mp4"
}
```

## Informações uteis:

### Dados default

### Postgres

- database: encoder
- user: postgres
- password: root
- host: micro-codeflix-postgres

### RabbitMQ

- user: rabbitmq
- password: rabbitmq
- host: micro-codeflix-rabbit
- ports:
    - 5672:5672
    - 15672:15672
- consumer name: micro-codeflix-videos-encoder
- consumer queue name: videos
- notification exchange: amq.direct
- notification routing key: jobs
- dead letter queue: dlx
    - type: fanout
