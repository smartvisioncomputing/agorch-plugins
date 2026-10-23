---
name: qualificacao
description: Qualifica as conversões novas do RD Station com nota de 0 a 3 (perfil, contato e interesse), atualiza o contato, move a negociação no funil, cria a tarefa de follow-up para o dono e avisa cada vendedor no Slack sem expor e-mail nem telefone. Use quando pedirem para qualificar ou triar leads, conversões e contatos novos, mover a negociação no funil, criar follow-up ou fazer o trabalho de SDR no RD Station.
---

# Qualificação de leads (RD Station)

## Achar os leads
- Use `search_contacts` para os contatos com conversão na janela pedida (hoje, se o usuário não disser outra). No máximo 50 por vez.

## Nota de qualificação (0 a 3)
Some um ponto por critério atendido, com base só no que está no RD Station:
1. **Perfil**: tem empresa e cargo de decisão (dono, sócio, diretor, gerente, head, C-level).
2. **Contato**: tem e-mail corporativo (não gmail, hotmail, outlook, yahoo) ou telefone.
3. **Interesse**: converteu em "fale com vendas", "demonstração" ou "orçamento", ou tem 2+ conversões.

- **3**: qualificado. **2**: nutrir com follow-up. **0–1**: não qualificado.
- Dado faltando conta como critério não atendido. Nunca complete dado com suposição.

## Gravar no RD Station (pede aprovação humana)
Mostre antes o que vai mudar e peça uma aprovação por lead:
- Qualificado: `update_contact` com a nota e o motivo; `move_deal_stage` para a primeira etapa do funil de vendas; `create_task` de follow-up para o dono, vencendo no próximo dia útil.
- Nutrir: `update_contact` com a nota e `create_task` em 3 dias úteis. Não mova etapa.
- Não qualificado: só `update_contact` com o motivo.
- Se o usuário pediu só relatório, não grave nada.

## Avisar o vendedor
- Uma mensagem por vendedor com os qualificados dele: nome, empresa, nota e motivo. Ache o canal com `slack_search_channels`. `slack_send_message` pede aprovação humana.
- Lead sem dono: diga isso no relatório e não poste.
- Não poste e-mail nem telefone do lead no Slack.
