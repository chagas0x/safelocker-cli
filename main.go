package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

func encrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}

func main() {
	var rootCommand = &cobra.Command{}

	var encryptCmd = &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt your password AES-256",
		Long:  `This command encrypt your password and store to database.`,
		Run: func(cmd *cobra.Command, args []string) {
			bytes := make([]byte, 32) 
			if _, err := rand.Read(bytes); err != nil {
				panic(err.Error())
			}
			key := hex.EncodeToString(bytes) 
			readerName := bufio.NewReader(os.Stdin)
			fmt.Print("enter name for your password: ")
			name, err := readerName.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}

			name = strings.TrimSpace(name)

			fmt.Print("enter your password to encrypt: ")
			readerPass := bufio.NewReader(os.Stdin)
			password, err := readerPass.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}

			password = strings.TrimSpace(password)

			encrypted := encrypt(password, key)
			fmt.Printf("encrypted : %s\n", encrypted)

			currentTime := time.Now()
			dataEncrypted := currentTime.Format("2006-01-02")
			dbDriver := "mysql"
			dbUser := os.Getenv("DBUSER")
			dbPass := os.Getenv("DBPASS")
			dbName := os.Getenv("DBNAME")

			db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
			if err != nil {
				log.Fatal(err)
			}

			queryResponse, err := db.Exec("INSERT INTO `sf` (`id`, `userName`, `passEncrypted`, `randomKey`, `dateEncrypted`) VALUES (?,?,?,?,?)", "0", name, encrypted, key, dataEncrypted)
			if err != nil {
				log.Fatal(err)
			}

			if queryResponse != nil {
				fmt.Println("Stored!")
			} else {
				fmt.Println("Exit!")
			}
			defer db.Close()
		},
	}

	var decryptCmd = &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt your password AES-256",
		Long:  `This command decrypt your password stored in database.`,
		Run: func(cmd *cobra.Command, args []string) {
			dbDriver := "mysql"
			dbUser := os.Getenv("DBUSER")
			dbPass := os.Getenv("DBPASS")
			dbName := os.Getenv("DBNAME")

			db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
			if err != nil {
				log.Fatal(err)
			}

			type ResultList struct {
				ID     string
				Name   string
				Pass   string
				Random string
				Date   string
			}

			rows, err := db.Query(`SELECT * FROM sf`)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			var queryStore []ResultList
			var id string
			var name string
			var pass string
			var random string
			var date string

			for rows.Next() {
				err := rows.Scan(&id, &name, &pass, &random, &date)
				if err != nil {
					log.Fatal(err)
				}
				queryStore = append(queryStore, ResultList{ID: id, Name: name, Pass: pass, Random: random, Date: date})
			}

			fmt.Println("\nDatabase entries: ")
			fmt.Println(queryStore)
			fmt.Println("")
			readerID := bufio.NewReader(os.Stdin)
			fmt.Print("select ID to decrypt password: ")
			ID, err := readerID.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			ID = strings.TrimSpace(ID)

			var storePass string
			db.QueryRow("SELECT passEncrypted FROM sf WHERE id=?", ID).Scan(&storePass)

			var randomKey string
			db.QueryRow("SELECT randomKey FROM sf WHERE id=?", ID).Scan(&randomKey)

			fmt.Println("Encrypted pass:", storePass)

			decrypted := decrypt(storePass, randomKey)
			fmt.Printf("decrypted : %s\n", decrypted)

			fmt.Println("Decrypted!")
			defer db.Close()
		},
	}

	rootCommand.AddCommand(encryptCmd)
	rootCommand.AddCommand(decryptCmd)
	rootCommand.Execute()
}
