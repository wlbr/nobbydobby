#!/bin/sh


## declare an array variable
declare -a arr=('{"firstname": "Uwe", "lastname": "Tube", "email": "uwe@tube.de"}'
    '{"firstname": "Bernd", "lastname": "das Brot", "email": "bernd@brot.de"}'
    '{"firstname": "Anna", "lastname": "Blume", "email": "anna@blume.de"}'
    '{"firstname": "Alice", "lastname": "Example", "email": "alice@example.org"}'
    '{"firstname": "Bob", "lastname": "Example", "email": "bob@example.org"}'
    '{"firstname": "Charlie", "lastname": "Example", "email": "charlie@example.org"}'
    '{"firstname": "Dave", "lastname": "Dummy", "email": "dave@dummy.org"}'
    '{"firstname": "Aragorn", "lastname": "Arathornson", "email": "Aragorn@dunedain.org"}'
    '{"firstname": "Gimli", "lastname": "Gloinsson", "email": "Gimli@erebor.com"}'
    '{"firstname": "Frodo", "lastname": "Baggins", "email": "frodo@beutlin.name"}'
    '{"firstname": "Samwise", "lastname": "Gamgee", "email": "samwise@gamgee.org"}'
    '{"firstname": "Pippin", "lastname": "Took", "email": "pippin@took.org"}'
    '{"firstname": "Merry", "lastname": "Brandybuck", "email": "merry@buck.org"}'
    '{"firstname": "Legolas", "lastname": "Greenleaf", "email": "legolas@sindar.org"}'
    )

## loop through above array (quotes are important if your elements may contain spaces)
for i in "${arr[@]}"
do
   curl -X POST 'http://localhost:8080/register' -d "$i"
done

