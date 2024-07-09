+++ 
date = 2024-07-09T23:04:13+01:00
title = "deployando site estático no S3 pt1: CDK's BucketDeployment sucks"
categories = ["technical"]
lang = "pt-br"
+++

Você tem lá um CDK para um site estático, s3, cloudfront, cert manager etc. Agora precisa fazer um deploy, e naturalmente busca uma solução na mesma vibe, o [`BucketDeployment`](https://docs.aws.amazon.com/cdk/api/v2/docs/aws-cdk-lib.aws_s3_deployment.BucketDeployment.html).

E até que funciona, faz até [invalidação com CloudFront!](https://github.com/aws/aws-cdk/blob/e5740c01a5f524b099258820b3206045873b6732/packages/aws-cdk-lib/aws-s3-deployment/lib/bucket-deployment.ts#L95-L101)

Porém chega uma hora que você percebe: isso leva mais tempo que deveria.

Aí olhando no código (e explorando o CloudFormation no console) percebe que faz muito mais do que você imaginava, cria um CRD, que cria (dentre outras coisas) um bucket e um lambda, que é o que faz efetivamente o trabalho de copiar.

Então se copia 1 vez para esse s3 temporário, e depois para o bucket correto. Para mais infos ver [o código do lambda](https://github.com/ohde/aws-cdk/blob/9ae6b1161d86016cfa459f25a295c0ee4fbb4343/packages/%40aws-cdk/aws-s3-deployment/lib/lambda/index.py).

O que é claramente um desperdício de tempo.

A solução é simplesmente usar a cli, usando `aws s3 sync` para copiar os arquivos se você quiser (e deletar os que não quer).

Ficaria algo como

```sh
aws s3 sync out s3://$SEU_BUCKET --delete
```

Se quiser checar antes de commitar com o deploy (recomendado), rode com `--dryrun` antes, ele vai listar os arquivos como por exemplo:

```
(dryrun) upload: out/_next/static/3BRZtz_wC32BCqy0h8gHL/_buildManifest.js to s3://$BUCKET_NAME/_next/static/3BRZtz_wC32BCqy0h8gHL/_buildManifest.js
(dryrun) delete: s3://$BUCKET_NAME/_next/static/99ixjLUK2BVi2N5LwU41p/_buildManifest.js
```

# yay!!?! not so fast
Tudo as mil maravilhas, porém estranhei que a lista de arquivos era sempre muito maior. Não mudei nada no meu site estático, por que está querendo subir todos os arquivos?!


Bom, [aparentemente eles checam o timestamp e o tamanho do arquivo](https://stackoverflow.com/questions/43529926/how-does-aws-s3-sync-determine-if-a-file-has-been-updated). Se você está fazendo um build novo, o timestamp dos arquivos também vai ser novo. Não sei se tem um jeito de falar pra ferramenta de build pra não mexer nos arquivos novos, e aí em CI cachear esses arquivos. E se tem, não sei se quero seguir nessa rota.

Uma outra sugestão é falar para ignorar timestamps e usar apenas o tamanho do arquivo, com a flag `--size-only`.

Ah, mas se eu trocar uma string de "CAVALO" para "CAMELO"? Não vai funcionar, sinto muito.

Além desse exemplo trivial, ainda tem o caso [levantado no Github onde](https://github.com/aws/aws-cli/issues/5216#issuecomment-1722248324) o seu index.html não muda de tamanho, especialmente no caso de SPA onde o html é o mesmo, só mudando o hash dos javascript inclusos.

Nesse caso acho mais fácil simplesmente dar um `aws s3 cp` no `index.html`.

# checksum no s3 sync? não pra você
Idealmente seria possível usar um checksum como gerado via MD5 para checar o contéudo, é uma [issue que existe desde 2014](https://github.com/aws/aws-cli/issues/599) (FELIZ BODAS DE ESTANHO!).

A implementação existe em [layers inferiores](https://github.com/aws/aws-cli/issues/6750), mas não chegou ainda na `aws s3`, e sabe-se-lá quando vai chegar.

Me perguntei como ferramentas como hugo fazendo deploy. Aparentemente [segundo a documentação](https://github.com/gohugoio/hugo/blob/439f07eac4706eb11fcaea259f04b3a4e4493fa1/docs/content/en/hosting-and-deployment/hugo-deploy.md?plain=1#L107-L110) o hugo checa o md5 dos arquivos, [aqui o código](https://github.com/gohugoio/hugo/blob/439f07eac4706eb11fcaea259f04b3a4e4493fa1/deploy/deploy.go#L692)


"Ah mas como pega o MD5? Tem que baixar o arquivo?". Não, tal como dito anteriormente e segundo [esse comentário no código do hugo](https://github.com/gohugoio/hugo/blob/439f07eac4706eb11fcaea259f04b3a4e4493fa1/deploy/deploy.go#L584-L591), S3 já retorna isso.

O que me faz pensar que deve ser mais fácil só usar a cli do hugo para deploys :^)

# bônus
Ainda tem um grande detalhe não necessariamente relacionado que não falei: você não pode só usar `--delete` e seguir com seu dia, já que possivelmente há alguém com o site aberto, com o manifest (do nextjs no caso) apontando para certos arquivos em caso de lazy load, que por sua vez não vão existir, quebrando a navegação do pobre coitado!

Então idealmente a solução seria:
1. Vê os arquivos que vão ser deletados
2. Muda a expiração deles pra sei lá, daqui 1 semana
3. Copia pro bucket os arquivos novos (`s3 sync`)
4. Muda o seu código/do framework para quando carregar um chunk inexistente, refreshar a página


BREAKING NEWS: fui informado que o NextJS já faz isso, mas não chequei. Trust but don't verify.


Inspiração: [Stop Using the CDK S3 Deployment Module | by Riccardo Giorato | AWS in Plain English](https://aws.plainenglish.io/stop-using-the-cdk-s3-deployment-module-4f31f36b4f21).
