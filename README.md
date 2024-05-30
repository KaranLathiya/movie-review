
# Movie review system

Movie review system for create/update/delete movies and movie reviews.

# Features
- Two different roles there one is Admin and User.
- Only Admin can create/update/delete movies.
- User can create/update/delete movie reviews.
- Can fetch movies and movie reviews.
- Can serach movies with different filters like movie title(full text search), movie director name and different sorting with recently created, movie title , movie average rating.
- Limiting in movie review. (maximum 3 revies can be added in 10 minutes)
- Can serach movie revies with comment(full text search),

# Tech Stack 
- GO 1.21
- CockroachDB 23.1
- Dbmate
- JWT (json web token)

## Run Locally

Prerequisites you need to set up on your local computer:

- [Golang](https://go.dev/doc/install)
- [Cockroach](https://www.cockroachlabs.com/docs/releases/)
- [Dbmate](https://github.com/amacneil/dbmate#installation)

1. Clone the project

```bash
  git clone https://github.com/KaranLathiya/movie-review.git
  cd movie-review
```

2. Copy the .env.example file to new .config/.env file and set env variables in .env:

```bash
  mkdir .config
  cp .env.example .config/.env
```

3. Create `.env` file in current directory and update below configurations:
   - Add Cockroach database URL in `DATABASE_URL` variable.
4. Run `dbmate up` to create database schema or Run `dbmate migrate` to migrate database schema.
5. Run `go run main.go` to run the programme.

## API Documentation:

After executing run command, open your favorite browser and type below URL to open graphQL playground.
```
http://localhost:8000/
```


