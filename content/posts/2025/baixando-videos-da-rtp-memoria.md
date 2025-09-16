+++ 
date = 2025-09-16T11:22:08+01:00
title = "Baixando Vídeos da RTP Memória"
categories = ["technical"]
+++

Dica rápida.

Por ex, baixando o vídeo https://arquivos.rtp.pt/conteudos/cantos-de-trabalho-2/

Abre o dev tools, procura pelo arquivo `manifest.mpd`.
Nesse caso, usa é `https://streaming-vod.rtp.pt/drm-dash/nas2.share/mcm/arquivo/mp4/380/380914df55af34fb9e5bc6202abfbd92arq.mp4/manifest.mpd`

Baixa usando `yt-dlp` e a flag `--allow-u` pra aceitar arquivos DRM.

Depois instala `bento4` (com brew é `brew install bento4`), vai estar disponível a cli `mp4decrypt`

Decripte o arquivo usando a chave `b57fd50083e83bfd85bbc31a32cf47e1:3d5f6bcf52c213ac76f5b896ee95be4e`


Exemplo:
```bash
mp4decrypt --key b57fd50083e83bfd85bbc31a32cf47e1:3d5f6bcf52c213ac76f5b896ee95be4e manifest\ \[manifest\].fv1-x3.mp4 video.mp4
```

Essa chave peguei [daqui](https://forum.videohelp.com/threads/416939-DRM-Protected-on-RTP-Arquivos#post2760881),
mas você pode sempre conseguir você mesmo usando a extensão [WidevineProxy2](https://forum.videohelp.com/threads/416316-%5BRelease%5D-WidevineProxy2-Extension-Bypass-HMAC-1-timetokens-Lic-wrapping).
