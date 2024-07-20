# Go Temppc Dist
Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade. Esse sistema deverá implementar OTEL(Open Telemetry) e Zipkin.
 
## Requisitos - Serviço A (responsável pelo input):

* O sistema deve receber um input de 8 dígitos via POST, através do schema:  { "cep": "29902555" }
* O sistema deve validar se o input é valido (contem 8 dígitos) e é uma STRING
* Caso seja válido, será encaminhado para o Serviço B via HTTP

- Caso não seja válido, deve retornar:
    - *Código HTTP:* **422**
    - *Mensagem:* **invalid zipcode**

---


## Requisitos - Serviço B (responsável pela orquestração):
* O sistema deve receber um CEP válido de 8 digitos
* O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin juntamente com o nome da localização.
* O sistema deve responder adequadamente nos seguintes cenários:

- Em caso de sucesso:
    - *Código HTTP:* **200**
    - *Response Body:*
    ```json 
    {
        "city": "São Paulo",
        "temp_C": 28.5,
        "temp_F": 28.5,
        "temp_K": 28.5
    } 
    ```
- Em caso de falha, caso o CEP não seja válido (com formato correto):
    - *Código HTTP:* **422**
    - *Mensagem:* **invalid zipcode**
- ​Em caso de falha, caso o CEP não seja encontrado:
    - *Código HTTP:* **404**
    - *Mensagem:* **can not find zipcode**

---
## OTEL + Zipkin
Após a implementação dos serviços, adicione a implementação do OTEL + Zipkin:
- Implementar tracing distribuído entre **Serviço A** - **Serviço B**
- Utilizar span para medir o tempo de resposta do serviço de busca de CEP e busca de temperatura

---

# Arquivos de Testes (_.http_)
- [Request API Local](./api/api.http)

---
# Executar localmente
Será necessário realizar o start das duas aplicações (A & B) para realizar o testes localmente.

### Serviço A [_validator_]
- port:8080
```bash
cd ./cmd/validator/
```
```bash
go run main.go
```

### Serviço B [_temppc_]
- port:8090
```bash
cd ./cmd/temppc/
```
```bash
go run main.go
```