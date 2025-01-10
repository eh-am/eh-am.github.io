+++ 
date = 2024-10-24T10:53:52Z
title = "Glassmorphism com CSS"
description = ""
categories = ["technical"]
+++

(É esperado que esteja rodando num navegador novo. Para ver os exemplos use a feature do browser "Ver código fonte" como um bom old schooler :)  )

Estava vendo um [vídeo de como fazer Glassmorphism no Photoshop](https://www.youtube.com/watch?v=aZOIbAeBVFA) e pensei: dá para fazer com CSS?

Sim, dá:

<div class="sand-image">
  <div class="glass-box">
    <span class="wrap"></span>
  </div>
</div>

<style>
/* source: https://stackoverflow.com/questions/65935742/button-with-transparent-background-and-rotating-gradient-border/65936035#65936035 */
/** 
 * It works as follows:
 * Create two elements
 * One refers to the border, which is actually just an element with linear gradient from top to bottom
 * Then a mask element, its mask contains 2 layers
 *  - one for the entire content
 *  - another for the content + padding
 *  By using mask-composite: exclude (think xor), it will
 *  only show the contents that are exclusive (ie not in both layers)
 *  in this case, it will only show the padding area
 *
 * */
/* ALTERNATIVE:
 * One could use https://hype4.academy/tools/glassmorphism-generator
 * the downside is that the border is slightly less interesting, but it's easier to implement
 * */
.sand-image {
  background: url('./img-94671052.1920.jpg');
  background-size: cover;
  background-position-y: 100%;
  width: 100%;
  height: 500px;
  position: relative;
}

.glass-box {
  position: absolute;
  inset: 0;
  margin: auto;
  border-radius: 1rem;

  width: 50%;
  height: 50%;
  backdrop-filter: blur(10px);
  background: linear-gradient(
    160deg,
    rgba(255, 255, 255, 0.15),
    rgba(255, 255, 255, 0)
  );
}

.glass-box .wrap {
  position: absolute;
  z-index: -1;
  inset: 0;
  border-radius: 1rem;
  padding: 1.5px; /* the thickness of the border */
  /* the below will do the magic */
  mask:
    linear-gradient(#fff 0 0) content-box,
    /* this will cover only the content area (no padding) */
      linear-gradient(#fff 0 0); /* this will cover all the area */

  -webkit-mask-composite: xor; /* needed for old browsers until the below is more supported */
  mask-composite: exclude; /* this will exclude the first layer from the second so only the padding area will be kept visible */
}


.glass-box .wrap::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;

  /* 180deg = to bottom = default so we can leave it empty */
  background: linear-gradient(180deg, #ffffff73 0%, #ffffff10 50%);
}
</style>


## Como funciona?

Há várias níveis diferentes de Glassmorphism.
![](https://media.nngroup.com/media/editor/2024/05/31/strokesgradients.jpg)
Imagem via [The Nielsen Norman Group](https://www.nngroup.com/articles/glassmorphism/)

O mais básico consiste em apenas um blur no fundo, que é trivial de se implementar usando `backdrop-filter: blur`. A dificuldade mesmo fica na borda.

Se for possível, a solução mais fácil é usar esse [Glassmorphism Generator](https://hype4.academy/tools/glassmorphism-generator):


<div class="sand-image">
  <div class="glass-box via-generator"></div>
</div>

<style>
.via-generator {
  background: rgba( 255, 255, 255, 0.25 );
  box-shadow: 0 8px 32px 0 rgba( 31, 38, 135, 0.37 );
  backdrop-filter: blur( 1.5px );
  -webkit-backdrop-filter: blur( 1.5px );
  border-radius: 10px;
  border: 1px solid rgba( 255, 255, 255, 0.18 );
}
</style>

No meu caso, gerou o seguinte CSS:
```css
.via-generator {
  background: rgba( 255, 255, 255, 0.25 );
  box-shadow: 0 8px 32px 0 rgba( 31, 38, 135, 0.37 );
  backdrop-filter: blur( 1.5px );
  -webkit-backdrop-filter: blur( 1.5px );
  border-radius: 10px;
  border: 1px solid rgba( 255, 255, 255, 0.18 );
}
```

Como podemos perceber, há 2 problemas:

O primeiro é essa box-shadow que tá mais pro azulado. É facilmente corrigível só mudando a cor para algo mais para o lado do preto. Tipo `0 8px 32px 0 rgba(22, 25, 57, 0.37)`.

O segundo é referente a borda. Ela é uniforme, então não dá aquela impressão de profundidade.

Aí que mora o perigo. 

Primeiro, precisamos fazer um gradiente com borda. Que é possível usando `border-image`.



<div class="sand-image">
  <div class="glass-box border-image"></div>
</div>


<style>
.border-image {
  border: 1px solid;
  border-image: linear-gradient(rgba(255, 255, 255, 1), rgba(255, 255, 255, 0.2)) 1;
  /* DOES NOT WORK! */
  border-radius: 15px;
}
</style>

Porém nesse caso, `border-radius` não funciona. Ou melhor, funciona, mas ignora que deveria ser em volta da imagem.

[Uma técnica comum](https://stackoverflow.com/a/53037637) para este caso é adicionar uma borda falsa usando `background-image`.

<div class="sand-image">
  <div class="glass-box background-image"></div>
</div>


<style>
.background-image {
  border: solid 10px transparent;
  background-image:
    linear-gradient(#fff, #fff),
    linear-gradient(rgba(255, 255, 255, 1), rgba(255, 255, 255, 0.2));
  background-origin: border-box;
  background-clip: content-box, border-box;
  border-radius: 15px;
}
</style>

Obviamente fica um lixo, porque requer que o centro seja uma cor sólida (o primeiro linear-gradient).

Mas estamos no caminho certo. A ideia é ter um elemento que crie a borda usando background-image e gradiente. Porém para tal, será um elemento INTEIRO com gradiente. Então precisamos de outro elemento que DESFAÇA esse gradiente.

Parece complicado (e é), mas podemos fazer isso usando [CSS Mask](https://developer.mozilla.org/en-US/docs/Web/CSS/mask).


## CSS Mask

Só ilustrando, a `mask` funciona similar ao Photoshop e outras ferramentas. Basicamente só mostra o que está dentro dessa máscara.

Por exemplo:

<div class="sand-image">
  <div class="glass-box mask-1"></div>
</div>

<style>
.mask-1{
  -webkit-mask-image: url(./blob.svg);
  mask-image: url("./blob.svg");
  mask-repeat: no-repeat;
  mask-position: center center;
  mask-size: cover;
}
</style>

Aqui ainda estamos usando aquele mesmo retângulo de vidro, mas aplicando uma máscara com essa forma estranha (btw gerada usando [SVG Shape Generator](https://www.softr.io/tools/svg-shape-generator)).


Também podemos aplicar contra o próprio background:

<div class="sand-image mask-2">
</div>

<style>
.mask-2{
  -webkit-mask-image: url(./blob.svg);
  mask-image: url("./blob.svg");
  mask-repeat: no-repeat;
  mask-position: center center;
  mask-size: cover;
}
</style>

Mas acho que um dos efeitos mais legais é usar contra uma máscara com gradiente:

<div class="sand-image mask-3">
</div>

<style>
.mask-3{
  mask-image: linear-gradient(to bottom, black, transparent 100%);
}
</style>

BTW a sintaxe do linear gradiente é um pouco enjoada. Então vou deixar mais um exemplo:


```css
.mask-4 {
  mask-image: linear-gradient(to top, transparent 5%, black 50%);
}
```

Isso diz que:
* o linear gradiente começa de baixo para cima
* só inicia a transparência a partir dos 5% iniciais da imagem (ie. joga fora os 5% iniciais)
* na metade final da imagem ela é preta, aka a máscara está full


<div class="sand-image mask-4">
</div>

<style>
.mask-4{
  mask-image: linear-gradient(to top, transparent 5%, black 50%);
}
</style>

Enfim, esse caso do gradiente é útil para ter imagens E um texto embaixo, tipo um título ou algo parecido.


Mas voltando ao problema original, vamos criar, dentro da div `glass-box` uma div `wrap`. Essa div `wrap` também tem um pseudo elemento com um gradiente de background que será a borda. Aí o pulo do gato: usamos mask com `mask-composite: exclude` que funciona com um xor, para EXCLUIR tudo que não for a falsa borda. Para isso, usamos duas "layers" na máscara: uma referente ao `content-box`, e outra referente a tudo. 

Para ilustração (não é o caso real), imagine que tem 2 divs: uma vermelha, ligeiramente maior, e outra verde. E onde há overlap, as cores se misturam, então onde está roxo é porque tem verde + vermelho:

<div class="wrap-fake">
</div>

<style>
.wrap-fake {
  position: relative;
  width: 100%;
  height: 500px;
  
  background: green;
}

.wrap-fake::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  border: 15px solid red;
  background: purple;
}
</style>

Agora adaptando, o purple é removido, já que está na máscara. E a borda vermelha agora é um linear-gradient.


## Final
Ficará assim:

<span class="wrap"></span>

<div class="sand-image">
  <div class="glass-box">
    <span class="wrap"></span>
  </div>
</div>

Lembrando que para ver os códigos clique em Visualizar Código fonte no seu browser :)

## Referências

* [mask - CSS: Cascading Style Sheets | MDN](https://developer.mozilla.org/en-US/docs/Web/CSS/mask)
* [Apply effects to images with CSS's mask-image property  |  Articles  |  web.dev](https://web.dev/articles/css-masking#masking_with_a_gradient)
* [Glassmorphism CSS Generator | SquarePlanet | SquarePlanet](https://hype4.academy/tools/glassmorphism-generator)
* [Glassmorphism: Definition and Best Practices](https://www.nngroup.com/articles/glassmorphism/)
* [html - Button with transparent background and rotating gradient border - Stack Overflow](https://stackoverflow.com/questions/65935742/button-with-transparent-background-and-rotating-gradient-border/65936035)
* [kriptonian_ comments on Does anyone know how can I give this border gradient, with glass morphism](https://old.reddit.com/r/tailwindcss/comments/164b88f/does_anyone_know_how_can_i_give_this_border/jy7sw51/?share_id=Qvx6Xvoc6sz-fM1jrUDgZ)
