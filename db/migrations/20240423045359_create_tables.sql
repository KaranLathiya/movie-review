-- migrate:up

CREATE TYPE IF NOT EXISTS public.role AS ENUM (
	'ADMIN',
	'USER');

CREATE TABLE IF NOT EXISTS public.users (
	id VARCHAR NOT NULL DEFAULT unique_rowid(),
	email VARCHAR(320) NOT NULL,
	password VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp():::TIMESTAMP,
	updated_at TIMESTAMP NULL,
	first_name VARCHAR(255) NOT NULL,
	last_name VARCHAR(255) NOT NULL,
	"role" movie_review.public."role" NOT NULL DEFAULT 'USER':::movie_review.public."role",
	full_name VARCHAR(511) NULL AS (concat(first_name, ' ':::STRING, last_name)) STORED,
	CONSTRAINT pk_users_id PRIMARY KEY (id ASC),
	UNIQUE INDEX uk_users_email (email ASC)
);

CREATE TABLE IF NOT EXISTS public.movie (
	id VARCHAR NOT NULL DEFAULT unique_rowid(),
	title VARCHAR(255) NOT NULL,
	director_id VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp():::TIMESTAMP,
	updated_at TIMESTAMP NULL,
	average_rating FLOAT8 NULL,
	description VARCHAR(1000) NOT NULL,
	updated_by VARCHAR NULL,
	CONSTRAINT pk_movie_id PRIMARY KEY (id ASC),
	CONSTRAINT fk_movie_director_id FOREIGN KEY (director_id) REFERENCES public.users(id),
	CONSTRAINT fk_movie_updated_by FOREIGN KEY (updated_by) REFERENCES public.users(id),
	UNIQUE INDEX uc_movie_title (title ASC)
);

CREATE TABLE IF NOT EXISTS public.review (
	id VARCHAR NOT NULL DEFAULT unique_rowid(),
	movie_id VARCHAR NOT NULL,
	reviewer_id VARCHAR NOT NULL,
	rating FLOAT8 NOT NULL,
	comment VARCHAR(1000) NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp():::TIMESTAMP,
	updated_at TIMESTAMP NULL,
	comment_tsvector TSVECTOR NULL AS (to_tsvector('english':::STRING, comment)) STORED,
	CONSTRAINT pk_review_id PRIMARY KEY (id ASC),
	CONSTRAINT fk_review_movie_id FOREIGN KEY (movie_id) REFERENCES public.movie(id),
	CONSTRAINT fk_review_reviewer_id FOREIGN KEY (reviewer_id) REFERENCES public.users(id),
	UNIQUE INDEX uc_review_movie_id_and_reviewer_id (movie_id ASC, reviewer_id ASC),
	INVERTED INDEX review_comment_tsvector_idx (comment_tsvector)
);

-- migrate:down

DROP TABLE IF EXISTS public.review;

DROP TABLE IF EXISTS public.movie;

DROP TABLE IF EXISTS public.users;

DROP TYPE IF EXISTS public.role;