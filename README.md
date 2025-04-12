# gator – A Command-Line RSS Feed Aggregator Using PostgreSQL

This project is CLI built with GO to:

- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post
- RSS feeds are a way for websites to publish updates to their content. You can use this project to keep up with your favorite blogs, news sites, podcasts, and more!

It was developed using [goose](https://github.com/pressly/goose) for database migrations, `psql` as the terminal client, and [sqlc](https://sqlc.dev/) to generate type-safe ORM functions.

## Run the CLI

To run the CLI, it's necessary to have [Go](https://go.dev/) installed and a PostgreSQL database connection set up.

Create a `~/.gatorconfig.json` file in your home directory with your PostgreSQL connection URL and a `current_user_name` field, which will be used by the CLI commands.

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Clone the repository locally and access the repo:

```shell
git clone https://github.com/lealre/gator.git
```

Install the CLI:

```shell
go install
```

Run the CLI:

```shell
gator register 'your_username'
```

> [!NOTE]
> To run the commands without installing, use `go run . <command>`.
