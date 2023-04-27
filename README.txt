------------------------------------------------------
# Intruções para execução do projeto
------------------------------------------------------

Foi utilizado Golang `go1.19.4 linux/amd64` para o desenvolvimento do
projeto.

RUN_TESTS:        `./run.sh`
RUN_SINGLE_TEST:  `./main examples/test.yail`
REPL:             `./repl.sh`
BUILD:            `go build ./main.go`

Os testes de exemplo estão na pasta examples/ e o executável main
pode ser executado com um ficheiro de testes como argumento.

Em ./examples/ todos os ficheiros que tenham:

*.yail.out)   Contêm a AST, possíveis Errors e Tokens do Lexer.
*.yail)       Contêm o código fonte da linguagem YAIL.
*.error.yail) Contêm o código fonte da linguagem YAIL com erros.

------------------------------------------------------
# Notas Importantes
------------------------------------------------------

Q) É obrigatório instalar o golang para executar o programa?

Não, assumindo que o seu OS é Linux. No entanto, caso, por algum motivo
seja necessário, deixo aqui o script para a sua instalação:

./install-go.sh

Importante ver o script e altere de acordo o seu OS.

Notas:

1) Caso queira executar o programa com um ficheiro de testes, basta
   executar o executável main com o caminho para o ficheiro de testes
   como argumento. Exemplo:

   ./main examples/test.yail

2) Caso exista algum erro no executável devido a incompatibilidades
   (assumindo que o OS é linux), terá que instalar o golang e executar
   o seguinte comando:

    # Build for Linux (amd64)
    GOOS=linux GOARCH=amd64 go build ./main.go

    # Build for macOS (amd64)
    GOOS=darwin GOARCH=amd64 go build ./main.go

    # Vou assumir que não é necessário fazer build para windows, caso
    sim:
    GOOS=windows GOARCH=amd64 go build ./main.go (os scripts em bash deixam
    de funcionar)...

    Pode passar na mesma os testes de exemplo como argumento, como dito
    no passo 1), ou executar o REPL

------------------------------------------------------
# Estrutura do projeto
------------------------------------------------------

+ast/         - Código relativamente à AST
+examples/    - Testes de exemplo para a linguagem
+lexer/       - Código relativamente ao lexer da linguagem
+parser/      - Código relativamente ao parser da linguagem
+repl/        - Código relativamente ao REPL da linguagem
+token/       - Código relativamente aos tokens da linguagem
  go.mod
  main
  main.go
  README.txt
  repl.sh
  run.sh

------------------------------------------------------
# Dependencias
------------------------------------------------------

- golang
- find
- bash
