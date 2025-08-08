# Load Tester

Ferramenta para realizar testes de carga em serviços web.

## Como usar

Após realizar o gitclone do projeto, levantar o ambiente executando: 
```
docker compose up -d
```

Executar os testes abaixo para utilizar a ferramenta:

- Teste Básico:
```
docker run --rm pos-stress-test-stress-test --url=http://facebook.com --requests=6000 --concurrency=400
```

- Teste com mais requests para forçar erros:
```
docker run --rm pos-stress-test-stress-test --url=http://facebook.com --requests=6000 --concurrency=400
```

Exemplo de saída do relatório:

<img width="1410" height="320" alt="stress-test" src="https://github.com/user-attachments/assets/e8843f3b-b4ce-4b1e-a0a0-54738bb50b78" />

