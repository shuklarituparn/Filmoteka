DROP DATABASE filmotek;
DROP TABLE actor_movies;
DROP TABLE actors;
DROP TABLE movies;
CREATE DATABASE filmotek;
CREATE TABLE actors (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    gender TEXT,
    birth_date TEXT
);
CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    release_year BIGINT NOT NULL,
    genre TEXT NOT NULL,
    description VARCHAR(1000),
    rating TEXT NOT NULL
);
CREATE TABLE actor_movies (
    actor_id BIGINT,
    movie_id BIGINT,
    PRIMARY KEY (actor_id, movie_id)
);
ALTER TABLE actor_movies
ADD CONSTRAINT fk_actor_id FOREIGN KEY (actor_id) REFERENCES actors(id);
ALTER TABLE actor_movies
ADD CONSTRAINT fk_movie_id FOREIGN KEY (movie_id) REFERENCES movies(id);
INSERT INTO movies (title, release_year, description, genre, rating)
VALUES (
        'The Shawshank Redemption',
        1994,
        'Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.',
        'Drama',
        9.3
    ),
    (
        'The Godfather',
        1972,
        'The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.',
        'Crime',
        9.2
    ),
    (
        'The Dark Knight',
        2008,
        'When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.',
        'Action',
        9.0
    ),
    (
        '12 Angry Men',
        1957,
        'A jury holdout attempts to prevent a miscarriage of justice by forcing his colleagues to reconsider the evidence.',
        'Drama',
        8.9
    ),
    (
        'Schindler''s List',
        1993,
        'In German-occupied Poland during World War II, industrialist Oskar Schindler gradually becomes concerned for his Jewish workforce after witnessing their persecution by the Nazis.',
        'Biography',
        8.9
    ),
    (
        'The Lord of the Rings: The Return of the King',
        2003,
        'Gandalf and Aragorn lead the World of Men against Sauron''s army to draw his gaze from Frodo and Sam as they approach Mount Doom with the One Ring.',
        'Adventure',
        8.9
    ),
    (
        'Pulp Fiction',
        1994,
        'The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.',
        'Crime',
        8.9
    ),
    (
        'Fight Club',
        1999,
        'An insomniac office worker and a devil-may-care soapmaker form an underground fight club that evolves into something much, much more.',
        'Drama',
        8.8
    ),
    (
        'Forrest Gump',
        1994,
        'The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75, whose only desire is to be reunited with his childhood sweetheart.',
        'Drama',
        8.8
    ),
    (
        'Inception',
        2010,
        'A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O.',
        'Action',
        8.8
    );
INSERT INTO actors (first_name, last_name, gender, birth_date)
VALUES ('Tim', 'Robbins', 'male', '1958-10-16'),
    ('Robert', 'Downey', 'male', '1965-04-04'),
    ('Morgan', 'Freeman', 'male', '1937-06-01'),
    ('Henry', 'Fonda', 'male', '1905-05-16'),
    ('Marlon', 'Brando', 'male', '1924-04-03'),
    ('Al', 'Pacino', 'male', '1940-04-25'),
    ('Christian', 'Bale', 'male', '1974-01-30'),
    ('Liam', 'Neeson', 'male', '1952-06-07'),
    ('John', 'Travolta', 'male', '1954-02-18'),
    ('Tom', 'Hanks', 'male', '1956-07-09');
INSERT INTO actor_movies (movie_id, actor_id)
VALUES (1, 1),
    (1, 2),
    (1, 3),
    (2, 4),
    (3, 5),
    (3, 6),
    (4, 7),
    (5, 8),
    (6, 9),
    (7, 10);
