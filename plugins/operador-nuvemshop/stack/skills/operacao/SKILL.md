---
name: operacao
description: Opera a loja Nuvemshop no dia a dia. Monta o resumo diário de vendas no Slack (pedidos pagos, faturamento, ticket médio e mais vendidos), avisa o estoque baixo e muda preço ou cria cupom sempre com aprovação. Use quando o lojista pedir resumo ou relatório de vendas do dia, pedidos, faturamento, ticket médio, produtos acabando, estoque baixo, reposição, mudar preço, promoção, desconto ou cupom.
---

# Operação da loja

## Resumo diário
- Use `list_orders` com a data pedida (ontem, se o usuário não disser). Conte só pedidos pagos; liste à parte os pendentes e cancelados.
- Calcule com `calcular`: quantidade, faturamento e ticket médio (faturamento ÷ pedidos pagos). Mostre em reais com duas casas.
- Mensagem no Slack em tópicos: pedidos pagos, faturamento, ticket médio, os 3 produtos mais vendidos e pendências. `slack_send_message` pede aprovação humana.
- Nunca poste nome, e-mail, telefone ou endereço de comprador.

## Estoque baixo
- Use `list_products` e considere estoque baixo `<= 5` unidades por variante, salvo outro limite dado pelo lojista. Ignore produtos sem controle de estoque.
- Liste produto, variante, estoque e vendas dos últimos 7 dias quando houver. Não altere estoque por conta própria.

## Preço (aprovação humana)
- Confirme o produto e a variante com `list_products` antes de mudar.
- Calcule o preço novo com `calcular` e mostre antes: produto, variante, preço atual → novo e a variação em %.
- Recuse queda de mais de 30% ou preço abaixo do custo informado sem confirmação explícita do lojista.
- Só então chame `update_product`.

## Cupom (aprovação humana)
- Mostre antes: código, tipo (percentual ou valor), valor, validade, limite de uso e valor mínimo do carrinho.
- Padrões, se o lojista não disser: validade de 7 dias, um uso por cliente, sem valor mínimo. Desconto acima de 30% exige confirmação explícita.
- Só então chame `create_coupon`.
