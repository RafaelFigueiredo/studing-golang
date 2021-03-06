# Blockchain com Golang
Por conta de um trabalho freelancer que peguei recentemente me deparei com o seguinte problema: 
> "Cada transação deve ser registrada no banco de dados de forma permanente e tenho que garantir que após a inserção esses registros não foram alterados"


O banco de dados utilizado no projeto é o PostgreSQL, então talvez uma solução pudesse ser simplesmente restringir os privilégios do usuário no banco de dados para "apenas leitura" mas nada impediria que um usuário com privilégios alterasse um dos registros, mas vamos com calma na teoria da conspiração e começar do começo.
Com essa implementando simples de uma blockchain poderemos saber quando um registro é modificado após a inserção mesmo que o usuário tenha acesso pleno ao banco de dados.


***NOTA: Neste primeiro artigo ainda não vamos utilizar o banco de dados***


## Conceitos Preliminares
### O que é hash e conhecendo SHA1?
Um hash é uma string gerada através de um algoritmo de criptografia em uma unica direção, como por exemplo MD5 ou SHA1 e representa de **forma única** seu conteúdo original, por exemplo 
```
SHA1("Rafael") = 3e05c90f8530b1ba72519824415d05e08cf5718b
```
Sempre que eu codificar a palavra "Rafael" ela vai gerar "3e05c90f8530b1ba72519824415d05e08cf5718b", porém não é possível através do hash descobrir que a palavra original era "Rafael", por isso dizemos que essa é uma **criptografia unidirecional**.
MD5 não é recomendado pois tem brechas de segurança que permitem que o mesmo conteúdo gera o mesmo hash, o que seria extremamente preocupante caso a informação de uma transação como `{To: 'Antonio da Feira', Amount: 3}` pudesse ter o mesmo identificador que `{To: 'Fundacao L...', Amount: 30000000}`, **não use MD5**.




### O que é Blockchain?
Uma blockchain ou "cadeia de blocos" é um estrutura lógica em que cada registro,"bloco", além do seu conteúdo possui um hash gerado a partir das informações contidas nele, e também o hash do registro anterior
```go
// Estrutura de um bloco
type Block struct{
  Id        int
  Data      string
  ...
  LastHash  []byte
  Hash      []byte
}`
```
* **Data**: A informação que queremos armazenar, pode ser um conjunto complexo de dados no formato JSON como `"{Amount: 1.034, Currency: 'BRL', etc...}"`
* **Hash**: Gerado utilizando SHA1 e representa todas as informações do bloco `SHA1(Id + Data + LastHash + ...)`
* **LastHash**: Se este é o bloco N, então LastHash indica o hash do bloco anterior, ou sejna, N-1


Essa referência a informação do bloco anterior que é a grande sacada da blockchain, assim mesmo que alterássemos a informação de um bloco e atualizasse-mos seu hash para parecer que nada aconteceu, o bloco seguinte utiliza esse dado para gerar seu próprio hash e seria necessário alterar o hash de todos os outros blocos da cadeia, o que devido a prova de esforço é computacionalmente impossível.
### O que é prova de esforço e mineração?
Simplesmente gerar um hash leva alguns milissegundos, então precisamos estabelecer um critério a ser atingido para que o hash seja aceito, e fazemos isso estabelecendo por exemplo que ele deva começar com uma certa quantidade de zeros, `00e5c90f8530b1ba72519824415d05e08cf5718b`.
O processo de **minerar** consiste em iteradamente modificar o conteúdo do bloco até que ele gere um hash que atende ao requisito da prova de esforço. Para que o hash mude a cada iteração introduzimos a variável `ntimes`(n-vezes), que é incrementada a cada iteração.
No nosso código controlamos essa quantidade de zeros para o hash ser aceito com a variável `dificult`, com o valor padrão 2, assim temos um tempo de processamento de alguns segundos por bloco, se utilizássemos uma dificuldade de 3 zeros levaria em torno de 10 minutos por bloco o que seria inviável por conta da frequência com que queremos gravar um registro.




## Let's Code
Acho que um bom ponto de partida é como queremos utilizar a biblioteca.


```go
func main(){  
  // Init
  var blockChain Blockchain
  blockChain.Init(2)                  // Inicia com dificuldade 2 
  
  // Cria e adiciona um bloco
  b := Block{Data: "{Action:'BUY'}"}
  blockChain.AddBlock(&b)
  //Imprime o hash de um bloco
  fmt.Println("Block hash: ", b.Hash)


  //Imprime o estatus da blockchain
  fmt.Println("Chain integrity ok?", blockChain.IsIntegre())


  // Imprime o ultimo bloco adicionado
  fmt.Println("LastBlock?", blockChain.LastBlock())
}
```




**No topo do nosso código**
Para esse projeto vamos precisar da biblioteca `crypto/sha1` para gerar os hashs, essa biblioteca retorna o hash como um array de bytes `[]byte`e por isso também vamos precisar da biblioteca `bytes` para podermos fazer comparações e gerar sequências de bytes. A biblioteca `time` é utilizada para registrar a data-hora da criação de cada bloco.
```go 
package main


