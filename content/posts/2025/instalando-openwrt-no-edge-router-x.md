+++ 
date = 2025-05-19T11:22:08+01:00
title = "Instalando openwrt no Edge Router X"
categories = ["technical"]
+++

Primeiro de conversa, fui tentar acessar o annas-archive e apareceu bloqueado.
Easy, basta mudar o DNS no roteador para usar o da google, CloudFare etc.

Descobri então, que o roteador da Woo (marca low budget da Nos) bloqueia as mudanças de DNS!

O que é engraçado, porque eu justamente escolhi a Woo porque tinha suporte à bridge mode,
o que insinua que é mais "aberto" do que outros provedores, como a Vodafone, meu antigo,
que não suportava bridge mode.

Ok, é possível mudar o DNS em cada device, mas dá trabalho né?

Então decidi tirar do armário um Edge Router X que havia comprado alguns anos, e usá-lo como roteador principal da rede domiciliar.

Para todos os efeitos, o Edge Router X (abreviado como `erx`) estava brickado. Talvez tenha tentado fazer algum update errado ou algo assim, então começaremos por aqui.

# Reinstalando o software de fábrica
Ele tem um sisteminha bem inteligente chamado de [TFTP Recovery](https://help.uisp.com/hc/en-us/articles/22591244564887-EdgeRouter-TFTP-Recovery), que basicamente é
um servidor TFTP (FTP só que via UDP), no qual você consegue subir uma imagem de recovery.

Primeiro, baixamos a imagem correta do device, no meu caso `ER-e50.recovery.v2.0.6.5208541.190708.0508.16de5fdde.img.signed`.

Ligamos um cabo de rede diretamente na **porta 0** do erx.

Para ativar o TFTP Recovery, desliga o roteador, depois segurando o botão de reset (com um clips ou algo similar), liga a energia e
observa as luzes até que estejam assim:

![](https://help.uisp.com/hc/article_attachments/22591244560535)

Aí já pode soltar o botão de reset.

Agora, o servidor TFTP estará disponível na porta 69 do ip 192.168.1.20. Não me recordo se ele tem um servidor DHCP nele,
logo talvez seja necessário na sua máquina colocar um ip fixo, qualquer coisa na casa dos 192.168.1.x funciona.

No terminal, primeiro rode `tftp` para abrir o cliente de tftp. Então rode

```bash
tftp> connect 192.168.1.20
tftp> binary
tftp> put ER-e50.recovery.v2.0.6.5208541.190708.0508.16de5fdde.img.signed
```

E aguarde, após algum tempo (creio que as luzes do erx mudam quando o processo esteja terminado), e conseguirá acessar a interface via http://192.168.1.1,
caso esteja usando ip fixo, não se esqueca de mudar o gateway para 192.168.1.1


# Instalando openwrt
A tal da EdgeOS, a OS padrão do Edge Router X não me anima muito. Sei lá, achei complicado e prefiro instalar o openwrt, que é open source e managed pela comunidade.

Os passos aqui são um pouco chatos, então vamos lá:

Nesse repositório, baixe o arquivo `initramfs-factory.tar`, eu baixei o da pasta `Version 22.03`

Então copie via `scp` para o erx essa imagem

```bash
scp -O -o HostKeyAlgorithms=ssh-rsa openwrt-ramips-mt7621-ubnt_edgerouter-x-initramfs-factory.tar ubnt@192.168.1.1:/tmp/
```

Pelo meu entendimento, tive que usar `-o HostKeyAlgorithms=ssh-rsa` porque o algoritmo usado no servidor [SSH é velho e considerado não seguro](https://unix.stackexchange.com/a/682921),
mas como é local neste caso é ok.

Então acesse o erx via ssh e instale a imagem:

```bash
ssh ubnt@192.168.1.1

add system image /tmp/openwrt-ramips-mt7621-ubnt_edgerouter-x-initramfs-factory.tar

reboot
```

Após reiniciar, mude o cabo de rede da porta 0 para porta 1.

De novo, creio que será necessário colocar um IP fixo, faça-o.

Agora podemos realizar um update do sistema, para tal, baixe a imagem oficial usando a página [firmware-selector](https://firmware-selector.openwrt.org/)

De novo copie usando `scp`, entre via ssh e rode:

```bash
sysupgrade -n -F /tmp/openwrt-24.10.0-ramips-mt7621-ubnt_edgerouter-x-squashfs-sysupgrade.bin
```

Agora a LuCI (a interface) deve estar disponível na porta 192.168.1.1, e creio que também tem um DHCP server,
então um IP estático não é necessário mais.

# Colocando bridge mode
Provavelmente ambos os ips do roteador da Woo E do erx estão como 192.168.1.1.

Ligue o cabo de rede diretamente no roteador da Woo, então vá para 192.168.1.1 no browser, o usuário/senha padrão deve ser `admin/admin`.
Acesse Internet Conectivity, clique na aba Bridge Mode. Ou acesse diretamente o link [https://192.168.1.1/2.0/gui/#/internetConnectivity/bridgemode](https://192.168.1.1/2.0/gui/#/internetConnectivity/bridgemode)

Habilite Bridge Mode E alguma porta específica, no meu caso, usei a quarta.

Confirme. Coloque um cabo de rede entre a porta 4 do roteador da Woo e a porta 0 do erx.

Pronto! Agora todo o roteamento é feito via erx.

BTW, ele permite colocar bridge mode POR porta, para que você consiga ligar diretamente TV ou whatever sem passar por bridge mode, interessante.

# Mudando o DNS
Ligue o cabo de rede no erx. Acesse a interface do LuCI.

Aba Network -> Interfaces. Em WAN clique edit, na aba Advanced Settings há uma seção "Use custom DNS servers". Coloque os
servidores DNS desejados.

(Em algum momento eu instale um unbound ou algo parecido e rode meu próprio servidor DNS, não hoje.)


# Downsides

Terminado todos esses passos, vamos falar dos problemas.

O principal deles é que o Wifi do roteador da Woo ainda está ativo, portanto é possível bypassar o erx.
Deixei assim porque o router possui Wifi6, que é melhor que o Wifi5 que eu tenho e que coloquei como access point.
Idealmente você deixaria um access point com Wifi6,7 diretamente ligado no erx, e desligaria o da operadora.

Essa questão de ambos os roteadores terem o mesmo IP é levemente confuso, mas estando conectado na rede do erx seria impossível
acessar o roteador da Woo anyway, então é ok. Mas se quiser mudar o IP do erx fique à vontade.
