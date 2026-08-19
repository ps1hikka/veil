use std::fs;
use std::io::{Read, Write};
use std::os::unix::net::UnixListener;

fn main() {
    let path = "/tmp/reality.sock";
    let _ = fs::remove_file(path);

    let listener = UnixListener::bind(path).expect("failed to bind socket");
    println!("reality listening on {}", path);

    for stream in listener.incoming() {
        match stream {
            Ok(mut stream) => {
                std::thread::spawn(move || {
                    let mut buf = [0u8; 4096];
                    loop {
                        match stream.read(&mut buf) {
                            Ok(0) => break,
                            Ok(n) => {
                                if stream.write_all(&buf[..n]).is_err() {
                                    break;
                                }
                            }
                            Err(_) => break,
                        }
                    }
                });
            }
            Err(e) => eprintln!("accept error: {}", e),
        }
    }
}