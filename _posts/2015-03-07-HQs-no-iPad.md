---
layout: post
title: Lendo HQs num caduco iPad2
---

Recentemente desenvolvi esse hábito de ler HQs. Marvel, DC, MSP, Vertigo, DarkHorse etc, o que achar na banca e for massa tá valendo. 
Porém, tem coisas que não dá para comprar por não terem edição, ou pelo preço ser abusivo. Nesse caso, meu velho companheiro de guerra iPad2 vem ao resgate.


# Participantes
## O Ipad
É um antigo iPad que comprei de um ex-chefe. Devidamente fiz um jailbreak (desculpem puristas).

## O Aplicativo
Dentre os pesquisados, o [ComicFlow](https://itunes.apple.com/br/app/comicflow/id409290355?mt=8) foi o vencedor. Bonitinho, funcional e compatível com iOS 7.1. Ah, claro, com código disponibilizado no [github](https://github.com/swisspol/ComicFlow).


## O problema
![Interface web do ComicFlow]({{ site.url }}/images/07_03_2015/comicFlow-webServer.png)

O ComicFlow tem uma ótima interface via web para enviar as HQs para o iPad. Porém vem com uma limitação de 50 arquivos, depois disso só comprando uma versão Premium Deluxe. Acho overkill resolver todos os problemas pagando, então resolvi lidar com o problema de outra forma.

# A solução

Primeiramente, instalei o OpenSSH via Cydia (se você fez jailbreak, saberá do que estou falando). 

Na minha máquina instalei um cliente de FTP, o Filezilla. Cuidado de onde forem baixar, pois do [SoundForge vem com malware](http://trac.filezilla-project.org/ticket/8888).

Certo. No iPad vou em Wi-fi para descobrir o IP dele na rede.

No cliente de FTP, coloco:
{% highlight bash %}
Host: IP do ipad
usuário: root
senha: alpine
{% endhighlight %}

Essa senha "alpine" é a padrão do root do iOS. 

Certo... como funciona essa árvore de diretórios?
![Árvore diretórios iPad]({{ site.url }}/images/07_03_2015/ipadtree.png)

Ok, não faço ideia de onde esteja a pasta do ComicFlow. Então abro uma janela no terminal e acesso via SSH o iPad:
```
ssh root@IP_DO_IPAD
```
Lembrando, a senha é ```alpine```.

Então uso uma velha ferramentinha do Unix para achar onde está a pasta do ComicFlow.
{% highlight bash %}
 find / -type d -name "ComicFlow*"
{% endhighlight %}

Algo como "por favor me ache um diretório que seu nome seja ComicFlow". Eu obtive o resultado:

{% highlight bash %}
/private/var/mobile/Applications/2F165C88-821E-46C4-8DAF-9EC50C47F90A/ComicFlow.app/ComicFlow
/private/var/mobile/Applications/2F165C88-821E-46C4-8DAF-9EC50C47F90A/ComicFlow.app/SC_Info/ComicFlow.sinf
/private/var/mobile/Applications/2F165C88-821E-46C4-8DAF-9EC50C47F90A/ComicFlow.app/SC_Info/ComicFlow.supp
{% endhighlight %}

Certo, já sei o caminho, agora voltando ao FTP é só encontrar a pasta e enviar seus arquivos em /Documents.

## Modo Bônus
Tá, agora preciso ficar toda vez rodando o find para achar a pasta do ComicFlow?
Por isso, criei um link simbólico para a pasta do ComicFlow, assim é só acessar em /ComicFlow!
{% highlight bash %}
ln -s /private/var/mobile/Applications/2F165C88-821E-46C4-8DAF-9EC50C47F90A/Documents/
{% endhighlight %}
## Modo 103%
Se achar que FTP gráfico é overkill, dá para instalar o rsync via Cydia e passar os arquivos via terminal mesmo:

{% highlight bash %}
rsync NOME_DO_ARQUIVO root@IP_DO_IPAD:/ComicFlow/
{% endhighlight %}


Agora é só ler um Animal Man da fase do Grant Morrisson :-)
![Animal Man]({{ site.url }}/images/07_03_2015/animalMan.PNG)



