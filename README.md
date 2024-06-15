# Desafio client-server-api

### server.go
Servidor com a rota **/cotacao**, responsável por:
 - buscar a cotação do dólar através da api **economia.awesomeapi.com.br**, 
 - salvá-la em um banco de dados **sqlite** e
 - por fim, retornar a cotação do dólar em formato **json**
 
### client.go
Cliente, que faz constantemente requições para o endpoint **/cotacao** e salva a última contação em um arquivo **cotacao.txt**

## Instruções para execução
Clone o projeto e entre na pasta raiz:
```
git clone https://github.com/usuario/repo.git
cd repo
```
Baixe as dependências:
```
go mod tidy
```
Execute os dois sistemas:
#### server.go
```
go run cmd/server/server.go
```
#### client.go
```
go run cmd/client/client.go <arg1>
```
- **arg1:** valor em inteiro que define tempo de intervalo entre as requisições em segundos
	- Utilize 0 ou menos caso queira fazer apenas uma requisição
	- Caso não seja informado, o valor será definido para padrão de 60 segundos
