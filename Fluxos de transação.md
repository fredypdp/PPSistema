- Fluxo de transação nacional:
    [] PSP do Pagador:
        [] Envia a mensagem de crédito (pacs.008) para o PP
        [] Debita o valor do pagamento na conta do usuário
    [x] PP:
        [x] Faz uma validação do pacs.008, e caso não dê erro, aprova a transação
        [x] Deduz a taxa (1.5%) no valor do pagamento:
            [x] São criados 2 pacs.008 a partir do pacs.008 original, um para o PSP do beneficiário (com o valor do pagamento depois da taxa), e outro para o PSP do PP, onde segue o fluxo de taxação
    [] PSP do Beneficiário:
        [] Recebe o pacs.008 (já taxado)
        [] Envia uma mensagem pacs.002 para o PSP do pagador (passa pelo PP) confirmando o recebimento do pacs.008
        [] Se creditou com sucesso o valor na conta do beneficiário
        [] Envia um pacs.002 para o PP
        [] Se não processou o pagamento ou teve uma devolução
        [] Envia um pacs.004 para o PP
    [] PP:
        [x] Recebe o pacs.002/pacs.004 e registra a transação
        [] Envia a mensagem para o PSP do pagador
    [] PSP do Pagador:
        [] Recebe o pacs.002/pacs.004
        [] Notifica o usuário o sucesso/rejeição da transação
- Fluxo de taxação:
    [x] Obs.: Apôs a criação dos 2 pacs.008, os dados do credor no pacs destinado ao PSP do PP são alterados pelos dados do PP
    [] PSP do PP:
        [] Recebe o pacs.008 com o valor adquirido da taxa
        [] Credita o valor na conta da empresa (PP)
        [] Envia uma mensagem pacs.002/pacs.004 para o PP
    [] PP:
        [] Recebe o pacs.002/pacs.004
        [] Analisa e registra a transação da taxa