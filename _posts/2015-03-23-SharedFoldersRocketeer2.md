---
layout: post
title: Shared Folders no Rocketeer2
---


Mais uma rápida. Usando o Rocketeer para Deploy em aplicações Laravel.

## O Objetivo
Manter o conteúdo de um diretório o mesmo para diferentes versões do seu deploy. Ideal para uploads, onde você não quer que seja sobescrito entre um deploy e outro.

Basicamente dentro do seu projeto terá uma pasta shared/ onde conterá seus arquivos/diretórios compartilhados, e então dentro da sua aplicação você terá [links simbólicos](http://en.wikipedia.org/wiki/Symbolic_link) para estes diretórios/arquivos no "shared/"


## Requisitos
Supõe-se que você já esteja conseguindo deployar e agora quer configurar pastas compartilhadas.


## Motivação pessoal
No meu caso, após instalar o [stapler](https://github.com/CodeSleeve/laravel-stapler), precisei compartilhar o diretório ```public/system```, pois é onde o Stapler sobe os arquivos.


## *Let's go, bub*
Primeiro, no remote.php do Rocketeer, adiciono o diretório que quero compartilhar no array 'shared'

{% highlight php %}
    'shared'         => [
        '{path.storage}/logs',
        '{path.storage}/sessions',
        '{path.public}/system'
    ],
{% endhighlight %}

Agora vem um pequeno "quirk". É preciso ter o diretório desejado sob versionamento no git, mas não quero que o seu conteúdo seja trackeado pelo git. Então para criar um diretório sempre vazio sob "as asas" do git, é só criar um .gitignore dentro desse seu diretório, com o seguinte conteúdo:

{% highlight bash %}
!.gitignore
{% endhighlight %}

Para ignorar tudo no diretório, exceto esse arquivo (o .gitignore).


(Créditos: [StackOverflow](http://stackoverflow.com/questions/115983/how-can-i-add-an-empty-directory-to-a-git-repository))

Aí é só commitar, mandar para o seu repositório externo (github/bitbucket provavelmente) e mandar um "php artisan deploy".


Só para ter uma noção do que é para acontecer:

{% highlight bash %}
|=> Setting permissions for MEU_CAMINHO/releases/20150323144405/public
$ cd MEU_CAMINHO/releases/20150323144405
$ chmod -R 755 MEU_CAMINHO/releases/20150323144405/public
$ chmod -R g+s MEU_CAMINHO/releases/20150323144405/public
$ chown -R USER:USER MEU_CAMINHO/releases/20150323144405/public

|=> Sharing file MEU_CAMINHO/releases/20150323144405/public/system
$ rm -rf MEU_CAMINHO/releases/20150323144405/public/system
$ ln -s MEU_CAMINHO/shared/public/system MEU_CAMINHO/releases/20150323144405/public/system-temp
$ mv -Tf MEU_CAMINHO/releases/20150323144405/public/system-temp MEU_CAMINHO/releases/20150323144405/public/system
{% endhighlight %}

Resumindo, linkou a pasta compartilhada do novo release com a pasta compartilhada do shared/


E é isso.
