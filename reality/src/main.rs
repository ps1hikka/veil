use std::io::{self, copy};
use std::net::{TcpListener, TcpStream};
use std::thread;

fn main() -> io::Result<()> {
    let listener = TcpListener::bind("0.0.0.0:443")?;
    println!("listening on 0.0.0.0:443 → gateway:8443");

    for stream in listener.incoming() {
        let client = stream?;
        thread::spawn(move || {
            if let Err(e) = pipe(client, "gateway:8443") {
                eprintln!("pipe error: {e}");
            }
        });
    }
    Ok(())
}

fn pipe(mut client: TcpStream, backend: &str) -> io::Result<()> {
    let mut backend = TcpStream::connect(backend)?;

    let mut client2 = client.try_clone()?;
    let mut backend2 = backend.try_clone()?;

    let t1 = thread::spawn(move || copy(&mut client, &mut backend));
    let t2 = thread::spawn(move || copy(&mut backend2, &mut client2));

    t1.join().unwrap()?;
    t2.join().unwrap()?;
    Ok(())
}