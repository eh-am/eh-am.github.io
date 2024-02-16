+++ 
date = 2024-02-16T19:04:52Z
title = "Windows Auto Login após Suspender"
description = ""
categories = ["life"]
+++

Este post talvez seja sobre a tal _death by thousand cuts_. Ou sobre _enshittification_. Ou Linux > Windows.
Interprete como quiser.

Em resumo, tenho um desktop que praticamente tornei num Console: fica ligado na televisão,
controle plugado, Steam inicializa sozinha em Big Picture Mode.

Claro, na maior parte do tempo o computador fica em hibernação.

Pra voltar da hibernação, aperto qualquer botão do controle. Inesperadamente o Windows reconhece
e "acorda".

Mas aí, sou apresentado com uma tela de login. Só existe um usuário, sem senha. Logo, só existe um único botão "Login"
que pode ser interagido.

Pressionar qualquer botão no controle não funciona. Apesar de literalmente só ter uma única ação possível. (Pra ser justo, até existem outras opções, mas são totalmente secundárias)

Então precisava pegar o teclado e apertar Enter. O que fazia com que o teclado ficasse plugado o tempo inteiro.

Decidi pesquisar como logar sozinho E A SOLUÇÃO VAI EXPLODIR SUA MENTE!!!!!!!1

[windows 10 - How do I get the desktop to show directly after wakeup from sleep? - Super User](https://superuser.com/questions/1066207/how-do-i-get-the-desktop-to-show-directly-after-wakeup-from-sleep/1115451#1115451)

Basicamente, precisa selecionar uma opção para Relogar automaticamente. Só que esta opção não existe
caso o usuário NÃO tenha senha.

Perceba bem o que isso significa. Você tem um (ou mais usuários) sem senha.
Por que alguém não teria senha? Porque ela não se importa com privacidade.
Ela não ter senha significa que qualquer um pode logar com seu usuário.

O caso perfeito para que logue automaticamente ao acordar.

Então agora você é obrigado a colocar uma senha. Bacana, agora a opção aparece e tudo se resolve.
Exceto que agora que o usuário tem senha, ao logar pela primeira vez (sei lá, energia acabou),
precisa da porcaria do teclado para colocar a senha. Inconveniente.

MAaaaaaaaas na verdade, eu menti. Dá para, após selecionar a opção, remover sua senha.
Aí a opção fica invisível, como era inicialmente. Porém seu estado é o que você acabou de colocar.

Uma comparação seria um interruptor, ligado, só que invisível. Para deixar o interruptor visível e poder mexer no estado do interruptor (ligá-lo / desligá-lo),
é necessário realizar alguma tarefa completamente não relacionada, como virar a chave da porta.

Jesus Amado.
