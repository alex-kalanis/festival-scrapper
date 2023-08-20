Task: Aktuální kapela na fesťáku

Jazyk: GO

Wtf:

- máme léto a fesťáky
- na nich hrají kapely
- fesťáky mají na stránkách rozpisy hrajících
- takže to sosnout a říct, kdo hraje teď a kdo bude
- interface na různé fesťáky - každý to má jinak

- stavy: aktuálně hraje / díra; další kapela / už nikdo
- potřeba: síť, HTML parser, porovnání času, output (napoprvé do CLI jenom dumpem)

Je to můj první program v GO, takže hodnotit s klidem...

Jak to běží:

```bash
$ go run scrapper.go --day=18 --hour=15 --min=45 --bands=5
```

Parametry jsou volitelné, pokud nic nepřijde, tak použije aktuální datum a čas a maximálně 2 kapely na aktivní stage
