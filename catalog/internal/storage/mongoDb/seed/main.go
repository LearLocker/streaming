package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/LearLocker/streaming/catalog/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// === ИЗМЕНИ НА СВОЙ URI ===
	uri := "mongodb://catalog-service-user:catalog-service-password@localhost:27017/catalog-service?authSource=admin"
	// Если запускаешь изнутри Docker-сети: mongo-catalog вместо localhost

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("catalog-service")

	// Создаём коллекции
	collections := []string{"movies", "genres", "rankings", "reviews"}
	for _, coll := range collections {
		err = db.CreateCollection(ctx, coll)
		if err != nil && !mongo.IsDuplicateKeyError(err) {
			log.Printf("Warning creating %s: %v", coll, err)
		} else {
			fmt.Printf("✅ Коллекция %s готова\n", coll)
		}
	}

	// Заполняем данными...
	seedGenres(db, ctx)
	seedRankings(db, ctx)
	seedMovies(db, ctx)

	fmt.Println("🎉 Все тестовые данные успешно загружены!")
}

func seedGenres(db *mongo.Database, ctx context.Context) {
	genres := []interface{}{
		model.Genre{1, "Action"}, model.Genre{2, "Drama"}, model.Genre{3, "Comedy"},
		model.Genre{4, "Sci-Fi"}, model.Genre{5, "Thriller"}, model.Genre{6, "Horror"},
		model.Genre{7, "Fantasy"}, model.Genre{8, "Romance"}, model.Genre{9, "Crime"},
		model.Genre{10, "Adventure"},
	}
	_, err := db.Collection("genres").InsertMany(ctx, genres)
	if err != nil {
		log.Printf("Genres error: %v", err)
	} else {
		fmt.Println("✅ Жанры добавлены")
	}
}

func seedRankings(db *mongo.Database, ctx context.Context) {
	rankings := []interface{}{
		model.Ranking{1, "Masterpiece"}, model.Ranking{2, "Excellent"}, model.Ranking{3, "Good"},
		model.Ranking{4, "Average"}, model.Ranking{5, "Mediocre"}, model.Ranking{6, "Bad"},
	}
	_, err := db.Collection("rankings").InsertMany(ctx, rankings)
	if err != nil {
		log.Printf("Rankings error: %v", err)
	} else {
		fmt.Println("✅ Рейтинги добавлены")
	}
}

func seedMovies(db *mongo.Database, ctx context.Context) {
	movies := make([]interface{}, 0, 30)
	genreList := []model.Genre{
		{1, "Action"}, {2, "Drama"}, {3, "Comedy"}, {4, "Sci-Fi"},
		{5, "Thriller"}, {6, "Horror"}, {7, "Fantasy"},
	}

	for i := 1; i <= 30; i++ {
		// Случайные жанры (2-4)
		numGenres := rand.Intn(3) + 2
		selectedGenres := make([]model.Genre, 0, numGenres)
		for j := 0; j < numGenres; j++ {
			selectedGenres = append(selectedGenres, genreList[rand.Intn(len(genreList))])
		}

		movie := model.Movie{
			ImdbID:     fmt.Sprintf("tt%07d", 1000000+i),
			Title:      fmt.Sprintf("Test Movie %d: The Great Adventure", i),
			PosterPath: fmt.Sprintf("/posters/movie%d.jpg", i),
			YouTubeID:  "dQw4w9wgxcQ",
			Genre:      selectedGenres,
			Ranking:    model.Ranking{RankingValue: rand.Intn(6) + 1, RankingName: "Excellent"},
			Review:     generateRandomReviews(i),
		}
		movies = append(movies, movie)
	}

	_, err := db.Collection("movies").InsertMany(ctx, movies)
	if err != nil {
		log.Printf("Movies error: %v", err)
	} else {
		fmt.Printf("✅ Добавлено 30 фильмов\n")
	}
}

func generateRandomReviews(movieID int) []model.Review {
	count := rand.Intn(6)
	reviews := make([]model.Review, 0, count)
	for i := 0; i < count; i++ {
		reviews = append(reviews, model.Review{
			ReviewID: movieID*10 + i,
			AuthorID: 100 + rand.Intn(900),
			Text:     "Отличный фильм! Рекомендую всем посмотреть.",
		})
	}
	return reviews
}
