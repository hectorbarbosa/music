CREATE TABLE IF NOT EXISTS public.songs(
    id SERIAL PRIMARY KEY,
    group_name varchar(50) NOT NULL,
    song_name varchar(200) NOT NULL,
    release_date date NOT NULL,
    song_text text NOT NULL,
    link text NOT NULL
);

CREATE INDEX group_name_idx on public.songs(group_name);
CREATE INDEX song_name_idx on public.songs(song_name);
CREATE INDEX release_date_idx on public.songs(release_date);
