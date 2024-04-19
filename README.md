# Full Cycle pos-graduation 

Access the course notes in [./Notes.org](./Notes.org) 

## Challenge 2

``` sh
git clone https://github.com/BuddhiLW/pos-go-expert
cd pos-go-expert/desafios/multithreading

go run main.go
```

Test the endpoint multiple times, in another terminal:

``` sh
curl http://localhost:8001/cep?cep=17800-970
```

This should give, something like this (after hitting multiple curls to that `cep`):

``` sh
brasilAPI || Estado: SP, Cidade: Adamantina, Bairro: Centro, Rua: Rua Deputado Salles Filho, 469
ViaCEP || Estado: SP, Cidade: Adamantina, Bairro: Centro, Rua: Rua Deputado Salles Filho
brasilAPI || Estado: SP, Cidade: Adamantina, Bairro: Centro, Rua: Rua Deputado Salles Filho, 469
brasilAPI || Estado: SP, Cidade: Adamantina, Bairro: Centro, Rua: Rua Deputado Salles Filho, 469
ViaCEP || Estado: SP, Cidade: Adamantina, Bairro: Centro, Rua: Rua Deputado Salles Filho
brasilAPI || Estado: SP, Cidade: Adamantina, Bairro: Centro, Rua: Rua Deputado Salles Filho, 469
ViaCEP || Estado: SP, Cidade: Adamantina, Bairro: Centro, Rua: Rua Deputado Salles Filho
brasilAPI || Estado: SP, Cidade: Adamantina, Bairro: Centro, Rua: Rua Deputado Salles Filho, 469
brasilAPI || Estado: SP, Cidade: Adamantina, Bairro: Centro, Rua: Rua Deputado Salles Filho, 469
```

