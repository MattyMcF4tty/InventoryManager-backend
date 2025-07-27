package database

import (
	"log"
	"os"

	env "github.com/joho/godotenv"
	supabase "github.com/supabase-community/supabase-go"
)

func Connect() *supabase.Client {
	err := env.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SECRET_KEY")

	client, err := supabase.NewClient(supabaseUrl, supabaseKey, nil)

	if err != nil {
		log.Fatalf("Error creating Supabase client: %v", err)
	}
	return client
}
