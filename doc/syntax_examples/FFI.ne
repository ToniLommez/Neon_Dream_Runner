// foreign permite a importação de funções de bibliotecas externas para serem utilizadas no código Neon.
//
// Sintaxe:
// foreign "<nome_biblioteca>" from "<nome_biblioteca_neon>" [using "<nome_biblioteca_conversao>"] {
//     <definições_de_funções>
// }
//
// - <nome_biblioteca>: O nome da biblioteca externa que contém as funções a serem importadas.
// - <nome_biblioteca_neon>: O nome da biblioteca Neon que fornece compatibilidade com a biblioteca externa.
// - <nome_biblioteca_conversao> (opcional): O nome da biblioteca Neon que fornece conversão automática de tipos entre Neon e a biblioteca externa.
//   Se especificado, os tipos serão convertidos automaticamente, sem a necessidade de conversões explícitas.
//
//
// Neste exemplo, a biblioteca "libc" é importada do contexto "c" usando a biblioteca de conversão "c_types".
// A função printf é definida e pode ser utilizada no código Neon.
//
// Observação: O uso da biblioteca de conversão "c_types" é opcional e depende das necessidades do seu código e da biblioteca externa utilizada.

foreign "libc" from "c" using "c_types"{
    fn printf
}

fn main() {
    printf "Hello, %s!\n" "Neon"
}

// Biblioteca de Conversão para tipos C
obj CTypes {
    fn to_int(v int) => i32 {
        v.raw: i32
    }

    fn to_string(v string) => *u8 {
        v.raw.ptr
    }
}

// Biblioteca de Compatibilidade do C
trait C {
    // Função de conversão de int para tipo C
    fn to_int(v int) => i32

    // Função de conversão de string para tipo C
    fn to_string(v string) => *u8
}