import (
  "fmt"
  "bytes"
  "time"
  "crypto/sha1"
  )
```
**Banco de Dados**
Por enquanto não estamos utilizando o banco de dados, mas podemos representá-lo por um mapa em que a chave é o id e o valor é cada bloco.
```go
// Esse mapa representa nosso banco de dados
var dbChain = map[int]*Block{}
```
Referenciando cada bloco através de um ponteiro, `*Block`, podemos acessar os dados de um bloco logo após adicioná-lo à cadeia sem a necessidade de pesquisarmos ela pelo id, `dbChain[Id]`
**Estrutura de um bloco**
```go
// Estrutura de um bloco
type Block struct{
  Id        int
  Data      string
  LastHash  []byte
  Datetime  time.Time
  NTimes     uint64
  Hash      []byte
}`
```
* **Id**: Auto-incrementado a cada bloco
* **Data**: Conteúdo do bloco
* **LastHash**: Hash do bloco anterior
* **Datetime**: Data-hora de inserção do bloco
* **NTimes**: Variável utilizada 
*  **Hash**: Hash do bloco
  
**Gera o hash para o conteúdo do bloco**
A função abaixo retorna um rash do conteúdo do bloco, observe que ela considera apenas o estado atual das variáveis.
```go
// Gera o hash para o conteúdo do bloco
func (bc Block) DoHash() []byte{
  bv := []byte(string(bc.Id) + bc.Data + string(bc.LastHash) + bc.Datetime.String() + string(bc.NTimes))
  hasher := sha1.New()
  hasher.Write(bv)
  sha_bytes := hasher.Sum(nil)
  return sha_bytes
}
```
**Blockchain**
Esta estrutura irá conter os métodos para:
* Adicionar um bloco novo `Add(b *Block)`;
* Validar se os dados estão íntegros `IsIntegre()`;
* `Mine(b *Block)` que trabalha em loop até que o hash gerado atenda aos requisitos da prova de esforço.
```go
// Representa a blockchain
type Blockchain struct{
  dificult int
}
```
A variável `dificult` vai determinar o requisito da prova de esforço, ou seja, quantos caracteres `0`(zero) vamos querer no início do nosso hash.


