# Seal

<em>`Seal` is a tool for simple, secure, fast file transfer between any two computers.</em>

### Build

on MacOS with brew

```
1. brew install go
2. git clone https://github.com/WuFangyu/Seal
3. cd Seal/src
4. go build -o <the directory that you want to place the binary>/seal
5. echo "alias seal=<path to the binary>" > ~/.bash_profile
6. source ~/.bash_profile
```

run `seal` Example output:

```
Host: 192.168.0.29:9911 Public Key: 0369ed664fcd98f9c0acc310cb6f0786436bbfa06a82edfc36d1caa261e8749d72
NAME:
   Seal - Seal is a tool for simple, fast and secure end-to-end file tranfer.

USAGE:
   seal [global options] command [command options] [arguments...]

COMMANDS:
   send     Send file to the destination.
   recv     Receive file from the remote.
   key      Manage cached key infomation.
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

### Installing

MacOS

-  Download 'seal' binary from [release](https://github.com/WuFangyu/Seal/releases)
- `echo "alias seal=<path to your downloaded binary>" > ~/.bash_profile`
- `source ~/.bash_profile`

Or build the project yourself :arrow_up:

### Usage

#### send file

`seal send -file <path to the file you want to send> -dest <remote address> (optional)-pubckey <hex key string of remote address>`

#### recv file

`seal recv -dir <path to a local directory>`

#### public keys

- list all cached key pairs: `seal key -list`

- add a new key: `seal key -put -name <remote address> -pubkey <ecc key string in hex>`

- clear all cached key pairs: `seal key -clear`

- get key string: `seal key -get <remote address>`

- remove a key pair: `seal key -remove <remote address>`

- update a key pair: `seal key -update -name <remote address> -pubkey <ecc key string in hex>`

- generate a new public key for the host: `seal key -genkey`


### Workflow

<em>prerequist</em>: sender needs to obtain an authentic copy of the receiver's public key 
(then use `seal key -put -name <name> -pubkey <authentic copy>` to cache the key) - no MITM attack :wink:

sender generates a random AES key & encrypt the file -> send encrypted AES key (by receiver's ecc public key)
-> receiver decrypts AES key (by ecc private key) -> sender sends encrypted file chunks in parallel over multiple TCP channels -> receiver decrypts the file

AES Key 256 bits, ECC Key 256 bits, ~2.8x faster than scp

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
