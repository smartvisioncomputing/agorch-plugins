---
name: qualificacao
description: Qualifica os leads novos do HubSpot com nota de 0 a 3 (perfil, contato e interesse), move os qualificados no funil, cria a tarefa de follow-up para o dono do contato e avisa cada vendedor no Slack sem expor e-mail nem telefone. Use quando pedirem para qualificar ou triar leads e contatos novos, separar SQL e MQL, atualizar a etapa do funil ou do pipeline, criar follow-up ou fazer o trabalho de SDR no HubSpot.
---

# Qualificação de leads (HubSpot)

## Achar os leads
- Use `search_crm_objects` em `contacts`, filtrando por `createdate` (hoje, se o usuário não disser outra janela) e `lifecyclestage` = `lead` ou vazio. No máximo 50 por vez.
- Leia empresa e negócio associados com `get_crm_objects` só quando precisar deles para a nota.

## Nota de qualificação (0 a 3)
Some um ponto por critério atendido, com base só no que está no CRM:
1. **Perfil**: tem empresa e cargo de decisão (dono, sócio, diretor, gerente, head, C-level).
2. **Contato**: tem e-mail corporativo (não gmail, hotmail, outlook, yahoo) ou telefone.
3. **Interesse**: pediu demonstração, preço ou contato, ou teve 2+ conversões.

- **3**: qualificado (SQL). **2**: nutrir com follow-up. **0–1**: não qualificado.
- Dado faltando conta como critério não atendido. Nunca complete dado com suposição.

## Gravar no CRM (pede aprovação humana)
Use `manage_crm_objects`, mostrando antes o que vai mudar:
- SQL: `lifecyclestage` = `salesqualifiedlead`; crie ou mova o negócio para a primeira etapa do pipeline de vendas; crie uma **tarefa** de follow-up para o dono, vencendo no próximo dia útil, com a nota e o motivo.
- Nutrir: `lifecyclestage` = `marketingqualifiedlead` e uma tarefa em 3 dias úteis.
- Não qualificado: só uma **nota** com o motivo. Não mova etapa.
- Se o usuário pediu só relatório, não grave nada.

## Avisar o vendedor
- Ache o dono com `search_owners` (pelo `hubspot_owner_id` do contato). Sem dono: diga isso no relatório e não poste.
- Uma mensagem por vendedor com os SQLs dele: nome, empresa, nota, motivo e o link do contato. `slack_send_message` pede aprovação humana.
- Não poste e-mail nem telefone do lead no Slack.
