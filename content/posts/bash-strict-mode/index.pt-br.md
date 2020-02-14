+++
date = "2020-02-08"
title = '"Bash Strict Mode" não-oficial'
slug = "bash-strict-mode"
categories = ["blog", "bash", "técnico"]
+++

Bash (e shell scripting no geral) não é fácil. É comum cometer equivocos se você não sabe o que está fazendo. Se você vem de um background tradicional de programação e só está interessado em juntar algumas linhas de código, haverão alguns comportamentos da linguagem que vão te confundir.

Para ajudar com essa, o [bash strict mode não-oficial](http://redsymbol.net/articles/unofficial-bash-strict-mode/) foi criado. Nesse post, passaremos por comportamentos confusos e como o modo estrito pode ajudar em cada caso (e suas manhas).

## errexit

### errexit: O problema

Dado o seguinte bash script:
{{< file "/posts/bash-strict-mode/errexit.sh" "bash" >}}

Assuma que o arquivo não exista. O script deve rodar por completo, ou deve falhar?

Pois o script roda normalmente!

Para abordar o problema, eu tendo ao princípio do "[Falhe Rápido](https://wiki.c2.com/?FailFast)":

(numa tradução direta)

> Isto é feito ao encontrar um erro tão sério que é possível que o estado do processo está corrupto ou inconsistente, e abortar imediatamente é a melhor maneira que nenhum dano adicional seja feito.

Parece bom. Então por que "Falhe quieto" o comportamento normal num shell script? Bem, pense que no contexto de um shell você NÃO quer abortar quando há um erro (imagine crashar o shell quando você dá um `cat` num arquivo que não existe). Me parece que esse comportamento simplesmente foi levado pra um shell non-interativo.

### errexit: A solução

Como podemos melhorar esse comportamento? Setando a flag `errexit`.

```sh
set -o errexit
```

ou a versão mais curta (e normalmente mais usada):

```sh
set -e
```

O que ela faz? A documentação diz que:

> Aborta imediatamente se uma pipeline (...) retorna um status não-zero.
>
> -- https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html#The-Set-Builtin

Voltando ao nosso example, nos faríamos assim:
{{< file "/posts/bash-strict-mode/errexit2.sh" "bash" >}}

Que por sua vez, falharia. Já que o arquivo não existe, `cat` retorna um exit code não-zero. Esse comportamento é descrito pelo seguinte
teste unitário [bats](https://github.com/sstephenson/bats):

{{< file "/posts/bash-strict-mode/errexit.bats" "bash" >}}

Isso funciona e vai definitivamente ([na minha opinião](http://mywiki.wooledge.org/BashFAQ/105))) te ajudar, mas fique esperto:

### Errexit: Manhas

#### Manha #1: Programas que retornam status não-zero

Nem todos os comandos retornam 0 quando rodam corretamente. Talvez o exemplo mais famoso seja `grep`. A documentação diz que:

> o código de saída é
>
> 1. 0 se a linha foi 'selecionada',
> 2. 1 se nenhuma linha foi selecionada,
> 3. e 2 se algum erro aconteceu.
>
> Porém, se -q ou --quiet ou --silent é usada e a linha é selecionada, o código de retorno é 0 mesmo que um erro tenha acontecido.

Portanto, no exemplo abaixo, `echo` nunca rodará.
{{< file "/posts/bash-strict-mode/grep_fail.sh" "bash" >}}

O que podemos fazer nesta situação? HÁ um pedaço no manual do bash na sessão de `errexit` que pode nos ajudar (formatada para claridade):

> O shell não aborta se o comando que falha é
>
> 1. parte de uma lista de comandos imediatamente seguida por um while ou until
> 2. parte de um teste numa declaração if
> 3. parte de algum comando executado num && ou || exceto o comando depois do último && ou ||
> 4. qualquer comando numa pipeline exceto o último, ou se o status de retorno é invertido com uma ! (negação)

No nosso caso, podemos simplesmente reescrever para cumpri com 2:

{{< file "/posts/bash-strict-mode/grep_correct.sh" "bash" >}}

Esse comportamento pode ser validado pelo seguinte teste bats:

{{< file "/posts/bash-strict-mode/grep.bats" "bash" >}}

#### Manha 2: E se você está ok com um comando falhando/retornando não-zero?

Nesse caso, simplesmente coloque um OR com `true`:

```sh
rm *.log || true
```

Já que (no exemplo acima) não queremos falhar caso não existam arquivos de log.

Vamos pensar um pouco porque isso funciona. A documentação (que já lemos num ponto anterior) diz:

> O shell não aborta se o comando que falha é
> (...)
>
> 3. part of any command executed in a && or || list except the command following the final && or ||

Como o comando seguindo o último `||` é `true`, então não há como a linha inteira falhar.

Uma outra opção seria desligar a flagg momentaneamente:

```sh
set +e
command_allowed_to_fail
set -e
```

A sintaxe `+` significa "remover" e `-` significa "adicionar" (vai entender)> Portanto, estamos simplesmente desabilitando essa feature entre nossas chamadas para `command_allowed_to_fail`!

#### Ponto bônus: Como sei qual comando falhou?

Não é uma manha específica ao `errexit`, mas muitas vezes precisamos descobrir onde um comando falhou.

{{< file "/posts/bash-strict-mode/which_command_failed.sh" "bash" >}}

Como podemos dizer qual comando é o problemático? (Ignore o erro óbvio)

- 1. `echo` tudo que você está fazendo \
     **Pros**: fácil e direto \
     **Cons**: entediante

- 2. `set -x`, que printará toda linha \
     **Pros**: simples de adicionar \
     **Cons**: você pode acabar expondo mais do que gostaria (imagina printando uma variável com alguma chave secreta, bem comum em ambiente de CI)

- 3. colocar uma `trap` para [printar o número da linha em que o comando falhou](https://intoli.com/blog/exit-on-errors-in-bash-scripts/) \
     **Pros**: pode ser adicionada globalmente \
     **Cons**: um pouco verboso

#### Mais alguma coisa?

Sim. Uma vez que você pegue o jeito, leia essa página [BashFAQ](http://mywiki.wooledge.org/BashFAQ/105) e as páginas linkadas.

## Pipefail

### Pipefail: O problema

{{< file "/posts/bash-strict-mode/pipefail_first.sh" "bash" >}}
Isso rodará normalmente!

Infelizmente `errexit` não é suficiente para nos salvar aqui. A documentação diz, de novo:

> O shell não aborta se o comando que falha é
> (...)
>
> 4. qualquer comando numa pipeline exceto o último, ou se o status de retorno é invertido com uma ! (negação)

> O código de retorno de uma pipeline é o código do último comando da pipeline

### Pipefail: A solução

Vamos setar `pipefail`:

> Se pipefail estiver habilitado, o código de retorno da pipeline é o value do último (à direita) comando a retornar com um status de não-zero, ou zero se todos os comandos rodarem com sucesso.

Em outras palavras, só vai retornar 0 se todas as partes da pipeline retornarem 0.

Ao contrário de `errexit`, `pipefail` só pode ser setado na sua forma completa:

```
set -o pipefail
```

Vamos arrumar o exemplo mostrado anterior:
{{< file "/posts/bash-strict-mode/pipefail_first_correct.sh" "bash" >}}

Ambos os comportamentos podem ser validados pelo seguinte teste bats:
{{< file "/posts/bash-strict-mode/pipefail_first.bats" "bash" >}}

### Pipefail: Manhas

#### Manha 1:

> o código de retorno da pipeline é o value do último (à direita) comando a retornar com um status de não-zero

{{< file "/posts/bash-strict-mode/pipefail_quirk_1.sh" "bash" >}}

O código de retorno de `cat` é `1` quando o arquivo não existe. E o do `xargs` é `123` "se qualquer invocação de um comando retorna com status 1-12".
Obviamente ambas as chamadas estão quebradas, mas qual código de retorno recebemos aqui?

A resposta é `123`, o que não é ideal.

Minha recomendação para este caso é simplesmente quebrar em diferentes instruções:

{{< file "/posts/bash-strict-mode/pipefail_quirk_1_correct.sh" "bash" >}}

Este comportamento pode ser confirmado pelo seguinte teste bats:

{{< file "/posts/bash-strict-mode/pipefail_quirk_1.bats" "bash" >}}

#### Manha 2:

Tome cuidado com o que você coloca num pipe:

{{< file "/posts/bash-strict-mode/pipefail_quirk_2.sh" "bash" >}}

Neste exemplo, estamos carregando um arquivo whitelist, passando para outro comando (aqui implementado como uma função) que passa para outro serviço (imagine uma CLI). Mesmo que o arquivo não exista, a pipeline não falha. Acabamos passando uma string vazia para `remove_hosts` o que poderia ter efeitos catrastróficos! (Nem tanto, só deletando mais do que você espera)

Idelamente, você quer falhar o mais cedo possível. O melhor a fazer é quebrar em instruções menores e ... just seja mais cuidadoso 4Head ¯\_(ツ)\_/¯

{{< file "/posts/bash-strict-mode/pipefail_quirk_2_correct.sh" "bash" >}}

Como sempre, o comportamento pode ser validado por:
{{< file "/posts/bash-strict-mode/pipefail_quirk_2.bats" "bash" >}}

Para mais exemplos, veja [Examples of why pipefail is really important to use](https://gist.github.com/yoramvandevelde/fb4d8aa6fa6d8b1eab6da81b62373d85).

## nounset

Por último mas não menos importante, este aqui é bem direto.

### nounset: O problema

{{< file "/posts/bash-strict-mode/nounset.sh" "bash" >}}

### nounset: A solução

> Trata variáveis e parâmetros não declarados que não os parâmetros especiais '@' e '\*' como erro quando realizando expansão de parâmetro. Uma mensagem de erro seja escrita na saída de erros, e um shell não interativo abortará.

{{< file "/posts/bash-strict-mode/nounset.bats" "bash" >}}

### nounset: Manhas

#### Manha 1:

Como mencionada na documentação , `@` and `*` são tratadas diferentemente:

{{< file "/posts/bash-strict-mode/nounset_quirk_1.sh" "bash" >}}

Então sempre verifique os argumentos você está recebendo são corretos:
{{< file "/posts/bash-strict-mode/nounset_quirk_1_correct.sh" "bash" >}}

{{< file "/posts/bash-strict-mode/nounset_quirk_1.bats" "bash" >}}

## Conclusão

Espero que este post seja suficiente para:

- mostrar como as expectativas que temos muitas vezes não são verdade;
- como o 'strict mode não-oficial' pode ajudar;
- e como o strict mode não é uma panacéia!

## Referências/Leituras recomendadas

- [The Let It Crash Philosophy Outside Erlang](http://stratus3d.com/blog/2020/01/20/applying-the-let-it-crash-philosophy-outside-erlang/)
- [Bash Strict Mode](http://redsymbol.net/articles/unofficial-bash-strict-mode/)
- [Dacav's Home: More shell inconsistencies](https://dacav.roundhousecode.com/blog/2019-10/02-more-shell-inconsistencies.html)
- [How to Exit When Errors Occur in Bash Scripts](https://intoli.com/blog/exit-on-errors-in-bash-scripts/)
- [Bash Errexit Inconsistency - Stratus3D](http://stratus3d.com/blog/2019/11/29/bash-errexit-inconsistency/)
- [Examples of why pipefail is really important to use](https://gist.github.com/yoramvandevelde/fb4d8aa6fa6d8b1eab6da81b62373d85)
