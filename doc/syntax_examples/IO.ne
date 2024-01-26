use (
    "io",
    "http"
)

fn main {
    // io nativo
    let value = input "Digite algo:"
    print "Valor digitado {}" value

    // leitura de arquivo
    let file = io.open "example.txt" "rw"
    print "Conteudo = {}" file.read
    file.write "Hello, Neon!"

    // io personalizado
    let reader = io.new_reader (io.open "example.txt" "w")
    let value = reader.input
    reader.print "Hello, Neon!"
}
