---
name: cobranca
description: Cobra e recupera faturas vencidas na Stripe. Lista e soma o total em aberto, usa o link da fatura ou cria um payment link, registra nota e tarefa de retorno no HubSpot, prepara o lembrete para o canal do financeiro no Slack conforme a régua de atraso e só faz reembolso a pedido e com aprovação. Use quando o usuário falar de fatura vencida ou em atraso, inadimplência, cobrança, saldo devedor, quanto falta receber, link de pagamento, lembrete de pagamento ou reembolso.
---

# Cobrança e recuperação

## Faturas vencidas
- Busque com `stripe_api_read` em `GET /v1/invoices` com `status=open` e `due_date[lt]=<agora menos N dias>` (N = 0 se o usuário não disser). Pagine até o fim; no máximo 100 faturas por execução.
- Some valores com `calcular`. Valores da Stripe vêm em centavos: divida por 100 e mostre em reais (ou na moeda da fatura).

## Payment link
- A fatura aberta já tem `hosted_invoice_url`: use esse link. Só crie payment link (`stripe_api_write`, `POST /v1/payment_links`) quando não houver fatura (cobrança avulsa).
- Toda chamada de `stripe_api_write` pede aprovação humana: mostre cliente, valor, moeda e o que será criado.

## Registro no CRM
- Ache o contato pelo e-mail do cliente com `search_crm_objects`.
- Com `manage_crm_objects` (aprovação humana), crie uma **nota** com fatura, valor, dias de atraso e link, e uma **tarefa** de retorno em 3 dias úteis para o dono do contato.
- Cliente fora do CRM: liste no relatório e não crie contato.

## Lembrete
- Prepare o texto do lembrete por cliente (cordial, valor, vencimento, link) e mande ao canal do financeiro com `slack_send_message` para revisão. O agente não envia e-mail ao cliente.
- Régua: 1–7 dias de atraso, lembrete gentil; 8–30, segundo aviso; mais de 30, escale para um humano e não mande lembrete.

## Reembolso
- Só com pedido explícito, com cliente e fatura identificados. Confirme na Stripe que a cobrança foi paga e não foi reembolsada.
- Use `stripe_api_write` em `POST /v1/refunds` com o `payment_intent` ou a `charge`. A aprovação é dupla: a do Agorch e a nativa da Stripe, que devolve uma URL de aprovação. Mostre a URL e espere; nunca tente contornar.
- Nunca reembolse mais que o valor pago nem crie reembolso por iniciativa própria.
