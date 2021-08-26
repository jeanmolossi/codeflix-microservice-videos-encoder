if [[ ! -f ".env" ]]; then
    cp .env.example .env
fi

if [[ ! -d ./tmp ]]; then
  mkdir ./tmp
fi

go mod tidy

docker-compose up -d