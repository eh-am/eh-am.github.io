+++ 
date = 2024-07-22T17:15:53+01:00
title = "copiando um disco com dd"
categories = ["technical"]
lang = "pt-br"
+++

# introdução obrigatória com estorinha comovente

Tenho um servidor paupérrimo, um Raspberry Pi 4 Model B Rev 1.1, 2GB de RAM.
Nada extraordinário, já que uso para Plex com Direct Play, sem necessidade de transcoding.

Além do mandatório cartão SD [^1], Por um bom tempo usei como volume para dados, um HD Externo
Seagate Backup Plus Slim Portable Drive 1TB (SRD00F1).

1TB é bastante, até que não é. Particularmente gostava de ter esse limite, porque me forçava a consumir as mídias baixadas e apagá-las, meio que numa imposição artificial anti hoarding.

Aí a Netflix começou a reclamar que compartilho senha com meus pais, e decidi parar de usar. O maior impacto foi a facilidade de acesso à episódios de Star Trek TNG/DS9. Ok, tem coisas que merecem ser baixadas e deixadas lá para sempre (cough cough Peep Show) :)

Na AmAzOn PrImE dAy TM comprei um Western Digital My Passport 260D, com 5TB de espaço. Não sei se suficiente para Data Hoarding, mas já melhor do que estava antes. Não ânseio muito por NAS e outras coisas do tipo, me parece um pouco too much, meio try hard. Gosto do espírito DYI de ter um RPi e um HD Externo.


# a migração
O HD antigo já estava quase cheio (~900GB). Parei todos os containers que estavam rodando (e eram os únicos apps com acesso àquele disco) para não ter um skew da cópia: o que acontece se copia arquivo A, e ele é modificado depois?.

Para copiar, até poderia simplesmente fazer um `rsync` para copiar os arquivos. Mas pensei: já que está quase cheio, mais fácil seria copiar o disco inteiro usando `dd`! [^2]

Pluguei o novo HD no RPi, que BTW é carinhosamente apelidado de `terok-nor` (MAS NÃO SOMOS APOLOGISTAS DA OCUPAÇÃO CARDASSIAN!), identifiquei qual device era qual (`fdisk -l`) e comecei copiar:

```bash
dukat@terok-nor:/mnt $ sudo dd if=/dev/sda of=/dev/sdb bs=100M status=progress
```

A sintaxe disso inicialmente me confundia, um mnemônico bom é:
`if` -> **input** file
`of` -> **output** file

Lembrando que em unix **TuDo É uM aRqUiVo**.

