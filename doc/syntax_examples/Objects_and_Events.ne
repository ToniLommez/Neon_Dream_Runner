obj Character {
    pub name: string
    gemCount: int = 0
    hasSpecialGem: bool = False

    fn collectGem {
        gemCount++
        print "Gema coletada! Total: {}", gemCount
    }

    fn collectSpecialGem {
        gemCount *= 2
        hasSpecialGem = True
        print "Gema especial coletada! Total: {}", gemCount
    }

    when gemEventTrigger {
        let chanceSpecialGem = random 0 10 // Cria uma chance aleatória para a gema especial.
        if chanceSpecialGem == 0 {
            collectSpecialGem
        } else {
            collectGem
        }
    }
}

fn acao_aleatoria {
    loop 10 { // Inicia a coleta de gemas.
        wait <| random 500 1500 // Simula o jogador
        trigger gemEventTrigger
    }
}

fn main {
    let! player Character
    play = player{"Player1", 0, 0}

    print "Bem-vindo ao jogo, {}!", char.name
  
  	!> acao_aleatoria player // Comando assincrono

    pulse until char.gemCount >= 10 {
        if player.gemCount % 2 == 0 {
            gerar_inimigo
        }
    }

    print "Parabéns, {}! Você coletou todas as 10 gemas e completou o jogo.", char.name
    char.free
}
