package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/yaml.v2"
)

func AddNewUser() error {
	if len(*new_password) == 0 {
		return fmt.Errorf("Password not specified")
	}

	if len(*new_password) < 6 {
		return fmt.Errorf("Password too short")
	}

	if len(*new_user) == 0 {
		return fmt.Errorf("Username not specified")
	}

	//Check if username is already used
	_, ok := usermap[*new_user]
	if ok {
		return fmt.Errorf("Username already used")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*new_password), COST)
	if err != nil {
		panic(err)
	}

	users = append(users, User{
		Username: *new_user,
		Password: string(hash),
	})

	data, err := yaml.Marshal(users)

	file, err := os.OpenFile(*userdb, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err

}