**Bloco Gênese e Init()**
A verificação de integridade compara o conteúdo de um bloco com o bloco anterior, então o bloco zero precisa ser inserido manualmente e faremos isso dentro da função `Init()`
```go
// Inicializa adicionando o bloco de gênese
func (bc *Blockchain) Init(){
  genesisBlock := Block{Id: 1, Data: "{}", LastHash: []byte(""), Datetime: time.Date(2018, time.May, 31, 1, 30, 5, 0, time.UTC), NTimes:0}
  bc.Mine(&genesisBlock)
  dbChain[0] = &genesisBlock
}
```
**Adiciona bloco**
Monta o bloco com base nas informações do bloco anterior, incrementando o Id e recuperando o hash anterior
```go
func (bc *Blockchain) AddBlock(b *Block) int{
  lastBlock := bc.LastBlock()
  b.Id = lastBlock.Id + 1
  b.LastHash = lastBlock.Hash
  b.Datetime = time.Now()
  b.NTimes=0


  bc.Mine(b)


  i := len(dbChain)
  dbChain[i]= b
  return i
}
```
**Last Block**
Retorna o último bloco da cadeia, essa função é útil pois precisamos incrementar o  `Id` e também do hash do último bloco.
```go
// Retorna o último bloco da cadeia
func (bc *Blockchain) LastBlock() *Block{
  return dbChain[len(dbChain)-1]
}
```
**Minerar**
A cada iteração incrementa NTime `block.NTimes +=1` e verifica se o hash atende o requisito da prova de esforço, no nosso caso começar com 2 zeros.
```go
// Gera um hash baseado na dificuldade da cadeia
func (bc *Blockchain) Mine(block *Block){
  fmt.Println("Mining block...")
  var hash []byte
  for !bytes.HasPrefix( hash, bytes.Repeat([]byte("0"), bc.dificult)  ) {
    fmt.Println(block)
    block.NTimes = block.NTimes + 1
    hash = block.DoHash()
  }
  block.Hash=hash
}
```
Aqui utilizamos a função `HasPrefix` da biblioteca `bytes` para verificar se o hash começa com `00`, e a função `Repeate` para gerar uma sequência de caracteres repetidos
A linha `fmt.Println(block)` é apenas para debug e pode ser removida.


**Checa a integridade da cadeia**
A integridade da cadeia faz duas verificações:
* Se o hash do bloco está correto e representa seu conteúdo, `bytes.Compare(dbChain[i].Hash, dbChain[i].DoHash())`; e
* Se o bloco referencia corretamente o bloco anterior, `bytes.Compare(dbChain[i].LastHash, dbChain[i-1].Hash)  !=  0`.


Como não é possível verificar o hash do bloco anterior para o primeiro bloco nós pulamos ele, `ìf i>0`
Note que para comparar duas cadeia de bytes, `[]byte` utilizamos a função `Compare` da biblioteca `bytes`, não é possível fazer a verificação diretamente `dbChain[i].LastHash != dbChain[i-1].Hash` pois teríamos um erro.
A função `Compare` recebe como argumento as duas cadeias de bytes que queremos comparar e retorna 0 se forem iguais, -1 se a primeira for menos  que a segunda, ou 1 se a primeira for maior.
```go
// Checa a integridade da cadeia
func (bc *Blockchain) IsIntegre() bool{
  fmt.Println("Checking integrity...")
  for i:=1; i< len(dbChain); i++ {
    //verifica se o conteudo do bloco foi alterado
    if bytes.Compare(dbChain[i].Hash, dbChain[i].DoHash()) !=0 {
      fmt.Println("Bloco modificado após a criação, corrompido")
      return false
    }


    //verifica se o bloco referência o bloco anterior
    if i>0 {  //ignora a verificação do bloco anterior para o bloco gênese
      if bytes.Compare(dbChain[i].LastHash, dbChain[i-1].Hash) != 0{
        fmt.Println("Bloco não referencia o bloco anterior, corrompido")
        return false
      }
    }
  }
  return true
}
```


## Referências
Esse código foi baseado no vídeo *Creating a blockchain with javascript* de [Simply Explained - Savjee](https://www.youtube.com/channel/UCnxrdFPXJMeHru_b4Q_vTPQ), Publicado em 18 de jul de 2017
https://youtu.be/zVqczFZr124
