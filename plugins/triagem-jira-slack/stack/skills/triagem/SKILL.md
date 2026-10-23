---
name: triagem
description: Faz a triagem dos tickets novos do Jira. Classifica prioridade e componente, sugere o responsável pelo histórico do componente, grava a triagem e o comentário com aprovação, posta o resumo no canal do Slack e aponta os tickets com SLA estourado. Use quando pedirem para triar, classificar ou priorizar tickets, chamados ou issues do Jira, atribuir responsável, resumir a fila no Slack ou cobrar SLA atrasado.
---

# Triagem de tickets

## Antes de tudo
1. Descubra o site com `get_accessible_atlassian_resources` e use o `cloudId` em todas as chamadas do Jira.
2. Busque os tickets com `search_jira_issues_using_jql`. Padrão para "tickets novos": `project = <CHAVE> AND created >= -24h ORDER BY created DESC`. Nunca busque mais de 50 por vez.

## Classificar cada ticket
- **Prioridade** (use os nomes do Jira do cliente; se não souber, leia um ticket com `get_jira_issue`):
  - Highest: sistema fora do ar, perda de dados, falha de pagamento ou segurança, muitos clientes afetados.
  - High: funcionalidade principal quebrada sem contorno.
  - Medium: falha com contorno, ou um cliente afetado.
  - Low: dúvida, ajuste visual, melhoria.
- **Componente**: escolha entre os componentes que já existem no projeto. Não crie componente novo.
- **Responsável sugerido**: quem mais resolveu tickets do mesmo componente nos últimos 30 dias (`component = "<X>" AND resolved >= -30d`). Se não houver histórico, deixe sem responsável e diga isso.

## Gravar
- `edit_jira_issue` e `add_or_edit_jira_issue_comment` pedem aprovação humana. Junte as mudanças e peça uma aprovação por ticket, mostrando: chave, prioridade atual → nova, componente, responsável sugerido e o motivo em uma linha.
- O comentário de triagem é curto: classificação, motivo e próximo passo. Nunca cole dados pessoais do cliente no comentário.
- Se o usuário pediu "só consultar" ou "não altere", não chame tools de escrita.

## Resumo no Slack
- Ache o canal com `slack_search_channels` quando o usuário der o nome.
- Uma mensagem só, em tópicos: total triado, quantos por prioridade, os Highest com chave e responsável, e os que ficaram sem dono.
- `slack_send_message` pede aprovação humana: mostre o texto final antes.

## SLA
- Padrão, se o time não informar outro: Highest 4 h, High 1 dia útil, Medium 3 dias úteis, Low sem SLA.
- "Estourado" é: aberto, sem resolução, e `created` mais antigo que o SLA. Cite chave, idade e responsável (ou "sem dono").
