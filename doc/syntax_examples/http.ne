use "http"

fn pong(w: http.Response, r: http.Request) => string {
    "Pong"
}

fn main {
    // Cria um servidor HTTP na porta 8080
    let! server = http.create_server

    server.port 8080

    server.get "/" ((w: http.Response, r: http.Request) => string {"Hello, Neon!"})
    server.get "/ping" pong

    !> server.ignite
}
