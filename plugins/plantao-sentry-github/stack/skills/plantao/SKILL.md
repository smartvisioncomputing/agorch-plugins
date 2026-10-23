---
name: plantao
description: Plantão de erros do Sentry para o GitHub. Busca e agrupa os erros novos, procura a issue existente para não duplicar, acha o dono do código pelo CODEOWNERS ou pelos commits do arquivo, abre a issue com stack e release (com aprovação) e avisa o canal no Slack. Use quando falarem de erro novo em produção, exceção, crash, alerta do Sentry, bug, plantão, on-call, abrir issue no GitHub ou descobrir quem é o dono do código.
---

# Plantão de erros

## Achar e agrupar
- Busque com `sentry_search_issues` (Sentry): erros não resolvidos, vistos pela primeira vez na janela pedida (última hora, se o usuário não disser), ordenados por eventos. No máximo 20 por execução.
- Agrupe pelo mesmo tipo de exceção + mesmo arquivo/função do topo da stack. Um grupo vira **uma** issue.
- Leia o detalhe com `get_sentry_resource`: mensagem, stack (frames do app, não de bibliotecas), release, ambiente, eventos e usuários afetados.

## Não duplicar
- Antes de abrir, procure com `github_search_issues` (GitHub) no repositório: o id curto do Sentry (ex.: `WEB-123`) no título ou no corpo, inclusive em issues fechadas há menos de 14 dias. Se existir, não abra outra: cite a existente.

## Dono do código
1. Leia `CODEOWNERS` com `get_file_contents` (tente `.github/CODEOWNERS`, `CODEOWNERS` e `docs/CODEOWNERS`) e aplique a regra mais específica ao arquivo do topo da stack.
2. Sem CODEOWNERS: `list_commits` do arquivo e use o autor mais frequente dos últimos 90 dias.
3. Sem nenhum dos dois: deixe sem dono e diga isso.

## Abrir a issue (aprovação humana)
`issue_write` com `method = "create"`. Mostre antes o título e o corpo.
- Título: `[Sentry <ID-CURTO>] <Exceção>: <mensagem curta>`.
- Corpo: link do Sentry, release, ambiente, eventos e usuários afetados, os 5 primeiros frames do app e uma menção ao dono (`@usuario` ou `@org/time`).
- Labels `bug` e `sentry`. Nunca cole tokens, cookies, e-mails ou dados de usuário que apareçam no evento.

## Avisar o canal
- Uma mensagem só com a lista: id do Sentry, título curto, eventos, dono e link da issue. `slack_send_message` pede aprovação humana.
