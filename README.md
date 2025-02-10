# safelocker-cli
SafeLocker é um software de gerenciamento de senhas criptografadas é como um cofre digital ultra seguro para todas as suas senhas, logins e informações confidenciais. Ele armazena suas credenciais de forma criptografada, protegendo-as de olhares curiosos e hackers maliciosos.

## Pré-requisitos
O software foi desenvolvido em ambiente Linux, mais precisamente na distribuição Arch Linux.

- MariaDB Server -> https://wiki.archlinux.org/title/MariaDB
- Golang -> https://go.dev/doc/install

## Instalação e configuração
### Database
```
CREATE DATABASE sf;
USE DATABASE sf;
CREATE TABLE `sf` (   `id` int(6) unsigned NOT NULL AUTO_INCREMENT,   `userName` varchar(100) NOT NULL,   `passEncrypted` varchar(300) NOT NULL, `randomKey` varchar(300) NOT NULL ,`dateEncrypted` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  PRIMARY KEY (`id`) ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
DESC table sf;
```

### Clone do repositório:
```
git clone https://github.com/0xc9e2/safelocker-cli.git
cd safelocker-cli
```

### Exportar variáveis de ambiente:
```
export DBUSER="database_user"
export DBPASS="database_pass"
export DBNAME="database_name"
```

### Instalação de dependência de código
```
go mod init sf
go mod tidy
go build sf
```
### Manual do CLI
```
./sf --help
```

### Encriptar senhas
```
./sf encrypt
```
### Decriptar senhas
```
./sf decrypt
```
## Disclaimer
O software SafeLocker é um software no modelo CLI desenvolvido na linguagem Golang. Ele é fornecido no estado atual, sem qualquer garantia expressa ou implícita de qualquer tipo.

## Author
Guilherme aka 0xc9e2
