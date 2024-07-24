+++ 
draft = false
date = 2024-07-24T17:19:29+01:00
title = "deployando site estático no S3 pt2: github actions não é seu amigo"
description = ""
slug = "" 
tags = []
categories = ["technical"]
+++

Esta é a parte 2, sem hard requirement na parte 1. Entretanto, se ainda quiser ver a parte 1, onde lido com AWS CDK [clique aqui](../s3-bucket-deployment-sucks).


# intro
Então, \*dá uma boa respirada\*, sabe aquela pipeline no Github Actions? Ela poderia ser melhor.

Melhor? Eu quis dizer **mais barata**. Se você usa apenas repositórios públicos - a paz de Cristo para você, te vejo na missa semana que vem - porque aparentemente os runners não [tem limitação de tempo](https://github.com/orgs/community/discussions/70492#discussioncomment-7362186) e nem de [artefatos](https://github.com/orgs/community/discussions/26438#discussioncomment-3251931).

Se procurarmos por [NEXTJS GITHUB ACTIONS DEPLOYMENT GOOGLE PESQUISAR](https://github.com/actions/starter-workflows/blob/main/pages/nextjs.yml) e afins, geralmente os Workflows são quebrados em 2 partes (removi as partes desnecessárias):

```yaml
jobs:
  # Build job
  build:
    runs-on: ubuntu-latest
    steps:
    (...)
      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: ./out

  # Deployment job
  deploy:
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```


O detalhe aqui é que o workflow é quebrado em 2: um faz build e sobe o artefato, e outro "baixa" esse artefato e o deploya.

O exemplo acima deploya direto para o Github Actions então talvez não seja claro, [esse aqui](https://gist.github.com/smcelhinney/9433bef18088b08e17a7349c454508f8) explicitamente sobe e baixa o artefato.

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    (...)
      - name: Upload artifacts
        uses: actions/upload-artifact@v1
        with:
          name: build_dir
          path: ./out
  publish-s3:
    needs: build
    runs-on: ubuntu-latest
    steps:
    - uses: actions/download-artifact@master
      with:
        name: build_dir
        path: ./out
```

# ineficácias iniciais

Então vamos analisar o que acontece. Obviamente pulando vários passos e simplificando.

1. Primeiro tem que achar um runner para o job `build`, que é rodado.
2. No final, ele zipa (tar?) o artefato etc.
3. Agora que job `build` terminou, acha um runner para rodar o job `publish-s3`.
4. Esse job tem que BAIXAR o zip que você acabou de dar upload, e aí faz o deploy e yadda yadda.

Visivelmente há 2 ineficácias, que na prática são a mesma:
1. Tem que esperar 2 runners serem scheduled, o que leva tempo.
2. Precisa zipar e fazer upload do artefato, e logo em seguida baixá-lo.

**Claramente duas ineficácias que podem ser resolvidas rodando tudo num job só :^)**


# outros problemas ensaboados

Outros problema que já toquei acima é que estes artefatos custam DIÑERO (para repos privados). Aí pra todo santo build (caso faça deploys, ou *dry-runs* para PRs) você vai lá e salva mais um artefato.

No meu caso era ainda pior, porque rodava uma [ferramenta para otimizar imagens para diferentes resoluções](https://next-export-optimize-images.vercel.app/), e o site possuía diversas imagens, chegando a dar uns 70MB cada build.

E pra piorar, você não precisa desse artefato depois do deploy. E não tem um jeito fácil de apagá-lo. Até dá para [apagar usando a API](https://github.com/actions/upload-artifact/issues/550#issuecomment-2059919616), mas eu não consegui e nem tive saco de debuggar.

Outra opção é nas próprias configurações do repositório (Actions -> General) colocar uma data de expiração de logs e artefatos. O downside é que mudou um, mudou o outro. É a solução mais fácil, só uma questão de achar um range bom para a retenção.


OOOOOOOOOutro problema é ilustrado por este artigo: [The Hidden Cost of Parallel Processing in GitHub Actions | by Wenqi Glantz | Better Programming](https://betterprogramming.pub/the-hidden-cost-of-parallel-processing-in-github-actions-63f25b2d5f6a). Que em resumo, o Github arredonda os tempos gastos de cada runner:

>GitHub rounds the minutes and partial minutes each job uses up to the nearest whole minute.

Fonte: [About billing for GitHub Actions - GitHub Docs](https://docs.github.com/en/billing/managing-billing-for-github-actions/about-billing-for-github-actions#about-billing-for-github-actions)


Então se você tem um Workflow que tem 2 jobs. Um job leva 1m10s e outro leva 1m05s.
O total será 4m, porque arredonda cada um deles para 2 e soma. Ao invés de SOMAR e depois arredondar. [Para mais infos ver essa discussão](https://github.com/orgs/community/discussions/8726).

It's always in the details :)

A solução, que já ficou clara, é usar "monolithic workflows". Workflows com só um job que faz tudo.

# múltiplos jobs, por quê?

Tá, mas por que tudo isso? Por que geralmente se quebram em vários jobs? Boa pergunta, umas das explicacões que encontrei é que caso queira rodar apenas um pedaço dele, é mais fácil e rápido. Por exemplo, o deploy falhou, mas o build em si está perfeito, então só precisa tentar fazer o deploy novamente.

Além é claro de paralelização, talvez você consiga rodar em paralelo testes + lint + build? Mas aí é outra investigação, no momento tendo a fazer tudo sequencialmente.



# Em resumo
1. Apague artefatos via API ou configure seu repositório para expirar os artefatos
2. Escreva o máximo possível num job, e só quebre quando for muito óbvio os ganhos de paralelização. Por exemplo, rodando a mesma pipeline num outro OS.
3. É no detalhe, no fine print, que o diabo trabalha.