Outra questão é que usei um `block size` de `100M`, sem nenhuma explicação específica.
Sempre soube que tunar este parâmetro é mais arte do que ciência, mas talvez devesse
ter dado mais importância. Esse blogpost mostra a diferença de escala de um block size mal tunado e um bom: [Tuning dd block size - tdg5](http://blog.tdg5.com/tuning-dd-block-size/)

Mas voltando, após um tempo o comando falhou:
```bash
2306867200 bytes (2.3 GB, 2.1 GiB) copied, 133 s, 17.4 MB/s
dd: error reading '/dev/sda1': Input/output error
22+0 records in
22+0 records out
2306867200 bytes (2.3 GB, 2.1 GiB) copied, 133.767 s, 17.2 MB/s
```

Ok, vamos rodar de novo. Não tenho os logs salvos, mas basicamente reclamava que o device
que estava usando, tipo `sda` não existia. Como assim?

Olhando nos logs do `dmesg`, percebi a mensagem de `over-current` change:

```bash
[1098323.479530] usb usb2-port3: over-current change #4
[1098323.695583] usb usb2-port4: over-current change #4
[1098323.796481] usb 1-1-port3: over-current change #2
[1098323.911561] usb usb2-port1: over-current change #7
[1098324.015713] usb 1-1-port4: over-current change #2
```

Ok, era algo que sabia que poderia acontecer. O raspberry pi simplesmente não tem
current suficiente para alimentar 2 HDs[^3]. Talvez se usar um HUB com alimentação externa funcionasse.
Mas não o tenho, então simplesmente movi os dois HDs para o meu Macbook.

A única diferença (que percebi) no `dd` que vem com o macOS é que `m` é a unidade para `MiB`.

Depois de MUITO tempo a cópia terminou.

Pluguei o novo HD no RPi, e rodei `fdisk -l`, ele deu um aviso engraçado
```
GPT PMBR size mismatch (1953525166 != 9767475199) will be corrected by write.
The backup GPT table is not on the end of the device.
```

"Corrected by write". Que write?

E eu nem fazia ideia que existia uma backup GPT table, que é escrita no final do device, mas pelo menos que bom que é backup.

BTW [GPT](https://en.wikipedia.org/wiki/GUID_Partition_Table) caso não saiba, é só uma tabela no disco informando cada partição e uns metadados associados à ela (imagino que tipo de file system, o UUID etc).

Rodei `parted` que supostamende veria corrigir o problema:
```bash
dukat@terok-nor:/mnt $ sudo parted -l
Warning: Not all of the space available to /dev/sda appears to be used, you can
fix the GPT to use all of the space (an extra 7813950033 blocks) or continue
with the current setting? 
Fix/Ignore? Fix                                                           
Model: WD My Passport 260D (scsi)
Disk /dev/sda: 5001GB
Sector size (logical/physical): 512B/4096B
Partition Table: gpt
Disk Flags: 

Number  Start   End     Size    File system  Name  Flags
 1      1049kB  1000GB  1000GB  ext4
```

(Não, não uso BTRFS ainda, como bom descendente de Ibéricos aqui se usa Azeite de Oliva 🤡)

Porém `fdisk` ainda acusou o tamanho anterior:
```bash
Disk /dev/sda: 4.55 TiB, 5000947302400 bytes, 9767475200 sectors
Disk model: My Passport 260D
Units: sectors of 1 * 512 = 512 bytes
Sector size (logical/physical): 512 bytes / 4096 bytes
I/O size (minimum/optimal): 4096 bytes / 4096 bytes
Disklabel type: gpt
Disk identifier: D7652DF4-F7C4-4411-AF10-17A3F5CE018B

Device     Start        End    Sectors   Size Type
/dev/sda1   2048 1953523711 1953521664 931.5G Linux filesystem
```

Dando uma [pesquisada](https://serverfault.com/a/1021192), existe uma ferramenta chamada [`growpart`](https://manpages.ubuntu.com/manpages/xenial/en/man1/growpart.1.html) que extend uma partição na tabela de partições para preencher todo o espaço disponível. Que é exatamente o que eu queria. Se tivesse que fazer algo mais complexo, melhor usar um `parted` com mais calma.

Primeiro acha onde está este binário e instala:
```bash
sudo apt-cache search growpart
sudo apt-get install cloud-guest-utils -y
```

Então expande a partição:
```bash
sudo growpart /dev/sda 1
CHANGED: partition=1 start=2048 old: size=1953521664 end=1953523711 new: size=9767473119 end=9767475166
```

E deu certo! `fdisk -l` Retorna

```bash
Disk /dev/sda: 4.55 TiB, 5000947302400 bytes, 9767475200 sectors
Disk model: My Passport 260D
Units: sectors of 1 * 512 = 512 bytes
Sector size (logical/physical): 512 bytes / 4096 bytes
I/O size (minimum/optimal): 4096 bytes / 4096 bytes
Disklabel type: gpt
Disk identifier: D7652DF4-F7C4-4411-AF10-17A3F5CE018B

Device     Start        End    Sectors  Size Type
/dev/sda1   2048 9767475166 9767473119  4.5T Linux filesystem
```

Ok agora também precisa expandir o filesystem, usando [`resize2fs`](https://linux.die.net/man/8/resize2fs). No debian vem junto quando se instala `cloud-guest-utils`.

```bash
dukat@terok-nor:/mnt $ sudo resize2fs /dev/sda1
resize2fs 1.47.0 (5-Feb-2023)
Resizing the filesystem on /dev/sda1 to 1220934139 (4k) blocks.
The filesystem on /dev/sda1 is now 1220934139 (4k) blocks long.
```


Confesso que não está 100% claro para mim a necessidade de crescer a partição E o file system. Claro que por implementação seu metadata pode estar em lugares diferentes, mas para mim ainda é apenas uma ação (aumentar filesystem), a partição deveria? vir automaticamente.


Anyway.


O mais legal é que como copiamos usando `dd`, o `UUID` da partição é o mesmo. Dessa forma `fstab` não precisa mudar :)



[^1]: Aparentemente é possível bootar sem um cartão SD. Ref: [PXE Booting Raspberry Pis – LTM Tech](https://ltm56.com/pxe-booting-raspberry-pis/)
[^2]: Teoricamente é mais rápido porque tem menos overhead de olhar pastas arquivos etc. Só copia todo o bloco de dados ignorando completamente o que se referem. Acho.
[^3]: Pelo visto conseguimos alimentar 2 HDs com USB 2.0, ou apenas 1 USB 3.0. Fonte: https://forum.libreelec.tv/thread/24395-raspberry-pi-4-wd-elements-4tb-external-drive-clicking-from-hdd/?postID=159188#post159188
