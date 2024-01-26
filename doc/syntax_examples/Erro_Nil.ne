fn bar(x: int) => any {
    => nil
}

// teste de nulidade
fn foo(x: any?, y: float, p: Person) {
    let! a1 = bar?: // a1 nao pode ser nulo, logo 0 por padrao
    let!? a2 = bar ?: 0 // a2 pode ser nulo, mas atribui-se 0

    let! b = p?.age // b nao pode ser nulo, logo 0 por padrao

    let z = p!.metodo?: // declaramos que p nao pode ser nulo por um caso de absurdo
}

// Funcao que pode disparar um erro ou retorna-lo
fn funcao_que_pode_dar_erro(x: int) => int, err? {
    // dispardo de erro (retorno de vazio padrao)
    if x == 1 { 
        error 
    }

    // disparo de erro built-in
    if x == 2 { 
        error reductio_ad_absurdum
    }

    // disparo + tipagem
    if x == 3 { 
        error is_three "deu ruim"
    }

    // disparo + retorno personalizado
    if x == 4 {
        error is_four "deu ruim" => 4 nil
    }

    // disparo + criacao de tipagem de erro
    if x == 5 {
        let e = err{type: is_five, msg: "deu ruim"}
        error e => 4, e
    }

    // disparo com criacao personalizada
    if x == 6 {
        error err{type: is_six, msg: "deu ruim"}
    }

    // disparo + execucao de bloco (blocos retornam valores)
    if x == 7 {
        error is_seven "deu ruim" => {
            funcao_aleatoria
            bar 2, nil
        }
    }

    8, err{type: ok, msg: "deu certo"}
}

fn funcao_com_varios_tratamentos => any, bool, int{
    // retorno full nil/default para erro
    let dll_handle = (load_dll dll_path)?

    // retorno para erro com escolha de retorno
    let dll_handle = (load_dll dll_path)? => nil, true, 23

    // ação pos erro + retorno
    let dll_handle = (load_dll dll_path)? => {
        execucao_de_acao
        nil, true, 23
    } 

    // ação pos erro 
    let? dll_handle = (load_dll dll_path)? {
        acao_bizarra_aleatoria
    }

    // ação pos erro + desempacotamento
    let dll_handle = (load_dll dll_path)? {
        case err
            of err.nilPointer => acao_aleatoria
            of err.fatal => die
        else => sei_la
    }
}
